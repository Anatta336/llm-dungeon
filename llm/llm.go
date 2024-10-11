package llm

// `llm` is for general purpose interacting with the LLM.
// Think of it like the driver for a database. It doesn't have business logic.

import (
	"context"
	"fmt"
	"log"
	colour "samdriver/dungeon/log"

	"github.com/ollama/ollama/api"
)

var client *api.Client
var ctx context.Context

type Request struct {
	System      string  `json:"system"`
	User        string  `json:"user"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
}

type Response struct {
	Result string `json:"result"`
}

func init() {
	var err error
	client, err = api.ClientFromEnvironment()
	if err != nil {
		log.Fatalf("error creating LLM API client: %v", err)
	}

	ctx = context.Background()
}

func (request *Request) Process() (Response, error) {
	isStream := false
	chatRequest := &api.ChatRequest{
		Model:     request.Model,
		Messages:  request.messages(),
		Stream:    &isStream,
		KeepAlive: &api.Duration{Duration: 3},
		Options: map[string]interface{}{
			"temperature": request.Temperature,
			"num_predict": 1024,
		},
	}

	var combinedText = ""
	err := client.Chat(ctx, chatRequest, func(response api.ChatResponse) error {
		if response.PromptEvalCount > 0 {
			fmt.Println(colour.Grey, "Prompt eval tokens:", response.PromptEvalCount, colour.Reset)
		}
		if response.EvalCount > 0 {
			fmt.Println(colour.Grey, "Output tokens:", response.EvalCount, colour.Reset)
		}

		combinedText += response.Message.Content

		return nil
	})

	if err != nil {
		return Response{}, fmt.Errorf("error with LLM: %v", err)
	}

	return Response{
		Result: combinedText,
	}, nil
}

func (request *Request) messages() []api.Message {
	return []api.Message{
		{
			Role:    "system",
			Content: request.System,
		},
		{
			Role:    "user",
			Content: request.User,
		},
	}
}

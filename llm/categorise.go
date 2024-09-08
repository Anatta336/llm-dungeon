package llm

import (
	"fmt"
	"log"
	"strings"

	"github.com/ollama/ollama/api"
)

type CategorizedUserMessage struct {
	UserMessage
	Category string `json:"category"`
}

func (userMessage *UserMessage) Categorise(useShort bool) (*CategorizedUserMessage, error) {
	message := api.Message{
		Role:    "user",
		Content: userMessage.Content,
	}

	var model string
	if useShort {
		model = "input-decide-short"
	} else {
		model = "input-decide"
	}

	isStream := false
	request := &api.ChatRequest{
		Model:    model,
		Messages: []api.Message{message},
		Stream:   &isStream,
	}

	var combinedResponse string = ""
	err := client.Chat(ctx, request, func(response api.ChatResponse) error {
		combinedResponse += response.Message.Content

		if response.EvalCount > 0 {
			log.Println("Eval tokens:", response.EvalCount)
		}

		if response.PromptEvalCount > 0 {
			log.Println("Prompt eval tokens:", response.PromptEvalCount)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error chatting with LLM API: %v", err)
	}

	words := strings.Fields(combinedResponse)
	if len(words) == 0 {
		return nil, fmt.Errorf("empty response received from LLM API")
	}

	log.Println("Response:", combinedResponse)

	finalWord := words[len(words)-1]
	finalWord = strings.Trim(finalWord, "\".,!? ")

	var categorizedMessage *CategorizedUserMessage = &CategorizedUserMessage{
		UserMessage: *userMessage,
		Category:    finalWord,
	}

	return categorizedMessage, nil
}

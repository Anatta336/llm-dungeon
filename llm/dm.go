package llm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
)

func (userMessage *UserMessage) DmProcess() (*DmResponseMessage, error) {
	dmResponseMessage := &DmResponseMessage{
		UserInput: userMessage.Content,
	}
	err := dmResponseMessage.requestDescription()
	if err != nil {
		return dmResponseMessage, err
	}

	err = dmResponseMessage.requestEncodedActions()
	if err != nil {
		return dmResponseMessage, err
	}

	return dmResponseMessage, nil
}

func (dmResponse *DmResponseMessage) requestDescription() error {
	isStream := false
	request := &api.GenerateRequest{
		Model:  "adjudicate-001",
		Prompt: dmResponse.UserInput,
		Stream: &isStream,
		KeepAlive: &api.Duration{
			Duration: time.Minute * 3,
		},
	}

	var combinedResponse = ""
	err := client.Generate(ctx, request, func(response api.GenerateResponse) error {
		if response.EvalCount > 0 {
			log.Println("Adjudicate. Eval tokens:", response.EvalCount)
		}

		if response.PromptEvalCount > 0 {
			log.Println("Adjudicate. Prompt eval tokens:", response.PromptEvalCount)
		}

		combinedResponse += response.Response

		return nil
	})

	if err != nil {
		return fmt.Errorf("error with Adjudicate LLM: %v", err)
	}

	thinkingStart := strings.Index(combinedResponse, "<thinking>")
	thinkingEnd := strings.Index(combinedResponse, "</thinking>")
	thoughts := func() string {
		if thinkingStart == -1 || thinkingEnd == -1 {
			return ""
		}
		return combinedResponse[thinkingStart+len("<thinking>") : thinkingEnd]
	}()

	outputStart := strings.Index(combinedResponse, "<output>")
	outputEnd := strings.Index(combinedResponse, "</output>")
	description := func() string {
		if outputStart == -1 || outputEnd == -1 {
			// Assume everything after the </thinking> is the description we want.
			return combinedResponse[thinkingEnd+len("</thinking>"):]
		}
		return combinedResponse[outputStart+len("<output>") : outputEnd]
	}()

	dmResponse.RawAdjudicate = combinedResponse
	dmResponse.AdjudicateThoughts = thoughts
	dmResponse.Description = description

	return nil
}

func (dmResponse *DmResponseMessage) requestEncodedActions() error {
	isStream := false
	request := &api.GenerateRequest{
		Model:  "action-encode-001",
		Prompt: dmResponse.Description,
		Stream: &isStream,
		KeepAlive: &api.Duration{
			Duration: time.Minute * 3,
		},
	}

	var combinedResponse = ""
	err := client.Generate(ctx, request, func(response api.GenerateResponse) error {
		if response.EvalCount > 0 {
			log.Println("Action encode. Eval tokens:", response.EvalCount)
		}

		if response.PromptEvalCount > 0 {
			log.Println("Action encode. Prompt eval tokens:", response.PromptEvalCount)
		}

		combinedResponse += response.Response

		return nil
	})

	if err != nil {
		return fmt.Errorf("error with action encode LLM: %v", err)
	}

	thinkingStart := strings.Index(combinedResponse, "<thinking>")
	thinkingEnd := strings.Index(combinedResponse, "</thinking>")
	thoughts := func() string {
		if thinkingStart == -1 || thinkingEnd == -1 {
			return ""
		}
		return combinedResponse[thinkingStart+len("<thinking>") : thinkingEnd]
	}()

	outputStart := strings.Index(combinedResponse, "<output>")
	outputEnd := strings.Index(combinedResponse, "</output>")
	actions := func() string {
		if outputStart == -1 || outputEnd == -1 {
			return ""
		}
		return combinedResponse[outputStart+len("<output>") : outputEnd]
	}()

	dmResponse.RawActionEncode = combinedResponse
	dmResponse.AdjudicateThoughts = thoughts
	dmResponse.Actions = actions

	return nil
}

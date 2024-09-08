package llm

import (
	"context"
	"log"

	"github.com/ollama/ollama/api"
)

var client *api.Client
var ctx context.Context

func init() {
	var err error
	client, err = api.ClientFromEnvironment()
	if err != nil {
		log.Fatalf("error creating LLM API client: %v", err)
	}

	ctx = context.Background()
}

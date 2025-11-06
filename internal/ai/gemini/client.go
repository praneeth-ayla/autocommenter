package ai

import (
	"context"
	"log"

	"google.golang.org/genai"
)

func NewClient(ctx context.Context) *genai.Client {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatalf("failed to create Gemini client: %v", err)
	}
	return client
}

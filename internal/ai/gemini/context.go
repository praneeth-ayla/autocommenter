package gemini

import (
	"context"
	"encoding/json"

	"github.com/praneeth-ayla/autocommenter/internal/contextstore"
	"github.com/praneeth-ayla/autocommenter/internal/prompt"
	"github.com/praneeth-ayla/autocommenter/internal/scanner"
	"google.golang.org/genai"
)

func (g *GeminiProvider) GenerateContextBatch(files []scanner.Data) ([]contextstore.FileDetails, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Build all parts for the batch
	var parts []*genai.Part
	for _, f := range files {
		promptText := prompt.BuildFileContextPrompt(f.Path, f.Content)
		parts = append(parts, &genai.Part{Text: promptText})
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: "Follow the JSON schema exactly"},
			},
		},
		ResponseMIMEType:   "application/json",
		ResponseJsonSchema: GenerateContextBatchSchema, // Use the predefined schema for validation.
	}

	input := []*genai.Content{
		{Parts: parts},
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		// "gemini-2.5-flash-lite",
		input,
		config,
	)
	if err != nil {
		return nil, err
	}

	raw := result.Text()

	var parsed struct {
		Files []contextstore.FileDetails `json:"files"`
	}

	// Unmarshal the raw JSON response into the parsed struct.
	err = json.Unmarshal([]byte(raw), &parsed)
	if err != nil {
		return nil, err
	}

	return parsed.Files, nil
}

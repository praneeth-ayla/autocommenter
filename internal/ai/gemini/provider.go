package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/praneeth-ayla/AutoCommenter/internal/ai/providerutil"
	"github.com/praneeth-ayla/AutoCommenter/internal/contextstore"
	"github.com/praneeth-ayla/AutoCommenter/internal/prompt"
	"github.com/praneeth-ayla/AutoCommenter/internal/scanner"
	"google.golang.org/genai"
)

type GeminiProvider struct{}

func New() *GeminiProvider {
	return &GeminiProvider{}
}

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
		ResponseJsonSchema: GenerateContextBatchSchema,
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

	err = json.Unmarshal([]byte(raw), &parsed)
	if err != nil {
		return nil, err
	}

	return parsed.Files, nil
}

func (g *GeminiProvider) GenerateComments(content string, contexts []contextstore.FileDetails) (string, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", err
	}

	var ctxParts strings.Builder
	for _, c := range contexts {
		j, _ := json.Marshal(c)
		ctxParts.Write(j)
		ctxParts.WriteByte('\n')
	}
	promptText := prompt.BuildGenerateCommentsForFilesPrompt(content, ctxParts.String())

	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: prompt.SystemInstructionComments},
			},
		},
		ResponseMIMEType: "text/plain",
	}

	input := []*genai.Content{
		{Parts: []*genai.Part{{Text: promptText}}},
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-lite",
		input,
		config,
	)
	if err != nil {
		return "", err
	}

	out := result.Text()

	out = providerutil.StripCodeFences(out)
	out = providerutil.EnsurePackageLine(out, content)

	const maxCommentBlocks = 40
	out = providerutil.PruneExcessiveComments(out, maxCommentBlocks)

	// Validate that no non-comment code lines were changed
	if changed, diff := providerutil.NonCommentCodeChanged(content, out); changed {
		// Return an explicit error so caller can retry/generate again with adjusted prompt
		return "", fmt.Errorf("commenting aborted: non-comment code was modified by the model. example diff snippet:\n%s", diff)
	}

	return out, nil
}

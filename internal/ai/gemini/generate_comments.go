package ai

import (
	"context"
	"encoding/json"

	"google.golang.org/genai"
)

type CommentedFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type CommentedFilesResponse struct {
	Files []CommentedFile `json:"files"`
}

func GenerateCommentsForFiles(ctx context.Context, client *genai.Client, files []FileContent) (CommentedFilesResponse, error) {
	prompt := BuildGenerateCommentsForFilesPrompt(files)

	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: "You are an expert project documentor"}},
		},
		ResponseMIMEType:   "application/json",
		ResponseJsonSchema: GenerateCommentsForFilesSchema,
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), config)
	if err != nil {
		return CommentedFilesResponse{}, err
	}

	var parsed CommentedFilesResponse
	err = json.Unmarshal([]byte(result.Text()), &parsed)
	if err != nil {
		return CommentedFilesResponse{}, err
	}

	return parsed, nil
}

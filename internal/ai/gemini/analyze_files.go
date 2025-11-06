package ai

import (
	"context"
	"encoding/json"

	"google.golang.org/genai"
)

type FileResponse struct {
	Files []string `json:"files"`
}

func AnalyzeFilesForComments(ctx context.Context, client *genai.Client, files []string) (FileResponse, error) {
	prompt := BuildAnalyzeFilesForCommentsPrompt(files)

	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: "You are a project build expert."}},
		},
		ResponseMIMEType:   "application/json",
		ResponseJsonSchema: AnalyzeFilesForCommentsSchema,
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), config)
	if err != nil {
		return FileResponse{}, err
	}

	var parsed FileResponse
	if err := json.Unmarshal([]byte(result.Text()), &parsed); err != nil {
		return FileResponse{}, err
	}

	return parsed, nil
}

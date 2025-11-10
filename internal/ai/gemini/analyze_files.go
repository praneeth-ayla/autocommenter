package ai

import (
	"context"
	"encoding/json"

	"github.com/praneeth-ayla/AutoCommenter/internal/scanner"
	"google.golang.org/genai"
)

// FileResponse defines the structure for the AI's response when analyzing files.
type FileResponse struct {
	Files []string `json:"files"`
}

// AnalyzeFilesForComments sends a list of files to the AI to determine which ones need comments.
func AnalyzeFilesForComments(ctx context.Context, client *genai.Client, files []scanner.FileInfo) (FileResponse, error) {
	// Build the prompt for the AI based on the provided file information.
	prompt := BuildAnalyzeFilesForCommentsPrompt(files)

	// Configure the AI model's behavior, including system instruction and response format.
	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: "You are a project build expert."}},
		},
		ResponseMIMEType:   "application/json",
		ResponseJsonSchema: AnalyzeFilesForCommentsSchema,
	}

	// Send the prompt to the Gemini model and get the content generation result.
	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), config)
	if err != nil {
		return FileResponse{}, err
	}

	var parsed FileResponse
	// Unmarshal the JSON response from the AI into the FileResponse struct.
	if err := json.Unmarshal([]byte(result.Text()), &parsed); err != nil {
		return FileResponse{}, err
	}

	return parsed, nil
}

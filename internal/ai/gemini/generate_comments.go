package ai

import (
	"context"
	"encoding/json"

	"google.golang.org/genai"
)

// CommentedFile represents a single file with its path and updated content.
type CommentedFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// CommentedFilesResponse defines the structure for the AI's response containing commented files.
type CommentedFilesResponse struct {
	Files []CommentedFile `json:"files"`
}

// GenerateCommentsForFiles sends file content to the AI to generate comments.
func GenerateCommentsForFiles(ctx context.Context, client *genai.Client, files []FileContent) (CommentedFilesResponse, error) {
	// Build the prompt for the AI based on the provided file contents.
	prompt := BuildGenerateCommentsForFilesPrompt(files)

	// Configure the AI model's behavior, including system instruction and response format.
	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: "You are an expert project documentor"}},
		},
		ResponseMIMEType:   "application/json",
		ResponseJsonSchema: GenerateCommentsForFilesSchema,
	}

	// Send the prompt to the Gemini model and get the content generation result.
	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), config)
	if err != nil {
		return CommentedFilesResponse{}, err
	}

	var parsed CommentedFilesResponse
	// Unmarshal the JSON response from the AI into the CommentedFilesResponse struct.
	err = json.Unmarshal([]byte(result.Text()), &parsed)
	if err != nil {
		return CommentedFilesResponse{}, err
	}

	return parsed, nil
}

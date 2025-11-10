package ai

import (
	"fmt"
	"log"

	"github.com/alpkeskin/gotoon"
	"github.com/praneeth-ayla/AutoCommenter/internal/prompt"
	"github.com/praneeth-ayla/AutoCommenter/internal/scanner"
)

// BuildAnalyzeFilesForCommentsPrompt constructs the AI prompt to analyze which files need comments.
func BuildAnalyzeFilesForCommentsPrompt(files []scanner.FileInfo) string {
	data := map[string]interface{}{
		"files": files,
	}

	// Encode the file information into a structured string format for the AI prompt.
	encoded, err := gotoon.Encode(
		data,
		gotoon.WithIndent(0),       // no extra spaces
		gotoon.WithDelimiter("\t"), // tabs tokenize better
	)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf(prompt.AnalyzeFilesForComments, encoded)
}

// BuildGenerateCommentsForFilesPrompt constructs the AI prompt to generate comments for given files.
func BuildGenerateCommentsForFilesPrompt(files []FileContent) string {
	data := map[string]interface{}{
		"files": files,
	}

	// Encode the file content into a structured string format for the AI prompt.
	encoded, err := gotoon.Encode(
		data,
		gotoon.WithIndent(0),       // no extra spaces
		gotoon.WithDelimiter("\t"), // tabs tokenize better
	)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf(prompt.GenerateCommentsForFiles, encoded)
}

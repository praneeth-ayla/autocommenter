package ai

import (
	"fmt"
	"log"

	"github.com/alpkeskin/gotoon"
	"github.com/praneeth-ayla/AutoCommenter/internal/prompt"
)

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

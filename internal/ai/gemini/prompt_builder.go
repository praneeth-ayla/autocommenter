package ai

import (
	"fmt"
	"log"

	"github.com/alpkeskin/gotoon"
	"github.com/praneeth-ayla/AutoCommenter/internal/prompt"
	"github.com/praneeth-ayla/AutoCommenter/internal/scanner"
)

func BuildAnalyzeFilesForCommentsPrompt(files []scanner.FileInfo) string {
	data := map[string]interface{}{
		"files": files,
	}

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

func BuildGenerateCommentsForFilesPrompt(files []FileContent) string {
	data := map[string]interface{}{
		"files": files,
	}

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

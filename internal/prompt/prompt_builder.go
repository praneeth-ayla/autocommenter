package prompt

import (
	"fmt"
	"log"
	"strings"

	"github.com/alpkeskin/gotoon"
)

func BuildGenerateCommentsForFilesPrompt(content string, contextData string) string {
	data := map[string]interface{}{
		"content": content,
		"context": contextData,
	}

	encoded, err := gotoon.Encode(
		data,
		gotoon.WithIndent(0),
		gotoon.WithDelimiter("\t"),
	)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf(GenerateCommentsForFiles, encoded)
}

func BuildFileContextPrompt(path string, content string) string {
	var b strings.Builder

	b.WriteString("Return JSON for this file using the schema fields path, file_name, summary, exports, and imports.\n")
	b.WriteString("Identify exports and imports from the content.\n\n")

	b.WriteString("Path: ")
	b.WriteString(path)
	b.WriteString("\n\n")

	b.WriteString("Content:\n")
	b.WriteString(content)

	return b.String()
}

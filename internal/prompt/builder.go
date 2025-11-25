package prompt

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alpkeskin/gotoon"
	"github.com/praneeth-ayla/AutoCommenter/internal/contextstore"
)

func BuildFileContextPrompt(path string, content string) string {
	return fmt.Sprintf(TemplateFileContext, path, content)
}

func BuildCommentPrompt(content string, contextData string) (string, error) {
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
		return "", fmt.Errorf("prompt encoding failed: %w", err)
	}

	return fmt.Sprintf(TemplateCommentsFile, encoded), nil
}

func BuildFixesPrompt(original string, aiOutput string) string {
	return fmt.Sprintf(TemplateApplyFixes, original, aiOutput)
}

func BuildReadmePrompt(contexts []contextstore.FileDetails, existingReadme string) (string, error) {
	var sb strings.Builder

	for _, c := range contexts {
		j, err := json.Marshal(c)
		if err != nil {
			return "", fmt.Errorf("context marshal error: %w", err)
		}
		sb.Write(j)
		sb.WriteByte('\n')
	}

	contextStr := sb.String()
	readmeStr := strings.TrimSpace(existingReadme)

	return fmt.Sprintf(TemplateReadme, contextStr, readmeStr), nil
}

package prompt

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/alpkeskin/gotoon"
	"github.com/praneeth-ayla/autocommenter/internal/contextstore"
)

// BuildFileContextPrompt constructs a prompt for analyzing the context of a file.
func BuildFileContextPrompt(path string, content string) string {
	return fmt.Sprintf(TemplateFileContext, path, content) // Uses a predefined template to format the prompt.
}

// BuildCommentPrompt constructs a prompt for generating code comments based on style and context.
func BuildCommentPrompt(style string, content string, contextData string) (string, error) {
	data := map[string]interface{}{
		"content": content,     // The code content to comment on.
		"context": contextData, // Additional context data for the AI.
	}

	encoded, err := gotoon.Encode(
		data,
		gotoon.WithIndent(0),       // No indentation for the encoded JSON.
		gotoon.WithDelimiter("\t"), // Use tab as a delimiter.
	)
	if err != nil {
		return "", fmt.Errorf("prompt encoding failed: %w", err) // Error if JSON encoding fails.
	}

	switch style {
	case "minimalist":
		return fmt.Sprintf(TemplateMinimalist, encoded), nil // Formats prompt for minimalist style.
	case "explanatory":
		return fmt.Sprintf(TemplateExplanatory, encoded), nil // Formats prompt for explanatory style.
	case "detailed":
		return fmt.Sprintf(TemplateDetailed, encoded), nil // Formats prompt for detailed style.
	case "docstring":
		return fmt.Sprintf(TemplateDocstring, encoded), nil // Formats prompt for docstring style.
	case "inline-only":
		return fmt.Sprintf(TemplateInlineOnly, encoded), nil // Formats prompt for inline-only style.
	default:
		return "", errors.New("unknown style: supported styles are minimalist, explanatory, detailed, docstring, inline-only") // Handles unsupported styles.
	}
}

// BuildFixesPrompt constructs a prompt to apply AI-generated fixes to original code.
func BuildFixesPrompt(original string, aiOutput string) string {
	return fmt.Sprintf(TemplateApplyFixes, original, aiOutput) // Uses a template to combine original code and AI suggestions.
}

// BuildReadmePrompt constructs a prompt for generating a project README file.
func BuildReadmePrompt(contexts []contextstore.FileDetails, existingReadme string, fileTree string) (string, error) {
	var sb strings.Builder // Efficiently builds the string for contexts.

	for _, c := range contexts {
		j, err := json.Marshal(c) // Marshals each file detail into JSON.
		if err != nil {
			return "", fmt.Errorf("context marshal error: %w", err) // Returns error if marshaling fails.
		}
		sb.Write(j)
		sb.WriteByte('\n') // Adds a newline after each JSON object.
	}

	contextStr := sb.String()                      // Converts the builder content to a string.
	readmeStr := strings.TrimSpace(existingReadme) // Removes leading/trailing whitespace from existing README.
	treeStr := strings.TrimSpace(fileTree)         // Removes leading/trailing whitespace from file tree.

	return fmt.Sprintf(TemplateReadme, contextStr, treeStr, readmeStr), nil // Formats the README prompt.
}

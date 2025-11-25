package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"strings"

	"github.com/praneeth-ayla/AutoCommenter/internal/ai/providerutil"
	"github.com/praneeth-ayla/AutoCommenter/internal/contextstore"
	"github.com/praneeth-ayla/AutoCommenter/internal/prompt"
	"google.golang.org/genai"
)

func (g *GeminiProvider) GenerateComments(content string, contexts []contextstore.FileDetails) (string, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", err
	}

	var ctxParts strings.Builder
	// Marshal file details into JSON and append to the context string builder.
	for _, c := range contexts {
		j, _ := json.Marshal(c)
		ctxParts.Write(j)
		ctxParts.WriteByte('\n')
	}
	// Build the prompt for generating comments, including content and context.
	promptText, err := prompt.BuildCommentPrompt(content, ctxParts.String())
	if err != nil {
		return "", err
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: prompt.SystemInstructionComments}},
		},
		ResponseMIMEType: "text/plain",
	}

	input := []*genai.Content{{Parts: []*genai.Part{{Text: promptText}}}}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash-lite", input, config)
	if err != nil {
		return "", err
	}

	out := result.Text()
	out = providerutil.StripCodeFences(out)
	out = providerutil.EnsurePackageLine(out, content)

	const maxCommentBlocks = 40
	out = providerutil.PruneExcessiveComments(out, maxCommentBlocks)

	if changed, diff := providerutil.NonCommentCodeChanged(content, out); changed {
		updatedCode, _ := applyAIFixes(content, out)
		fmt.Printf("commenting aborted: non comment code changed. diff:\n%s", diff)
		return updatedCode, nil
	}

	return out, nil
}

func applyAIFixes(original string, aiOutput string) (string, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", err
	}

	promptText := prompt.BuildFixesPrompt(original, aiOutput)

	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: prompt.SystemInstructionFixes}},
		},
		ResponseMIMEType: "text/plain",
	}

	input := []*genai.Content{{Parts: []*genai.Part{{Text: promptText}}}}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-pro", input, config)
	if err != nil {
		return "", err
	}

	fixed := result.Text()
	fixed = providerutil.StripCodeFences(fixed)
	fixed = providerutil.EnsurePackageLine(fixed, original)

	fset := token.NewFileSet()
	if _, err := parser.ParseFile(fset, "", fixed, parser.AllErrors); err != nil {
		return "", fmt.Errorf("ai returned invalid Go source: %w", err)
	}

	return fixed, nil
}

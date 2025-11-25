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

func (g *GeminiProvider) GenerateComments(content string, contexts []contextstore.FileDetails, style string) (string, error) {
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
	promptText, err := prompt.BuildCommentPrompt(style, content, ctxParts.String())
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

	// If non-comment code changed in the AI output, attempt to fix it using applyAIFixes.
	if changed, _ := providerutil.NonCommentCodeChanged(content, out); changed {
		// Try to apply AI fixes, validating after each attempt that non-comment code
		// hasn't been altered. If we cannot produce a safe fix, return an error and
		// do not change the file.
		const maxFixAttempts = 2
		var lastErr error
		var fixed string

		for attempt := 1; attempt <= maxFixAttempts; attempt++ {
			fixed, lastErr = applyAIFixes(content, out)
			if lastErr != nil {
				// applyAIFixes already returns parse errors; we can retry if < attempts
				out = fixed // try next round with whatever AI returned (if any)
				continue
			}

			// ensure we got non-empty output
			if strings.TrimSpace(fixed) == "" {
				lastErr = fmt.Errorf("applyAIFixes returned empty output on attempt %d", attempt)
				out = fixed
				continue
			}

			// make sure fixes did not change non-comment code
			if changed2, _ := providerutil.NonCommentCodeChanged(content, fixed); !changed2 {
				// success: fixed comments only (or produced safe output)
				return fixed, nil
			}

			// still changes non-comment code; prepare for another attempt
			out = fixed
			lastErr = fmt.Errorf("non-comment code still changed after attempt %d", attempt)
		}

		// if we reach here, attempts exhausted and we couldn't safely fix the code
		if lastErr == nil {
			lastErr = fmt.Errorf("ai fixes failed and changed non-comment code")
		}
		return "", fmt.Errorf("ai fixes unsafe: %w", lastErr)
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

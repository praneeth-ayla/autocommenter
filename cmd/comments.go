/* Copyright © 2025 NAME HERE <EMAIL ADDRESS> */
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/googleapis/gax-go/v2/apierror"
	"github.com/praneeth-ayla/AutoCommenter/internal/ai"
	"github.com/praneeth-ayla/AutoCommenter/internal/contextstore"
	"github.com/praneeth-ayla/AutoCommenter/internal/scanner"
	"github.com/spf13/cobra"
	"google.golang.org/genai"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
)

const (
	defaultRetryDelay = 60 * time.Second
	maxRetryAttempts  = 3
)

var commentsCmd = &cobra.Command{
	Use:   "comments",
	Short: "Generate and manage code comments",
	Long: `Commands that help you scan the project and add missing comments
by using AI to document your code automatically.

Example:
  AutoCommenter comments gen
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'comments gen' to generate comments for your files")
		err := HitRateLimit()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("OK")
		}
	},
}

var genCommentsCmd = &cobra.Command{
	Use:   "gen",
	Short: "Add comments to code files that need them",
	Long: `Scan the project and find files without proper comments.
Use AI to generate comments and write them back into the files.
Retries automatically if API rate limits occur.

Example:
  AutoCommenter comments gen
`,
	RunE: runGenerateComments,
}

func init() {
	rootCmd.AddCommand(commentsCmd)
	commentsCmd.AddCommand(genCommentsCmd)
}

// runGenerateComments is the main entry point for comment generation
func runGenerateComments(cmd *cobra.Command, args []string) error {
	// Initialize
	rootPath := scanner.GetProjectRoot()
	provider := ai.NewProvider("gemini")

	// Scan project
	fmt.Println("📁 Scanning project files...")
	files, err := scanner.Scan(rootPath)
	if err != nil {
		return fmt.Errorf("failed to scan files: %w", err)
	}

	filteredFiles := scanner.FilterFilesNeedingComments(files)
	if len(filteredFiles) == 0 {
		fmt.Println("✅ No files need comments!")
		return nil
	}

	fmt.Printf("📝 Found %d files needing comments\n", len(filteredFiles))

	// Load context
	fmt.Println("🔍 Loading project context...")
	context, err := contextstore.Load()
	if err != nil {
		return fmt.Errorf("failed to load context: %w", err)
	}
	allCtxSlice := contextstore.MapToSlice(context)

	// Process files
	fmt.Println("🤖 Generating comments...")
	successCount, errorCount := 0, 0

	for i, file := range filteredFiles {
		fmt.Printf("\n[%d/%d] Processing: %s\n", i+1, len(filteredFiles), file.Path)

		if err := processFile(file, provider, allCtxSlice); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			errorCount++
		} else {
			fmt.Printf("✅ Updated successfully\n")
			successCount++
		}
	}

	// Summary
	fmt.Println("\n" + strings.Repeat("─", 50) + "\n")
	fmt.Printf("📊 Summary: %d succeeded, %d failed\n", successCount, errorCount)

	if errorCount > 0 {
		return fmt.Errorf("completed with %d errors", errorCount)
	}

	return nil
}

// processFile handles comment generation for a single file with automatic retries
func processFile(file scanner.Info, provider ai.Provider, context []contextstore.FileDetails) error {
	fd := scanner.LoadSingle(file)
	var lastErr error

	for attempt := 1; attempt <= maxRetryAttempts; attempt++ {
		commented, err := provider.GenerateComments(fd.Content, context)
		if err == nil {
			// Success - write file
			return os.WriteFile(file.Path, []byte(commented), 0644)
		}

		lastErr = err

		// Check if this is a retryable rate limit error
		if retryDelay, isRateLimit := checkRateLimitError(err); isRateLimit {
			if attempt < maxRetryAttempts {
				fmt.Printf("⏳ Rate limit hit. Waiting %v before retry %d/%d...\n",
					retryDelay, attempt+1, maxRetryAttempts)
				time.Sleep(retryDelay)
				continue
			}
			return fmt.Errorf("rate limit exceeded after %d attempts", maxRetryAttempts)
		}

		// Non-retryable error
		return fmt.Errorf("generation failed: %w", err)
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// checkRateLimitError determines if an error is a rate limit error and extracts retry delay
func checkRateLimitError(err error) (retryDelay time.Duration, isRateLimit bool) {
	// Check error string for rate limit indicators
	errStr := err.Error()
	if strings.Contains(errStr, "RESOURCE_EXHAUSTED") ||
		strings.Contains(errStr, "429") ||
		strings.Contains(errStr, "Quota exceeded") {
		return extractRetryDelay(errStr), true
	}

	// Check structured API error
	var apiErr *apierror.APIError
	if errors.As(err, &apiErr) {
		if status := apiErr.GRPCStatus(); status != nil && status.Code() == codes.ResourceExhausted {
			// Try to extract RetryInfo from status details
			for _, detail := range status.Details() {
				if retryInfo, ok := detail.(*errdetails.RetryInfo); ok {
					if retryInfo.RetryDelay != nil {
						if d := retryInfo.RetryDelay.AsDuration(); d > 0 {
							return d, true
						}
					}
				}
			}
			// Rate limit confirmed but no retry info - use default
			return defaultRetryDelay, true
		}
	}

	return 0, false
}

// extractRetryDelay parses retry delay from error message
func extractRetryDelay(errStr string) time.Duration {
	// Match patterns like "retry in 1.5s" or "Please retry in 2.5s"
	re := regexp.MustCompile(`retry in ([0-9.]+)s`)
	if matches := re.FindStringSubmatch(errStr); len(matches) > 1 {
		if seconds, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return time.Duration(seconds * float64(time.Second))
		}
	}
	return defaultRetryDelay
}

func HitRateLimit() error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return err
	}

	// small prompt enough to trigger request but not heavy
	input := []*genai.Content{
		{Parts: []*genai.Part{{Text: "hi"}}},
	}

	for i := 0; i < 50; i++ { // spam 50 times
		fmt.Println("Request", i+1)

		_, err := client.Models.GenerateContent(
			ctx,
			"gemini-2.5-flash-lite",
			input,
			nil,
		)

		if err != nil {
			fmt.Println("Error on request", i+1, err)
			return err
		}

		fmt.Println("OK")
	}

	return nil
}

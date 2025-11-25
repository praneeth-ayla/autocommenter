package cmd

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/googleapis/gax-go/v2/apierror"
	"github.com/praneeth-ayla/AutoCommenter/internal/ai"
	"github.com/praneeth-ayla/AutoCommenter/internal/config"
	"github.com/praneeth-ayla/AutoCommenter/internal/contextstore"
	"github.com/praneeth-ayla/AutoCommenter/internal/prompt"
	"github.com/praneeth-ayla/AutoCommenter/internal/scanner"
	"github.com/praneeth-ayla/AutoCommenter/internal/ui"
	"github.com/spf13/cobra"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
)

const (
	defaultRetryDelay = 5 * time.Second
	maxRetryAttempts  = 3
	perRequestTimeout = 60 * time.Second
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
	},
}

var genCommentsCmd = &cobra.Command{
	Use:   "gen",
	Short: "Add comments to code files that need them",
	RunE:  runGenerateComments,
}

func init() {
	genCommentsCmd.SilenceUsage = true
	genCommentsCmd.SilenceErrors = true

	rootCmd.AddCommand(commentsCmd)
	commentsCmd.AddCommand(genCommentsCmd)
}

func runGenerateComments(cmd *cobra.Command, args []string) error {
	providerName, _ := config.GetProvider()
	provider, err := ai.NewProvider(providerName)
	if err != nil {
		return fmt.Errorf("provider init: %w", err)
	}

	commentStyle, err := ui.SelectOne("Select comment style:", prompt.Styles)
	if err != nil {
		return err
	}
	_ = commentStyle

	rootPath := scanner.GetProjectRoot()
	fmt.Println("Scanning project files...")
	files, err := scanner.Scan(rootPath)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	filteredFiles := scanner.FilterFilesNeedingComments(files)
	if len(filteredFiles) == 0 {
		fmt.Println("No files need comments")
		return nil
	}

	fmt.Printf("Found %d files needing comments\n", len(filteredFiles))

	fmt.Println("Loading project context...")
	ctxMap, err := contextstore.Load()
	if err != nil {
		return fmt.Errorf("load context: %w", err)
	}
	allCtxSlice := contextstore.MapToSlice(ctxMap)

	fmt.Println("Generating comments (this may take a while)...")
	successCount, errorCount := 0, 0

	for i, file := range filteredFiles {
		fmt.Printf("\n[%d/%d] %s\n", i+1, len(filteredFiles), file.Path)
		if err := processFile(file, provider, allCtxSlice); err != nil {
			fmt.Printf("  ✖ error: %v\n", err)
			errorCount++
		} else {
			fmt.Println("  ✓ updated")
			successCount++
		}
	}

	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Printf("Summary: %d succeeded, %d failed\n", successCount, errorCount)

	if errorCount > 0 {
		return fmt.Errorf("completed with %d errors", errorCount)
	}
	return nil
}

func processFile(file scanner.Info, provider ai.Provider, context []contextstore.FileDetails) error {
	fd := scanner.LoadSingle(file)
	var lastErr error

	for attempt := 1; attempt <= maxRetryAttempts; attempt++ {
		attemptTimeout := time.Duration(perRequestTimeout)
		ctx, cancel := contextWithTimeout(attemptTimeout)
		defer cancel()

		done := make(chan struct{})
		var commented string
		var err error

		go func() {
			commented, err = provider.GenerateComments(fd.Content, context)
			close(done)
		}()

		select {
		case <-ctx.Done():
			lastErr = fmt.Errorf("request timed out after %s", attemptTimeout)
		case <-done:
			if err == nil {
				return scanner.WriteFile(file.Path, commented)
			}
			lastErr = err
		}

		if delay, isRateLimit := checkRateLimitError(lastErr); isRateLimit {
			if attempt < maxRetryAttempts {
				sleepWithJitter(delay)
				continue
			}
			return fmt.Errorf("rate limit after %d attempts: %w", maxRetryAttempts, lastErr)
		}

		// Non-rate-limit error -> fail fast
		return fmt.Errorf("generation failed: %w", lastErr)
	}

	return fmt.Errorf("retries exhausted: %w", lastErr)
}

func contextWithTimeout(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}

func sleepWithJitter(base time.Duration) {
	j := time.Duration(rand.Int63n(int64(base / 2))) // jitter up to half of base
	time.Sleep(base + j)
}

func checkRateLimitError(err error) (time.Duration, bool) {
	if err == nil {
		return 0, false
	}
	errStr := err.Error()
	if strings.Contains(errStr, "RESOURCE_EXHAUSTED") ||
		strings.Contains(errStr, "429") ||
		strings.Contains(errStr, "Quota exceeded") {
		return extractRetryDelay(errStr), true
	}

	var apiErr *apierror.APIError
	if errors.As(err, &apiErr) {
		if status := apiErr.GRPCStatus(); status != nil && status.Code() == codes.ResourceExhausted {
			for _, detail := range status.Details() {
				if retryInfo, ok := detail.(*errdetails.RetryInfo); ok {
					if retryInfo.RetryDelay != nil {
						if d := retryInfo.RetryDelay.AsDuration(); d > 0 {
							return d, true
						}
					}
				}
			}
			return defaultRetryDelay, true
		}
	}

	return 0, false
}

func extractRetryDelay(errStr string) time.Duration {
	re := regexp.MustCompile(`retry in ([0-9.]+)s`)
	if matches := re.FindStringSubmatch(errStr); len(matches) > 1 {
		if seconds, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return time.Duration(seconds * float64(time.Second))
		}
	}
	return defaultRetryDelay
}

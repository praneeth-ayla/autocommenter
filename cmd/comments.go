package cmd

import (
	"fmt"
	"strings"

	"github.com/praneeth-ayla/AutoCommenter/internal/ai"
	"github.com/praneeth-ayla/AutoCommenter/internal/ai/providerutil"
	"github.com/praneeth-ayla/AutoCommenter/internal/config"
	"github.com/praneeth-ayla/AutoCommenter/internal/contextstore"
	"github.com/praneeth-ayla/AutoCommenter/internal/prompt"
	"github.com/praneeth-ayla/AutoCommenter/internal/scanner"
	"github.com/praneeth-ayla/AutoCommenter/internal/ui"
	"github.com/spf13/cobra"
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

func processFile(file scanner.Info, provider ai.Provider, ctx []contextstore.FileDetails) error {
	fd := scanner.LoadSingle(file)

	commented, err := providerutil.DoWithRetry[string](
		providerutil.MaxRetryAttempts,
		providerutil.PerRequestTimeout,
		func() (string, error) {
			return provider.GenerateComments(fd.Content, ctx)
		},
	)
	if err != nil {
		return err
	}

	return scanner.WriteFile(file.Path, commented)
}

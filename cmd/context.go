/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/praneeth-ayla/AutoCommenter/internal/ai"
	"github.com/praneeth-ayla/AutoCommenter/internal/contextstore"
	"github.com/praneeth-ayla/AutoCommenter/internal/scanner"
	"github.com/spf13/cobra"
)

// contextCmd represents the context command
var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Work with code context",
	Long: `This command helps you manage context data for your project.
Context helps AutoCommenter understand your code better.

To generate context run:

  AutoCommenter context gen

This will scan your project and build a context database.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run 'AutoCommenter context gen' to generate context")
	},
}

func init() {
	rootCmd.AddCommand(contextCmd)
	contextCmd.AddCommand(contextGenCmd)
}

var contextGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate context for your project",
	Long: `Scan your project and build context info for all supported files.

Usage example:
  AutoCommenter context gen

After generation you can run comment commands that use this context.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Generating Context")

		rootPath := getProjectRoot()
		files, err := scanner.Scan(rootPath)
		if err != nil {
			fmt.Println("scan error:", err)
			return fmt.Errorf("scan failed: %w", err)
		}

		fmt.Printf("File path: %v\n", rootPath)
		if len(files) == 0 {
			fmt.Println("No files found for context generation")
			return nil
		}

		provider := ai.NewProvider("gemini")
		batches := scanner.BatchByLines(files, 500)
		allContext := make(map[string]contextstore.FileDetails)

		for _, batch := range batches {
			batchData := scanner.Load(batch)

			ctx, err := provider.GenerateContextBatch(batchData)
			if err != nil {
				fmt.Println("context batch error:", err)
				continue
			}

			for _, item := range ctx {
				allContext[item.Path] = item
			}
		}

		if len(allContext) == 0 {
			fmt.Println("No context generated after processing batches")
			return nil
		}

		if err := contextstore.Save(allContext); err != nil {
			fmt.Println("save error:", err)
			return fmt.Errorf("context save failed: %w", err)
		}

		fmt.Println("context generation completed")
		return nil
	},
}

func getProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}

	for {
		// check if go.mod exists here
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)

		// reached filesystem root
		if parent == dir {
			return "." // fallback
		}

		dir = parent
	}
}

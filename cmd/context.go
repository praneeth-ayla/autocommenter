/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"sync"

	"github.com/praneeth-ayla/AutoCommenter/internal/ai"
	"github.com/praneeth-ayla/AutoCommenter/internal/contextstore"
	"github.com/praneeth-ayla/AutoCommenter/internal/scanner"
	"github.com/spf13/cobra"
)

// contextCmd represents the context command
var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Manage project context",
	Long: `Work with the stored context that helps generate accurate comments and readme doc.
Use this command to scan the project and collect useful information.

Example:
  AutoCommenter context gen
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run 'AutoCommenter context gen' to generate context")
	},
}

var contextGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Scan project and build context",
	Long: `Scan supported files in the project and store context data.
This improves the quality of generated comments and readme later.

Example:
  AutoCommenter context gen
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Generating Context")

		rootPath := scanner.GetProjectRoot()
		files, err := scanner.Scan(rootPath)
		if err != nil {
			fmt.Println("scan error:", err)
			return fmt.Errorf("scan failed: %w", err)
		}

		if len(files) == 0 {
			fmt.Println("No files found for context generation")
			return nil
		}

		provider := ai.NewProvider("gemini")
		batches := scanner.BatchByLines(files, 500)
		allContext := make(map[string]contextstore.FileDetails)

		var wg sync.WaitGroup
		var mu sync.Mutex

		for _, batch := range batches {
			wg.Add(1)
			go func(b []scanner.Info) {
				defer wg.Done()

				batchData := scanner.Load(b)
				ctx, err := provider.GenerateContextBatch(batchData)
				if err != nil {
					fmt.Println("context batch error:", err)
					return
				}

				mu.Lock()
				for _, item := range ctx {
					allContext[item.Path] = item
				}
				mu.Unlock()
			}(batch)
		}
		wg.Wait()

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

func init() {
	rootCmd.AddCommand(contextCmd)
	contextCmd.AddCommand(contextGenCmd)
}

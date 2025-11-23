package cmd

import (
	"fmt"
	"os"

	"github.com/praneeth-ayla/AutoCommenter/internal/ai"
	"github.com/praneeth-ayla/AutoCommenter/internal/contextstore"
	"github.com/praneeth-ayla/AutoCommenter/internal/scanner"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan project, generate context and apply comments",
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Println("Scanning project files")

		files, err := scanner.Scan(".")
		if err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}

		candidates := scanner.FilterFilesNeedingComments(files)
		if len(candidates) == 0 {
			fmt.Println("No files need comments")
			return nil
		}

		provider := ai.NewProvider("gemini")
		allContext := make(map[string]contextstore.FileDetails)

		batches := scanner.BatchByLines(candidates, 500)

		for _, batch := range batches {
			batchData := scanner.Load(batch)

			ctx, err := provider.GenerateContextBatch(batchData)
			if err != nil {
				fmt.Println("context batch error", err)
				continue
			}

			for _, item := range ctx {
				allContext[item.Path] = item
			}
		}

		if err := contextstore.Save(allContext); err != nil {
			return fmt.Errorf("context save failed: %w", err)
		}

		fmt.Println("context generation completed")

		allCtxSlice := contextstore.MapToSlice(allContext)

		for _, file := range candidates {
			fd := scanner.LoadSingle(file)
			commented, err := provider.GenerateComments(fd.Content, allCtxSlice)
			if err != nil {
				fmt.Println("comment generation error for", file.Path, err)
				continue
			}

			if err := os.WriteFile(file.Path, []byte(commented), 0644); err != nil {
				fmt.Println("file update failed for", file.Path, err)
				continue
			}

			fmt.Println("Updated", file.Path)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}

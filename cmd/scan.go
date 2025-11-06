/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	ai "github.com/praneeth-ayla/AutoCommenter/internal/ai/gemini"
	"github.com/praneeth-ayla/AutoCommenter/internal/scanner"
	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan and list files needing comments",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		files, err := scanner.Scanner(".")
		if err != nil {
			return fmt.Errorf("failed to scan: %w", err)
		}

		client := ai.NewClient(ctx)
		response, err := ai.AnalyzeFilesForComments(ctx, client, files)
		if err != nil {
			return fmt.Errorf("AI processing failed: %w", err)
		}

		for i, f := range response.Files {
			fmt.Println(f, i)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}

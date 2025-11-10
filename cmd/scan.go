/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

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

		files, err := scanner.Scanner("../fed-fresh")
		if err != nil {
			return fmt.Errorf("failed to scan: %w", err)
		}

		fmt.Println("Scanning files...")
		client := ai.NewClient(ctx)
		response, err := ai.AnalyzeFilesForComments(ctx, client, files)
		if err != nil {
			return fmt.Errorf("AI processing failed: %w", err)
		}

		fmt.Println("Scanning files Content...")
		fileContents := ai.ReadFileContent(response)

		commentedFiles, err := ai.GenerateCommentsForFiles(ctx, client, fileContents)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Updating files Content...")
		for _, f := range commentedFiles.Files {
			updateFiles(f)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}

func updateFiles(f ai.CommentedFile) {

	err := ioutil.WriteFile(f.Path, []byte(f.Content), 0644) // 0644 sets file permissions
	if err != nil {
		log.Fatal(err)
	}

	log.Println("File written successfully:", f.Path)

}

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

// scanCmd represents the scan command, used to scan files and generate comments.
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan and list files needing comments",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Scan the current directory for files.
		files, err := scanner.Scanner(".")
		if err != nil {
			return fmt.Errorf("failed to scan: %w", err)
		}

		fmt.Println("Scanning files...")
		// Initialize a new Gemini AI client.
		client := ai.NewClient(ctx)
		// Ask the AI to analyze which files require comments.
		response, err := ai.AnalyzeFilesForComments(ctx, client, files)
		if err != nil {
			return fmt.Errorf("AI processing failed: %w", err)
		}

		fmt.Println("Scanning files Content...")
		// Read the content of the identified files.
		fileContents := ai.ReadFileContent(response)

		// Generate comments for the files using the AI.
		commentedFiles, err := ai.GenerateCommentsForFiles(ctx, client, fileContents)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Updating files Content...")
		// Iterate and update each file with the new commented content.
		for _, f := range commentedFiles.Files {
			updateFiles(f)
		}

		return nil
	},
}

func init() {
	// Add the scan command as a subcommand to the root command.
	rootCmd.AddCommand(scanCmd)
}

// updateFiles writes the new content to the specified file path.
func updateFiles(f ai.CommentedFile) {

	err := ioutil.WriteFile(f.Path, []byte(f.Content), 0644) // 0644 sets file permissions
	if err != nil {
		log.Fatal(err)
	}

	log.Println("File written successfully:", f.Path)

}

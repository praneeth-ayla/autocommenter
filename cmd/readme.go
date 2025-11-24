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

// readmeCmd represents the readme command
var readmeCmd = &cobra.Command{
	Use:   "readme",
	Short: "Manage project README file",
	Long: `Generate or update README.md based on the scanned project context
and any existing README found in the project.

Example:
  AutoCommenter readme gen
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'readme gen' to generate readme.md for your project")
	},
}

var genReadmeCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate or update README.md",
	Long: `Create or replace the README.md file in the project root.
This uses collected project context and AI generation.

Actions:
  1. Load stored project context
  2. Merge with existing README if available
  3. Write a new README.md

Example:
  AutoCommenter readme gen
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		rootPath := scanner.GetProjectRoot()
		provider := ai.NewProvider("gemini")

		fmt.Println("Loading project context...")
		contextData, err := contextstore.Load()
		if err != nil {
			return fmt.Errorf("failed to load context: %w", err)
		}
		allCtxSlice := contextstore.MapToSlice(contextData)

		// check existing README file
		var existingReadme string
		readmePaths := []string{
			filepath.Join(rootPath, "README.md"),
			filepath.Join(rootPath, "readme.md"),
		}

		for _, path := range readmePaths {
			if data, err := os.ReadFile(path); err == nil {
				existingReadme = string(data)
				fmt.Println("Existing README found:", path)
				break
			}
		}

		fmt.Println("️Generating README...")
		newReadme, err := provider.GenerateReadme(allCtxSlice, existingReadme)
		if err != nil {
			return fmt.Errorf("README generation failed: %w", err)
		}

		outputPath := filepath.Join(rootPath, "README.md")

		err = scanner.WriteFile(outputPath, newReadme)
		if err != nil {
			return fmt.Errorf("failed to write README.md: %w", err)
		}

		fmt.Println(";) README.md updated:", outputPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(readmeCmd)
	readmeCmd.AddCommand(genReadmeCmd)
}

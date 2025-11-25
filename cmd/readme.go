/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/praneeth-ayla/autocommenter/internal/ai"
	"github.com/praneeth-ayla/autocommenter/internal/config"
	"github.com/praneeth-ayla/autocommenter/internal/contextstore"
	"github.com/praneeth-ayla/autocommenter/internal/scanner"
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

var (
	readmePath string // Flag for custom README path
)

var genReadmeCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate or update README.md",
	Long: `Create or replace the README.md file in the project root.
This uses collected project context and AI generation.

Actions:
  1. Load stored project context
  2. Merge with existing README if available
  3. Write a new README.md

Examples:
  AutoCommenter readme gen
  AutoCommenter readme gen --path docs/README.md
  AutoCommenter readme gen -p ./documentation/README.md
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		providerName, _ := config.GetProvider()
		provider, err := ai.NewProvider(providerName)
		if err != nil {
			fmt.Println("provider error:", err)
			return err
		}

		rootPath := scanner.GetProjectRoot()
		fmt.Println("Loading project context...")
		contextData, err := contextstore.Load()
		if err != nil {
			return fmt.Errorf("no project context found. Run: AutoCommenter context gen")
		}

		allCtxSlice := contextstore.MapToSlice(contextData)

		// Determine output path
		outputPath := filepath.Join(rootPath, "README.md")
		if readmePath != "" {
			outputPath = readmePath
			// If it's a relative path, make it absolute relative to project root
			if !filepath.IsAbs(outputPath) {
				outputPath = filepath.Join(rootPath, outputPath)
			}
		}

		// Ensure directory exists
		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
		}

		// check existing README file (check both default location and custom path)
		var existingReadme string
		readmePaths := []string{
			outputPath, // Check the target path first
			filepath.Join(rootPath, "README.md"),
			filepath.Join(rootPath, "readme.md"),
		}

		// get custom path from user using --path or -p flag
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

		err = scanner.WriteFile(outputPath, newReadme)
		if err != nil {
			return fmt.Errorf("failed to write README.md: %w", err)
		}

		fmt.Println(";) README.md updated:", outputPath)
		return nil
	},
}

func init() {
	genReadmeCmd.SilenceUsage = true
	// genReadmeCmd.SilenceErrors = true

	// Add path flag
	genReadmeCmd.Flags().StringVarP(&readmePath, "path", "p", "", "Custom path for README file (default: ./README.md)")

	rootCmd.AddCommand(readmeCmd)
	readmeCmd.AddCommand(genReadmeCmd)
}

/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/praneeth-ayla/AutoCommenter/internal/ai"
	"github.com/praneeth-ayla/AutoCommenter/internal/contextstore"
	"github.com/praneeth-ayla/AutoCommenter/internal/scanner"
	"github.com/spf13/cobra"
)

// commentsCmd represents the comments command
var commentsCmd = &cobra.Command{
	Use:   "comments",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("comments called")
	},
}

// init initializes the command tree by adding subcommands.
func init() {
	rootCmd.AddCommand(commentsCmd)
	commentsCmd.AddCommand(genCommentsCmd)
}

var genCommentsCmd = &cobra.Command{
	Use:   "gen",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		rootPath := getProjectRoot()
		fmt.Println("Scanning for files")
		files, err := scanner.Scan(rootPath)
		if err != nil {
			return err
		}

		fmt.Println("Filtering files")
		filteredFiles := scanner.FilterFilesNeedingComments(files)

		fmt.Println("Loading project context")
		context, err := contextstore.Load()
		if err != nil {
			return err
		}

		provider := ai.NewProvider("gemini")
		allCtxSlice := contextstore.MapToSlice(context)

		fmt.Println("Comment generating")
		for _, file := range filteredFiles {
			fd := scanner.LoadSingle(file)

			// Single attempt only. If error occurs, return it immediately.
			commented, err := provider.GenerateComments(fd.Content, allCtxSlice)
			if err != nil {
				return fmt.Errorf("comment generation error for %s: %w", file.Path, err)
			}

			if err := os.WriteFile(file.Path, []byte(commented), 0644); err != nil {
				return fmt.Errorf("file update failed for %s: %w", file.Path, err)
			}

			fmt.Printf("Updated %s\n", file.Path)
		}

		return nil
	},
}

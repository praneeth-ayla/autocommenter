/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/praneeth-ayla/AutoCommenter/internal/scanner"
	"github.com/spf13/cobra"
)

// scanCmd represents the scan command, used to scan files and generate comments.
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan and list files needing comments",
	RunE: func(cmd *cobra.Command, args []string) error {

		files, err := scanner.Scanner(".")
		if err != nil {
			return fmt.Errorf("failed to scan: %w", err)
		}
		filteredFiles := scanner.FilterFilesNeedingComments(files)

		fmt.Println("Total Files: ", len(files))
		fmt.Println("Scanning files...")

		for i, file := range filteredFiles {
			fmt.Println(i, file)
		}

		return nil
	},
}

func init() {
	// Add the scan command as a subcommand to the root command.
	rootCmd.AddCommand(scanCmd)
}

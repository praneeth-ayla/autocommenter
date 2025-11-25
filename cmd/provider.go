/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/praneeth-ayla/autocommenter/internal/ai"
	"github.com/praneeth-ayla/autocommenter/internal/config"
	"github.com/praneeth-ayla/autocommenter/internal/ui"
	"github.com/spf13/cobra"
)

// providerCmd represents the base command for provider-related subcommands.
var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Manage AI provider setting",
	Long:  `Set or view which AI provider is used for context generation.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'AutoCommenter provider set' to interactively set provider")
	},
}

// providerSetCmd represents the command to interactively set the AI provider.
var providerSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Interactively select an AI provider",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Prompt the user to select an AI provider from the supported list.
		selectedProvider, err := ui.SelectOne("Select AI Provider:", ai.SupportedProviders)
		if err != nil {
			return err // Return the error if selection fails.
		}

		// Save the selected provider to the configuration.
		err = config.SetProvider(selectedProvider)
		if err != nil {
			return fmt.Errorf("could not save provider: %w", err) // Wrap error for context.
		}

		fmt.Println("Provider updated to", selectedProvider)
		return nil
	},
}

// providerGetCmd represents the command to display the current AI provider.
var providerGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Show current AI provider",
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve the current AI provider name. The error is ignored as per original logic.
		name, _ := config.GetProvider()
		fmt.Println("Current provider:", name)
	},
}

// init initializes the provider commands and adds them to the root command.
func init() {
	// Add the providerCmd as a subcommand to the root command.
	rootCmd.AddCommand(providerCmd)
	// Add providerSetCmd as a subcommand to providerCmd.
	providerCmd.AddCommand(providerSetCmd)
	// Add providerGetCmd as a subcommand to providerCmd.
	providerCmd.AddCommand(providerGetCmd)
}

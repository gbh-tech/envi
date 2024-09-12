package cmd

import (
	"github.com/gbh-tech/envi/pkg/providers"
	"github.com/spf13/cobra"
)

func NewWerfCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "werf",
		Short: "Generate .env file from Werf",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize Werf provider
			werfProvider := providers.NewWerfProvider()
			// Use the provider to generate .env file
			return werfProvider.GenerateEnvFile()
		},
	}

	// Add flags specific to Werf command
	cmd.Flags().StringP("environment", "e", "", "Target environment")
	cmd.Flags().Bool("secrets", true, "Include secret files")
	cmd.Flags().StringSliceP("to-path", "p", []string{".env"}, "Paths to generate .env files")

	return cmd
}

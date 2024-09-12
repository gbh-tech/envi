package cmd

import (
	"github.com/gbh-tech/envi/pkg/providers"
	"github.com/spf13/cobra"
)

func NewOnePasswordCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "op",
		Short: "Generate .env file from 1Password",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize 1Password provider
			opProvider := providers.NewOnePasswordProvider()
			// Use the provider to generate .env file
			return opProvider.GenerateEnvFile()
		},
	}

	// Add flags specific to 1Password command
	cmd.Flags().StringP("vault", "v", "", "Target vault")
	cmd.Flags().StringSliceP("item", "i", nil, "Target secret items")
	cmd.Flags().StringSliceP("to-path", "p", []string{".env"}, "Paths to generate .env files")

	return cmd
}

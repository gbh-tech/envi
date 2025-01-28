package werf

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

type Options struct {
	Environment string
	Secrets     bool
	Values      string
	Path        []string
}

var Command = &cobra.Command{
	Use:     "werf",
	Short:   "Generate .env file from Werf",
	Example: "envi werf [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		options := parseCommandFlags(cmd)

		GenerateEnvFile(options)
	},
}

func parseCommandFlags(cmd *cobra.Command) Options {
	environment, _ := cmd.Flags().GetString("environment")
	secrets, _ := cmd.Flags().GetBool("secrets")
	values, _ := cmd.Flags().GetString("values")
	path, _ := cmd.Flags().GetStringArray("path")

	return Options{
		Environment: environment,
		Secrets:     secrets,
		Values:      values,
		Path:        path,
	}
}

func init() {
	Command.Flags().StringP(
		"values",
		"v",
		"",
		"Add custom values",
	)
	Command.Flags().BoolP(
		"secrets",
		"s",
		true,
		"Include secret files",
	)
	Command.Flags().StringP(
		"environment",
		"e",
		"",
		"Target environment",
	)

	if err := Command.MarkFlagRequired("environment"); err != nil {
		log.Fatalf("Error marking 'environment' flag as required: %v", err)
	}
}

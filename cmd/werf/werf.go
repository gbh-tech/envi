package werf

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

type Options struct {
	Environment string
	Secrets     bool
	Values      string
	ValueFiles  []string
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
	values, _ := cmd.Flags().GetString("extra-value")
	valueFiles, _ := cmd.Flags().GetStringArray("extra-values-file")
	path, _ := cmd.Flags().GetStringArray("path")

	return Options{
		Environment: environment,
		Secrets:     secrets,
		Values:      values,
		ValueFiles:  valueFiles,
		Path:        path,
	}
}

func init() {
	Command.Flags().String(
		"extra-value",
		"",
		"Adds extra value to Werf using --set flag",
	)
	Command.Flags().StringArray(
		"extra-values-file",
		[]string{},
		"Adds extra YAML value file to Werf using --values flag",
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

package werf

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

type Options struct {
	Development bool
	Environment string
	Secrets     bool
	Values      string
	ValueFiles  []string
	Path        []string
	Dir         string
	Overwrite   bool
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
	development, _ := cmd.Flags().GetBool("development")
	environment, _ := cmd.Flags().GetString("environment")
	secrets, _ := cmd.Flags().GetBool("secrets")
	values, _ := cmd.Flags().GetString("extra-value")
	valueFiles, _ := cmd.Flags().GetStringArray("extra-values-file")
	path, _ := cmd.Flags().GetStringArray("path")
	dir, _ := cmd.Flags().GetString("dir")
	overwrite, _ := cmd.Flags().GetBool("overwrite")

	return Options{
		Development: development,
		Environment: environment,
		Secrets:     secrets,
		Values:      values,
		ValueFiles:  valueFiles,
		Path:        path,
		Dir:         dir,
		Overwrite:   overwrite,
	}
}

func init() {
	Command.Flags().BoolP(
		"development",
		"d",
		false,
		"Enable development mode for werf render command",
	)
	Command.Flags().StringP(
		"dir",
		"w",
		"",
		"Directory to search for werf config files",
	)
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

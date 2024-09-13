package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	werf "github.com/gbh-tech/envi/pkg/providers/werf"
)

var WerfCommand = &cobra.Command{
	Use:     "werf",
	Short:   "Generate .env file from Werf",
	Example: "envi werf [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		options := configureWerfFlags(cmd)

		err := werf.GenerateEnvFile(options)
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
	},
}

func configureWerfFlags(cmd *cobra.Command) werf.WerfOptions {
	environment, _ := cmd.Flags().GetString("environment")
	secrets, _ := cmd.Flags().GetBool("secrets")
	values, _ := cmd.Flags().GetString("values")
	path, _ := cmd.Flags().GetStringArray("path")

	return werf.WerfOptions{
		Environment: environment,
		Secrets:     secrets,
		Values:      values,
		Path:        path,
	}
}

func init() {
	RootCmd.AddCommand(WerfCommand)

	WerfCommand.PersistentFlags().StringP("environment", "e", "", "Target environment")
	WerfCommand.MarkPersistentFlagRequired("environment")
	WerfCommand.PersistentFlags().StringP("values", "v", "", "Add custom values")
	WerfCommand.PersistentFlags().BoolP("secrets", "s", true, "Include secret files")
}

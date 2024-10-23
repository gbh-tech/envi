package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	manual "github.com/gbh-tech/envi/pkg/providers/manual"
)

var ManualCommand = &cobra.Command{
	Use:     "manual",
	Aliases: []string{"new"},
	Short:   "Generate .env file from manual input",
	Example: "envi manual [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		options := configureManual(cmd)

		err := manual.GenerateEnvFile(options)
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
	},
}

func configureManual(cmd *cobra.Command) manual.ManualOptions {
	path, _ := cmd.Flags().GetStringArray("path")
	value, _ := cmd.Flags().GetStringArray("value")

	return manual.ManualOptions{
		Value: value,
		Path:  path,
	}
}

func init() {
	RootCmd.AddCommand(ManualCommand)

	ManualCommand.PersistentFlags().StringArrayP("value", "v", []string{}, "Manual values to add")

	// Required flags
	if err := ManualCommand.MarkPersistentFlagRequired("value"); err != nil {
		log.Fatalf("Error marking 'value' flag as required: %v", err)
	}
}

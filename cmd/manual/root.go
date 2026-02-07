package manual

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

type Options struct {
	Values    []string
	Paths     []string
	Overwrite bool
	NextJS    bool
}

var Command = &cobra.Command{
	Use:     "manual",
	Aliases: []string{"new"},
	Short:   "Generate .env file from manual input",
	Example: "envi manual [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		options := parseCommandFlags(cmd)
		generateEnvFileFromManualInput(options)
	},
}

func parseCommandFlags(cmd *cobra.Command) Options {
	path, _ := cmd.Flags().GetStringArray("path")
	value, _ := cmd.Flags().GetStringArray("values")
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	nextjs, _ := cmd.Flags().GetBool("nextjs")

	return Options{
		Values:    value,
		Paths:     path,
		Overwrite: overwrite,
		NextJS:    nextjs,
	}
}

func init() {
	Command.Flags().StringArrayP(
		"values",
		"v",
		[]string{},
		"Manual-specified values to include",
	)

	if err := Command.MarkFlagRequired("values"); err != nil {
		log.Fatalf(
			"Error marking 'value' flag as required: %v",
			err,
		)
	}
}

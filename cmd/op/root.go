package op

import (
	"github.com/charmbracelet/log"
	"github.com/gbh-tech/envi/pkg/utils"
	"github.com/spf13/cobra"
)

type Options struct {
	Vault string
	Items []string
	Path  []string
}

var Command = &cobra.Command{
	Use:     "op",
	Aliases: []string{"onepassword"},
	Short:   "Generate .env file from 1Password",
	Example: "envi op [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		options := parseCommandFlags(cmd)

		token := utils.GetEnvironment("OP_SERVICE_ACCOUNT_TOKEN")
		client := NewClient(token)

		client.GenerateEnvFile(options)
	},
}

func parseCommandFlags(cmd *cobra.Command) Options {
	vault, _ := cmd.Flags().GetString("vault")
	item, _ := cmd.Flags().GetStringArray("item")
	path, _ := cmd.Flags().GetStringArray("path")

	return Options{
		Vault: vault,
		Items: item,
		Path:  path,
	}
}

func init() {
	Command.Flags().StringP(
		"vault",
		"v",
		"",
		"Vault ID",
	)
	Command.Flags().StringArrayP(
		"item",
		"i",
		[]string{""},
		"Item ID",
	)

	if err := Command.MarkFlagRequired("vault"); err != nil {
		log.Fatalf(
			"Error marking 'vault' flag as required: %v",
			err,
		)
	}
	if err := Command.MarkFlagRequired("item"); err != nil {
		log.Fatalf(
			"Error marking 'item' flag as required: %v",
			err,
		)
	}
}

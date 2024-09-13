package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	op "github.com/gbh-tech/envi/pkg/providers/onepassword"
)

var OnePasswordClient op.OnePasswordClient

var OnePasswordCommand = &cobra.Command{
	Use:     "op",
	Aliases: []string{"onepassword"},
	Short:   "Generate .env file from 1Password",
	Example: "envi op [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		options := configureOnepasswordFlags(cmd)

		token := os.Getenv("OP_SERVICE_ACCOUNT_TOKEN")
		OnePasswordClient, err := op.NewClient(token)
		if err != nil {
			log.Fatalf("Error initializing OnePassowrd client: %v", err)
		}

		err = OnePasswordClient.GenerateEnvFile(options)
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
	},
}

func configureOnepasswordFlags(cmd *cobra.Command) op.OnepasswordOptions {
	vault, _ := cmd.Flags().GetString("vault")
	item, _ := cmd.Flags().GetStringArray("item")
	path, _ := cmd.Flags().GetStringArray("path")

	return op.OnepasswordOptions{
		Vault: vault,
		Items: item,
		Path:  path,
	}
}

func init() {
	RootCmd.AddCommand(OnePasswordCommand)

	OnePasswordCommand.PersistentFlags().StringP("vault", "v", "", "Vault ID")
	OnePasswordCommand.MarkPersistentFlagRequired("vault")
	OnePasswordCommand.PersistentFlags().StringArrayP("item", "i", []string{""}, "Item ID")
	OnePasswordCommand.MarkPersistentFlagRequired("item")
}

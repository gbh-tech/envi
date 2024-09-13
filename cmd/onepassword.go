package cmd

import (
	"fmt"
	"github.com/gbh-tech/envi/pkg/providers"
	"github.com/spf13/cobra"
)

var Vault string
var Item []string
var ToPath []string

func OnePasswordCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "op",
		Short: "Generate .env file(s) using 1Password as provider",
		Run: func(cmd *cobra.Command, args []string) {
			client := providers.NewOnePasswordProvider()
			items, err := client.GetItems(Vault)
			if err != nil {
				return
			}
			fmt.Println(items)
		},
	}

	OnePasswordCommand := configureFlags(command)
	return OnePasswordCommand
}

func configureFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(
		&Vault,
		"vault",
		"v",
		"",
		"Target vault to fetch secret item",
	)
	cmd.Flags().StringSliceVarP(
		&Item,
		"item",
		"i",
		nil,
		"Target secret(s) from which to generate the .env file",
	)
	cmd.Flags().StringSliceVarP(
		&ToPath,
		"to-path",
		"p",
		[]string{".env"},
		"Path(s) to generate the dot env (.env) file to",
	)

	return cmd
}

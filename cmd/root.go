package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/gbh-tech/envi/cmd/manual"
	"github.com/gbh-tech/envi/cmd/op"
	"github.com/gbh-tech/envi/cmd/werf"
	"github.com/spf13/cobra"
)

var Version string = "dev"

var RootCmd = &cobra.Command{
	Use:     "envi",
	Short:   "A CLI tool for generating .env files from various sources",
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, show usage and return error
		_ = cmd.Help()
		log.Fatalf("requires at least one subcommand")
	},
}

func init() {
	RootCmd.AddCommand(manual.Command)
	RootCmd.AddCommand(op.Command)
	RootCmd.AddCommand(werf.Command)

	RootCmd.PersistentFlags().StringArrayP(
		"path",
		"p",
		[]string{".env"},
		"Target file path",
	)
	RootCmd.PersistentFlags().BoolP(
		"overwrite",
		"o",
		false,
		"Overwrite existing file",
	)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

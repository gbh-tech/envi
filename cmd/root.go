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

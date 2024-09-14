package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "envi",
	Short:   "A CLI tool for generating .env files from various sources",
	Version: "1.0.0",
}

func init() {
	RootCmd.PersistentFlags().StringArrayP("path", "p", []string{".env"}, "Target repository")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

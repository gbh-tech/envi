package main

import (
	"fmt"
	"github.com/gbh-tech/envi/cmd"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "envi",
		Short:   "A CLI tool for generating .env files from various sources",
		Version: "1.0.0",
	}

	rootCmd.AddCommand(cmd.OnePasswordCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

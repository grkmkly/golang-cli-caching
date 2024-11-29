package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "root",
	Run: func(cmd *cobra.Command, args []string) {
		fncrootCmd()
	},
}

func fncrootCmd() {
	fmt.Println("Hello World")
}

func init() {
	rootCmd.AddCommand(caching)
	rootCmd.AddCommand(cachingFile)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Print("Error")
		os.Exit(1)
	}
}

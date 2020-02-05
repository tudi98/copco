package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parses a problem",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("parse called")
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)
}

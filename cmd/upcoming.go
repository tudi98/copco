package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var upcomingCmd = &cobra.Command{
	Use:   "upcoming",
	Short: "Shows upcoming contests.",
	Run: func(cmd *cobra.Command, args []string) {
		upcoming()
	},
}

func init() {
	rootCmd.AddCommand(upcomingCmd)

}

func upcoming() {
	fmt.Println("Upcoming contests...")
}

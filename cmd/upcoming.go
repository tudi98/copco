package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tudi98/copco/pkg/parser"
)

var upcomingCmd = &cobra.Command{
	Use:   "upcoming",
	Short: "Show upcoming contests",
	Run: func(cmd *cobra.Command, args []string) {
		upcoming()
	},
}

func init() {
	rootCmd.AddCommand(upcomingCmd)
}

func upcoming() {
	fmt.Println("***** Upcoming Contests *****")
	p := parser.CodeforcesParser{}
	upcomingContests := p.GetUpcoming()
	for _, v := range upcomingContests {
		fmt.Println(v)
	}
	fmt.Println("*****************************")
}

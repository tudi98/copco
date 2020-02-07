package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tudi98/copco/pkg/parser"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parses a problem",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		parseProblem(args[0])
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)
}

func parseProblem(url string) {
	var p parser.Parser = parser.CodeforcesParser{}
	problem := p.GetProblem(url)
	fmt.Printf("%+v", problem)
}

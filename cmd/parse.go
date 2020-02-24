package cmd

import (
	"log"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tudi98/copco/pkg/parser/codeforces"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parses a problem or a contest",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		parse(args[0])
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)
}

func parse(onlineJudgeUrl string) {
	u, err := url.ParseRequestURI(onlineJudgeUrl)
	if err != nil {
		log.Fatal(err)
	}
	if strings.Contains(u.Host, "codeforces.com") {
		codeforces.Parse(onlineJudgeUrl)
	} else {
		log.Fatal("Site not supported.")
	}
}

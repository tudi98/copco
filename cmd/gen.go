package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tudi98/copco/code_generator"
	"io/ioutil"
	"log"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate source code using custom code",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gen(args[0])
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}

func gen(file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Error while opening %s", file)
	}
	sourceCode, err := code_generator.GenerateSolution(string(content))
	err = ioutil.WriteFile(file, []byte(sourceCode), 0644)
	if err != nil {
		log.Fatalf("Error while writing to %s", file)
	}
}

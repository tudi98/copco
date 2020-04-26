package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
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
	sourceCode, err := generateSourceCode(string(content))
	err = ioutil.WriteFile(file, []byte(sourceCode), 0644)
	if err != nil {
		log.Fatalf("Error while writing to %s", file)
	}
}

func generateSourceCode(sourceCode string) (string, error) {
	r := regexp.MustCompile("#copco \".+\"")
	lines := strings.Split(sourceCode, "\n")
	sourceCode = ""
	for _, line := range lines {
		if r.MatchString(line) {
			codePath := strings.Split(line, "\"")[1]
			codePath = strings.ReplaceAll(codePath, "/", string(os.PathSeparator))
			codePath = viper.GetString("COPCO_CUSTOM_CODE") + string(os.PathSeparator) + codePath
			content, err := ioutil.ReadFile(codePath)
			if err != nil {
				return "", fmt.Errorf("could not find custom code file %s", codePath)
			}
			sourceCode += string(content) + "\n"
		} else {
			sourceCode += line + "\n"
		}
	}
	return sourceCode, nil
}

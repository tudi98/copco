package code_generator

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func GenerateSolution(sourceCode string) (string, error) {
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

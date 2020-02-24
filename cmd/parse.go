package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

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
	p := parser.CodeforcesParser{}
	problem := p.GetProblem(url)

	separator := "/"
	if runtime.GOOS == "windows" {
		separator = "\\"
	}

	problemPath := os.Getenv("COPCO_PATH") + separator + "codeforces" +
		separator + problem.ContestId + separator + problem.Name
	templatePath := os.Getenv("COPCO_TEMPLATE")

	if _, err := os.Stat(problemPath); os.IsNotExist(err) {
		err := os.MkdirAll(problemPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Error when creating %s", problemPath)
		}
	}

	from, err := os.Open(templatePath)
	if err != nil {
		log.Fatalf("Error when opening template %s", templatePath)
	}
	defer from.Close()

	sourcePath := problemPath + separator + "main.cpp"
	to, err := os.OpenFile(sourcePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("Error when creating %s", sourcePath)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}

	testsPath := problemPath + separator + "tests"
	if _, err := os.Stat(testsPath); os.IsNotExist(err) {
		err := os.MkdirAll(testsPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Error when creating %s", testsPath)
		}
	}

	for i, v := range problem.Inputs {
		filePath := testsPath + separator + fmt.Sprintf("%d.in", i)
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatalf("Error while creating %s", filePath)
		}
		_, err = file.WriteString(v)
		if err != nil {
			log.Fatalf("Error while writing to %s", filePath)
		}
		file.Close()
	}

	for i, v := range problem.Outputs {
		filePath := testsPath + separator + fmt.Sprintf("%d.ok", i)
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatalf("Error while creating %s", filePath)
		}
		_, err = file.WriteString(v)
		if err != nil {
			log.Fatalf("Error while writing to %s", filePath)
		}
		file.Close()
	}

	jsonPath := problemPath + separator + "problem.json"
	file, err := os.OpenFile(jsonPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("Error when creating %s", jsonPath)
	}

	jsonVal, err := json.Marshal(problem)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.WriteString(string(jsonVal))
	if err != nil {
		log.Fatalf("Error when writing to %s", jsonPath)
	}

	file.Close()
}

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

	problem_path := os.Getenv("COPCO_PATH") + separator + "codeforces" +
		separator + "problemset" + separator + problem.Id
	template_path := os.Getenv("COPCO_TEMPLATE")

	if _, err := os.Stat(problem_path); os.IsNotExist(err) {
		os.MkdirAll(problem_path, os.ModePerm)
	} else {
		log.Fatal(err)
	}

	from, err := os.Open(template_path)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(problem_path+separator+"main.cpp", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}

	for i, v := range problem.Inputs {
		file, err := os.OpenFile(problem_path+separator+fmt.Sprintf("%d.in", i), os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		file.WriteString(v)
		file.Close()
	}

	for i, v := range problem.Outputs {
		file, err := os.OpenFile(problem_path+separator+fmt.Sprintf("%d.ok", i), os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		file.WriteString(v)
		file.Close()
	}

	file, err := os.OpenFile(problem_path+separator+"problem.json", os.O_RDWR|os.O_CREATE, 0666)
	json_val, err := json.Marshal(problem)
	if err != nil {
		log.Fatal(err)
	}
	file.WriteString(string(json_val))
	file.Close()
}

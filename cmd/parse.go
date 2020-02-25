package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tudi98/copco/pkg/atcoder"
	"github.com/tudi98/copco/pkg/codeforces"
	"github.com/tudi98/copco/pkg/models"
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

func parse(Url string) {
	var parser models.ParserInterface
	var onlineJudge string

	Url = strings.TrimRight(Url, "/")

	u, err := url.ParseRequestURI(Url)
	if err != nil {
		log.Fatal(err)
	}

	if strings.Contains(u.Host, "codeforces.com") {
		parser = codeforces.Parser{}
		onlineJudge = codeforces.OnlineJudge
	} else if strings.Contains(u.Host, "atcoder.jp") {
		parser = atcoder.Parser{}
		onlineJudge = atcoder.OnlineJudge
	} else {
		log.Fatal("site not supported")
	}

	if parser.ValidateProblemUrl(Url) {
		createProblem(Url, parser, onlineJudge)
	} else if parser.ValidateContestUrl(Url) {
		createContest(Url, parser, onlineJudge)
	}
}

func createContest(Url string, parser models.ParserInterface, onlineJudge string) {
	contest, err := parser.ParseContest(Url)
	if err != nil {
		log.Fatal(err)
	}
	for _, problemUrl := range contest.Urls {
		createProblem(problemUrl, parser, onlineJudge)
	}
	color.Green("Done!")
}

func createProblem(Url string, parser models.ParserInterface, onlineJudge string) {
	separator := "/"
	if runtime.GOOS == "windows" {
		separator = "\\"
	}

	problem, err := parser.ParseProblem(Url)
	if err != nil {
		log.Fatal(err)
	}

	problemPath := os.Getenv("COPCO_PATH") + separator + onlineJudge + separator + problem.ContestId + separator + problem.Name
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

	fmt.Printf("%s...", Url)
	color.Green("OK!")
}

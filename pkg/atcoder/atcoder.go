package atcoder

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/gocolly/colly"
	"github.com/tudi98/copco/pkg/models"
)

func Parse(url string) {
	r1, _ := regexp.Compile("atcoder\\.jp/contests/.*/tasks/.*")
	r2, _ := regexp.Compile("atcoder\\.jp/contests/.*/tasks")
	r3, _ := regexp.Compile("atcoder\\.jp/contests/.*")
	if r1.MatchString(url) {
		createProblem(url)
	} else if r2.MatchString(url) {
		createContest(url)
	} else if r3.MatchString(url) {
		createContest(url + "/tasks")
	} else {
		log.Fatal("Invalid url")
	}
}

func createContest(url string) {
	contest := parseContest(url)
	for _, problemUrl := range contest.Urls {
		createProblem(problemUrl)
	}
	color.Green("Done!")
}

func createProblem(url string) {
	fmt.Printf("%s...", url)

	separator := "/"
	if runtime.GOOS == "windows" {
		separator = "\\"
	}

	problem := parseProblem(url)

	problemPath := os.Getenv("COPCO_PATH") + separator + "atcoder" + separator + problem.ContestId + separator + problem.Name
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

	color.Green("OK!")
}

func parseContest(url string) models.Contest {
	contest := models.Contest{}

	urlArray := strings.Split(url, "/")
	contest.ContestId = urlArray[len(urlArray)-1]

	c := colly.NewCollector(
		colly.AllowedDomains("atcoder.jp"),
	)

	c.OnHTML("table > tbody > tr > td:first-child > a[href]", func(e *colly.HTMLElement) {
		contest.Urls = append(contest.Urls, "https://atcoder.jp"+e.Attr("href"))
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}

	return contest
}

func parseProblem(url string) models.Problem {
	problem := models.Problem{}

	urlArray := strings.Split(url, "/")
	problem.ContestId = urlArray[len(urlArray)-3]
	problem.ProblemId = urlArray[len(urlArray)-1]

	c := colly.NewCollector(
		colly.AllowedDomains("atcoder.jp"),
	)

	c.OnHTML("div.col-sm-12", func(e *colly.HTMLElement) {
		if html, _ := e.DOM.Html(); !strings.Contains(html, "<span class=\"h2\"") {
			return
		}
		problem.Name = e.ChildText("span.h2")
		text := strings.Split(e.ChildText("p"), " ")
		timeLimit, err := strconv.ParseFloat(text[2], 32)
		if err != nil {
			log.Fatal(err)
		}
		timeLimit *= 1000
		problem.TimeLimit = int(timeLimit)
		problem.MemoryLimit, err = strconv.Atoi(text[7])
		if err != nil {
			log.Fatal(err)
		}
		problem.MemoryLimit = problem.MemoryLimit * 1024 * 1024
	})

	c.OnHTML("div#task-statement > span > span > div > section", func(e *colly.HTMLElement) {
		if strings.Contains(e.ChildText("h3"), "Sample Input") {
			problem.Inputs = append(problem.Inputs, e.ChildText("pre"))
		} else if strings.Contains(e.ChildText("h3"), "Sample Output") {
			problem.Outputs = append(problem.Outputs, e.ChildText("pre"))
		}
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}

	return problem
}

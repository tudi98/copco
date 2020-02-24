package codeforces

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
	"github.com/tudi98/copco/pkg/parser/models"
)

func Parse(url string) {
	r1, _ := regexp.Compile("codeforces\\.com/problemset/problem/[0-9]+/.*")
	r2, _ := regexp.Compile("codeforces\\.com/contest/[0-9]+/problem/.*")
	r3, _ := regexp.Compile("codeforces\\.com/contest/[0-9]+")
	if r1.MatchString(url) || r2.MatchString(url) {
		createProblem(url)
	} else if r3.MatchString(url) {
		createContest(url)
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

	problemPath := os.Getenv("COPCO_PATH") + separator + "codeforces" + separator + problem.ContestId + separator + problem.Name
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
		colly.AllowedDomains("codeforces.com"),
	)

	c.OnHTML("table.problems > tbody > tr > td:first-child > a[href]", func(e *colly.HTMLElement) {
		contest.Urls = append(contest.Urls, "https://codeforces.com"+e.Attr("href"))
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
	if strings.Contains(url, "problemset") {
		problem.ContestId = urlArray[len(urlArray)-2]
	} else {
		problem.ContestId = urlArray[len(urlArray)-3]
	}
	problem.ProblemId = problem.ContestId + urlArray[len(urlArray)-1]

	c := colly.NewCollector(
		colly.AllowedDomains("codeforces.com"),
	)

	c.OnHTML("div.problem-statement", func(e *colly.HTMLElement) {
		problem.Name = e.ChildText("div.header > div.title")
		text := strings.Split(e.ChildText("div.header > div.time-limit"), "test")[1]
		timeLimit, err := strconv.ParseFloat(strings.Split(text, " ")[0], 32)
		timeLimit *= 1000
		problem.TimeLimit = int(timeLimit)
		if err != nil {
			log.Fatal(err)
		}
		text = strings.Split(e.ChildText("div.header > div.memory-limit"), "test")[1]
		problem.MemoryLimit, err = strconv.Atoi(strings.Split(text, " ")[0])
		if err != nil {
			log.Fatal(err)
		}
		problem.MemoryLimit = problem.MemoryLimit * 1024 * 1024
	})

	c.OnHTML("div.input > pre", func(e *colly.HTMLElement) {
		html, err := e.DOM.Html()
		if err != nil {
			log.Fatal(err)
		}
		r := strings.NewReplacer("<br/>", "\n", "<br />", "\n", "<br>", "\n")
		html = r.Replace(html)
		problem.Inputs = append(problem.Inputs, html)
	})

	c.OnHTML("div.output > pre", func(e *colly.HTMLElement) {
		html, err := e.DOM.Html()
		if err != nil {
			log.Fatal(err)
		}
		r := strings.NewReplacer("<br/>", "\n", "<br />", "\n", "<br>", "\n")
		html = r.Replace(html)
		problem.Outputs = append(problem.Outputs, html)
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}

	return problem
}

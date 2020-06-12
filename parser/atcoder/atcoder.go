package atcoder

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/tudi98/copco/parser/models"
)

const OnlineJudge = "atcoder"

type Parser struct{}

func (p Parser) ValidateContestUrl(url string) bool {
	r1 := regexp.MustCompile("atcoder\\.jp/contests/.+/tasks")
	r2 := regexp.MustCompile("atcoder\\.jp/contests/.+")
	if r1.MatchString(url) || r2.MatchString(url) {
		return true
	}
	return false
}

func (p Parser) ValidateProblemUrl(url string) bool {
	r1 := regexp.MustCompile("atcoder\\.jp/contests/.+/tasks/.+")
	if r1.MatchString(url) {
		return true
	}
	return false
}

func (p Parser) ParseContest(url string) (models.Contest, error) {
	contest := models.Contest{}

	urlArray := strings.Split(url, "/")

	if urlArray[len(urlArray)-1] == "tasks" {
		contest.ContestId = urlArray[len(urlArray)-2]
	} else {
		url += "/tasks"
		contest.ContestId = urlArray[len(urlArray)-1]
	}

	contest.ContestUrl = url

	c := colly.NewCollector(
		colly.AllowedDomains("atcoder.jp"),
	)

	// Set a delay between requests
	c.Limit(&colly.LimitRule{
		DomainGlob: "atcoder.jp/*",
		Delay:      1 * time.Second,
	})

	c.OnHTML("table > tbody > tr > td:first-child > a[href]", func(e *colly.HTMLElement) {
		contest.Urls = append(contest.Urls, "https://atcoder.jp"+e.Attr("href"))
	})

	err := c.Visit(url)
	if err != nil {
		return contest, err
	}

	return contest, nil
}

func (p Parser) ParseProblem(url string) (models.Problem, error) {
	problem := models.Problem{}

	urlArray := strings.Split(url, "/")
	problem.ContestId = urlArray[len(urlArray)-3]
	problem.ProblemId = urlArray[len(urlArray)-1]
	problem.ProblemUrl = url

	var parsingErr error
	parsingErr = nil

	c := colly.NewCollector(
		colly.AllowedDomains("atcoder.jp"),
	)

	c.OnHTML("div.col-sm-12", func(e *colly.HTMLElement) {
		if html, _ := e.DOM.Html(); !strings.Contains(html, "<span class=\"h2\"") {
			return
		}
		problem.Name = e.ChildText("span.h2")
		text := strings.Split(e.ChildText("p"), " ")
		if len(text) < 8 {
			parsingErr = fmt.Errorf("could not parse time limit and memory limit")
			return
		}
		timeLimit, err := strconv.ParseFloat(text[2], 32)
		if err != nil {
			parsingErr = err
			return
		}
		timeLimit *= 1000
		problem.TimeLimit = int(timeLimit)
		problem.MemoryLimit, err = strconv.Atoi(text[7])
		if err != nil {
			parsingErr = err
			return
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

	c.OnError(func(r *colly.Response, err error) {
		parsingErr = err
	})

	err := c.Visit(url)
	if err != nil {
		parsingErr = fmt.Errorf("error when parsing %s", url)
	}

	return problem, parsingErr
}

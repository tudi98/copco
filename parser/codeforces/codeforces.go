package codeforces

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/tudi98/copco/parser/models"
)

const OnlineJudge = "codeforces"

type Parser struct{}

func (p Parser) ValidateContestUrl(url string) bool {
	r := regexp.MustCompile("codeforces\\.com/contest/[0-9]+")
	if r.MatchString(url) {
		return true
	}
	return false
}

func (p Parser) ValidateProblemUrl(url string) bool {
	r1 := regexp.MustCompile("codeforces\\.com/problemset/problem/[0-9]+/.+")
	r2 := regexp.MustCompile("codeforces\\.com/contest/[0-9]+/problem/.+")
	if r1.MatchString(url) || r2.MatchString(url) {
		return true
	}
	return false
}

func (p Parser) ParseContest(url string) (models.Contest, error) {
	contest := models.Contest{}

	urlArray := strings.Split(url, "/")
	contest.ContestId = urlArray[len(urlArray)-1]
	contest.ContestUrl = url

	c := colly.NewCollector(
		colly.AllowedDomains("codeforces.com"),
		colly.URLFilters(
			regexp.MustCompile(strings.ReplaceAll(url, ".", "\\.")),
		),
	)

	var pErr error
	pErr = nil

	c.OnResponse(func(r *colly.Response) {
		if r.Request.URL.Path == "/" {
			pErr = errors.New("no such contest")
		}
	})

	c.OnHTML("table.problems > tbody > tr > td:first-child > a[href]", func(e *colly.HTMLElement) {
		contest.Urls = append(contest.Urls, "https://codeforces.com"+e.Attr("href"))
	})

	err := c.Visit(url)
	if err != nil {
		return models.Contest{}, err
	}

	return contest, pErr
}

func (p Parser) ParseProblem(url string) (models.Problem, error) {
	problem := models.Problem{}

	urlArray := strings.Split(url, "/")
	if strings.Contains(url, "problemset") {
		problem.ContestId = urlArray[len(urlArray)-2]
	} else {
		problem.ContestId = urlArray[len(urlArray)-3]
	}
	problem.ProblemId = problem.ContestId + urlArray[len(urlArray)-1]
	problem.ProblemUrl = url

	var pErr error
	pErr = nil

	c := colly.NewCollector(
		colly.AllowedDomains("codeforces.com"),
	)

	// Set a delay between requests
	c.Limit(&colly.LimitRule{
		DomainGlob: "codeforces.com/*",
		Delay:      1 * time.Second,
	})

	c.OnResponse(func(r *colly.Response) {
		if r.Request.URL.Path == "/" {
			pErr = errors.New("no such problem")
		}
	})

	c.OnHTML("div.problem-statement", func(e *colly.HTMLElement) {
		problem.Name = e.ChildText("div.header > div.title")
		text := strings.Split(e.ChildText("div.header > div.time-limit"), "test")[1]
		timeLimit, err := strconv.ParseFloat(strings.Split(text, " ")[0], 32)
		if err != nil {
			pErr = err
			return
		}
		timeLimit *= 1000
		problem.TimeLimit = int(timeLimit)
		text = strings.Split(e.ChildText("div.header > div.memory-limit"), "test")[1]
		problem.MemoryLimit, err = strconv.Atoi(strings.Split(text, " ")[0])
		if err != nil {
			pErr = err
			return
		}
		problem.MemoryLimit = problem.MemoryLimit * 1024 * 1024
	})

	c.OnHTML("div.input > pre", func(e *colly.HTMLElement) {
		html, err := e.DOM.Html()
		if err != nil {
			pErr = err
			return
		}
		r := strings.NewReplacer("<br/>", "\n", "<br />", "\n", "<br>", "\n")
		html = r.Replace(html)
		problem.Inputs = append(problem.Inputs, html)
	})

	c.OnHTML("div.output > pre", func(e *colly.HTMLElement) {
		html, err := e.DOM.Html()
		if err != nil {
			pErr = err
			return
		}
		r := strings.NewReplacer("<br/>", "\n", "<br />", "\n", "<br>", "\n")
		html = r.Replace(html)
		problem.Outputs = append(problem.Outputs, html)
	})

	c.OnError(func(r *colly.Response, err error) {
		pErr = err
	})

	err := c.Visit(url)
	if err != nil {
		return models.Problem{}, err
	}

	if pErr != nil {
		problem = models.Problem{}
	}

	return problem, pErr
}

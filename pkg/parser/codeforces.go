package parser

import (
	"github.com/tudi98/copco/pkg/models"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type CodeforcesParser struct{}

func (CodeforcesParser) GetUpcoming() []string {
	contests := make([]string, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("codeforces.com"),
	)

	// TODO: check html when there is a live contest
	c.OnHTML("div.contestList > div.datatable > div > table > tbody > tr[data-contestid]", func(e *colly.HTMLElement) {
		// TODO: convert time to local time
		contests = append(contests, e.ChildText("td:nth-child(1)")+" "+e.ChildText("td:nth-child(3)")+" MSK")
	})

	c.Visit("https://codeforces.com/contests?complete=true")

	return contests
}

func (CodeforcesParser) GetProblem(url string) models.Problem {
	problem := models.Problem{}

	// TODO: validate url

	urlArray := strings.Split(url, "/")
	if strings.Contains(url, "problemset") {
		problem.ContestId = urlArray[len(urlArray)-2]
	} else {
		problem.ContestId = urlArray[len(urlArray)-3]
	}
	problem.Id = problem.ContestId + urlArray[len(urlArray)-1]

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

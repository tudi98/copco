package parser

import (
	"strings"

	"github.com/gocolly/colly"
)

type CodeforcesParser struct{}

func (CodeforcesParser) GetUpcoming() []string {
	contests := make([]string, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("codeforces.com"),
	)

	c.OnHTML("div.contestList > div.datatable > div > table > tbody > tr[data-contestid]", func(e *colly.HTMLElement) {
		// TODO convert time to local time
		contests = append(contests, e.ChildText("td:nth-child(1)")+" "+e.ChildText("td:nth-child(3)")+" MSK")
	})

	c.Visit("https://codeforces.com/contests?complete=true")

	return contests
}

func (CodeforcesParser) GetProblem(url string) Problem {
	problem := Problem{}

	// TODO: validate url

	url_array := strings.Split(url, "/")

	if strings.Contains(url, "problemset") {
		problem.Id = url_array[len(url_array)-2] + url_array[len(url_array)-1]
	} else {
		problem.Id = url_array[len(url_array)-3] + url_array[len(url_array)-1]
	}

	c := colly.NewCollector(
		colly.AllowedDomains("codeforces.com"),
	)

	c.OnHTML("div.problem-statement", func(e *colly.HTMLElement) {
		text := strings.Split(e.ChildText("div.header > div.time-limit"), "test")[1]
		problem.TimeLimit = strings.Split(text, " ")[0]
		text = strings.Split(e.ChildText("div.header > div.memory-limit"), "test")[1]
		problem.MemoryLimit = strings.Split(text, " ")[0]
	})

	c.OnHTML("div.input > pre", func(e *colly.HTMLElement) {
		problem.Inputs = append(problem.Inputs, e.Text)
	})

	c.OnHTML("div.output > pre", func(e *colly.HTMLElement) {
		problem.Outputs = append(problem.Outputs, e.Text)
	})

	c.Visit(url)

	return problem
}

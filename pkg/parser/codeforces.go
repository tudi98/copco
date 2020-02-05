package parser

import "github.com/gocolly/colly"

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

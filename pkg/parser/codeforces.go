package parser

type CodeforcesParser struct{}

func (CodeforcesParser) GetUpcoming() []string {
	return []string{"Contest 1 11:00", "Contest 2 12:00"}
}

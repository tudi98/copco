package parser

type Parser interface {
	GetUpcoming() []string
	GetProblem(url string) Problem
}

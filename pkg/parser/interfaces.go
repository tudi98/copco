package parser

import "github.com/tudi98/copco/pkg/models"

type Parser interface {
	GetUpcoming() []string
	GetProblem(url string) models.Problem
}

package models

type Problem struct {
	ProblemId   string
	ProblemUrl  string
	Name        string
	ContestId   string
	TimeLimit   int
	MemoryLimit int
	Inputs      []string
	Outputs     []string
}

type Contest struct {
	ContestId  string
	ContestUrl string
	Urls       []string
}

type ParserInterface interface {
	ValidateContestUrl(url string) bool
	ValidateProblemUrl(url string) bool
	ParseContest(url string) (Contest, error)
	ParseProblem(url string) (Problem, error)
}

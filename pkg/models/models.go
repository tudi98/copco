package models

type Problem struct {
	ProblemId   string
	Name        string
	ContestId   string
	TimeLimit   int
	MemoryLimit int
	Inputs      []string
	Outputs     []string
}

type Contest struct {
	ContestId string
	Urls      []string
}

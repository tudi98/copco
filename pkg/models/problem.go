package models

type Problem struct {
	Id          string
	Name        string
	ContestId   string
	TimeLimit   int
	MemoryLimit int
	Inputs      []string
	Outputs     []string
}

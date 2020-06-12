package atcoder

import (
	"github.com/tudi98/copco/parser/models"
	"reflect"
	"testing"
)

func TestParser_ValidateProblemUrl(t *testing.T) {
	testCases := []struct {
		in  string
		out bool
	}{
		{"", false},
		{"https://atcoder.jp/contests/agc044/tasks/agc044_a", true},
		{"atcoder.jp/contests/agc044/tasks/agc044_a", true},
		{"https://atcoder.com/contests/agc044/tasks/agc044_a", false},
		{"https://atcoder.jp/contests/agc044/tasks/", false},
	}
	parser := Parser{}
	for _, tc := range testCases {
		out := parser.ValidateProblemUrl(tc.in)
		if out != tc.out {
			t.Errorf("Error when validating '%s' got: %t, want: %t", tc.in, out, tc.out)
		}
	}
}

func TestParser_ValidateContestUrl(t *testing.T) {
	testCases := []struct {
		in  string
		out bool
	}{
		{"", false},
		{"https://atcoder.jp/contests/agc044", true},
		{"https://atcoder.jp/contests/agc044/tasks", true},
		{"atcoder.jp/contests/agc044", true},
		{"https://atcoder.com/contests/agc044", false},
		{"https://atcoder.com/contests/", false},
	}
	parser := Parser{}
	for _, tc := range testCases {
		out := parser.ValidateContestUrl(tc.in)
		if out != tc.out {
			t.Errorf("Error when validating '%s' got: %t, want: %t", tc.in, out, tc.out)
		}
	}
}

func TestParser_ParseProblem(t *testing.T) {
	testCases := []struct {
		in  string
		out models.Problem
	}{
		{
			"https://atcoder.jp/contests/agc044/tasks/agc044_a",
			models.Problem{
				ProblemId:   "agc044_a",
				ProblemUrl:  "https://atcoder.jp/contests/agc044/tasks/agc044_a",
				Name:        "A - Pay to Win",
				ContestId:   "agc044",
				TimeLimit:   2000,
				MemoryLimit: 1073741824,
				Inputs:      []string{"5\n11 1 2 4 8\n11 1 2 2 8\n32 10 8 5 4\n29384293847243 454353412 332423423 934923490 1\n900000000000000000 332423423 454353412 934923490 987654321"},
				Outputs:     []string{"20\n19\n26\n3821859835\n23441258666"},
			},
		},
		{
			"https://atcoder.jp/contests/agc044/tasks/agc044_z",
			models.Problem{},
		},
	}
	parser := Parser{}
	for _, tc := range testCases {
		out, _ := parser.ParseProblem(tc.in)
		if !reflect.DeepEqual(out, tc.out) {
			t.Errorf("Error when parsing '%s'\n got: %+v,\nwant: %+v", tc.in, out, tc.out)
		}
	}
}

func TestParser_ParseContest(t *testing.T) {
	testCases := []struct {
		in  string
		out models.Contest
	}{
		{
			"https://atcoder.jp/contests/agc044",
			models.Contest{
				ContestId:  "agc044",
				ContestUrl: "https://atcoder.jp/contests/agc044/tasks",
				Urls: []string{
					"https://atcoder.jp/contests/agc044/tasks/agc044_a",
					"https://atcoder.jp/contests/agc044/tasks/agc044_b",
					"https://atcoder.jp/contests/agc044/tasks/agc044_c",
					"https://atcoder.jp/contests/agc044/tasks/agc044_d",
					"https://atcoder.jp/contests/agc044/tasks/agc044_e",
					"https://atcoder.jp/contests/agc044/tasks/agc044_f",
				},
			},
		},
		{
			"https://atcoder.jp/contests/",
			models.Contest{},
		},
	}
	parser := Parser{}
	for _, tc := range testCases {
		out, _ := parser.ParseContest(tc.in)
		if !reflect.DeepEqual(out, tc.out) {
			t.Errorf("Error when parsing '%s' \n got: %+v, \nwant: %+v", tc.in, out, tc.out)
		}
	}
}

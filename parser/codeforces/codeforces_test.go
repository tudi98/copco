package codeforces

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
		{"https://codeforces.com/contest/1344/problem/C", true},
		{"codeforces.com/contest/1344/problem/C", true},
		{"https://codeforces.com/contest/1344/problem/", false},
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
		{"https://codeforces.com/contest/1344", true},
		{"codeforces.com/contest/1344", true},
		{"https://codeforces.com/contest/", false},
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
			"https://codeforces.com/contest/1344/problem/C",
			models.Problem{
				ProblemId:   "1344C",
				ProblemUrl:  "https://codeforces.com/contest/1344/problem/C",
				Name:        "C. Quantifier Question",
				ContestId:   "1344",
				TimeLimit:   1000,
				MemoryLimit: 268435456,
				Inputs: []string{
					"2 1\n1 2\n",
					"4 3\n1 2\n2 3\n3 1\n",
					"3 2\n1 3\n2 3\n",
				},
				Outputs: []string{
					"1\nAE\n",
					"-1\n",
					"2\nAAE\n",
				},
			},
		},
		{
			"https://codeforces.com/contest/1344/problem/Z",
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
			"https://codeforces.com/contest/1344",
			models.Contest{
				ContestId:  "1344",
				ContestUrl: "https://codeforces.com/contest/1344",
				Urls: []string{
					"https://codeforces.com/contest/1344/problem/A",
					"https://codeforces.com/contest/1344/problem/B",
					"https://codeforces.com/contest/1344/problem/C",
					"https://codeforces.com/contest/1344/problem/D",
					"https://codeforces.com/contest/1344/problem/E",
					"https://codeforces.com/contest/1344/problem/F",
				},
			},
		},
		{
			"https://codeforces.com/contest/",
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

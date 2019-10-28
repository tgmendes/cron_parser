package cron

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	testCases := []struct {
		desc      string
		expr      string
		expOutput string
	}{
		{
			desc: "given example",
			expr: "*/15 0 1,15 * 1-5 /usr/bin/find",
			expOutput: `minute        0 15 30 45
hour          0
day of month  1 15
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   1 2 3 4 5
command       /usr/bin/find
`,
		}, {
			desc: "run all",
			expr: "* * * * * /usr/bin/findAll",
			expOutput: fmt.Sprintf(`minute        %s
hour          %s
day of month  %s
month         %s
day of week   %s
command       /usr/bin/findAll
`, getInterval(0, 59), getInterval(0, 23), getInterval(1, 31), getInterval(1, 12), getInterval(0, 6)),
		}, {
			desc: "time steps",
			expr: "*/15 */4 */5 */4 */3 /usr/bin/findSteps",
			expOutput: `minute        0 15 30 45
hour          0 4 8 12 16 20
day of month  1 6 11 16 21 26 31
month         1 5 9
day of week   0 3 6
command       /usr/bin/findSteps
`,
		}, {
			desc: "run at specific time",
			expr: "30 4 5 10 5 /usr/bin/findTime",
			expOutput: `minute        30
hour          4
day of month  5
month         10
day of week   5
command       /usr/bin/findTime
`,
		}, {
			desc: "run at specific times",
			expr: "0,5,10 0,4,6 1,15,31 4,8,12 1,3,5 /usr/bin/findTimes",
			expOutput: `minute        0 5 10
hour          0 4 6
day of month  1 15 31
month         4 8 12
day of week   1 3 5
command       /usr/bin/findTimes
`,
		}, {
			desc: "run at time intervals",
			expr: "0-10 18-23 1-15 4-10 1-1 /usr/bin/findIntervals",
			expOutput: `minute        0 1 2 3 4 5 6 7 8 9 10
hour          18 19 20 21 22 23
day of month  1 2 3 4 5 6 7 8 9 10 11 12 13 14 15
month         4 5 6 7 8 9 10
day of week   1
command       /usr/bin/findIntervals
`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parsed, err := Parse(tC.expr)

			if err != nil {
				t.Fatalf("unexpected error occurred %s", err.Error())
			}

			if parsed != tC.expOutput {
				t.Logf("unexpected result\n")
				t.Logf("want:\n%s", tC.expOutput)
				t.Logf("got:\n%s", parsed)
				t.FailNow()
			}
		})
	}
}

func TestOutOfRangeErrors(t *testing.T) {
	testCases := []struct {
		desc string
		expr string
	}{
		{
			desc: "minutes upper",
			expr: "60 0 1,15 * 1-5 /usr/bin/find",
		}, {
			desc: "minutes lower",
			expr: "-1 0 1,15 * 1-5 /usr/bin/find",
		},
		{
			desc: "hours upper",
			expr: "15 24 1,15 * 1-5 /usr/bin/find",
		}, {
			desc: "hours lower",
			expr: "15 -1 1,15 * 1-5 /usr/bin/find",
		},
		{
			desc: "day of month upper",
			expr: "15 0 32 * 1-5 /usr/bin/find",
		}, {
			desc: "day of month lower",
			expr: "15 0 -1 * 1-5 /usr/bin/find",
		},
		{
			desc: "month upper",
			expr: "15 0 1,15 13 1-5 /usr/bin/find",
		}, {
			desc: "month lower",
			expr: "15 0 1,15 0 1-5 /usr/bin/find",
		},
		{
			desc: "day of week upper",
			expr: "15 24 1,15 * 7 /usr/bin/find",
		}, {
			desc: "day of week lower",
			expr: "15 -1 1,15 * -1 /usr/bin/find",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parsed, err := Parse(tC.expr)

			if parsed != "" {
				t.Fatalf("expected output to be empty")
			}

			if err == nil {
				t.Fatalf("expected an error to occur")
			}

		})
	}
}

func TestOddErrors(t *testing.T) {
	testCases := []struct {
		desc string
		expr string
	}{
		{
			desc: "out of bounds days of month",
			expr: "15 0 1,32 * 1-5 /usr/bin/find",
		}, {
			desc: "with text",
			expr: "a 0 1,15 * 1-5 /usr/bin/find",
		}, {
			desc: "with bad interval",
			expr: "15 0 1,15 * 5-1 /usr/bin/find",
		}, {
			desc: "with bad interval (non number)",
			expr: "15 0 1,15 * 1-five /usr/bin/find",
		}, {
			desc: "with mixed blocks types",
			expr: "15 0 1,15-3 * 1-5 /usr/bin/find",
		}, {
			desc: "with string in blocks",
			expr: "*/five 0 1,15 * 1-5 /usr/bin/find",
		}, {
			desc: "with unexpected number of arguments",
			expr: "15 0 1,15 * 1-5",
		}, {
			desc: "with multiple forward slashes",
			expr: "*/15/21 0 1,15 * 1-5",
		}, {
			desc: "with forward slashe without star",
			expr: "12/21 0 1,15 * 1-5",
		}, {
			desc: "with multiple dashes",
			expr: "*/15 0 1,15 * 1-5-8",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parsed, err := Parse(tC.expr)

			// THEN
			if parsed != "" {
				t.Log("want: ")
				t.Logf("got: %s", parsed)
				t.Fatalf("expected output to be empty")
			}

			if err == nil {
				t.Fatalf("expected an error to occur")
			}

		})
	}
}

func getInterval(start, end int) string {
	intervals := []string{}
	for i := start; i <= end; i++ {
		intervals = append(intervals, strconv.Itoa(i))
	}

	return strings.Join(intervals, " ")
}

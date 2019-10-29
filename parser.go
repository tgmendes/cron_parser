package cron

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var allowedRanges = map[string][]int{
	"minute":     []int{0, 59},
	"hour":       []int{0, 23},
	"dayOfMonth": []int{1, 31},
	"month":      []int{1, 12},
	"dayOfWeek":  []int{0, 6},
}

// Parse takes a cron expression as an input, and will return a table
// with the times at which it will run.
func Parse(expr string) (string, error) {
	splitExpr := strings.Split(expr, " ")
	if len(splitExpr) != 6 {
		return "", errors.New("malformed cron expression")
	}

	minute, err := parseBlock(splitExpr[0], "minute")
	if err != nil {
		return "", err
	}

	hour, err := parseBlock(splitExpr[1], "hour")
	if err != nil {
		return "", err
	}

	dayOfMonth, err := parseBlock(splitExpr[2], "dayOfMonth")
	if err != nil {
		return "", err
	}

	month, err := parseBlock(splitExpr[3], "month")
	if err != nil {
		return "", err
	}

	dayOfWeek, err := parseBlock(splitExpr[4], "dayOfWeek")
	if err != nil {
		return "", err
	}

	command := splitExpr[5]

	output := fmt.Sprintf(`minute        %s
hour          %s
day of month  %s
month         %s
day of week   %s
command       %s
`, minute, hour, dayOfMonth, month, dayOfWeek, command)

	return output, nil
}

func parseBlock(block, dateType string) (string, error) {
	allowedRange := allowedRanges[dateType]
	if block == "*" {
		return getRange(allowedRange[0], allowedRange[1], 1), nil
	}

	// Cron expression with steps
	everyNthBlock := strings.Split(block, "/")

	// e.g. 12/15/19
	if len(everyNthBlock) > 2 {
		return "", fmt.Errorf("invalid cron block %s", block)
	}

	if len(everyNthBlock) == 2 {
		// e.g. 23/12 (not standard cron)
		if everyNthBlock[0] != "*" {
			return "", fmt.Errorf("invalid cron block %s", block)
		}

		step, err := strconv.Atoi(everyNthBlock[1])
		if err != nil {
			return "", err
		}

		return getRange(allowedRange[0], allowedRange[1], step), nil
	}

	// Cron expression with intervals
	interval := strings.Split(block, "-")

	// e.g. 2-4-6
	if len(interval) > 2 {
		return "", fmt.Errorf("invalid cron block %s", block)
	}

	if len(interval) == 2 {
		start, err := strconv.Atoi(interval[0])
		if err != nil {
			return "", err
		}
		end, err := strconv.Atoi(interval[1])
		if err != nil {
			return "", err
		}

		if end < start {
			return "", fmt.Errorf("invalid interval block for %s", dateType)
		}

		return getRange(start, end, 1), nil
	}

	// Cron expression with specific time value(s)
	timeValues := strings.Split(block, ",")
	times := []string{}

	for _, t := range timeValues {
		tInt, err := strconv.Atoi(t)
		if err != nil {
			return "", err
		}

		if tInt < allowedRange[0] || tInt > allowedRange[1] {
			return "", fmt.Errorf("%s exceed allowed value of %d", dateType, allowedRange[1])
		}

		times = append(times, t)
	}

	return strings.Join(times, " "), nil
}

func getRange(start, stop, step int) string {
	r := []string{}
	for i := start; i <= stop; i = i + step {
		r = append(r, strconv.Itoa(i))
	}

	return strings.Join(r, " ")
}

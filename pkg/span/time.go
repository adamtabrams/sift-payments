package span

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Time struct {
	Years  []int
	Months []int
}

func Parse(years, months string) (*Time, error) {
	y, err := parseTimeSpan(years, strconv.Atoi)
	if err != nil {
		return nil, err
	}

	m, err := parseTimeSpan(months, parseMonths)
	if err != nil {
		return nil, err
	}

	return &Time{
		Years:  y,
		Months: m,
	}, nil
}

func parseTimeSpan(text string, parseFn func(string) (int, error)) ([]int, error) {
	if text == "" {
		return nil, nil
	}

	if strings.Contains(text, ",") {
		res := []int{}
		for _, v := range strings.Split(text, ",") {
			m, err := parseTimeSpan(v, parseFn)
			if err != nil {
				return nil, err
			}

			res = append(res, m...)
		}

		return res, nil
	}

	if strings.Contains(text, "-") {
		v := strings.Split(text, "-")
		start, err := parseFn(v[0])
		if err != nil {
			return nil, err
		}

		end, err := parseFn(v[1])
		if err != nil {
			return nil, err
		}

		res := []int{}
		for i := start; i <= end; i++ {
			res = append(res, i)
		}

		return res, nil
	}

	m, err := parseFn(text)
	if err != nil {
		return nil, err
	}

	return []int{m}, nil
}

func parseMonths(month string) (int, error) {
	for i := 1; i <= 12; i++ {
		m := time.Month(i)
		switch {
		case month == fmt.Sprint(i):
			return i, nil
		case len(month) > 2 && strings.EqualFold(month[:3], m.String()[:3]):
			return i, nil
		}
	}

	return 0, errors.Errorf("not a valid month: %v", month)
}

func (t *Time) HasYears() bool {
	return len(t.Years) > 0
}

func (t *Time) HasMonths() bool {
	return len(t.Months) > 0
}

func (t *Time) Contains(currentYear int, date time.Time) bool {
	year, month, _ := date.Date()
	var yearsMatch, monthsMatch bool

	for _, y := range t.Years {
		if y == year {
			yearsMatch = true
		}
	}

	for _, m := range t.Months {
		if m == int(month) {
			monthsMatch = true
		}
	}

	switch {
	case t.HasYears() && t.HasMonths():
		return yearsMatch && monthsMatch
	case t.HasYears():
		return yearsMatch
	case t.HasMonths():
		return year == currentYear && monthsMatch
	default:
		return year == currentYear
	}
}

func (t *Time) Key(date time.Time) string {
	y, m, _ := date.Date()

	if t.HasYears() && !t.HasMonths() {
		return fmt.Sprintf("%v", y)
	}

	return fmt.Sprintf("%v-%v", y, int(m))
}

func (t *Time) Name(key string) any {
	if !strings.Contains(key, "-") {
		name, _ := strconv.Atoi(key)

		return name
	}

	subkeys := strings.Split(key, "-")
	year := subkeys[0]
	v, _ := strconv.Atoi(subkeys[1])
	month := time.Month(v)

	if t.HasYears() && t.HasMonths() {
		return fmt.Sprintf("%v %v", month, year)
	}

	return month.String()
}

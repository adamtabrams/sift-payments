package sift

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Record struct {
	Name   string
	Amount int
	Date   time.Time
}

type Header struct {
	ID, Name, Date, Amount string
}

func (h *Header) ParseRecords(dateFormat string, table [][]string) (map[string]Record, error) {
	res := make(map[string]Record, len(table[1:]))

	iID, iName, iDate, iAmount, err := h.index(table[0])
	if err != nil {
		return nil, err
	}

	for _, row := range table[1:] {
		amount, err := parseAmount(row[iAmount])
		if err != nil {
			return nil, err
		}

		date, err := time.Parse(dateFormat, row[iDate])
		if err != nil {
			return nil, fmt.Errorf("cannot parse time %s: %s", date, err)
		}

		key := row[iID]
		res[key] = Record{
			Name:   row[iName],
			Amount: amount,
			Date:   date,
		}
	}

	return res, nil
}

func (h *Header) index(row []string) (id, name, date, amount int, err error) {
	found := 0

	for i, x := range row {
		switch x {
		case h.ID:
			id = i
			found++
		case h.Name:
			name = i
			found++
		case h.Date:
			date = i
			found++
		case h.Amount:
			amount = i
			found++
		}
	}

	if found != 4 {
		return 0, 0, 0, 0, fmt.Errorf("cannot parse header indexes")
	}

	return id, name, date, amount, err
}

func parseAmount(s string) (int, error) {
	s = strings.Trim(s, `"`)
	s = strings.ReplaceAll(s, ",", "")

	if strings.HasPrefix(s, "-$") {
		s = "-" + strings.TrimPrefix(s, "-$")
	}

	s = strings.TrimPrefix(s, "$")
	i := strings.Index(s, ".")

	amount, err := strconv.Atoi(s[:i])
	if err != nil {
		return 0, fmt.Errorf("cannot parse amount %s: %s", s, err)
	}

	return amount, nil
}

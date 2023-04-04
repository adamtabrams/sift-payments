package pay

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type RecordMap map[string]Record

type Record struct {
	Name   string
	Amount int
	Date   time.Time
}

func (r Record) Payment(category string) Payment {
	return Payment{
		Name:     r.Name,
		Amount:   r.Amount,
		Category: category,
	}
}

func RecordMapFromBytes(config *Config, csvData []byte) (RecordMap, error) {
	r := csv.NewReader(bytes.NewReader(csvData))

	table, err := r.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse csv")
	}

	return RecordMapFromTable(config, table[0], table[1:])
}

func RecordMapFromTable(config *Config, header []string, rowList [][]string) (RecordMap, error) {
	rm := make(RecordMap, len(rowList))

	index, err := indexFromHeader(config, header)
	if err != nil {
		return nil, err
	}

	for _, row := range rowList {
		amount, err := parseAmount(row[index.amount])
		if err != nil {
			return nil, err
		}

		date, err := parseDate(config, row[index.date])
		if err != nil {
			return nil, err
		}

		key := row[index.id]
		rm[key] = Record{
			Name:   row[index.name],
			Amount: amount,
			Date:   date,
		}
	}

	return rm, nil
}

func parseAmount(s string) (int, error) {
	s = strings.Trim(s, `"`)
	s = strings.ReplaceAll(s, ",", "")

	if strings.HasPrefix(s, "-") {
		s = "-" + strings.TrimPrefix(s, "-$")
	}

	s = strings.TrimPrefix(s, "$")
	i := strings.Index(s, ".")

	amount, err := strconv.Atoi(s[:i])
	if err != nil {
		return 0, errors.Wrapf(err, "cannot parse amount: %v", s)
	}

	return amount, nil
}

func parseDate(config *Config, s string) (time.Time, error) {
	t, err := time.Parse(config.Format.Date, s)
	if err != nil {
		return t, errors.Wrapf(err, "cannot parse time: %v", s)
	}

	return t, nil
}

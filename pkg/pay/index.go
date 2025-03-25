package pay

import (
	"github.com/pkg/errors"
)

type index struct {
	id     int
	name   int
	amount int
	date   int
}

func indexFromHeader(config *Config, header []string) (*index, error) {
	index := &index{}
	numFound := 0

	for i, col := range header {
		switch col {
		case config.Header.ID:
			index.id = i
			numFound++
		case config.Header.Name:
			index.name = i
			numFound++
		case config.Header.Amount:
			index.amount = i
			numFound++
		case config.Header.Date:
			index.date = i
			numFound++
		}
	}

	if numFound < 4 {
		return nil, errors.Errorf("cannot index header: %+v following config: %+v", header, config.Header)
	}

	return index, nil
}

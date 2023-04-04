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

// TODO fix initialization
// func indexFromHeader(config *Config, header []string) (index *index, err error) {
func indexFromHeader(config *Config, header []string) (*index, error) {
	numFound := 0
	index := &index{}
	// fmt.Printf("debug: %+v\n", header)

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
		// TODO improve error message
		return nil, errors.Errorf("cannot use config to create valid index: %+v", index)
	}

	return index, nil
}

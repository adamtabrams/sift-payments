package pay

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Paths      PathsConfig
	Header     HeaderConfig
	Format     FormatConfig
	Categories []string
}

type PathsConfig struct {
	RulesFile  string
	RecordsDir string
}

type HeaderConfig struct {
	ID     string
	Name   string
	Date   string
	Amount string
}

type FormatConfig struct {
	Date string
}

// TODO set default paths

func ConfigFromBytes(yamlData []byte) (c *Config, err error) {
	err = yaml.Unmarshal(yamlData, &c)
	if err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal config")
	}

	return c, nil
}

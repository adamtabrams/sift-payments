package pay

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
)

func ConfigFromFile(filePath string) (*Config, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read config file: %v", filePath)
	}

	config, err := ConfigFromBytes(bytes)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) RecordMapFromDir() (RecordMap, error) {
	dirPath := c.Paths.RecordsDir
	globPath := filepath.Join(dirPath, "*.csv")

	pathList, err := filepath.Glob(globPath)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read records dir: %v", dirPath)
	}

	rm := make(RecordMap, len(pathList))
	for _, path := range pathList {
		recordMap, err := c.RecordMapFromFile(path)
		if err != nil {
			return nil, err
		}

		maps.Copy(rm, recordMap)
	}

	return rm, nil
}

func (c *Config) RecordMapFromFile(filePath string) (RecordMap, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read records file: %v", filePath)
	}

	recordMap, err := RecordMapFromBytes(c, bytes)
	if err != nil {
		return nil, err
	}

	return recordMap, nil
}

func (c *Config) RuleListFromFile() (RuleList, error) {
	filePath := c.Paths.RulesFile

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read rules file: %v", filePath)
	}

	ruleList, err := RuleListFromBytes(bytes)
	if err != nil {
		return nil, err
	}

	return ruleList, nil
}

func (c *Config) AddRuleToFile(rule Rule) error {
	file, err := os.OpenFile(c.Paths.RulesFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrap(err, "cannot open rules files")
	}
	defer file.Close()

	bytes, err := rule.Bytes()
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)
	if err != nil {
		return errors.Wrap(err, "cannot append to rules files")
	}

	return nil
}

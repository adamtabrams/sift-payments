package pay

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
)

func ConfigFromFile() (*Config, error) {
	// TODO support XDG
	const filePath = "config.yaml"

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

func RecordMapFromDir(config *Config) (RecordMap, error) {
	dirPath := config.Paths.RecordsDir

	fileInfoList, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read records dir: %v", dirPath)
	}

	rm := make(RecordMap)
	for _, fileInfo := range fileInfoList {
		// TODO consider using glob
		if !strings.HasSuffix(fileInfo.Name(), ".csv") {
			continue
		}
		// TODO handle invalid files
		// TODO improve path handling
		recordMap, err := RecordMapFromFile(config, dirPath + "/" + fileInfo.Name())
		if err != nil {
			return nil, err
		}

		maps.Copy(rm, recordMap)
	}

	return rm, nil
}

func RecordMapFromFile(config *Config, filePath string) (RecordMap, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read records file: %v", filePath)
	}

	recordMap, err := RecordMapFromBytes(config, bytes)
	if err != nil {
		return nil, err
	}

	return recordMap, nil
}

func RuleListFromFile(config *Config) (RuleList, error) {
	filePath := config.Paths.RulesFile

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

// TODO implement
func AddRuleToFile(config *Config, rule Rule) error {
	file, err := os.OpenFile(config.Paths.RulesFile, os.O_APPEND|os.O_WRONLY, 0600)
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

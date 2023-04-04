package pay

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type RuleListMap map[string]RuleList

type RuleList []Rule

type Rule struct {
	Name     string
	Amount   int `yaml:",omitempty"`
	Category string
}

func RuleListFromBytes(yamlData []byte) (rl RuleList, err error) {
	err = yaml.Unmarshal(yamlData, &rl)
	if err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal rules")
	}

	return rl, nil
}

func NewRule(name, category string, amount int) Rule {
	return Rule{
		Name:     KeyFromRecordName(name),
		Amount:   amount,
		Category: category,
	}
}

func (r Rule) Bytes() ([]byte, error) {
	list := RuleList{r}

	bytes, err := yaml.Marshal(list)
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal rule")
	}

	return bytes, nil
}

// TODO this mutates
func (rlm RuleListMap) AppendRule(rule Rule) {
	// key := rlm.Key(rule.Name)
	key := rule.Name
	existing := rlm[key]
	rlm[key] = append(existing, rule)
}

func (rl RuleList) Map() RuleListMap {
	rlm := make(RuleListMap)

	for _, rule := range rl {
		rlm.AppendRule(rule)
	}

	return rlm
}

func (rlm RuleListMap) Category(record Record) (category string, ok bool) {
	// key := rlm.Key(record.Name)
	key := KeyFromRecordName(record.Name)
	// fmt.Println("debug1 - " + record.Name)
	// fmt.Println("debug2 - " + key)
	ruleList := rlm[key]

	if len(ruleList) == 0 {
		return "", false
	}

	for _, rule := range ruleList {
		if rule.Amount == record.Amount {
			return rule.Category, true
		}
		if rule.Amount == 0 {
			category = rule.Category
			ok = true
		}
	}

	return category, ok
}

// TODO make private
func KeyFromRecordName(name string) string {
	r := regexp.MustCompile(`\S*[0-9]\S*`)
	s := r.ReplaceAllString(name, "")
	s = strings.TrimSpace(s)

	return s
}

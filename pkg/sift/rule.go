package sift

import (
	"strings"
)

type RuleTable map[string][]Rule

type Rule struct {
	Name     string
	Amount   int `yaml:",omitempty"`
	Category string
}

func NewRule(name, category string, amount int) Rule {
	return Rule{
		Name:     getKey(name),
		Amount:   amount,
		Category: category,
	}
}

func (rt RuleTable) Append(rule Rule) {
	key := rule.Name
	existing := rt[key]
	rt[key] = append(existing, rule)
}

func (rt RuleTable) Category(rec Record) (string, bool) {
	key := getKey(rec.Name)
	rules := rt[key]

	if len(rules) == 0 {
		return "", false
	}

	var res string
	var ok bool
	for _, r := range rules {
		if r.Amount == rec.Amount {
			return r.Category, true
		}
		if r.Amount == 0 {
			res = r.Category
			ok = true
		}
	}

	return res, ok
}

func getKey(recordName string) string {
	var res []string

	for word := range strings.SplitSeq(recordName, " ") {
		if len(word) == 0 {
			continue
		}
		if strings.ContainsAny(word, "0123456789") {
			word = "#"
		}
		res = append(res, word)
	}

	return strings.Join(res, " ")
}

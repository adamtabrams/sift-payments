package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/adamtabrams/sift-payments/pkg/pay"
	"gopkg.in/yaml.v3"
)

func main() {
	config, err := pay.ConfigFromFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	ruleList, err := pay.RuleListFromFile(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	ruleListByName := ruleList.Map()

	recordByID, err := pay.RecordMapFromDir(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	flags := parseFlags()
	categoriesByNum := categoriesByNum(config)

	paymentListByTime := make(pay.PaymentListMap)

	for _, record := range recordByID {
		if !timeMatches(flags, record.Date) {
			continue
		}

		category, ok := ruleListByName.Category(record)
		if !ok {
			rule := pay.NewRule(record.Name, "skipped", 0)
			save := false

			// TODO replace tesing logic for conditional
			if len(flags) == 0 || flags[0] != "skip" {
				rule, save = promptUser(categoriesByNum, record)
			}

			if save {
				pay.AddRuleToFile(config, rule)
			}

			ruleListByName.AppendRule(rule)
			category = rule.Category
		}

		payment := record.Payment(category)
		key := keyFromTime(flags, record.Date)
		paymentList := paymentListByTime[key]
		paymentListByTime[key] = append(paymentList, payment)
	}

	summary := paymentListByTime.Summary()

	bytes, err := yaml.Marshal(summary)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\n%s\n", bytes)

	// TODO handle printing for category
}

func parseFlags() []string {
	return nil
}

func timeMatches(flags []string, date time.Time) bool {
	return true
}

func keyFromTime(flags []string, date time.Time) string {
	return "2023"
}

// TODO use STDERR
func promptUser(categoriesByNum map[int]string, record pay.Record) (pay.Rule, bool) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\nCould not sort expense of %v dollars for:\n%v\n\n", record.Amount, record.Name)

	// TODO consider replacing with function that prints from list
	cBytes, _ := yaml.Marshal(categoriesByNum)
	fmt.Printf("%s", cBytes)

	fmt.Printf("\nSelect a category (blank to skip): ")
	numRaw, _ := reader.ReadString('\n')
	num := strings.TrimSpace(numRaw)

	category := "skipped"
	if num != "" {
		n, _ := strconv.Atoi(num)
		if c, ok := categoriesByNum[n]; ok {
			category = c
		}
	}

	if category == "skipped" {
		return pay.NewRule(record.Name, "skipped", 0), false
	}

	fmt.Printf("\nAlso require amount to match %v dollars (y/N): ", record.Amount)
	matchAmountRaw, _ := reader.ReadString('\n')
	matchAmount := strings.TrimSpace(matchAmountRaw)

	amount := 0
	if matchAmount == "y" {
		amount = record.Amount
	}

	rule := pay.NewRule(record.Name, category, amount)
	rBytes, _ := yaml.Marshal(rule)
	fmt.Printf("\n%s", rBytes)

	fmt.Printf("\nSave to rules (y/N): ")
	saveRaw, _ := reader.ReadString('\n')
	save := strings.TrimSpace(saveRaw)

	return rule, save == "y"
}

func categoriesByNum(config *pay.Config) map[int]string {
	categoryList := config.Categories
	categoryMap := make(map[int]string, len(categoryList))

	for i, category := range config.Categories {
		categoryMap[i+1] = category
	}

	return categoryMap
}

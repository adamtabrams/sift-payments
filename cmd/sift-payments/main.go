package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/adamtabrams/sift-payments/pkg/pay"
	"github.com/adamtabrams/sift-payments/pkg/span"
	"gopkg.in/yaml.v3"
)

var version = "v0.0.1"

func main() {
	log.SetFlags(log.Lshortfile)

	var printVersion, skipPrompt bool
	flag.BoolVar(&printVersion, "v", false, "prints the current version")
	flag.BoolVar(&skipPrompt, "s", false, "skip user prompts")

	var yearFlag, monthFlag, configPath string
	flag.StringVar(&yearFlag, "y", "", "view payments for a given year")
	flag.StringVar(&monthFlag, "m", "", "view payments for a given month")
	flag.StringVar(&configPath, "config", "config.yaml", "path to config file")

	// var category string
	// flag.StringVar(&category, "c", "", "view payments for a given category")

	flag.Parse()
	if printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	timespan, err := span.Parse(yearFlag, monthFlag)
	if err != nil {
		log.Fatal(err)
	}

	config, err := pay.ConfigFromFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	ruleList, err := config.RuleListFromFile()
	if err != nil {
		log.Fatal(err)
	}

	recordByID, err := config.RecordMapFromDir()
	if err != nil {
		log.Fatal(err)
	}

	categoriesByNum := categoriesByNum(config)
	ruleListByName := ruleList.Map()
	paymentListByTime := make(pay.PaymentListMap)
	currentYear := time.Now().Year()

	for _, record := range recordByID {
		if !timespan.Contains(currentYear, record.Date) {
			continue
		}

		category, ok := ruleListByName.Category(record)
		if !ok {
			rule := pay.NewRule(record.Name, "skipped", 0)
			save := false

			if !skipPrompt {
				rule, save = promptUser(categoriesByNum, record)
			}

			if save {
				err := config.AddRuleToFile(rule)
				if err != nil {
					log.Fatal(err)
				}
			}

			ruleListByName.AppendRule(rule)
			category = rule.Category
		}

		payment := record.Payment(category)
		key := timespan.Key(record.Date)
		paymentList := paymentListByTime[key]
		paymentListByTime[key] = append(paymentList, payment)
	}

	// Get summaries for each time period
	summaryMap := paymentListByTime.SummaryMap()

	// Get all keys
	keyList := make([]string, 0, len(summaryMap))
	for key := range summaryMap {
		keyList = append(keyList, key)
	}

	// Sort all keys
	sort.Strings(keyList)

	// Convert name, create map, marshal text
	for _, key := range keyList {
		v := make(map[any]pay.Summary, 1)
		name := timespan.Name(key)
		v[name] = summaryMap[key]

		bytes, err := yaml.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\n%s\n", bytes)
	}

	// TODO handle printing for category
}

func promptUser(categoriesByNum map[int]string, record pay.Record) (pay.Rule, bool) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Fprintf(os.Stderr, "\nCould not sort expense of %v dollars for:\n%v\n\n", record.Amount, record.Name)

	// TODO consider replacing with function that prints from list
	cBytes, _ := yaml.Marshal(categoriesByNum)
	fmt.Fprintf(os.Stderr, "%s", cBytes)

	fmt.Fprintf(os.Stderr, "\nSelect a category (blank to skip): ")
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

	fmt.Fprintf(os.Stderr, "\nAlso require amount to match %v dollars (y/N): ", record.Amount)
	matchAmountRaw, _ := reader.ReadString('\n')
	matchAmount := strings.TrimSpace(matchAmountRaw)

	amount := 0
	if matchAmount == "y" {
		amount = record.Amount
	}

	rule := pay.NewRule(record.Name, category, amount)
	rBytes, _ := yaml.Marshal(rule)
	fmt.Fprintf(os.Stderr, "\n%s", rBytes)

	fmt.Fprintf(os.Stderr, "\nSave to rules (y/N): ")
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

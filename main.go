package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/adamtabrams/sift-payments/pkg/sift"
	"github.com/alexflint/go-arg"
	"gopkg.in/yaml.v3"
)

// var version = "v0.0.2"

func (Flags) Version() string {
	return "v0.0.2"
}

type Flags struct {
	Years    Range  `arg:"-y,--,required" help:"view specific year(s)"`
	Months   Range  `arg:"-m,--" help:"view specific month(s)"`
	Category string `arg:"-c,--" help:"view specific category"`
	Path     string `arg:"-p,--" help:"path to config file" default:"config.yaml"`
	NoPrompt bool   `arg:"-n,--no-prompt" help:"skip user prompts"`
}

type Range []int

type ConfigFile struct {
	RulesFile  string
	RecordsDir string
	DateFormat string
	Header     *sift.Header
	Categories []string
}

func main() {
	f := Flags{}
	arg.MustParse(&f)
	log.SetFlags(log.Lshortfile)

	configYaml, err := os.ReadFile(f.Path)
	assertNil(err, "cannot read config file %s", f.Path)

	config := &ConfigFile{}
	err = yaml.Unmarshal(configYaml, &config)
	assertNil(err, "cannot unmarshal configs")

	rulesYaml, err := os.ReadFile(config.RulesFile)
	assertNil(err, "cannot read rules file %s", config.RulesFile)

	var rules []sift.Rule
	err = yaml.Unmarshal(rulesYaml, &rules)
	assertNil(err, "cannot unmarshal rules")

	ruleTable := make(sift.RuleTable)
	for _, r := range rules {
		ruleTable.Append(r)
	}

	records, err := readRecordsDir(config)
	assertNil(err, "cannot parse records dir %s", config.RecordsDir)

	keyLayout := "2006"
	outLayout := "2006"
	maps.DeleteFunc(records, func(_ string, v sift.Record) bool {
		return !slices.Contains(f.Years, v.Date.Year())
	})

	if f.Months != nil {
		keyLayout = "2006-01"
		outLayout = "January 2006"
		maps.DeleteFunc(records, func(_ string, v sift.Record) bool {
			return !slices.Contains(f.Months, int(v.Date.Month()))
		})
	}

	summaries := make(map[string]sift.Summary)
	categoryRes := []sift.Record{}
	count := 0

	for _, r := range records {
		count++
		// c := getCategory(config, f.NoPrompt, ruleTable)
		c := getCategory(config, f.NoPrompt, ruleTable, r, count, len(records))
		if f.Category == c {
			categoryRes = append(categoryRes, r)
			continue
		}

		key := r.Date.Format(keyLayout)
		s := summaries[key]
		summaries[key] = s.Add(r, c)
	}

	if f.Category != "" {
		slices.SortFunc(categoryRes, func(a, b sift.Record) int {
			return a.Date.Compare(b.Date)
		})
		for _, r := range categoryRes {
			fmt.Printf("\n%s  %v\n%s\n", r.Date.Format("2006-01-02"), r.Amount, r.Name)
		}

		return
	}

	keys := slices.Collect(maps.Keys(summaries))
	slices.Sort(keys)

	for _, k := range keys {
		t, err := time.Parse(keyLayout, k)
		assertNil(err, "cannot parse summary key %s", k)

		wrapped := map[string]sift.Summary{
			t.Format(outLayout): summaries[k],
		}
		summaryYaml, err := yaml.Marshal(wrapped)
		assertNil(err, "cannot build final summary")

		fmt.Printf("\n%s", summaryYaml)
	}
}

// func getCategory(c *ConfigFile, noPrompt bool, rt sift.RuleTable, rec sift.Record) string {
func getCategory(c *ConfigFile, noPrompt bool, rt sift.RuleTable, rec sift.Record, count, total int) string {
	category, ok := rt.Category(rec)
	if ok {
		return category
	}

	if noPrompt {
		rule := sift.NewRule(rec.Name, "skipped", 0)
		rt.Append(rule)
		return rule.Category
	}

	// rule, save := promptUser(c.Categories, rec)
	rule, save := promptUser(c.Categories, rec, count, total)

	// TODO: removed asking to save rules?
	if save {
		err := addRuleToFile(c.RulesFile, rule)
		assertNil(err, "error saving rule to file %s", c.RulesFile)
	}

	rt.Append(rule)
	return rule.Category
}

// TODO: add way to save rule for a specific date
func promptUser(categories []string, rec sift.Record, count, total int) (sift.Rule, bool) {
// func promptUser(categories []string, rec sift.Record) (sift.Rule, bool) {
	fmt.Fprintf(os.Stderr, "\nprogress: %v / %v", count, total)
	fmt.Fprintf(os.Stderr, "\nname: %s\namount: %v\n\n", rec.Name, rec.Amount)

	for i, c := range categories {
		fmt.Fprintf(os.Stderr, "%v - %s\n", i+1, c)
	}

	fmt.Fprintf(os.Stderr, "\nselect a category (blank to skip): ")
	var input string
	fmt.Scanln(&input)

	category := "skipped"
	if input != "" {
		n, err := strconv.Atoi(input)
		assertNil(err, "invalid input %s", input)
		category = categories[n-1]
	}

	if category == "skipped" {
		return sift.NewRule(rec.Name, "skipped", 0), false
	}

	fmt.Fprintf(os.Stderr, "\nrequire amount to match %v dollars (y/N): ", rec.Amount)
	fmt.Scanln(&input)

	amount := 0
	if strings.ToLower(input) == "y" {
		amount = rec.Amount
	}

	rule := sift.NewRule(rec.Name, category, amount)
	ruleYaml, err := yaml.Marshal(rule)
	assertNil(err, "error marshalling rule")

	fmt.Fprintf(os.Stderr, "\n%s\nsave rule (Y/n): ", ruleYaml)
	fmt.Scanln(&input)
	save := strings.ToLower(input) != "n"

	return rule, save
}

func readRecordsDir(config *ConfigFile) (map[string]sift.Record, error) {
	globPath := filepath.Join(config.RecordsDir, "*.csv")
	paths, err := filepath.Glob(globPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create glob: %s", err)
	}

	res := make(map[string]sift.Record, len(paths))
	for _, p := range paths {
		recordCSV, err := os.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %s", p, err)
		}

		reader := strings.NewReader(string(recordCSV))
		table, err := csv.NewReader(reader).ReadAll()
		if err != nil {
			return nil, fmt.Errorf("failed to parse csv %s: %s", p, err)
		}

		r, err := config.Header.ParseRecords(config.DateFormat, table)
		if err != nil {
			return nil, err
		}

		maps.Copy(res, r)
	}

	return res, nil
}

func addRuleToFile(filePath string, rule sift.Rule) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("cannot open rules files: %s", err)
	}
	defer file.Close()

	ruleYaml, err := yaml.Marshal([]sift.Rule{rule})
	if err != nil {
		return fmt.Errorf("cannot marshal rule: %s", err)
	}

	_, err = file.Write(ruleYaml)
	if err != nil {
		return fmt.Errorf("cannot append to rules files: %s", err)
	}

	return nil
}

func (r *Range) UnmarshalText(b []byte) error {
	s := string(b)

	if !strings.Contains(s, "-") {
		n, err := strconv.Atoi(s)
		*r = Range{n}
		return err
	}

	bounds := strings.Split(s, "-")
	if len(bounds) != 2 {
		return fmt.Errorf("invalid input")
	}

	min, err := strconv.Atoi(bounds[0])
	if err != nil {
		return err
	}

	max, err := strconv.Atoi(bounds[1])
	if err != nil {
		return err
	}

	if max <= min {
		return fmt.Errorf("invalid range")
	}

	for i := min; i <= max; i++ {
		*r = append(*r, i)
	}

	return nil
}

func assertNil(err error, format string, a ...any) {
	if err != nil {
		msg := fmt.Sprintf(format, a...)
		log.Fatalf("%s: %s", msg, err)
	}
}

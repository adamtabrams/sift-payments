## About

Sift is a simple way to generate a spending report.

I wanted an easy, private, offline way to review monthly transactions and spending.
Sift builds rules by prompting the user to categorize transactions.
Those rules are stored in a file and used to categorize future transactions.


## Setup

If `go` is already installed, just run `make` from the top of the repo.
That will build an executable at `bin/sift`.
Tests can be run with `make test`.


## Examples

Process records for January 2023:
`./sift -y 2023 -m 1`

Process records for January through March 2023:
`./sift -y 2023 -m 1-3`

Process records for 2022 and 2023:
`./sift -y 2022-2023`

Process records for each month in 2022 and 2023:
`./sift -y 2022-2023 -m 1-12`

View all transactions in the subscriptions category for Jan 2023:
`./sift -y 2023 -m 1 -c subscriptions`

By default, transactions not matching any existing rules will prompt the user.
Rules can match either on just the transaction ID or both ID and amount.
Using `--no-prompt` or `-n` puts non-matching transactions in the `skipped` category.

See all flags:
`./sift --help`


## Config

Sift is requires a `config.yaml` file like this:
```
rulesfile: rules.yaml
recordsdir: records
dateformat: 01/02/06

header:
  id: Transaction ID
  name: Description
  date: Date
  amount: Amount

categories:
  - salary
  - rent
  - food
  - insurance
  - subscriptions
```

New rules will be saved to the value of `rulefiles`.
CSV files will be read from the value of `recordsdir`.
Data in the date column of CSV files must match the value of `dateformat`.
- Format must reference `Jan 02 2006`.
- Refer to [golang time pkg](https://pkg.go.dev/time) if needed.
The values in `header` specify which columns to use.
Each value in `categories` is a budget category option to select.

Check out the `sample` dir for more info.

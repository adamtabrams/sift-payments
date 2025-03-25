#!/bin/sh

assert_eq() {
  [ $? != 0 ] && echo "EXIT - $name" >&2
  [ "$expect" = "$actual" ] && echo "pass - $name" >&2 && return
  echo "FAIL - $name" >&2
  expectFile=$(mktemp)
  actualFile=$(mktemp)
  echo "$expect" >"$expectFile"
  echo "$actual" >"$actualFile"
  diff --side-by-side "$expectFile" "$actualFile"
  rm "$expectFile" "$actualFile"
}

expect=''
actual=$(../bin/sift -n -y 2022)
name='year empty'
assert_eq
actual=$(../bin/sift -n -y 2023 -m 12)
name='year month empty'
assert_eq
actual=$(../bin/sift -n -y 2022-2023 -m 5-11)
name='year-range month-range empty'
assert_eq

expect='
February 2023:
    income: 0
    expenses: 1295
    total: -1295
    categories:
        insurance: -87
        rent: -1200
        subscriptions: -8'
actual=$(../bin/sift -n -y 2023 -m 2)
name='year month'
assert_eq
actual=$(../bin/sift -n -y 2022-2023 -m 2)
name='year-range month'
assert_eq

expect='
"2023":
    income: 4000
    expenses: 3846
    total: 154
    categories:
        food: -9
        insurance: -221
        rent: -3600
        salary: 4000
        subscriptions: -16'
actual=$(../bin/sift -n -y 2023)
name='year'
assert_eq

expect='
January 2023:
    income: 4000
    expenses: 1286
    total: 2714
    categories:
        insurance: -78
        rent: -1200
        salary: 4000
        subscriptions: -8

February 2023:
    income: 0
    expenses: 1295
    total: -1295
    categories:
        insurance: -87
        rent: -1200
        subscriptions: -8

March 2023:
    income: 0
    expenses: 1265
    total: -1265
    categories:
        food: -9
        insurance: -56
        rent: -1200'
actual=$(../bin/sift -n -y 2023 -m 1-3)
name='year month-range'
assert_eq
actual=$(../bin/sift -n -y 2022-2023 -m 1-6)
name='year-range month-range'
assert_eq

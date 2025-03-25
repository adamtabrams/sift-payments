#!/bin/sh

assert_eq() {
    [ $? != 0 ] && echo "EXIT - $name" >&2
    if [ "$expect" = "$actual" ]; then
        echo "pass - $name" >&2
    else
        echo "FAIL - $name" >&2
        expectFile=$(mktemp)
        actualFile=$(mktemp)
        echo "$expect" > "$expectFile"
        echo "$actual" > "$actualFile"
        diff --side-by-side "$expectFile" "$actualFile"
        rm "$expectFile" "$actualFile"
    fi
}

expect=""
actual=$(../sift-payments -s -m Dec)
name="month empty"
assert_eq
actual=$(../sift-payments -s -y 2022)
name="year empty"
assert_eq
actual=$(../sift-payments -s -y 2022 -m Dec)
name="year month empty"
assert_eq
actual=$(../sift-payments -s -y 2021,2022 -m 1-11)
name="year month range empty"
assert_eq

expect="
February:
    income: 0
    expenses: 1295
    total: -1295
    categories:
        insurance: -87
        rent: -1200
        subscriptions: -8"
actual=$(../sift-payments -s -m Feb)
name="month single"
assert_eq

expect="
January:
    income: 4000
    expenses: 1286
    total: 2714
    categories:
        insurance: -78
        rent: -1200
        salary: 4000
        subscriptions: -8


March:
    income: 0
    expenses: 1265
    total: -1265
    categories:
        food: -9
        insurance: -56
        rent: -1200"
actual=$(../sift-payments -s -m 1,3)
name="month comma"
assert_eq

expect="
January:
    income: 4000
    expenses: 1286
    total: 2714
    categories:
        insurance: -78
        rent: -1200
        salary: 4000
        subscriptions: -8


February:
    income: 0
    expenses: 1295
    total: -1295
    categories:
        insurance: -87
        rent: -1200
        subscriptions: -8


March:
    income: 0
    expenses: 1265
    total: -1265
    categories:
        food: -9
        insurance: -56
        rent: -1200"
actual=$(../sift-payments -s -m 1-march)
name="month dash"
assert_eq
actual=$(../sift-payments -s)
name="default"
assert_eq

expect="
2023:
    income: 4000
    expenses: 3846
    total: 154
    categories:
        food: -9
        insurance: -221
        rent: -3600
        salary: 4000
        subscriptions: -16"
actual=$(../sift-payments -s -y 2023)
name="year single"
assert_eq

expect="
February 2023:
    income: 0
    expenses: 1295
    total: -1295
    categories:
        insurance: -87
        rent: -1200
        subscriptions: -8"
actual=$(../sift-payments -s -y 2023 -m Feb)
name="year month single"
assert_eq

expect="
January 2023:
    income: 4000
    expenses: 1286
    total: 2714
    categories:
        insurance: -78
        rent: -1200
        salary: 4000
        subscriptions: -8


March 2023:
    income: 0
    expenses: 1265
    total: -1265
    categories:
        food: -9
        insurance: -56
        rent: -1200"
actual=$(../sift-payments -s -y 2022-2023 -m 1,3)
name="year range month comma"
assert_eq

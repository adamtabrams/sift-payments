##@ Build

bin/sift: *.go pkg/*/*.go
	go build -o bin/sift

.PHONY: build
build:
	make bin/sift

.PHONY: test
test:
	cd sample && ./tests.sh

.PHONY: coverage
coverage:
	go test -coverprofile /dev/null ./...

.PHONY: lint
lint:
	golangci-lint run

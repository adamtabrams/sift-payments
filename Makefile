##@ Build

.PHONY: build
build:
	go build ./cmd/sift-payments

.PHONY: integration-test
integration-test:
	cd sample && ./tests.sh

.PHONY: test
test:
	go test ./...

.PHONY: coverage
coverage:
	go test -coverprofile /dev/null ./...

.PHONY: lint
lint:
	golangci-lint run

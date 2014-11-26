# ex : shiftwidth=2 tabstop=2 softtabstop=2 :
SHELL := /bin/sh
GOPROCS := 4
SRC := *.go

.PHONY: all
all: get-deps test vet histogram

histogram: $(SRC)
	go build

.PHONY: get-deps
get-deps:
	go get -d -v ./...

.PHONY: clean
clean:
	go clean -i ./...

.PHONY: format
format:
	go fmt ./...

coverage.out: $(SRC)
	go test -coverprofile=coverage.out

.PHONY: cov
cov: coverage.out
	go tool cover -func=coverage.out

.PHONY: htmlcov
htmlcov: coverage.out
	go tool cover -html=coverage.out

.PHONY: test
test:
	go test

.PHONY: run
run: all
	./histogram -duration=900 -threads=5000

.PHONY: vet
vet:
	go vet -x

.PHONY: install
install: cov
	go install

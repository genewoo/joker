.PHONY: build test clean test-coverage coverage-report

build:
	go build -o bin ./... 

test:
	go test ./...

test-coverage:
	go test -coverprofile=coverage.out ./...

coverage-report: test-coverage
	go tool cover -html=coverage.out -o coverage.html

clean:
	go clean
	rm -f coverage.out coverage.html
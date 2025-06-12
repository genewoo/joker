.PHONY: build test clean test-coverage coverage-report

build:
	go build -o bin ./cmd/joker

test:
@for attempt in 1 2; do \
if go test ./...; then \
break; \
fi; \
echo "Tests failed (attempt $$attempt)."; \
if [ "$$attempt" -eq 2 ]; then \
exit 1; \
fi; \
echo "Retrying..."; \
done

test-coverage:
	go test -coverprofile=coverage.out ./...

coverage-report: test-coverage
	go tool cover -html=coverage.out -o coverage.html

clean:
	go clean
	rm -f coverage.out coverage.html

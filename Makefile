.PHONY: dev-run swagger run-app test

SWAG=go run github.com/swaggo/swag/cmd/swag@v1.16.4

dev-run:
	$(SWAG) init -g ./cmd/api/main.go --output docs
	go run cmd/api/main.go

swagger:
	$(SWAG) init -g ./cmd/api/main.go --output docs

run-app:
	go run cmd/api/main.go

COVERAGE_EXCLUDE=mocks|main.go|docs|test|config.go
COVERAGE_THRESHOLD = 50

test:
	go test ./... -coverprofile=coverage.tmp -coverpkg=./... -covermode=atomic -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
    if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
        echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
        exit 1; \
    else \
        echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
    fi

docker-build:
	docker build -t bookmark_service:dev .
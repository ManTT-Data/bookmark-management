.PHONY: dev-run swagger run-app test
GIT_TAG := $(shell git describe --tags --exact-match --abbrev=0 2>/dev/null)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
IMG_TAG := latest

IMG_NAME=traman26/teaman-bookmark

ifneq ($(GIT_TAG),)
	IMG_TAG := $(GIT_TAG)
endif

export IMG_TAG

dev-run:
	swag init -g ./cmd/api/main.go --output docs
	go run cmd/api/main.go

swagger:
	swag init -g ./cmd/api/main.go --output docs

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

COVERAGE_FOLDER=./test-output

docker-test:
	mkdir -p $(COVERAGE_FOLDER)
	docker buildx build --build-arg COVERAGE_EXCLUDE="${COVERAGE_EXCLUDE}" --build-arg COVERAGE_THRESHOLD="${COVERAGE_THRESHOLD}" --progress=plain --target test -t test:test --output ./test-output .
	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

docker-build:
	docker build -t $(IMG_NAME):$(IMG_TAG) .


DOCKER_USERNAME ?=
DOCKER_PASSWORD ?=

docker-login:
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin

docker-release:
	docker push $(IMG_NAME):$(IMG_TAG)
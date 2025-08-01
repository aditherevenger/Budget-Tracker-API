# Go parameters
BINARY_NAME=budget-tracker-api
COVERAGE_FILE=coverage.out
COVERAGE_THRESHOLD=80

.PHONY: all build run test coverage clean

all: build

build:
	@echo "Building the application..."
	go build -o $(BINARY_NAME) main.go

run:
	@echo "Running the application..."
	./$(BINARY_NAME)

test:
	@echo "Running tests with coverage..."
	go test ./... -coverprofile=$(COVERAGE_FILE) -covermode=atomic

coverage: test
	@echo "Checking coverage..."
	@total_coverage=$$(go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print substr($$3, 1, length($$3)-1)}'); \
	coverage_int=$${total_coverage%.*}; \
	if [ $$coverage_int -lt $(COVERAGE_THRESHOLD) ]; then \
		echo "Coverage ($$total_coverage%) is below threshold ($(COVERAGE_THRESHOLD)%)."; \
		exit 1; \
	else \
		echo "Coverage ($$total_coverage%) meets the threshold ($(COVERAGE_THRESHOLD)%)."; \
	fi

clean:
	@echo "Cleaning up..."
	go clean
	rm -f $(BINARY_NAME) $(COVERAGE_FILE)

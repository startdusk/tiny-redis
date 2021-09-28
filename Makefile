APP_NAME=tinyredis
VERSION="v0.0.1"

all: fmt vet test testrace

run: build
	@echo "run app..."
	@chmod +x ./$(APP_NAME)
	@./$(APP_NAME)

build: all
	@rm -rf $(APP_NAME)-cli &>/dev/null
	@echo "build app..."
	@go build ./cmd/...

fmt:
	@echo "fmt code..."
	@go fmt ./...

vet:
	@echo "vet code..."
	@go vet ./...

test:
	@echo "run testing..."
	@go test ./...

testrace:
	@echo "run test race..."
	@go test -race -cpu 1,4 -timeout 7m ./...

.PHONY: \
		fmt \
		vet \
		test \
		testrace \
		build \
		run

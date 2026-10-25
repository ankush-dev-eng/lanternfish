.PHONY: all build run test vet fmt clean docker-build

APP_NAME=lanternfish
BIN_DIR=bin
PORT?=8080

all: vet test build

build:
	@mkdir -p $(BIN_DIR)
	go build -ldflags="-s -w" -o $(BIN_DIR)/$(APP_NAME) ./cmd/lanternfish

run:
	PORT=$(PORT) go run ./cmd/lanternfish

test:
	go test -v ./...

vet:
	go vet ./...

fmt:
	gofmt -w -s .

clean:
	rm -rf $(BIN_DIR) coverage.out

docker-build:
	docker build -t $(APP_NAME):latest .

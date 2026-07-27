BINARY=devagent
GOFLAGS=-ldflags="-s -w"

.PHONY: build run test lint clean docker

build:
	go build $(GOFLAGS) -o $(BINARY) .

run: build
	./$(BINARY)

test:
	go test ./... -v

lint:
	go vet ./...

clean:
	rm -f $(BINARY)

docker:
	docker build -t devagent .

docker-up:
	docker compose up -d

docker-down:
	docker compose down

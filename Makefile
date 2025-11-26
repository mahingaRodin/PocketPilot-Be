.PHONY: test build

test:
	go test -v ./... -cover

build:
	go build -o bin/pocketpilot ./cmd/api

docker-build:
	docker build -t pocketpilot-api:latest .

docker-test:
	docker run --rm -p 8080:8080 pocketpilot-api:latest
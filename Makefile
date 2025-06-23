SHELL := /bin/bash

.PHONY: lint test run docker-run docker-test

lint:
	./scripts/lint.sh

test:
	./scripts/test.sh

run:
	go run ./cmd/main.go

docker-run:
	docker build -t graphql-backend:latest .
	docker network create backend-net || true
	docker run -d --rm --name api --network=backend-net -p 8080:8080 graphql-backend:latest

docker-test:
	docker run -it --rm --network=backend-net -v $(PWD):/app -w /app --env-file tests/.env.dist golang:1.23 ./scripts/test.sh

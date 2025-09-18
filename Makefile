#!make
.PHONY: up down rmv swag test

up:
	docker-compose -f ./deployments/docker-compose.yml up -d

down:
	docker-compose -f ./deployments/docker-compose.yml down

rmv:
	docker volume rm $$(docker volume ls -q)

swag:
	swag init -g ./cmd/main.go -o ./api

test:
	go test -v -race -covermode=atomic -coverprofile=coverage.out $$(go list ./... | grep domain)
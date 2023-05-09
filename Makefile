build:
	go mod download && go build -o ./.bin/app ./cmd/app/main.go

run: build
	./.bin/app

swag:
	swag init -g internal/app/app.go

.DEFAULT_GOAL := run
.PHONY: build, run
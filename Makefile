SHELL := /bin/bash

.PHONY: docker-run

docker-run:
	@docker run --rm -it \
		-v $(PWD):/app \
		-w /app \
		--name syspector \
		golang:latest \
		go run ./cmd/example/main.go


.PHONY: test up all

all: test up

test:
	go test -v ./...

up:
	docker-compose up --build
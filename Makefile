all: lint build test

build:
	go build ./...

install: 
	./scripts/make-install.sh

lint:
	gometalinter --exclude=vendor --exclude=repos --disable-all --enable=golint --enable=vet --enable=gofmt ./...
	find . -name '*.go' | xargs gofmt -w -s

test: 
	 go test -cover ./check ./handlers

start:
	 go run cmd/web/main.go

start-dev:
	go run cmd/web/main.go -http 127.0.0.1:8000 -dev

misspell:
	find . -name '*.go' -not -path './vendor/*' | xargs misspell -error

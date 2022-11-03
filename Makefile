container_name_api = gmail-sender
COLORIZE_PASS=sed ''/PASS/s//$$(printf "$(GREEN)PASS$(RESET)")/''
COLORIZE_FAIL=sed ''/FAIL/s//$$(printf "$(RED)FAIL$(RESET)")/''
SHELL=/bin/bash

.PHONY: bin build test start stop

bin:
	go build -a -tags "netgo" -installsuffix netgo -ldflags="-s -w -extldflags \"-static\"" -o build/bin/ ./...

build: bin
	docker build -t $(container_name_api):latest -f build/Dockerfile .

test: 
	gofmt -l .
	go vet -v ./...
	staticcheck ./...
	go test -v ./...  | $(COLORIZE_PASS) | $(COLORIZE_FAIL)

start:
	docker compose -f deployment/compose-local.yml up -d

stop:
	docker compose -f deployment/compose-local.yml down

clean:
	rm -rf build/bin/*

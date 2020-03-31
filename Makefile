.PHONY: default install test
all: default install test

VERSION=v1.0.0

gosec:
	go get github.com/securego/gosec/cmd/gosec

sec:
	@gosec ./...
	@echo "[OK] Go security check was completed!"

proxy:
	export GOPROXY=https://goproxy.cn


default: proxy
	go fmt ./...&&revive .&&goimports -w .&&golangci-lint run --enable-all

install: proxy
	go install -ldflags="-s -w" ./...


test: proxy
	go test ./...

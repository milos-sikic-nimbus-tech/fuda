SHELL := /bin/bash
ENV := set -a; [ -f .env ] && . ./.env; set +a;

.PHONY: dev dev-api dev-client check check-api check-client fmt build client docker desktop desktop-macos desktop-windows clean

VERSION ?= dev
GITHUB_CLIENT_ID ?=
ARCH ?= $(shell go env GOARCH)
DESKTOP_LDFLAGS := -X main.version=$(VERSION) -X main.githubClientID=$(GITHUB_CLIENT_ID)

dev:
	$(MAKE) -j2 dev-api dev-client

dev-api:
	go -C api build -o ../bin/fuda ./cmd/fuda
	$(ENV) ./bin/fuda

dev-client: client/node_modules
	pnpm -C client dev

check: check-api check-client

check-api:
	cd api && golangci-lint run ./... && go test ./...

check-client: client/node_modules
	pnpm -C client lint
	pnpm -C client fmt:check
	pnpm -C client typecheck
	pnpm -C client test

fmt: client/node_modules
	cd api && golangci-lint fmt ./...
	pnpm -C client fmt

client: client/node_modules
	pnpm -C client build

build: client
	go -C api build -o ../bin/fuda ./cmd/fuda

desktop: desktop-macos desktop-windows

desktop-macos: client
	CGO_ENABLED=1 GOOS=darwin GOARCH=$(ARCH) MACOSX_DEPLOYMENT_TARGET=13.0 CGO_LDFLAGS=-mmacosx-version-min=13.0 \
		go -C api build -trimpath -ldflags="-s -w $(DESKTOP_LDFLAGS)" -o ../bin/fuda-desktop ./cmd/fuda-desktop

desktop-windows: client
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 \
		go -C api build -trimpath -ldflags="-s -w -H=windowsgui $(DESKTOP_LDFLAGS)" -o ../bin/fuda-desktop.exe ./cmd/fuda-desktop

docker:
	docker build -t fuda .

clean:
	rm -rf bin .fuda-cache
	find api/internal/web/dist -mindepth 1 ! -name .gitkeep -delete

client/node_modules: client/package.json client/pnpm-lock.yaml
	pnpm -C client install --frozen-lockfile
	touch client/node_modules

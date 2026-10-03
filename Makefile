SHELL := /bin/bash
ENV := set -a; [ -f .env ] && . ./.env; set +a;

.PHONY: dev dev-api dev-client webhook check check-api check-client fmt build client docker clean

dev:
	$(MAKE) -j2 dev-api dev-client

dev-api:
	go -C api build -o ../bin/fuda ./cmd/fuda
	$(ENV) ./bin/fuda

dev-client: client/node_modules
	pnpm -C client dev

webhook:
	$(ENV) scripts/webhook.sh $(or $(HOST),github)

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

docker:
	docker build -t fuda .

clean:
	rm -rf bin .fuda-cache
	find api/internal/web/dist -mindepth 1 ! -name .gitkeep -delete

client/node_modules: client/package.json client/pnpm-lock.yaml
	pnpm -C client install --frozen-lockfile
	touch client/node_modules

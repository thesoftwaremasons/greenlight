
include .envrc
current_time = $(shell date --iso-8601=seconds)
git_description=$(shell git describe --always --dirty)
linker_flags = '-s -X main.buildTime=${current_time} main.version=${git_description}'



help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y]

.PHONY: run/api
run/api:
	go run ./cmd/api -db-dsn=${GREENLIGHT_DB_DSN}

.PHONY: build/api
build/api:
	@echo 'Building api'
	go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/api

.PHONY: db/psql
db/psql:
	${GREENLIGHT_DB_DSN}

.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}'
	migrate create -seq -ext .sql -dir ./migrations ${name}

.PHONY: db/migrations/up
db/migrations/up:	confirm
	@echo	'Running up migrations ....'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up

.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...
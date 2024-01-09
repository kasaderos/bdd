# general
CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
PACKAGE=core-system/cmd/app

all: format build

build: bindir
	go build -o ${BINDIR}/app ${PACKAGE}

init: 
	make .bin-deps
	make migrate

run:
	go run ${PACKAGE}

# linters
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.51.1
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}

lint: install-lint
	${LINTBIN} run

bindir:
	mkdir -p ${BINDIR}

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

# binaries
LOCAL_BIN:=$(CURDIR)/bin

.PHONY: .bin-deps
.bin-deps:
	$(info Installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest
 	GOBIN=$(LOCAL_BIN) go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# migrations
MIGRATE:=$(LOCAL_BIN)/goose
MIGRATE_DATABASE_DRIVER=$(DB_DRIVER)
MIGRATE_DATABASE_HOST=$(DB_HOST)
MIGRATE_DATABASE_PORT=$(DB_PORT)
MIGRATE_DATABASE_URL="postgresql://postgres:passwd@$(MIGRATE_DATABASE_HOST):$(MIGRATE_DATABASE_PORT)/postgres?sslmode=disable"
.PHONY: migrate
migrate:
	$(MIGRATE) -dir migrations $(MIGRATE_DATABASE_DRIVER) $(MIGRATE_DATABASE_URL) up

.PHONY: migrate_up
migrate_up:
	$(MIGRATE) -dir migrations $(MIGRATE_DATABASE_DRIVER) $(MIGRATE_DATABASE_URL) up-by-one

.PHONY: migrate_down
migrate_down:
	$(MIGRATE) -dir migrations $(MIGRATE_DATABASE_DRIVER) $(MIGRATE_DATABASE_URL) down

.PHONY: reset_migrations
reset_migrations:
	$(MIGRATE) -dir migrations $(MIGRATE_DATABASE_DRIVER) $(MIGRATE_DATABASE_URL) reset

# scripts
.PHONY: go_generate
go_generate:
	go generate ./...

.PHONY: sqlc
sqlc:
	sqlc -f scripts/sqlc/sqlc.yaml generate

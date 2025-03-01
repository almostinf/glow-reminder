include .env
export $(shell sed 's/=.*//' .env)

CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d+\.\d+)/; print $$1;')

SMARTIMPORTSVER=v0.2.0
SMARTIMPORTSBIN=${BINDIR}/smartimports_${GOVER}

LINTVER=v1.57.1
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}

MOCKGENVER=v0.4.0
MOCKGENBIN=${BINDIR}/mockgen_${GOVER}

GOOSEVER=v3.21.1
GOOSEBIN=${BINDIR}/goose_${GOVER}

GOSWAGGERVER=v0.31.0
GOSWAGGERBIN=${BINDIR}/go_swagger_${GOVER}

# ==============================================================================
# Main commands

bindir:
	mkdir -p ${BINDIR}
.PHONY: bindir

precommit: format lint
	echo "OK"
.PHONY: precommit

dev-env-up:
	docker-compose up -d
.PHONY: dev-env-up

dev-env-down:
	docker-compose down
	docker rmi glow_reminder
.PHONY: dev-env-down

# ==============================================================================
# Tools commands

install-mockgen: bindir
	test -f ${MOCKGENBIN} || \
		(GOBIN=${BINDIR} go install go.uber.org/mock/mockgen@${MOCKGENVER} && \
			go get go.uber.org/mock/mockgen/model && \
		mv ${BINDIR}/mockgen ${MOCKGENBIN})
.PHONY: install-mockgen

gen-mocks: install-mockgen
	 go generate -run mockgen ./...
.PHONY: gen-mocks

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})
.PHONY: install-lint

lint: install-lint
	${LINTBIN} run --fix
.PHONY: lint

install-smartimports: bindir
	test -f ${SMARTIMPORTSBIN} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@${SMARTIMPORTSVER} && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTSBIN})
.PHONY: install-smartimports

format: install-smartimports
	${SMARTIMPORTS}
.PHONY: format

install-goose: bindir
	test -f ${GOOSEBIN} || \
		(GOBIN=${BINDIR} go install github.com/pressly/goose/v3/cmd/goose@${GOOSEVER} && \
		mv ${BINDIR}/goose ${GOOSEBIN})
.PHONY: install-goose

create-migration: install-goose
	${GOOSEBIN} -dir ${GOOSE_MIGRATION_DIR} create $(name) sql
.PHONY: create-migration

down-migration: install-goose
	${GOOSEBIN} -dir ${GOOSE_MIGRATION_DIR} ${GOOSE_DRIVER} ${GOOSE_DBSTRING} down
.PHONY: down-migration

up-migration: install-goose
	${GOOSEBIN} -dir ${GOOSE_MIGRATION_DIR} ${GOOSE_DRIVER} ${GOOSE_DBSTRING} up
.PHONY: up-migration

install-go-swagger:
	test -f ${GOSWAGGERBIN} || \
		(GOBIN=${BINDIR} go install github.com/go-swagger/go-swagger/cmd/swagger@${GOSWAGGERVER} && \
		mv ${BINDIR}/swagger ${GOSWAGGERBIN})
.PHONY: install-go-swagger

gen-client:
	${GOSWAGGERBIN} generate client -f api/swagger.yaml -t pkg/glow_reminder
.PHONY: gen-client

# ==============================================================================
# Tests commands

test:
	go test -v -race -count=1 ./...
.PHONY: test

test-100:
	go test -v -race -count=100 ./...
.PHONY: test-100
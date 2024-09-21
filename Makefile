PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin
$(shell [ -f bin ] || mkdir -p $(PROJECT_BIN))
PATH := $(PROJECT_BIN):$(PATH)

.PHONY: dc
dc:
	docker-compose up --remove-orphans --build

.PHONY: cleandc
cleandc:
	rm -rf pg_volume
	docker-compose up --remove-orphans --build

.PHONY: postgres-init
postgres-init:
	docker run --name postgres -p 5433:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=admin -d postgres:15-alpine

.PHONY: postgres-drop
postgres-drop:
	docker stop postgres
	docker remove postgres

.PHONY: postgres
postgres:
	docker exec -it postgres psql

.PHONY: create-db
create-db:
	docker exec -it postgres createdb --username=postgres --owner=postgres demo

.PHONY: drop-db
drop-db:
	docker exec -it postgres dropdb demo

.PHONY: docker
docker:
	docker build -t template .
	docker run --rm \
		--name template \
		--network host \
		-p 9000:9000 \
		-e DB_PASSWORD=$(DB_PASSWORD) \
		template

# ----------------------------------- TESTING -----------------------------------
.PHONY: tests
tests:
	go clean -testcache && go test ./...
# ---------------------------------- PROFILING ----------------------------------
.PHONY: cpuprof
cpuprof:
	( PPROF_TMPDIR=${PPROFDIR} go tool pprof -http :8081 -seconds 20 http://127.0.0.1:9000/debug/pprof/profile )

.PHONY: memprof
memprof:
	( PPROF_TMPDIR=${PPROFDIR} go tool pprof -http :8081 http://127.0.0.1:9000/debug/pprof/heap )

# ---------------------------------- LINTING ------------------------------------
GOLANGCI_LINT_VERSION = v1.60.3
GOLANGCI_LINT = $(PROJECT_BIN)/golangci-lint

.PHONY: .install-golangci-lint
.install-golangci-lint:
	[ -f $(PROJECT_BIN)/golangci-lint ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) $(GOLANGCI_LINT_VERSION)

.PHONY: lint
lint: .install-golangci-lint
	$(GOLANGCI_LINT) run ./... --config=./.golangci.yml

.PHONY: lint-fast
lint-fast: .install-golangci-lint
	$(GOLANGCI_LINT) run ./... --fast --config=./.golangci.yml

# ---------------------------------- MIGRATIONS ---------------------------------

.PHONY: new-migration
new-migration:
	migrate create -ext sql -dir ./migrations involvio_pg


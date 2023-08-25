BIN_DIR = bin
PROTO_DIR = proto
DB_DIR = db
SERVER_DIR = server
CLIENT_DIR = client

include .env
export

COVERAGE_DIR ?= .coverage

ifeq ($(OS), Windows_NT)
	SHELL := powershell.exe
	.SHELLFLAGS := -NoProfile -Command
	SHELL_VERSION = $(shell (Get-Host | Select-Object Version | Format-Table -HideTableHeaders | Out-String).Trim())
	OS = $(shell "{0} {1}" -f "windows", (Get-ComputerInfo -Property OsVersion, OsArchitecture | Format-Table -HideTableHeaders | Out-String).Trim())
	PACKAGE = $(shell (Get-Content go.mod -head 1).Split(" ")[1])
	CHECK_DIR_CMD = if (!(Test-Path $@)) { $$e = [char]27; Write-Error "$$e[31mDirectory $@ doesn't exist$${e}[0m" }
	HELP_CMD = Select-String "^[a-zA-Z_-]+:.*?\#\# .*$$" "./Makefile" | Foreach-Object { $$_data = $$_.matches -split ":.*?\#\# "; $$obj = New-Object PSCustomObject; Add-Member -InputObject $$obj -NotePropertyName ('Command') -NotePropertyValue $$_data[0]; Add-Member -InputObject $$obj -NotePropertyName ('Description') -NotePropertyValue $$_data[1]; $$obj } | Format-Table -HideTableHeaders @{Expression={ $$e = [char]27; "$$e[36m$$($$_.Command)$${e}[0m" }}, Description
	RM_F_CMD = Remove-Item -erroraction silentlycontinue -Force
	RM_RF_CMD = ${RM_F_CMD} -Recurse
	SERVER_BIN = ${SERVER_DIR}.exe
	CLIENT_BIN = ${CLIENT_DIR}.exe
else
	SHELL := bash
	SHELL_VERSION = $(shell echo $$BASH_VERSION)
	UNAME := $(shell uname -s)
	VERSION_AND_ARCH = $(shell uname -rm)
	ifeq ($(UNAME),Darwin)
		OS = macos ${VERSION_AND_ARCH}
	else ifeq ($(UNAME),Linux)
		OS = linux ${VERSION_AND_ARCH}
	else
		$(error OS not supported by this Makefile)
	endif
	PACKAGE = $(shell head -1 go.mod | awk '{print $$2}')
	CHECK_DIR_CMD = test -d $@ || (echo "\033[31mDirectory $@ doesn't exist\033[0m" && false)
	HELP_CMD = grep -E '^[a-zA-Z_-]+:.*?\#\# .*$$' $(firstword $(MAKEFILE_LIST)) | awk 'BEGIN {FS = ":.*?\#\# "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	RM_F_CMD = rm -f
	RM_RF_CMD = ${RM_F_CMD} -r
	SERVER_BIN = ${SERVER_DIR}
	CLIENT_BIN = ${CLIENT_DIR}
endif

postgres: ## Start postgresql docker container 
	docker-compose up $(DATABASE_NAME) -d

createdb: ## Create database from .env file
	docker-compose exec $(DATABASE_NAME) createdb --username=$(DATABASE_USERNAME) --owner=$(DATABASE_USERNAME) $(DATABASE_NAME)
dropdb: ## Drop database from .env file
	docker-compose exec $(DATABASE_NAME) dropdb --username=$(DATABASE_USERNAME) $(DATABASE_NAME) 

# Migrate database up & down lvl
MIGRATE_COMMAND:=migrate -path ./db/migrations --database "${DATABASE_DRVIER}://$(DATABASE_USERNAME):$(DATABASE_PASSWORD)@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" --verbose
migrateup: ## Migrate database up to the latest revision
	$(MIGRATE_COMMAND) up
migratedown: ## Migrate database down to the oldest revision
	$(MIGRATE_COMMAND) down
migrateup1: ## Migrate database up 1 revision
	$(MIGRATE_COMMAND) up 1
migratedown1: ## Migrate database down 1 revision
	$(MIGRATE_COMMAND) down 1

db_docs: ## Generate database documentation
	dbdocs build db/bank.dbml
db_schema: ## Generate database .sql schema
	dbml2sql --postgres -o db/bank.sql db/bank.dbml

sqlc: # Generate sqlc type-safe boilerplate go code from sql queries
	@-rm -r ${DB_DIR}/sqlc/*.sql.go
	sqlc generate

proto: ## Generate protobuf code
	@-rm -r ${PROTO_DIR}/pb/*.go
	@-rm -r docs/swagger/*.swagger.json
	protoc --proto_path=${PROTO_DIR} -I${PROTO_DIR} --go_out=. --go-grpc_out=. \
		--go_opt=module=${PACKAGE} --go-grpc_opt=module=${PACKAGE}  \
		--grpc-gateway_out=${PROTO_DIR}/pb \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		--openapiv2_out docs/swagger \
   		--openapiv2_opt logtostderr=true --openapiv2_opt allow_merge=true \
   		--openapiv2_opt merge_file_name=api_spec \
		${PROTO_DIR}/*.proto

evans: ## Run evans CLI tool connected to the GRPC server
	evans -r repl -p 2137

test: ## Run tests
	go test -coverprofile $(COVERAGE_DIR)/combined.txt --cover ./...

cover: ## Run tests and generate coverage report
	@-rm -r $(COVERAGE_DIR)
	@mkdir $(COVERAGE_DIR)
	make test
	go tool cover -html=$(COVERAGE_DIR)/combined.txt -o cover.html

server: ## Run server API
	go run main.go

install-prerequisites-cli: ## Install CLI tools like sqlc, migrate, dbdocs, protobuf
	brew install act
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	npm install -g dbdocs @dbml/cli
	brew tap ktr0731/evans && brew install evans

build: ## Build the binary file
	go build -o main main.go

bump: all ## Update packages version
	go get -u ./...

PR-check: ## Local Github Action run of the Pull Request steps 
	(cd ../	&& act pull_request)


.DEFAULT_GOAL := help
.PHONY: postgres createdb dropdb migrateup migratedown migratedown1 migrateup1 sqlc test server db_docs db_schema help proto build bump

about: ## Display info related to the build
	@echo "OS: ${OS}"
	@echo "Shell: ${SHELL} ${SHELL_VERSION}"
	@echo "Protoc version: $(shell protoc --version)"
	@echo "Go version: $(shell go version)"
	@echo "Go package: ${PACKAGE}"
	@echo "Openssl version: $(shell openssl version)"


help: ## Show this help
	@${HELP_CMD}
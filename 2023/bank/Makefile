include .env
export

COVERAGE_DIR ?= .coverage


postgres:
	docker-compose up $(DATABASE_NAME) -d

# Create & drop database
createdb:
	docker-compose exec $(DATABASE_NAME) createdb --username=$(DATABASE_USERNAME) --owner=$(DATABASE_USERNAME) $(DATABASE_NAME)
dropdb:
	docker-compose exec $(DATABASE_NAME) dropdb --username=$(DATABASE_USERNAME) $(DATABASE_NAME) 

# Migrate database up & down lvl
MIGRATE_COMMAND:=migrate -path ./db/migrations --database "postgresql://$(DATABASE_USERNAME):$(DATABASE_PASSWORD)@localhost:5432/pg-bank?sslmode=disable" --verbose
migrateup:
	$(MIGRATE_COMMAND) up
migratedown:
	$(MIGRATE_COMMAND) down
migrateup1:
	$(MIGRATE_COMMAND) up 1
migratedown1:
	$(MIGRATE_COMMAND) down 1


sqlc:
	sqlc generate

# Test and coverage
test:
	go test -coverprofile $(COVERAGE_DIR)/combined.txt --cover ./...

cover:
	@-rm -r $(COVERAGE_DIR)
	@mkdir $(COVERAGE_DIR)
	make test
	go tool cover -html=$(COVERAGE_DIR)/combined.txt -o cover.html

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown migratedown1 migrateup1 sqlc test server
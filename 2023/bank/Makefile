include .env
export

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

sqlc:
	sqlc generate

# Test and coverage
test:
	go test --cover ./...

cover: test
	go tool cover -html cover.out -o cover.html
	open cover.html

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test
killall:
	docker kill $(docker ps -aq)
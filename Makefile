###########
#Migration#
###########


dbConnectionString := postgresql://postgres:admin@127.0.0.1:5435/tennis?sslmode=disable
testDb := postgresql://postgres:admin@127.0.0.1:5436/tennistest?sslmode=disable

.PHONY: ensure-migrate
ensure-migrate: ## Ensure that the migrate binary is installed
	@which migrate > /dev/null || echo "Please install \`migrate\` by running \`make install-migrate\`"

.PHONY: install-migrate
install-migrate:
	@echo "To install migrate either run \'brew install golang-migrate\' or check the repository https://github.com/golang-migrate/migrate"

.PHONY: create-migration
create-migration: ensure-migrate
	migrate create -ext sql -dir db/migrations -seq $(name)

.PHONY: test-migrate-up
test-migrate-up: ensure-migrate
	migrate -source file://db/migrations -database ${dbConnectionString} up

.PHONY: test-migrate-down
test-migrate-down: ensure-migrate
	migrate -source file://db/migrations -database ${dbConnectionString} down

.PHONY: test-migrations
test-migrations: ensure-migrate test-migrate-up test-migrate-down test-migrate-up

#######
#Setup#
#######
.PHONY: start-dev-env
start-dev-env: ensure-migrate
	@docker compose up tennisdb -d  
	@sleep 1
	@migrate -source file://db/migrations -database ${dbConnectionString} up
	@go run main.go

.PHONY: end-dev-env
end-dev-env: 
	@docker compose stop tennisdb
	@docker rm tennisdb

.PHONY: restart-dev-env
restart-dev-env: ensure-migrate end-dev-env start-dev-env

######
#Test#
######
.PHONY: test
test: run-test-db
	@{ \
	trap 'docker compose stop tennistestdb 2> /dev/null; docker rm tennistestdb 2> /dev/null; exit 1' ERR; \
	go test ./... -v -p 1; \
	docker compose stop tennistestdb 2> /dev/null; \
	docker rm tennistestdb 2> /dev/null; \
	}

# .PHONY: test-refresh-token
# test-refresh-token: run-test-db run-example-data
# 	@{ \
# 	trap 'docker compose stop tennistestdb 2> /dev/null; docker rm tennistestdb 2> /dev/null; exit 1' ERR; \
# 	go test ./... -v -p 1; \
# 	docker compose stop tennistestdb 2> /dev/null; \
# 	docker rm tennistestdb 2> /dev/null; \
# 	}
# 
# .PHONY: run-example-data
# run-example-data: run-test-db 
# 	psql ${testDb} -f ./tests/example-user.sql -q


.PHONY: run-test-db 
run-test-db: ensure-migrate
	@echo "\n\r\tStarting docker container\n\r"
	@docker compose up -d tennistestdb 2> /dev/null || echo "Please make sure you have docker running and the port 5436 is free"
	@sleep 2
	@echo "\n\r\tApplying migrations to testdb\n\r"
	@migrate -source file://db/migrations -database ${testDb} up >> /dev/null
	@echo ""

######
#GEN
####
.PHONY: sqlc-gen
sqlc-gen:
	sqlc generate

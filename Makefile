include .env

postgres:
	docker compose up -d

createdb:
	docker exec -it postgres createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${POSTGRES_DB}

dropdb:
	docker exec -it postgres dropdb ${POSTGRES_DB}

migrateup:
	migrate -path db/migration -database 'postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable' -verbose up

migratedown:
	migrate -path db/migration -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose down

sqlc:
	sqlc generate

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/mortum5/simple-bank/db/sqlc Store

server:
	go run main.go

lint:
	golangci-lint --color auto -v run --fix 

cloc:
	gocloc .
	
test:
	go test -v -count=1 -race ./...

stop:
	docker compose down

.PHONY: postgres createdb dropdb migrateup migratedown sqlc mock server lint test stop
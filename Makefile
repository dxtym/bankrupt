DB_URL=postgresql://root:secret@localhost:5432/bankrupt?sslmode=disable

postgres:
	docker run --name postgres --network bank-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root bankrupt

dropdb:
	docker exec -it postgres dropdb bankrupt

migrateup:
	migrate -path db/migration -database "${DB_URL}" -verbose up

migrateup1:
	migrate -path db/migration -database "${DB_URL}" -verbose up 1

migrateup2:
	migrate -path db/migration -database "${DB_URL}" -verbose up 2

migratedown:
	migrate -path db/migration -database "${DB_URL}" -verbose down

migratedown1:
	migrate -path db/migration -database "${DB_URL}" -verbose down 1

migratedown2:
	migrate -path db/migration -database "${DB_URL}" -verbose down 2

dbdocs:
	dbdocs build doc/db.dbml

dbml:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	@/snap/bin/go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/dxtym/bankrupt/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 dbdocs dbml sqlc test server mock

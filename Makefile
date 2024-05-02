postgres:
	docker run --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root bankrupt

dropdb:
	docker exec -it postgres dropdb bankrupt

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bankrupt?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bankrupt?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test
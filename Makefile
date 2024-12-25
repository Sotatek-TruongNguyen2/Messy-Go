postgres: 
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=postgres POSTGRES_PASSWORD=postgres -d postgres:12-alpine
createdb: 
	docker exec -it postgres createdb --username=postgres --owner=postgres simple_bank
dropdb:
	docker exec -it postgres drop simple_bank

migrateup: 
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown: 
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc: 
	sqlc generate

test: 
	go test -v ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc
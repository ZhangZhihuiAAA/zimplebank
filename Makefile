DB_URL=postgresql://root:aaa@localhost:5432/zimple_bank?sslmode=disable

postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=aaa -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root zimple_bank

dropdb:
	docker exec -it postgres16 dropdb zimple_bank

initupdatetype:
	sed -i 's/"numeric(32, 6)"/numeric(32, 6)/g' db/migration/*.sql

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate
	# manually change pgtype.Numeric to float64
	sed -i 's/pgtype.Numeric/float64/g' db/sqlc/*.go
	goimports -w db/sqlc/*.go
	gofmt -w db/sqlc/*.go
	sed -i 's/\t/    /g' db/sqlc/*.go

initschema4githubtest:
	ls -1 db/migration/*.up.sql | xargs -I{} psql postgresql://root:aaa@postgres:5432/zimple_bank?sslmode=disable-f {}

test:
	go test -v -count=1 -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/ZhangZhihuiAAA/zimplebank/db/sqlc Store

proto:
	rm -rf pb/*.go
	rm -rf doc/swagger/*.json
	rm -rf doc/statik/*.go
	protoc --proto_path=proto \
	       --go_out=pb --go_opt=paths=source_relative \
	       --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	       --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	       --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=zimple_bank \
	       proto/*.proto
	statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: postgres createdb dropdb initupdatetype migrateup migratedown migrateup1 migratedown1 new_migration db_docs db_schema sqlc initschema4githubtest test server mock proto evans redis
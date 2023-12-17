postgres:
	docker run --name postgres16 --network zbank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=aaa -d postgres:alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root zimple_bank

dropdb:
	docker exec -it postgres16 dropdb zimple_bank

initupdatetype:
	sed -i 's/"numeric(32, 6)"/numeric(32, 6)/g' db/schema/0000_init_schema.up.sql

initup:
	psql postgresql://root:aaa@localhost:5432/zimple_bank?sslmode=disable -f db/schema/0000_init_schema.up.sql

initdown:
	psql postgresql://root:aaa@localhost:5432/zimple_bank?sslmode=disable  -f db/schema/0000_init_schema.down.sql

sqlc:
	sqlc generate
	# manually change pgtype.Numeric to float64
	sed -i 's/pgtype.Numeric/float64/g' db/sqlc/*.sql.go
	goimports -w db/sqlc/*.go
	gofmt -w db/sqlc/*.go
	sed -i 's/\t/    /g' db/sqlc/*.go

test:
	go test -v -count=1 -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/ZhangZhihuiAAA/zimplebank/db/sqlc Store

.PHONY: postgres createdb dropdb initupdatetype initup initdown sqlc test server mock
postgres:
	docker rm postgres16
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=aaa -d postgres:alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root zimple_bank

dropdb:
	docker exec -it postgres16 dropdb zimple_bank

settimezone:
	psql postgresql://root:aaa@localhost:5432/zimple_bank?sslmode=disable -f db/schema/0000_set_timezone.sql

initup:
	psql postgresql://root:aaa@localhost:5432/zimple_bank?sslmode=disable -f db/schema/0001_init_schema.up.sql

initdown:
	psql postgresql://root:aaa@localhost:5432/zimple_bank?sslmode=disable  -f db/schema/0001_init_schema.down.sql

initdb: createdb settimezone initup

sqlc:
	sqlc generate
	# manually change pgtype.Numeric to float64
	sed -i 's/pgtype.Numeric/float64/g' db/sqlc/*.go
	goimports -w db/sqlc/*.go
	gofmt -w db/sqlc/*.go
	sed -i 's/\t/    /g' db/sqlc/*.go

test:
	go test -v -count=1 -cover ./...

.PHONY: postgres createdb dropdb settimezone initup initdown initdb sqlc test
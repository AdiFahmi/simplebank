migrateup:
	migrate -path db/migration -database "mysql://app:app@tcp(127.0.0.1:3390)/app" -verbose up

migratedown:
	migrate -path db/migration -database "mysql://app:app@tcp(127.0.0.1:3390)/app" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

testclean:
	go clean -testcache

server:
	go run main.go

dryrun: migratedown migrateup testclean test

.PHONY: migrateup migratedown sqlc test testclean server drynrun

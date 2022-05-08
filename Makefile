migrateup:
	migrate -path db/migration -database "mysql://root:root@tcp(127.0.0.1:3390)/app" -verbose up

migratedown:
	migrate -path db/migration -database "mysql://root:root@tcp(127.0.0.1:3390)/app" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

singletest:
	go test -run TestGetAccountAPI -v

testclean:
	go clean -testcache

server:
	go run main.go

autoreload:
	gin --appPort 8088 --all -i main.go

dryrun: migratedown migrateup testclean test

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/adifahmi/simplebank/db/sqlc Store

.PHONY: migrateup migratedown sqlc test testclean server drynrun mock

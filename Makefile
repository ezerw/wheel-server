migration:
	migrate create -ext ".sql" -dir db/migrations -seq new_migration
.PHONY: migration-create

migrate:
	migrate -path db/migrations \
		-database "mysql://wheel:secret@tcp(127.0.0.1:3306)/wheel?parseTime=true&loc=Pacific%2FAuckland" \
		-verbose up
.PHONY: migrate

sqlc:
	sqlc generate
.PHONY: sqlc

test:
	go test -v -cover ./...
.PHONY: test
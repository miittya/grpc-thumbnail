.SILENT:

migrate:
	go run ./cmd/migrator --storage-path=./storage/thumbnail.db --migrations-path=./migrations

generate-proto:
	protoc -I proto proto\proto\thumbnail\thumbnail.proto --go_out=.\proto\gen\go --go_opt=paths=source_relative --go-grpc_out=.\proto\gen\go --go-grpc_opt=paths=source_relative
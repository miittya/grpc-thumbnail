.SILENT:

migrate:
	go run ./server/cmd/migrator --storage-path=./storage/thumbnail.db --migrations-path=./server/migrations

generate-proto:
	protoc -I proto proto\proto\thumbnail\thumbnail.proto --go_out=.\proto\gen\go --go_opt=paths=source_relative --go-grpc_out=.\proto\gen\go --go-grpc_opt=paths=source_relative
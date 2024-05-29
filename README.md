# gRPC Thumbnail Proxy Service
## Сборка и запуск сервера
`go run ./server/cmd/grpc-thumbnail --config=./server/config/local.yaml`
### При первом запуске необходимо сначала применить миграции!
`make migrate`

Или:

`mkdir storage`

`go run ./server/cmd/migrator --storage-path=./storage/thumbnail.db --migrations-path=./server/migrations`


## Сборка клиента
`go build -o client ./client/cmd/grpc-thumbnail`
## Запуск клиента
`cd client`

`./grpc-thumbnail *video url* *video url*...`
### Запуск клиента с флагом async
`cd client`

`./grpc-thumbnail --async *video url* *video url*...`
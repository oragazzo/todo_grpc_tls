gen_proto:
	protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. proto/todo.proto

gen_certs:
	bash ./cert/gen.sh

run_server:
	go run cmd/server/main.go -env-path=cmd/server/.env

run_client:
	go run cmd/client/main.go -env-path=cmd/client/.env


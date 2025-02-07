gen_proto:
	protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. proto/todo.proto

gen_certs:
	bash ./cert/gen.sh

run_server:
	go run server/main.go -port 8080 -cert_file ./cert/server-cert.pem -key_file ./cert/server-key.pem -ca_file ./cert/ca-cert.pem

run_client:
	go run client/main.go -address 127.0.0.1:8080 -cert_file ./cert/server-cert.pem


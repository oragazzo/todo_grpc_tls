syntax = "proto3"

option go_package = "todo_grpc/pb"

service "TodoService" {
	rpc SayHello(HelloRequest) returns (HelloResponse)
}

message HelloRequest {
	string greeting = 1;
}

message HelloResponse {
	string repply = 1;
}
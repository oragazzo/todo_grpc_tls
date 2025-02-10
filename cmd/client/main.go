package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	pb "github.com/oragazzo/todo_grpc_tls/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	port = ":8080"
)

func CreateTodo(client pb.TodoServiceClient) {
	reqCreateTodo := &pb.CreateTodoRequest{
		Title:     "Recording",
		Completed: false,
	}

	resCreateTodo, err := client.CreateTodo(context.Background(), reqCreateTodo)
	if err != nil {
		log.Fatal("Error to create todo", err.Error())
	}

	fmt.Printf("Created Todo: ID=%d, Title=%s, Completed=%t\n", resCreateTodo.Id, resCreateTodo.Title, resCreateTodo.Completed)
}

func GetTodos(client pb.TodoServiceClient) {
	reqGetTodos := &pb.GetTodosRequest{}

	resGetTodos, err := client.GetTodos(context.Background(), reqGetTodos)
	if err != nil {
		log.Fatal("Error to get todos", err.Error())
	}

	for _, todo := range resGetTodos.TodoList {
		fmt.Printf("Todo: ID=%d, Title=%s, Completed=%t\n", todo.Id, todo.Title, todo.Completed)
	}
}

func UpdateTodo(client pb.TodoServiceClient) {
	reqUpdateTodo := &pb.UpdateTodoRequest{
		Id:        1,
		Title:     "Studying",
		Completed: true,
	}

	resUpdateTodo, err := client.UpdateTodo(context.Background(), reqUpdateTodo)
	if err != nil {
		log.Fatal("Error to update todo", err.Error())
	}

	fmt.Printf("Updated Todo: ID=%d, Title=%s, Completed=%t\n", resUpdateTodo.Id, resUpdateTodo.Title, resUpdateTodo.Completed)
}

func DeleteTodo(client pb.TodoServiceClient) {
	reqDeleteTodo := &pb.DeleteTodoRequest{
		Id: 2,
	}

	resDeleteTodo, err := client.DeleteTodo(context.Background(), reqDeleteTodo)
	if err != nil {
		log.Fatal("Error to delete todo", err.Error())
	}

	fmt.Printf("Deleted Todo: ID=%d", reqDeleteTodo.Id)
	fmt.Println(resDeleteTodo)
}

var (
	caCert     = "cert/ca.crt"
	clientCert = "cert/client.crt"
	clientKey  = "cert/client.key"
)

func main() {
	// Load client's certificates and create credentials
	tlsConfig := loadTLSConfig()
	creds := credentials.NewTLS(tlsConfig)

	conn, err := grpc.NewClient("localhost"+port, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalln("Error in dial", err.Error())
	}
	defer conn.Close()

	// Create a new TodoService client
	client := pb.NewTodoServiceClient(conn)

	CreateTodo(client)
	// GetTodos(client)
	// UpdateTodo(client)
	// DeleteTodo(client)
}

func loadTLSConfig() *tls.Config {
	// Load certificate of the CA who signed server's certificate
	caCert, err := os.ReadFile(caCert)
	if err != nil {
		log.Fatalf("Failed to load CA certificate: %v", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		log.Fatal("Failed to append CA certificate to pool")
	}

	// Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		log.Fatalf("Failed to load client certificate and key: %v", err)
	}

	// Create and return the TLS configuration
	return &tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{clientCert},
		ServerName:   "localhost",
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	pb "github.com/oragazzo/todo_grpc_tls/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/oragazzo/todo_grpc_tls/internal/config"
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

	fmt.Printf("Created Todo: ID=%d, Title=%s, Completed=%t, CreatedAt=%s\n", resCreateTodo.Id, resCreateTodo.Title, resCreateTodo.Completed, resCreateTodo.CreatedAt.AsTime())
}

func GetTodos(client pb.TodoServiceClient) {
	reqGetTodos := &pb.GetTodosRequest{}

	resGetTodos, err := client.GetTodos(context.Background(), reqGetTodos)
	if err != nil {
		log.Fatal("Error to get todos", err.Error())
	}

	for _, todo := range resGetTodos.TodoList {
		fmt.Printf("Todo: ID=%d, Title=%s, Completed=%t, CreatedAt=%s, UpdatedAt=%s\n", todo.Id, todo.Title, todo.Completed, todo.CreatedAt.AsTime(), todo.UpdatedAt.AsTime())
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

	fmt.Printf("Updated Todo: ID=%d, Title=%s, Completed=%t, UpdatedAt=%s\n", resUpdateTodo.Id, resUpdateTodo.Title, resUpdateTodo.Completed, resUpdateTodo.UpdatedAt.AsTime())
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

func loadEnvFile(envPath string) error {
	if envPath != "" {
		// If absolute path is provided, use it directly
		if filepath.IsAbs(envPath) {
			return godotenv.Load(envPath)
		}
		// Convert relative path to absolute
		absPath, err := filepath.Abs(envPath)
		if err != nil {
			return fmt.Errorf("failed to resolve path: %v", err)
		}
		return godotenv.Load(absPath)
	}

	// Default behavior: try to load .env from the current directory
	return godotenv.Load()
}

func main() {
	// Parse command line flags
	envPath := flag.String("env-path", "", "Path to the environment file")
	flag.Parse()

	// Load environment variables
	if err := loadEnvFile(*envPath); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Load TLS configuration from environment
	tlsCfg, err := config.NewTLSConfigFromEnv()
	if err != nil {
		log.Fatalf("Failed to load TLS config: %v", err)
	}

	// Create TLS credentials
	creds := credentials.NewTLS(config.LoadTLSConfig(*tlsCfg))

	// Get server address from environment
	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = "localhost:8080" // default value
	}

	// Create gRPC connection
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a new TodoService client
	client := pb.NewTodoServiceClient(conn)

	// CreateTodo(client)
	GetTodos(client)
	// UpdateTodo(client)
	// DeleteTodo(client)
}

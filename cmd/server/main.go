package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/joho/godotenv"
	"github.com/oragazzo/todo_grpc_tls/internal/config"
	"github.com/oragazzo/todo_grpc_tls/internal/database"
	"github.com/oragazzo/todo_grpc_tls/internal/server"
	pb "github.com/oragazzo/todo_grpc_tls/proto"
)

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

	// Connect to the database
	db, err := database.Connect(os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Load TLS configuration from environment
	tlsCfg, err := config.NewTLSConfigFromEnv()
	if err != nil {
		log.Fatalf("Failed to load TLS config: %v", err)
	}

	// Create gRPC server with TLS credentials
	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(config.LoadTLSConfig(*tlsCfg))),
	)

	// Create and register TodoServer
	todoServer := server.NewTodoServer(db)
	pb.RegisterTodoServiceServer(grpcServer, todoServer)

	// Start listening on port
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Server listening on port localhost%v", port)

	// Serve gRPC requests
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

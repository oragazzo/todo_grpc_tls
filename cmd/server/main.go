package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
	pb "github.com/oragazzo/todo_grpc/proto"
)

const (
	port = ":8080"
)

type Todo struct {
	gorm.Model
	Title     string `gorm:"not null"`
	Completed bool   `gorm:"not null"`
}

type server struct {
	db *gorm.DB
	pb.UnimplementedTodoServiceServer
}

func (s *server) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	todo := &Todo{
		Title:     req.Title,
		Completed: req.Completed,
	}

	s.db.Create(todo)

	return &pb.CreateTodoResponse{
		Id:        int64(todo.ID),
		Title:     todo.Title,
		Completed: todo.Completed,
	}, nil
}

func (s *server) GetTodos(ctx context.Context, req *pb.GetTodosRequest) (*pb.GetTodosResponse, error) {
	var todos []*Todo

	s.db.Find(&todos)

	var resTodos []*pb.GetTodosResponse_Todo
	for _, todo := range todos {
		resTodo := &pb.GetTodosResponse_Todo{
			Id:        int64(todo.ID),
			Title:     todo.Title,
			Completed: todo.Completed,
		}
		resTodos = append(resTodos, resTodo)
	}

	return &pb.GetTodosResponse{
		TodoList: resTodos,
	}, nil
}

func (s *server) UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*pb.UpdateTodoResponse, error) {
	var todo Todo

	s.db.First(&todo, req.Id)
	if todo.ID == 0 {
		return nil, fmt.Errorf("Todo not found")
	}

	todo.Title = req.Title
	todo.Completed = req.Completed

	s.db.Save(&todo)

	return &pb.UpdateTodoResponse{
		Id:        int64(todo.ID),
		Title:     todo.Title,
		Completed: todo.Completed,
	}, nil
}

func (s *server) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*pb.DeleteTodoResponse, error) {
	var todo Todo

	s.db.First(&todo, req.Id)
	if todo.ID == 0 {
		return nil, fmt.Errorf("Todo not found")
	}

	s.db.Delete(&todo)

	return &pb.DeleteTodoResponse{}, nil
}

func main() {
	// Load environment variable
	godotenv.Load()

	// Connect to the PostgreSQL database
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Database connecting is success")

	// Automigrate the Todo struct to the PostgreSQL database
	db.AutoMigrate(&Todo{})

	// Create a new server
	grpcServer := grpc.NewServer()

	// Register TodoService server
	pb.RegisterTodoServiceServer(grpcServer, &server{db: db})

	// Start listening on port
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("error to listen", err.Error())
	}
	log.Printf("Server listening on port localhost%v", port)

	// Serve gRPC requests
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalln("error when serve grpc", err.Error())
	}
}

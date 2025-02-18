package server

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/oragazzo/todo_grpc_tls/internal/models"
	pb "github.com/oragazzo/todo_grpc_tls/proto"
)

// TodoServer implements the TodoService gRPC server
type TodoServer struct {
	db *gorm.DB
	pb.UnimplementedTodoServiceServer
}

// NewTodoServer creates a new TodoServer instance
func NewTodoServer(db *gorm.DB) *TodoServer {
	return &TodoServer{
		db: db,
	}
}

func (s *TodoServer) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	now := time.Now()
	todo := &models.Todo{
		Title:     req.Title,
		Completed: req.Completed,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.db.Create(todo).Error; err != nil {
		return nil, fmt.Errorf("failed to create todo: %v", err)
	}

	return &pb.CreateTodoResponse{
		Id:        int64(todo.ID),
		Title:     todo.Title,
		Completed: todo.Completed,
		CreatedAt: timestamppb.New(todo.CreatedAt),
	}, nil
}

func (s *TodoServer) GetTodos(ctx context.Context, req *pb.GetTodosRequest) (*pb.GetTodosResponse, error) {
	var todos []*models.Todo

	if err := s.db.Find(&todos).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch todos: %v", err)
	}

	var resTodos []*pb.GetTodosResponse_Todo
	for _, todo := range todos {
		resTodos = append(resTodos, &pb.GetTodosResponse_Todo{
			Id:        int64(todo.ID),
			Title:     todo.Title,
			Completed: todo.Completed,
			CreatedAt: timestamppb.New(todo.CreatedAt),
			UpdatedAt: timestamppb.New(todo.UpdatedAt),
		})
	}

	return &pb.GetTodosResponse{
		TodoList: resTodos,
	}, nil
}

func (s *TodoServer) UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*pb.UpdateTodoResponse, error) {
	var todo models.Todo

	if err := s.db.First(&todo, req.Id).Error; err != nil {
		return nil, fmt.Errorf("todo not found: %v", err)
	}

	todo.Title = req.Title
	todo.Completed = req.Completed

	if err := s.db.Save(&todo).Error; err != nil {
		return nil, fmt.Errorf("failed to update todo: %v", err)
	}

	return &pb.UpdateTodoResponse{
		Id:        int64(todo.ID),
		Title:     todo.Title,
		Completed: todo.Completed,
		UpdatedAt: timestamppb.New(todo.UpdatedAt),
	}, nil
}

func (s *TodoServer) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*pb.DeleteTodoResponse, error) {
	var todo models.Todo

	if err := s.db.First(&todo, req.Id).Error; err != nil {
		return nil, fmt.Errorf("todo not found: %v", err)
	}

	if err := s.db.Delete(&todo).Error; err != nil {
		return nil, fmt.Errorf("failed to delete todo: %v", err)
	}

	return &pb.DeleteTodoResponse{}, nil
}

package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/oragazzo/todo_grpc/proto"
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

func main() {
	// Set up a connection to the server
	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

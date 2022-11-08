package main

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/ts1106/go-grpc-todo-app/todo"
	"google.golang.org/grpc"
)

type TodoServer struct {
	todo.UnimplementedTodoServer
	todoItems []*todo.Item
}

func (s *TodoServer) GetItem(ctx context.Context, req *todo.GetRequest) (*todo.Item, error) {
	for _, item := range s.todoItems {
		if item.Id == req.GetId() {
			return item, nil
		}
	}
	return &todo.Item{}, nil
}

func (s *TodoServer) ListItem(req *todo.ListRequest, stream todo.Todo_ListItemServer) error {
	for _, item := range s.todoItems {
		if item.Id >= req.LowId && item.Id <= req.HighId {
			if err := stream.Send(item); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *TodoServer) CreateItem(stream todo.Todo_CreateItemServer) error {
	var createCount int32
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&todo.Summary{
				CreateCount: createCount,
			})
		}
		if err != nil {
			return err
		}
		s.todoItems = append(s.todoItems, item)
		createCount++
	}
}

func (s *TodoServer) ProgressChat(stream todo.Todo_ProgressChatServer) error {
	for {
		note, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := stream.Send(note); err != nil {
			return err
		}
	}
}

func newServer() *TodoServer {
	return &TodoServer{
		todoItems: []*todo.Item{},
	}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	todo.RegisterTodoServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}

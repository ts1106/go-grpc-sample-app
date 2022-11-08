package main

import (
	"context"
	"io"
	"log"

	"github.com/ts1106/go-grpc-todo-app/todo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createItem(ctx context.Context, client todo.TodoClient) {
	items := []*todo.Item{
		{Id: 1, Title: "first task"},
		{Id: 2, Title: "Second task"},
		{Id: 3, Title: "Third task"},
	}
	stream, err := client.CreateItem(ctx)
	if err != nil {
		log.Fatalf("client.CreateItem failed: %v", err)
	}
	for _, item := range items {
		if err := stream.Send(item); err != nil {
			log.Fatalf("CreateItem stream.Send failed: %v", err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("CreateItem stream.CloseAndRecv failed: %v", err)
	}
	log.Printf("Summary: %v", reply)
}

func getItem(ctx context.Context, client todo.TodoClient) {
	item, err := client.GetItem(ctx, &todo.GetRequest{Id: 1})
	if err != nil {
		log.Fatalf("client.GetItem failed: %v", err)
	}
	log.Println(item)
}

func listItem(ctx context.Context, client todo.TodoClient) {
	stream, err := client.ListItem(ctx, &todo.ListRequest{LowId: 1, HighId: 2})
	if err != nil {
		log.Fatalf("client.ListItem failed: %v", err)
	}
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("ListItem stream.Recv failed: %v", err)
		}
		log.Printf("Got Item: %v", item)
	}
}

func progressChat(ctx context.Context, client todo.TodoClient) {
	waitc := make(chan struct{})
	stream, err := client.ProgressChat(ctx)
	if err != nil {
		log.Fatalf("client.ProgressChat failed: %v", err)
	}
	go func() {
		for {
			note, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("ProgressChat stream.Recv failed: %v", err)
			}
			log.Printf("Task(Id: %d) Progress is %d %%", note.Id, note.Per)
		}
	}()
	notes := []*todo.ProgressNote{
		{Id: 1, Per: 50},
		{Id: 2, Per: 30},
		{Id: 3, Per: 80},
	}
	for _, note := range notes {
		if err := stream.Send(note); err != nil {
			log.Fatalf("ProgressChat stream.Send failed: %v", err)
		}
	}
	stream.CloseSend()
	<-waitc
}

func main() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	ctx := context.Background()
	client := todo.NewTodoClient(conn)

	createItem(ctx, client)
	getItem(ctx, client)
	listItem(ctx, client)
	progressChat(ctx, client)
}

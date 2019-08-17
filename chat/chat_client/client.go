package main

import (
	"context"
	"fmt"
	"github.com/benibana2001/gcp-go/chat/chatpb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	fmt.Println("Welcome to Chat !")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("You could not connect: %v\n", err)
	}

	defer cc.Close()

	client := chatpb.NewMessageServiceClient(cc)

	// Create a Stream
	stream, err := client.SendMessage(context.Background())
	if err != nil {
		fmt.Printf("error was occured while creating stream: %v\n", err)
		return
	}

	// ---
	ch := make(chan struct{})

	go receive(stream)

	go func(){
		for {
			var message string
			// 入力を待ち受ける
			fmt.Println("Please Input Message")
			fmt.Scanf("%v", &message)
			send(stream, message)
		}
	}()

	<-ch
}

func send(stream chatpb.MessageService_SendMessageClient, msg string) {
	// Create a request
	request := chatpb.SendMessageRequest{
		UserId:  "Taro",
		Message: msg,
	}

	// Send
	fmt.Printf("Send message: %v\n", msg)
	stream.Send(&request)

}

func receive(stream chatpb.MessageService_SendMessageClient) {
	for {

		fmt.Println("Waiting Message...")

		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("Error while receiving: %v\n", err)
			break
		}
		fmt.Printf("Received: %v\n", res.GetUserId())
		fmt.Printf("Received: %v\n", res.GetMessage())
	}
}


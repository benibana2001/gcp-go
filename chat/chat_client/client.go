package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/benibana2001/gcp-go/chat/chatpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
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
	//transferStream, err := client.TransferMessage(context.Background(), &chatpb.Null{})
	//if err != nil {
	//	fmt.Printf("error was occured while creating stream: %v\n", err)
	//	return
	//}

	//postStream, err := client.PostMessage(context.Background())
	//if err != nil {
	//	fmt.Printf("error was occured while creating stream: %v\n", err)
	//	return
	//}

	testStream, err := client.Test(context.Background())

	// ---
	go receive(testStream)

	for {
		//var message string
		// 入力を待ち受ける
		//fmt.Println("Please Input Message")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		msg := input.Text()
		//fmt.Scanf("%v", &message)
		post(testStream, msg)
	}
}

func post(stream chatpb.MessageService_TestClient, msg string) {
	// Create a request
	request := chatpb.PostRequest{
		Name:    "Taro",
		Content: msg,
	}

	// Send
	//fmt.Printf("Send message: %v\n", msg)
	if err := stream.Send(&request); err != nil {
		fmt.Printf("Error was occured while Send: %v\n", err)
	}
}

func receive(stream chatpb.MessageService_TestClient) {
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("Error while receiving: %v\n", err)
			break
		}
		fmt.Printf(">>>>: [%v]%v\n", res.Name, res.Content)
	}
}

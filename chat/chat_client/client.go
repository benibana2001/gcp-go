package main

import (
	"context"
	"fmt"
	"github.com/benibana2001/gcp-go/chat/chatpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	fmt.Println("Welcome to Chat !")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("You could not connect: %v\n", err)
	}

	defer cc.Close()

	c := chatpb.NewMessageServiceClient(cc)

	Send(c)
}

var request01 = chatpb.SendMessageRequest{
	UserId:  "Taro",
	Message: "Konichiha !",
}
var request02 = chatpb.SendMessageRequest{
	UserId:  "Taro",
	Message: "Ima Nani Shite Masuka ?",
}

func Send(c chatpb.MessageServiceClient) {

	// Create a Stream
	stream, err := c.SendMessage(context.Background())
	if err != nil {
		fmt.Printf("error was occured while creating stream: %v\n", err)
		return
	}

	requests := []*chatpb.SendMessageRequest{
		&request01,
		&request02,
	}

	waitc := make(chan struct{})
	// send a bunch of message
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive a bunch of message
	go func() {
		for {
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
		close(waitc)
	}()

	//
	<-waitc

}

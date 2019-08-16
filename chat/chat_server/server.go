package main

import (
	"fmt"
	"github.com/benibana2001/gcp-go/chat/chatpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type server struct{}

// メッセージの双方向通信
func (*server) SendMessage(stream chatpb.MessageService_SendMessageServer) error {
	fmt.Printf("SendMessage function was invoked!\n")

	// リクエストを待ち受ける
	for {
		req, err := stream.Recv()

		// 最後までメッセージを読んだら終了
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Printf("Error occured while Reading stream: %v", err)
		}

		// 送信されてきたデータをパース
		userId := req.GetUserId()
		message := req.GetMessage()

		fmt.Println(userId + ": " + message)

		sendErr := stream.Send(&chatpb.SendMessageResponse{
			UserId: "[FROM]: " + userId,
			Message: "[MESSAGE]: " + message,
		})
		if sendErr != nil {
			fmt.Printf("Error while sending data to client: %v", err)
			return err
		}
	}
}

func main() {
	fmt.Println("Server is listening...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	chatpb.RegisterMessageServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

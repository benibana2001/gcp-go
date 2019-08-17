package main

import (
	"fmt"
	"github.com/benibana2001/gcp-go/chat/chatpb"
	"google.golang.org/grpc"
	"io"
	"net"
)

type server struct {
	contents []message
}

type message struct {
	author string
	content string
}

// メッセージの双方向通信
func (s *server) SendMessage(stream chatpb.MessageService_SendMessageServer) error {
	// リクエストを待ち受ける
	for {
		req, err := stream.Recv()

		// 最後までメッセージを読んだら終了
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Printf("Error occured while Reading stream: %v\n", err)
		}

		// 送信されてきたデータをパース
		author := req.GetName()
		content := req.GetContent()

		fmt.Println(author + ": " + content)

		// メッセージをserverに追加
		s.contents = append(s.contents, message{
			author: author,
			content: content,
		})
	}
}

func (s *server) GetMessage(req *chatpb.Null, stream chatpb.MessageService_GetMessageServer) error {
	preNum := len(s.contents)
	currentNum := 0

	for {
		currentNum = len(s.contents)
		if currentNum > preNum {
			latest := s.contents[len(s.contents)-1]
			err := stream.Send(&chatpb.Message{
				Name: latest.author,
				Content: latest.content,
			})
			if err != nil {
				fmt.Printf("Error was occured while GetMessage: %v\n", err)
			}
		}
		preNum = currentNum
	}

}

func main() {
	fmt.Println("Server is listening...")

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		fmt.Printf("Failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	chatpb.RegisterMessageServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}



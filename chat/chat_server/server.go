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

/*
func a(s *server, stream chatpb.MessageService_PostMessageServer) error{
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

func b(s *server, req *chatpb.Null, stream chatpb.MessageService_TransferMessageServer) error {
	preNum := len(s.contents)
	currentNum := 0

	for {
		currentNum = len(s.contents)
		if currentNum > preNum {
			latest := s.contents[len(s.contents)-1]
			err := stream.Send(&chatpb.TransferResult{
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

// メッセージの双方向通信
func (s *server) PostMessage(stream chatpb.MessageService_PostMessageServer) error {
	err := a(s, stream)
	if err != nil {
		return err
	}
	return nil
}

func (s *server) TransferMessage(req *chatpb.Null, stream chatpb.MessageService_TransferMessageServer) error {
	err := b(s, req, stream)
	if err != nil {
		return err
	}
	return nil
}
*/

func (s *server) Test(stream chatpb.MessageService_TestServer) error {
	preNum := len(s.contents)
	currentNum := 0

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

		currentNum = len(s.contents)
		if currentNum > preNum {
			latest := s.contents[len(s.contents)-1]
			err := stream.Send(&chatpb.TransferResult{
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



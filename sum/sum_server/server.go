package main

import (
	"context"
	"fmt"
	"github.com/benibana2001/gcp-go/sum/sumpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}


func (*server) Sum(ctx context.Context, req *sumpb.SumRequest) (*sumpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with %v\n", req)
	firstNumber := req.GetSum().GetFirstNum()
	secondNumber := req.GetSum().GetSecondNum()
	amount := firstNumber + secondNumber

	res := &sumpb.SumResponse{
		Result: amount,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello Sum Server!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	sumpb.RegisterSumServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package main

import (
	"context"
	"fmt"
	"github.com/benibana2001/gcp-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type server struct{}


func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with %v\n", req)
	firstNumber := req.GetSum().GetFirstNum()
	secondNumber := req.GetSum().GetSecondNum()
	amount := firstNumber + secondNumber

	res := &calculatorpb.SumResponse{
		Result: amount,
	}
	return res, nil
}

// Client Streaming
func (*server) ComputeAverage(stream calculatorpb.CalculateService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverate function was invoked with a streaming request\n")
	var result int32
	var i int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// finish running server
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: result / i,
			})
		}
		if err != nil {
			fmt.Printf("Error while reading stream: %v", err)
		}
		num := req.GetNumber()
		result += num
		i++
	}
}


func (*server) DecompositManyTimes(req *calculatorpb.DecompositManyTimeRequest, stream calculatorpb.CalculateService_DecompositManyTimesServer) error {
	fmt.Printf("DecompositManyTimes function was invoked with %v\n", req)
	primeNumber := req.GetPrimeNumber()

	var k int32 = 2
	for primeNumber > 1 {
		if primeNumber % k == 0 {
			res := &calculatorpb.DecompositManyTimesResponse{
				Result: k,
			}
			primeNumber = primeNumber / k
			stream.Send(res)
		}else {
			k = k + 1
		}
	}
	return nil
}

func main() {
	fmt.Println("Hello Sum Server!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculateServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

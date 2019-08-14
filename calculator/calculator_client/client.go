package main

import (
	"context"
	"fmt"
	"github.com/benibana2001/gcp-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	fmt.Println("Hello sumClient!")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculateServiceClient(cc)
	//doUnary(c)

	doServerStreaming(c)
}

func doServerStreaming(c calculatorpb.CalculateServiceClient)  {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &calculatorpb.DecompositManyTimeRequest{
		PrimeNumber: 1000,
	}

	resStream, err := c.DecompositManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while callng Greet RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("error while reading stream: %v",err)
		}
		fmt.Printf("Response from DecompositManyTimes: %v\n", msg.GetResult())
	}
}

func doUnary(c calculatorpb.CalculateServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")

	req := &calculatorpb.SumRequest{
		Sum: &calculatorpb.Sum{
			FirstNum: 10,
			SecondNum: 2,
		},
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while callng Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}

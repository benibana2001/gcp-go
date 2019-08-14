package main

import (
	"context"
	"fmt"
	"github.com/benibana2001/gcp-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	fmt.Println("Hello sumClient!")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewSumServiceClient(cc)
	doUnary(c)
}

func doUnary(c calculatorpb.SumServiceClient) {
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

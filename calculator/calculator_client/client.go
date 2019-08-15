package main

import (
	"context"
	"fmt"
	"github.com/benibana2001/gcp-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
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

	//doServerStreaming(c)

	//doClientStreming(c)

	doBiDiStreaming(c)
}

func doBiDiStreaming(c calculatorpb.CalculateServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	// create a stream
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		fmt.Printf("error while creating stream: %v", err)
		return
	}

	// create a requests
	requests := []*calculatorpb.FindMaximumRequest{
		&calculatorpb.FindMaximumRequest{
			Number: 10,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 16,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 4,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 20,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 22,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 1,
		},
	}

	waitc := make(chan struct{})

	// send a bunch of message
	go func(){
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
			if err == io.EOF{
				break
			}
			if err != nil {
				fmt.Printf("Error while receiving: %v\n", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}

func doClientStreming(c calculatorpb.CalculateServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	reqests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: 10,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 12,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 18,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 20,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 30,
		},
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		fmt.Printf("error while calling LongServer: %v", err)
	}

	for _, req := range reqests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("error while receiving LongServer: %v", err)
	}
	fmt.Printf("ComputeAverate Response: %v\n", res)
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

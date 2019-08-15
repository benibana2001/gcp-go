package main

import (
	"context"
	"fmt"
	"github.com/benibana2001/gcp-go/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	fmt.Println("Hello I'm a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	//doUnary(c)

	//doServerStreaming(c)

	//doClientStreaming(c)

	doBiDiStreaming(c)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	// create a stream
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		fmt.Printf("error while creating stream: %v", err)
		return
	}
	requests := []*greetpb.GreetEveryOneRequest{
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Taro",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jiro",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Saburo",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Shiro",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Goro",
			},
		},
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
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	//
	<-waitc
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Taro",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jiro",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Saburo",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Shiro",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Goro",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		fmt.Printf("error while calling LongServer: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("error while receiving LongServer: %v", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)
}
func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Yusuke",
			LastName:  "Iwase",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		fmt.Printf("error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("error while reading stream: %v", err)
		}
		fmt.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Yusuke",
			LastName:  "Iwase",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while callng Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}

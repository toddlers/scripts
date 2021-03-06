package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/toddlers/learngrpc/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello I'm a client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect :%v", err)
	}
	defer conn.Close()
	c := greetpb.NewGreetServiceClient(conn)
	// doUnary(c)
	doUnaryWithDeadline(c, 5)  //ok
	doUnaryWithDeadline(c, 20) // error
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDirectionalStreaming(c)

}

func doBiDirectionalStreaming(c greetpb.GreetServiceClient) {
	log.Println("Starting BiDirectional streaming")
	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Suresh",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Sonu",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Sam",
			},
		},
	}
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while calling GreetEveryone: %v", err)
		return
	}
	// to block
	waitc := make(chan struct{})
	//send messages
	go func() {
		for _, req := range requests {
			log.Printf("Sending request :  %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	//receive messages
	go func() {
		for {
			resStream, err := stream.Recv()
			if err == io.EOF {
				// server closed stream
				break
			}
			if err != nil {
				log.Fatalf("Error while reading the stream : %v", err)
				break
			}
			log.Printf("response from greet : %v\n", resStream.GetResult())
			time.Sleep(1000 * time.Millisecond)
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	log.Println("Starting client streaming")
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Suresh",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Sonu",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Sam",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}
	for _, req := range requests {
		log.Printf("Sending reque :  %v", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response LongGreet: %v", err)
	}
	log.Printf("LongGreet Response : %v", resp)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	log.Println("Starting server streaming RPC...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Suresh",
			LastName:  "Kumar",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling greet : %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// server closed stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading the stream : %v", err)
		}
		log.Printf("response from greet : %v", msg.GetResult())
	}
}

func doUnary(c greetpb.GreetServiceClient) {
	log.Println("Unary RPC")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "",
			LastName:  "Kumar",
		},
	}
	// erro call
	doUnaryCall(c, req)

	// correct call
	req = &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "suresh",
			LastName:  "Kumar",
		},
	}
	doUnaryCall(c, req)
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	log.Println("UnaryWithDeadline RPC")
	req := &greetpb.GreetWithDealineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Suresh",
			LastName:  "Kumar",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := c.GreetWithDealine(ctx, req)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC(user error)
			log.Printf("Error from server %s: %s", respErr.Message(), respErr.Code())
			// log.Println(respErr.Code())
			if respErr.Code() == codes.DeadlineExceeded {
				log.Println("Deadline execeeded")
			} else {
				log.Printf("Unexpected error : %v", respErr)
			}
		} else {
			log.Fatalf("Error while calling GreetWithDeadline: %v\n", err)
			return
		}
		return
	}
	log.Printf("Response from GreetWithDeadline : %v", res.Result)
}

func doUnaryCall(c greetpb.GreetServiceClient, req *greetpb.GreetRequest) {
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC(user error)
			log.Printf("Error from server %s: %s", respErr.Message(), respErr.Code())
			// log.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				log.Println("Empty first name")
			}

		} else {
			log.Fatalf("Error while calling greet : %v\n", err)
			return
		}
	}
	log.Printf("response from greet : %v\n", res.GetResult())
}

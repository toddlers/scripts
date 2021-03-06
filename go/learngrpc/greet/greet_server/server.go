package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/toddlers/learngrpc/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Printf("Received request : %v", req)
	firstName := req.GetGreeting().FirstName
	if firstName == "" {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Received empty firstname",
		)
	}
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetWithDealine(ctx context.Context, req *greetpb.GreetWithDealineRequest) (*greetpb.GreetWithDealineResponse, error) {
	log.Printf("Received request GreetWithDealine : %v", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			log.Println("The client cancelled the request")
			return nil, status.Error(codes.Canceled, "The client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := req.GetGreeting().FirstName
	result := "Hello " + firstName
	res := &greetpb.GreetWithDealineResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	log.Printf("Greet many times request : %v", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number" + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)

	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	log.Println("clien streaming request")
	result := ""
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// done reading from the client
			return stream.SendAndClose(
				&greetpb.LongGreetResponse{
					Result: result,
				},
			)
		}
		if err != nil {
			log.Fatalf("Error reading client stream : %v", err)
			return err
		}
		firstName := msg.Greeting.FirstName
		result += "Hello " + firstName + "! \n"
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	log.Println("Greet everyone request")
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error reading client stream : %v", err)
			return err
		}
		firstName := msg.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "\n"
		log.Printf("Sending result : %v", result)
		err = stream.Send(
			&greetpb.GreetEveryoneResponse{
				Result: result,
			},
		)
		if err != nil {
			log.Fatalf("Error sending data to client stream : %v", err)
			return err
		}
	}
}
func main() {
	fmt.Println("Hello World")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}

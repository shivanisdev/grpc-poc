package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"example.com/grpc-poc/pb"

	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Greetings: %s", r.GetMessage())

	r1, _ := client.PrintAgeByYear(ctx, &pb.YearRequest{Year: 1993})

	log.Printf("Your Age is: %d", r1.GetAge())

	stream, err := client.GreetsStream(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes: %v", err)
	}
	for {
		// Receive streamed messages
		res, err := stream.Recv()
		if err == io.EOF {
			break // End of stream
		}
		if err != nil {
			log.Fatalf("Error while receiving stream: %v", err)
		}

		// Print the response from the server
		fmt.Printf("Response from server: %s\n", res.GetMessage())
	}
}

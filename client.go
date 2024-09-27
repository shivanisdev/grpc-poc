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

	// create video streaming client

	videoClient := pb.NewVideoStreamClient(conn)

	req := &pb.VideoRequest{VideoId: "test"}

	videoStream, err := videoClient.StreamVideo(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling StreamVideo: %v", err)
	}

	file, err := os.Create("videos/downloaded_sample.mp4")
	if err != nil {
		log.Fatalf("could not create file: %v", err)
	}
	defer file.Close()

	for {
		chunk, err := videoStream.Recv()
		if err == io.EOF {
			log.Println("File received successfully")
			break
		}
		if err != nil {
			log.Fatalf("error receiving chunk: %v", err)
		}

		_, err = file.Write(chunk.ChunkData)
		if err != nil {
			log.Fatalf("error writing to file: %v", err)
		}
	}
}

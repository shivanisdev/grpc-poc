package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"example.com/grpc-poc/pb"
	"google.golang.org/grpc"
)

// Server is used to implement helloworld.GreeterServer
type Server struct {
	pb.UnimplementedGreeterServer
	pb.UnimplementedVideoStreamServer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello !! " + in.GetName()}, nil
}

func (s *Server) PrintAgeByYear(ctx context.Context, r *pb.YearRequest) (*pb.AgeResponse, error) {
	age := 2024 - r.GetYear()
	return &pb.AgeResponse{Age: age}, nil
}

func (s *Server) GreetsStream(r *pb.HelloRequest, stream pb.Greeter_GreetsStreamServer) error {
	name := r.GetName()
	for i := 0; i < 5; i++ { // Simulate streaming 5 greetings
		response := &pb.HelloReply{
			Message: fmt.Sprintf("Hello %s, message number %d", name, i+1),
		}
		if err := stream.Send(response); err != nil {
			return fmt.Errorf("errorfit: %v", err)
		}
		time.Sleep(1 * time.Second) // Simulate delay
	}
	return nil
}

func (s *Server) StreamVideo(r *pb.VideoRequest, stream pb.VideoStream_StreamVideoServer) error {
	fileName := fmt.Sprintf("./videos/%s.mp4", r.VideoId)
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("could not open video file: %v", err)
	}
	defer file.Close()

	buffer := make([]byte, 1024*64)
	seqNum := 0

	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read from file: %v", err)
		}

		err = stream.Send(&pb.VideoChunk{
			ChunkData:      buffer[:n],
			SequenceNumber: int32(seqNum),
		})
		if err != nil {
			return fmt.Errorf("failed to send chunk: %v", err)
		}

		seqNum++
	}

	return nil
}

func main() {
	// Listen on a port (e.g., :50051)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	fmt.Println("Starting Server at 50051")
	// Register the Greeter service on the gRPC server
	pb.RegisterGreeterServer(grpcServer, &Server{})

	// Register the Video Streaming service on the gRPC ser
	pb.RegisterVideoStreamServer(grpcServer, &Server{})

	// Start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

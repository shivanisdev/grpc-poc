package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"example.com/grpc-poc/pb"
	"google.golang.org/grpc"
)

// Server is used to implement helloworld.GreeterServer
type Server struct {
	pb.UnimplementedGreeterServer
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

	// Start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

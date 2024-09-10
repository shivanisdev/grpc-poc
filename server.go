package main

import (
	"context"
	"fmt"
	"log"
	"net"

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
	return &pb.AgeResponse{Age: 2024 - r.GetYear()}, nil
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

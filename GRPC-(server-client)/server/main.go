package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	pb "server/proto/gen"

	// Enable Gzip compression support for gRPC
	_ "google.golang.org/grpc/encoding/gzip"
)

// server implements the gRPC Calculator and Greeter services.
type server struct {
	pb.UnimplementedCalculatorServer
	pb.UnimplementedGreeterServer
}

// Add handles the addition of two numbers and demonstrates gRPC metadata usage.
func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	// Extract incoming metadata (headers) from client
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("No metadata found")
	}
	log.Println("Metadata:", md)

	// Read authorization header from incoming metadata
	authHeader, ok := md["authorization"]
	if !ok {
		log.Println("No metadata found")
	}
	log.Println(authHeader)

	// Set custom response headers
	responseHeaders := metadata.Pairs("key1", "value1")
	err := grpc.SendHeader(ctx, responseHeaders)
	if err != nil {
		log.Println("Couldnot set headers")
	}

	// Set custom response trailers
	responseTrailers := metadata.Pairs("key1", "value1")
	err = grpc.SetTrailer(ctx, responseTrailers)
	if err != nil {
		log.Println("Couldnot set trailers")
	}

	// Calculate sum
	sum := req.A + req.B

	log.Println("Sum :", sum)

	return &pb.AddResponse{
		Sum: sum,
	}, nil
}

// Greet handles returning a personalized greeting message.
func (s *server) Greet(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: fmt.Sprintf("Hello %s, nice to meet you,", req.Name),
	}, nil
}

func main() {
	cert := "cert.pem"
	key := "key.pem"

	// Start TCP listener
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Load TLS credentials
	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		log.Fatalln("Failed to load credentials:", err)
	}

	// Initialize gRPC server with TLS credentials
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	// Register gRPC service implementations
	pb.RegisterCalculatorServer(grpcServer, &server{})
	pb.RegisterGreeterServer(grpcServer, &server{})

	// Start serving incoming RPC requests
	log.Println("Server is running on port", port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}

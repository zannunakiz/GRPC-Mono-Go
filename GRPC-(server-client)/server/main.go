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

	_ "google.golang.org/grpc/encoding/gzip"
)

type server struct {
	pb.UnimplementedCalculatorServer
	pb.UnimplementedGreeterServer
}

func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("No metadata found")
	}
	log.Println("Metadata:", md)
	// authHeader, ok := md[0]
	authHeader, ok := md["authorization"]
	if !ok {
		log.Println("No metadata found")
	}
	log.Println(authHeader)

	// Set response Headers
	responseHeaders := metadata.Pairs("key1", "value1")
	err := grpc.SendHeader(ctx, responseHeaders)
	if err != nil {
		log.Println("Couldnot set headers")
	}

	responseTrailers := metadata.Pairs("key1", "value1")
	err = grpc.SetTrailer(ctx, responseTrailers)
	if err != nil {
		log.Println("Couldnot set trailers")
	}

	sum := req.A + req.B

	log.Println("Sum :", sum)

	return &pb.AddResponse{
		Sum: sum,
	}, nil
}

func (s *server) Greet(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: fmt.Sprintf("Hello %s, nice to meet you,", req.Name),
	}, nil
}

func main() {
	cert := "cert.pem"
	key := "key.pem"

	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		log.Fatalln("Failed to load credentials:", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	pb.RegisterCalculatorServer(grpcServer, &server{})
	pb.RegisterGreeterServer(grpcServer, &server{})

	log.Println("Server is running on port", port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}

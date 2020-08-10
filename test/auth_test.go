package test

import (
	"context"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func init() {}

// Test scenario: The server side registers the interceptor, constructs the client to add the authentication part, and sends the request
// Authentication to print OK error correctly client returns GRPC-related error

var (
	port        = ""
	gGrpcServer = new(Server)
)

func TestCreateGrpcServer(t *testing.T) {
	CreateGrpcServer()
}

func CreateGrpcServer() *grpc.Server {
	c, err := credentials.NewServerTLSFromFile("server.pem", "server.key")
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
	}
	listen, err := net.Listen("tcp", ":"+port)
	grpcserver := grpc.NewServer(
		grpc.Creds(c),
	)
	pb.RegisterGreeterServer(grpcserver, gGrpcServer)
	if err != nil {
		log.Fatalf("net listen err: %v", err)
	}
	go grpcserver.Serve(listen)
	return grpcserver
}

type Server struct {
	expectedMd map[string]string
}

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{}, nil
}

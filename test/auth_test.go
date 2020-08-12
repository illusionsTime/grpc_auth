package test

import (
	"context"
	"crypto/tls"
	server_sdk "grpc_auth/authinterceptor"
	"grpc_auth/sdk"
	"net"
	"runtime"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func init() {}

// Test scenario: The server side registers the interceptor, constructs the client to add the authentication part, and sends the request
// Authentication to print OK error correctly client returns GRPC-related error

var (
	port          = "50051"
	localserver   = "127.0.0.1"
	gGrpcServer   = new(Server)
	gMockProvider = new(AuthProvider)
)

func TestCreateGrpcServer(t *testing.T) {
	CreateGrpcServer()
	server_sdk.RegisterDefaultAuthProvider(gMockProvider)
}

func CreateGrpcServer() *grpc.Server {
	c, err := credentials.NewServerTLSFromFile("server.pem", "server.key")
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
	}
	listen, err := net.Listen("tcp", ":"+port)
	grpcserver := grpc.NewServer(
		grpc.Creds(c),
		grpc.UnaryInterceptor(server_sdk.NewUnaryServerInterceptor),
	)
	pb.RegisterGreeterServer(grpcserver, gGrpcServer)
	if err != nil {
		log.Errorf("net listen err: %v", err)
	}
	go grpcserver.Serve(listen)
	return grpcserver
}

type Server struct {
	expectedMd map[string]string
}

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "call sayHello check md got error,The client is not authorized")
	}
	checkMd := false
	vs := md.Get("subject")
	for i := 0; i < len(vs); i++ {
		if vs[i] == "subject001" {
			checkMd = true
			break
		}
	}
	if !checkMd {
		return nil, status.Error(codes.Unauthenticated, "call sayHello check md value got error,The client is not authorized")
	}
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func StartGrpcClient(cre credentials.PerRPCCredentials) error {
	runtime.GOMAXPROCS(runtime.NumCPU())
	conn, err := grpc.Dial(localserver+":"+port,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			ServerName:         "",
			InsecureSkipVerify: true,
		})),
		grpc.WithPerRPCCredentials(cre))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: "zmy"})
	if err != nil {
		return err
	}
	log.Errorf("Greeting: %s", r.GetMessage())
	return nil
}

func TestCase1(t *testing.T) {
	gMockProvider.expectedToken = "000000"
	gMockProvider.expectedType = "Bearer"
	err := StartGrpcClient(sdk.NewBearerAuth("000000"))
	time.Sleep(1 * time.Second)
	if err == nil {
		t.Logf("test ok!")
	} else {
		t.Errorf("test with error,error is %v", err)
	}
}

type AuthProvider struct {
	expectedToken string
	expectedType  string
	expectedMap   map[string]string
	expectedMd    map[string]string
}

func (s *AuthProvider) CheckAuth(info server_sdk.CallBackInfo) (map[string]string, error) {
	m := make(map[string]string)
	m["subject"] = "subject001"
	if s.expectedType != info.TokenType {
		return m, status.Error(codes.Unauthenticated, "info token type is error,The client is not authorized")
	}
	if s.expectedToken != "" && s.expectedToken != info.Token {
		return m, status.Error(codes.Unauthenticated, "info token check is error,The client is not authorized")
	}
	if len(s.expectedMap) > 0 {
		for k, v := range s.expectedMap {
			v1, ok := info.TokenCtx[k]
			if !ok || v1 != v {
				return m, status.Error(codes.Unauthenticated, "info token map check  error,The client is not authorized")
			}
		}
	}
	return m, nil
}

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	hellopb "github.com/TadayoshiOtsuka/grpc_sample/src/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	hellopb.RegisterGreetingServiceServer(server, NewMyServer())
	reflection.Register(server)

	go func() {
		log.Printf("start gRPC server! port: %v", port)
		server.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	server.GracefulStop()
}

type MyServer struct {
	hellopb.UnimplementedGreetingServiceServer
}

func (s *MyServer) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func NewMyServer() *MyServer {
	return &MyServer{}
}

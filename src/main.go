package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

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

func NewMyServer() *MyServer {
	return &MyServer{}
}

func (s *MyServer) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func (s *MyServer) HelloServerStream(req *hellopb.HelloRequest, stream hellopb.GreetingService_HelloServerStreamServer) error {
	resCount := 5

	for i := 0; i < resCount; i++ {
		msg := fmt.Sprintf("[%d] Hello, %s!", i, req.GetName())

		if err := stream.Send(
			&hellopb.HelloResponse{Message: msg},
		); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}

	return nil
}

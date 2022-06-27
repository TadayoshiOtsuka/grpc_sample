package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	pb "github.com/TadayoshiOtsuka/grpc_sample/src/pkg/grpc"

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

	pb.RegisterGreetingServiceServer(server, NewMyServer())
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
	pb.UnimplementedGreetingServiceServer
}

func NewMyServer() *MyServer {
	return &MyServer{}
}

func (s *MyServer) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func (s *MyServer) HelloServerStream(req *pb.HelloRequest, stream pb.GreetingService_HelloServerStreamServer) error {
	resCount := 5

	for i := 0; i < resCount; i++ {
		msg := fmt.Sprintf("[%d] Hello, %s!", i, req.GetName())

		if err := stream.Send(
			&pb.HelloResponse{Message: msg},
		); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}

	return nil
}

func (s *MyServer) HelloClientStream(stream pb.GreetingService_HelloClientStreamServer) error {
	nameList := make([]string, 0)

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			msg := fmt.Sprintf("Hello, %v!", nameList)
			return stream.SendAndClose(&pb.HelloResponse{Message: msg})
		}
		if err != nil {
			return err
		}
		nameList = append(nameList, req.GetName())
	}
}

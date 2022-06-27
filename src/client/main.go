package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/TadayoshiOtsuka/grpc_sample/src/pkg/grpc"
)

func main() {
	fmt.Println("start gRPC client")

	scanner := bufio.NewScanner(os.Stdin)

	address := "localhost:8080"
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("connection failed")
		return
	}

	defer conn.Close()

	client := pb.NewGreetingServiceClient(conn)

	for {
		fmt.Println("1: send unary req")
		fmt.Println("2: send server stream req")
		fmt.Println("3: send client stream req")
		fmt.Println("4: send bi stream req")
		fmt.Println("5: exit")
		fmt.Print("please Enter >")
		scanner.Scan()
		in := scanner.Text()

		switch in {
		case "1":
			Hello(client, scanner)
		case "2":
			HelloServerStream(client, scanner)
		case "3":
			HelloClientStream(client, scanner)
		case "5":
			fmt.Println("bye.")
			goto M
		}
	}
M:
}

func Hello(client pb.GreetingServiceClient, scanner *bufio.Scanner) {
	fmt.Println("Please Enter Your Name")
	scanner.Scan()
	name := scanner.Text()

	req := &pb.HelloRequest{
		Name: name,
	}

	res, err := client.Hello(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.GetMessage())
}

func HelloServerStream(client pb.GreetingServiceClient, scanner *bufio.Scanner) {
	fmt.Println("Please Enter Your Name")
	scanner.Scan()
	name := scanner.Text()
	req := &pb.HelloRequest{
		Name: name,
	}

	stream, err := client.HelloServerStream(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("all res already received")
			break
		}

		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println(res)
	}
}

func HelloClientStream(client pb.GreetingServiceClient, scanner *bufio.Scanner) {
	stream, err := client.HelloClientStream(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	sendCount := 5
	fmt.Printf("Please enter %d names.\n", sendCount)
	for i := 0; i < sendCount; i++ {
		scanner.Scan()
		name := scanner.Text()

		if err := stream.Send(&pb.HelloRequest{
			Name: name,
		}); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("stream send, %s\n", name)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.GetMessage())
}

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

	hellopb "github.com/TadayoshiOtsuka/grpc_sample/src/pkg/grpc"
)

func main() {
	fmt.Println("start gRPC client")

	scanner := bufio.NewScanner(os.Stdin)

	address := "localhost:8080"
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal("connection failed")
		return
	}

	defer conn.Close()

	client := hellopb.NewGreetingServiceClient(conn)

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

		case "5":
			fmt.Println("bye.")
			goto M
		}
	}
M:
}

func Hello(client hellopb.GreetingServiceClient, scanner *bufio.Scanner) {
	fmt.Println("Please Enter Your Name")
	scanner.Scan()
	name := scanner.Text()

	req := &hellopb.HelloRequest{
		Name: name,
	}

	res, err := client.Hello(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.GetMessage())

}

func HelloServerStream(client hellopb.GreetingServiceClient, scanner *bufio.Scanner) {
	fmt.Println("Please Enter Your Name")
	scanner.Scan()
	name := scanner.Text()
	req := &hellopb.HelloRequest{
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
		}

		fmt.Println(res)
	}
}

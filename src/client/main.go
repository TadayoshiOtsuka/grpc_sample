package main

import (
	"bufio"
	"context"
	"fmt"
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
		fmt.Println("1: send req")
		fmt.Println("2: exit")
		fmt.Print("please Enter >")

		scanner.Scan()
		in := scanner.Text()
		switch in {
		case "1":
			Hello(client, scanner)

		case "2":
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

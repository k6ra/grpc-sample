package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/k6ra/grpc-sample/echo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("[echo] ")
}

func main() {
	target := "localhost:50051"
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dif not connect: %s", err)
	}
	defer conn.Close()
	client := pb.NewEchoServiceClient(conn)
	msg := os.Args[1]
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.Echo(
		ctx,
		&pb.EchoRequest{Message: msg},
	)
	if err != nil {
		log.Println(err)
	}
	log.Println(r.GetMessage())
}

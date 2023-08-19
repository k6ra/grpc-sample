package main

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	pb "github.com/k6ra/grpc-sample/downloader/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("[file] ")
}

func main() {
	target := "localhost:50051"
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s\n", err)
	}
	defer conn.Close()
	c := pb.NewFileServiceClient(conn)
	name := os.Args[1]
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stream, err := c.Download(ctx, &pb.FileRequest{Name: name})
	if err != nil {
		log.Fatalf("cloud not download: %s\n", err)
	}
	var blob []byte
	for {
		c, err := stream.Recv()
		if err == io.EOF {
			log.Printf("done %d bytes\n", len(blob))
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		blob = append(blob, c.GetData()...)
	}

	file, err := os.Create(name)
	if err != nil {
		log.Fatalf("cloud not create file: %s\n", err)
	}
	_, err = file.Write(blob)
	if err != nil {
		log.Fatalf("file write failed: %s\n", err)
	}
}

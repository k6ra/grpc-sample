package main

import (
	"io"
	"log"
	"sync"

	pb "github.com/k6ra/grpc-sample/chat/proto"
)

var streams sync.Map

type chatService struct{}

func (c *chatService) Connect(stream pb.ChatService_ConnectServer) error {
	log.Println("connect", &stream)
	streams.Store(stream, struct{}{})
	defer func() {
		log.Println("disconnect", &stream)
		streams.Delete(stream)
	}()
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		streams.Range(func(key, value any) bool {
			stream := key.(pb.ChatService_ConnectServer)
			stream.Send(&pb.Post{
				Name:    req.GetName(),
				Message: req.GetMessage(),
			})
			return true
		})
	}
}

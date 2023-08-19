package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	pb "github.com/k6ra/grpc-sample/uploader/proto"
)

type fileService struct{}

func (s *fileService) Upload(stream pb.FileService_UploadServer) error {
	var blob []byte
	var name string
	for {
		c, err := stream.Recv()
		if err == io.EOF {
			log.Printf("done %d bytes\n", len(blob))
			break
		}
		if err != nil {
			return err
		}
		name = c.GetName()
		blob = append(blob, c.GetData()...)
	}
	fp := filepath.Join("./uploader/resource", name)
	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	_, err = file.Write(blob)
	if err != nil {
		return err
	}
	stream.SendAndClose(&pb.FileResponse{
		Size: int64(len(blob)),
	})
	return nil
}

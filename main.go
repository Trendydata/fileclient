package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/Trendydata/fileserver/pkg"
	"google.golang.org/grpc"
)

const (
	CHUNK_SIZE = 512
)

func sendFile(client pkg.FileClient, filename string) {
	ctx := context.Background()

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Failed to read file %s: %v", filename, err)
	}
	client.Upload(ctx, &pkg.Chunk{
		Meta: &pkg.Metadata{
			Name: fmt.Sprintf("%s", filename),
		},
		Data: data,
	})
	log.Printf("Sent %d bytes", len(data))
}

func sendFileStream(client pkg.FileClient, filename string) {
	ctx := context.Background()
	stream, err := client.UploadStream(ctx)
	if err != nil {
		log.Fatal("Failed to create stream: %v", err)
	}
	file, err := os.OpenFile(filename, os.O_RDONLY, 0444)
	if err != nil {
		log.Fatal("Failed to read file %s: %v", filename, err)
	}
	buf := make([]byte, CHUNK_SIZE)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalln(err)
		}
		stream.Send(&pkg.Chunk{
			Meta: &pkg.Metadata{
				Name: fmt.Sprintf("%s", filename),
			},
			Data: buf[:n],
		})
		log.Printf("Sent %d bytes", n)
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
}

func main() {
	conn, err := grpc.Dial(fmt.Sprintf(":%d", 8080), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer conn.Close()
	client := pkg.NewFileClient(conn)
	sendFileStream(client, os.Args[1])
}

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Trendydata/fileserver/pkg"
	"google.golang.org/grpc"
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
}

func main() {
	conn, err := grpc.Dial(fmt.Sprintf(":%d", 8080), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer conn.Close()
	client := pkg.NewFileClient(conn)
	sendFile(client, os.Args[1])
}

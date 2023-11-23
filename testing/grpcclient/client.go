package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"ports-service/internal/adapters/streamfromfile"
	"ports-service/internal/domain"

	"google.golang.org/grpc"

	"github.com/google/uuid"
	pb "ports-service/internal/gen/grpc"
)

type Port struct {
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float64 `json:"coordinates"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
	Code        string    `json:"code"`
}

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewPortServiceClient(conn)

	filePath := "C:\\Users\\tillk\\GolandProjects\\ports-service\\data\\ports.json"
	if err != nil {
		log.Fatal(err)
	}

	stream, err := client.StreamPorts(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Println("opening file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	decoder := json.NewDecoder(file)

	// Check the first token to determine if it's an array or an object.
	_, err = decoder.Token()
	if err != nil {
		log.Printf("Error reading the first JSON token: %v\n", err)
		return
	}

	// Iterate over each entry in the JSON object
	for decoder.More() {
		// Read the key
		key, err := decoder.Token()
		if err != nil {
			fmt.Println(err)
			return
		}

		var item domain.Port
		if err := decoder.Decode(&item); err != nil {
			if err != io.EOF {
				log.Printf("Error decoding object: %v\n", err)
			}
			break
		}

		streamfromfile.SetKey(&item, key.(string))

		fmt.Println("sending the following data: ", item)

		time.Sleep(time.Millisecond)

		id := uuid.New()

		req := &pb.StreamPortsRequest{
			Uuid: id.String(),
			Port: &pb.Port{
				Key:         item.Key,
				Name:        item.Name,
				City:        item.City,
				Country:     item.Country,
				Alias:       item.Alias,
				Regions:     item.Regions,
				Coordinates: item.Coordinates,
				Province:    item.Province,
				Timezone:    item.Timezone,
				Unlocs:      item.Unlocs,
				Code:        item.Code,
			},
		}

		stream.Send(req)
	}

	stream.CloseAndRecv()
}

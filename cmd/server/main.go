package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"ports-service/internal/adapters/database"
	"ports-service/internal/adapters/grpc"
	"ports-service/internal/adapters/streamfromfile"
	"ports-service/internal/domain"
	"syscall"
	"time"
)

func main() {
	var runGRPC = flag.Bool("grpc", true, "Whether to run gRPC server")
	var bufferSize = flag.Int("buffer", 100, "Size of buffered channel to limit memory usage")
	var filePath = flag.String("file", "data/ports.json", "Path to JSON file")
	var debugKey = flag.String("debugkey", "", "Key to lookup in the database")

	flag.Parse()

	//TODO: consider using slog package for structured logging
	log.Println("Run gRPC server:", *runGRPC)
	log.Println("Buffer size:", *bufferSize)
	log.Println("File path:", *filePath)
	log.Println("Debug key:", *debugKey)

	db := database.MemDB[domain.Port]{
		DB: make(map[string]domain.Port),
	}

	//TODO: move this to separate package
	if *debugKey != "" {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			ticker := time.NewTicker(1 * time.Second)
			for {
				select {
				case <-ctx.Done():
					ticker.Stop()
					return
				case <-ticker.C:
					port, ok := db.DB[*debugKey]
					if ok {
						log.Printf("Lookup key found %s: %v\n", *debugKey, port)
						return
					} else {
						log.Printf("No entry for key %s\n", *debugKey)
					}
				}
			}
		}()
	}

	repo := domain.StorePortRepository{Data: &db}

	if *runGRPC {
		portService := grpc.PortService{PortForShipsRepository: repo}
		grpc.StartServer("localhost:8080", portService, *bufferSize)
	} else {
		portService := streamfromfile.PortService{PortForShipsRepository: repo}
		ctx, cancel := context.WithCancel(context.Background())

		// Start streaming
		err := portService.StreamJSONfromFile(ctx, *filePath, *bufferSize)
		if err != nil {
			log.Fatalln(err)
		}

		// Wait for SIGINT (Ctrl+C)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		// Got SIGINT, cancel context
		cancel()
	}
}

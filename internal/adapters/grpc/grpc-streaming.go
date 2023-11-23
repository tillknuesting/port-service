package grpc

import (
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"ports-service/internal/domain"
	pb "ports-service/internal/gen/grpc"
	"syscall"
	"time"
)

type PortServiceServer struct {
	pb.UnimplementedPortServiceServer
	grpcStreamChan chan domain.Port // Channel for streaming data
	portService    PortService
}

func NewPortServiceServer(portService PortService, grpcStreamChan chan domain.Port) *PortServiceServer {
	return &PortServiceServer{
		portService:    portService,
		grpcStreamChan: grpcStreamChan,
	}
}

func (p *PortServiceServer) StreamPorts(server pb.PortService_StreamPortsServer) error {
	for {
		portData, err := server.Recv()
		if err == io.EOF {
			// End of stream
			return nil
		}
		if err != nil {
			return err // Handle the error appropriately
		}

		port := domain.Port{
			Key:         portData.Port.Key,
			Name:        portData.Port.Name,
			City:        portData.Port.City,
			Country:     portData.Port.Country,
			Alias:       portData.Port.Alias,
			Regions:     portData.Port.Regions,
			Coordinates: portData.Port.Coordinates,
			Province:    portData.Port.Province,
			Timezone:    portData.Port.Timezone,
			Unlocs:      portData.Port.Unlocs,
			Code:        portData.Port.Code,
		}

		// Process received data
		p.grpcStreamChan <- port // Assuming the types match
	}
}

func (p *PortService) StreamFromGRPC(ctx context.Context, bufferSize int, grpcStreamChan chan domain.Port) error {
	streamer := NewStreamer[domain.Port](grpcStreamChan)
	portStream, err := streamer.StreamObjects(ctx, bufferSize)
	if err != nil {
		return err
	}

	for port := range portStream {
		err := p.PortForShipsRepository.Store(ctx, port)
		if err != nil {
			return err
		}
	}
	return nil
}

func StartServer(address string, portService PortService, bufferzise int) error {
	// Create a channel for streaming data from the gRPC handler
	grpcStreamChan := make(chan domain.Port)

	// Create a new PortServiceServer with the channel and PortService
	server := NewPortServiceServer(portService, grpcStreamChan)

	// Start the streaming process
	go func() {
		if err := portService.StreamFromGRPC(context.Background(), bufferzise, grpcStreamChan); err != nil {
			log.Fatalf("Error streaming ports: %v", err)
		}
	}()

	// Listen on the specified address
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register PortServiceServer with the gRPC server
	pb.RegisterPortServiceServer(s, server)

	// Set up graceful shutdown with forced stop
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down gRPC server...")

		// Create a channel to signal the completion of a graceful shutdown
		done := make(chan struct{})

		go func() {
			s.GracefulStop()
			close(done)
		}()

		// Wait for either the completion of the graceful shutdown or a timeout
		select {
		case <-done:
			log.Println("gRPC server gracefully stopped")
		case <-time.After(5 * time.Second): // TODO: Make timeout configurable
			log.Println("Timeout reached, forcibly stopping gRPC server")
			s.Stop()
		}
	}()

	// Start the gRPC server
	log.Printf("Starting gRPC server on %s", address)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}

	return nil
}

type PortService struct {
	PortForShipsRepository domain.StorePortRepository
}

// Streamer is a generic type for streaming data from a JSON file.
// T is the type of data that will be streamed.
type Streamer[T any] struct {
	grpcStream <-chan T
}

func NewStreamer[T any](grpcStream <-chan T) *Streamer[T] {
	return &Streamer[T]{grpcStream: grpcStream}
}

func (s *Streamer[T]) StreamObjects(ctx context.Context, bufferSize int) (<-chan T, error) {
	ch := make(chan T, bufferSize)
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-s.grpcStream:
				if !ok {
					return // gRPC stream closed
				}
				ch <- item
			}
		}
	}()
	return ch, nil
}

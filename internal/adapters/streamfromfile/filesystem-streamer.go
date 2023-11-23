package streamfromfile

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"ports-service/internal/domain"
)

// FileStreamer is a generic type for streaming data from a JSON file.
// T is the type of data that will be streamed.
type FileStreamer[T any] struct {
	filePath string // Path to the JSON file.
}

// NewFileStreamer acts as a constructor for FileStreamer.
func NewFileStreamer[T any](filePath string) *FileStreamer[T] {
	// Initialize a new FileStreamer with the provided file path.
	return &FileStreamer[T]{filePath: filePath}
}

// StreamObjects streams objects of type T from a JSON file.
// The use of a buffered channel (make(chan T, bufferSize)) allows some degree of pre-fetching of data,
// but this buffering is controlled by the bufferSize parameter.
// The buffer size determines how many objects are held in memory after being read from the file
// but before being processed by the consumer of the channel.
func (fs *FileStreamer[T]) StreamObjects(ctx context.Context, bufferSize int) (<-chan T, error) {
	ch := make(chan T, bufferSize)
	go func() {
		defer close(ch)

		file, err := os.Open(fs.filePath)
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

			var item T
			if err := decoder.Decode(&item); err != nil {
				if err != io.EOF {
					log.Printf("Error decoding object: %v\n", err)
				}
				break
			}

			SetKey(&item, key.(string))

			//TODO: Remove this debug code
			//// Print the key and location
			//fmt.Printf("Key: %v, Location: %+v\n", key, item)

			select {
			case <-ctx.Done():
				return
			case ch <- item:
			}
		}
	}()

	return ch, nil
}

// SetKey sets the field "Key" to the given value on the passed
// struct pointer item, if the field exists. The struct type
// passed for T must be a pointer. Helpful to preserve metadata e.g. Key
// that would be lost otherwise.
func SetKey[T any](item *T, value string) {
	t := reflect.TypeOf(item).Elem()
	v := reflect.ValueOf(item).Elem()

	if field, ok := t.FieldByName("Key"); ok {
		v.FieldByIndex(field.Index).SetString(value)
	}
}

type PortService struct {
	PortForShipsRepository domain.StorePortRepository
}

// StreamJSONfromFile streams objects of type T from a JSON file. TODO: Perhaps move this to a service/application layer?
func (p PortService) StreamJSONfromFile(ctx context.Context, filePath string, bufferSize int) error {
	streamer := NewFileStreamer[domain.Port](filePath)
	portStream, err := streamer.StreamObjects(ctx, bufferSize)
	if err != nil {
		return fmt.Errorf("setting up JSON stream from filesystem: %w", err)
	}

	for port := range portStream {
		err := p.PortForShipsRepository.Store(ctx, port)
		if err != nil {
			return fmt.Errorf("set fails on Data from StorePortRepository: %w", err)
		}
	}
	return nil
}

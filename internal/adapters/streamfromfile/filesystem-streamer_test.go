package streamfromfile_test

import (
	"context"
	"os"
	"ports-service/internal/adapters/streamfromfile"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestObject struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// helper function to create a temporary file with JSON content
func createTempJSONFile(content string) (string, error) {
	file, err := os.CreateTemp("", "*.json")
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func TestStreamObjects_Success(t *testing.T) {
	filePath, err := createTempJSONFile(`{"obj1": {"key": "1", "value": "one"}, "obj2": {"key": "2", "value": "two"}}`)
	assert.NoError(t, err)
	defer os.Remove(filePath)

	fileStreamer := streamfromfile.NewFileStreamer[TestObject](filePath)
	ctx := context.Background()
	ch, err := fileStreamer.StreamObjects(ctx, 2)
	assert.NoError(t, err)

	// Read objects from channel
	obj1 := <-ch
	obj2 := <-ch

	assert.Equal(t, TestObject{Key: "obj1", Value: "one"}, obj1)
	assert.Equal(t, TestObject{Key: "obj2", Value: "two"}, obj2)
}

func TestStreamObjects_InvalidJSON(t *testing.T) {
	filePath, err := createTempJSONFile(`invalid json content`)
	assert.NoError(t, err)
	defer os.Remove(filePath)

	fileStreamer := streamfromfile.NewFileStreamer[TestObject](filePath)
	ctx := context.Background()
	ch, err := fileStreamer.StreamObjects(ctx, 2)
	assert.NoError(t, err)

	// Expect no objects and channel to be closed
	_, ok := <-ch
	assert.False(t, ok, "channel should be closed with no objects sent")
}

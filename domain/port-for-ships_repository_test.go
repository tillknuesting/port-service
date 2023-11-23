package domain_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"ports-service/internal/domain"
	"testing"
)

// MockPortRepository is a mock implementation of the PortRepository interface
type MockPortRepository struct {
	mock.Mock
}

// GetPort mocks the GetPort method
func (m *MockPortRepository) GetPort(ctx context.Context, key string) (domain.Port, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(domain.Port), args.Error(1)
}

// SavePort mocks the SavePort method
func (m *MockPortRepository) SavePort(ctx context.Context, port domain.Port) error {
	args := m.Called(ctx, port)
	return args.Error(0)
}

// TestGetPort tests the GetPort method
func TestGetPort(t *testing.T) {
	mockRepo := new(MockPortRepository)
	ctx := context.TODO()
	port := domain.Port{Key: "PORT123", Name: "PortName"}

	// Setting up the expectations
	mockRepo.On("GetPort", ctx, "PORT123").Return(port, nil)

	// Call the method
	result, err := mockRepo.GetPort(ctx, "PORT123")

	// Assert expectations
	mockRepo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, port, result)
}

// TestSavePort tests the SavePort method
func TestSavePort(t *testing.T) {
	mockRepo := new(MockPortRepository)
	ctx := context.TODO()
	port := domain.Port{Key: "PORT123", Name: "PortName"}

	// Setting up the expectations
	mockRepo.On("SavePort", ctx, port).Return(nil)

	// Call the method
	err := mockRepo.SavePort(ctx, port)

	// Assert expectations
	mockRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

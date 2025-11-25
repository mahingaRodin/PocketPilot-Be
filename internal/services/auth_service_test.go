package services

import (
	"pocketpilot-api/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) EmailExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func TestAuthService_Register(t *testing.T) {
    mockRepo := new(MockUserRepository)
    authService := NewAuthService(mockRepo, "test-secret-key")

    registerReq := &models.RegisterRequest{
        Email:     "test@example.com",
        Password:  "password123",
        FirstName: "John",
        LastName:  "Doe",
    }

    t.Run("Successful Registration", func(t *testing.T) {
        // Setup expectations
        mockRepo.On("EmailExists", "test@example.com").Return(false, nil)
        mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil).Run(func(args mock.Arguments) {
            user := args.Get(0).(*models.User)
            user.ID = "user-123"
            // Password should be hashed
            assert.NotEqual(t, "password123", user.PasswordHash)
            assert.Equal(t, "test@example.com", user.Email)
            assert.Equal(t, "John", user.FirstName)
            assert.Equal(t, "Doe", user.LastName)
        })

        // Execute
        authResponse, err := authService.Register(registerReq)

        // Assert
        require.NoError(t, err)
        require.NotNil(t, authResponse)
        assert.NotEmpty(t, authResponse.Token)
        assert.Equal(t, "user-123", authResponse.User.ID)
        assert.Equal(t, "test@example.com", authResponse.User.Email)

        mockRepo.AssertExpectations(t)
    })

    t.Run("Registration with Existing Email", func(t *testing.T) {
        mockRepo.On("EmailExists", "existing@example.com").Return(true, nil)

        registerReq.Email = "existing@example.com"
        authResponse, err := authService.Register(registerReq)

        assert.Error(t, err)
        assert.Nil(t, authResponse)
        assert.Equal(t, "email already registered", err.Error())

        mockRepo.AssertExpectations(t)
    })

    t.Run("Registration with EmailExists Error", func(t *testing.T) {
        mockRepo.On("EmailExists", "error@example.com").Return(false, assert.AnError)

        registerReq.Email = "error@example.com"
        authResponse, err := authService.Register(registerReq)

        assert.Error(t, err)
        assert.Nil(t, authResponse)

        mockRepo.AssertExpectations(t)
    })
}
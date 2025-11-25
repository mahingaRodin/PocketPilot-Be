package services

import (
	"pocketpilot-api/internal/models"
	"pocketpilot-api/internal/utils"
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

func TestAuthService_Login(t *testing.T) {
    mockRepo := new(MockUserRepository)
    authService := NewAuthService(mockRepo, "test-secret-key")

    loginReq := &models.LoginRequest{
        Email:    "test@example.com",
        Password: "password123",
    }

    t.Run("Successful Login", func(t *testing.T) {
        // Create a user with hashed password (simulate real hashing)
        hashedPassword, _ := utils.HashPassword("password123")
        mockUser := &models.User{
            ID:           "user-123",
            Email:        "test@example.com",
            PasswordHash: hashedPassword,
            FirstName:    "John",
            LastName:     "Doe",
        }

        mockRepo.On("GetUserByEmail", "test@example.com").Return(mockUser, nil)

        authResponse, err := authService.Login(loginReq)

        require.NoError(t, err)
        require.NotNil(t, authResponse)
        assert.NotEmpty(t, authResponse.Token)
        assert.Equal(t, "user-123", authResponse.User.ID)
        assert.Equal(t, "test@example.com", authResponse.User.Email)

        mockRepo.AssertExpectations(t)
    })

    t.Run("Login with Non-Existent User", func(t *testing.T) {
        mockRepo.On("GetUserByEmail", "nonexistent@example.com").Return(nil, nil)

        loginReq.Email = "nonexistent@example.com"
        authResponse, err := authService.Login(loginReq)

        assert.Error(t, err)
        assert.Nil(t, authResponse)
        assert.Equal(t, "invalid email or password", err.Error())

        mockRepo.AssertExpectations(t)
    })

    t.Run("Login with Wrong Password", func(t *testing.T) {
        // Create user with different password hash
        mockUser := &models.User{
            ID:           "user-123",
            Email:        "test@example.com",
            PasswordHash: "different-hash",
            FirstName:    "John",
            LastName:     "Doe",
        }

        mockRepo.On("GetUserByEmail", "test@example.com").Return(mockUser, nil)

        authResponse, err := authService.Login(loginReq)

        assert.Error(t, err)
        assert.Nil(t, authResponse)
        assert.Equal(t, "invalid email or password", err.Error())

        mockRepo.AssertExpectations(t)
    })

    t.Run("Login with Repository Error", func(t *testing.T) {
        mockRepo.On("GetUserByEmail", "error@example.com").Return(nil, assert.AnError)

        loginReq.Email = "error@example.com"
        authResponse, err := authService.Login(loginReq)

        assert.Error(t, err)
        assert.Nil(t, authResponse)

        mockRepo.AssertExpectations(t)
    })
}


func TestAuthService_GetUserProfile(t *testing.T) {
    mockRepo := new(MockUserRepository)
    authService := NewAuthService(mockRepo, "test-secret-key")

    t.Run("Successful GetUserProfile", func(t *testing.T) {
        mockUser := &models.User{
            ID: "user-123",
            Email: "test@example.com",
            FirstName: "Rosine",
            LastName: "Smith",
        }
        mockRepo.On("GetUserByID", "user-123").Return(mockUser, nil)
        user, err := authService.GetUserProfile("user-123")

        require.NoError(t, err)
        require.NotNil(t, user)
        assert.Equal(t, "user-123", user.ID)
        assert.Equal(t, "test@example.com", user.Email)

        mockRepo.AssertExpectations(t)
    })

        t.Run("Profile Not Found", func(t *testing.T) {
        mockRepo.On("GetUserByID", "nonexistent-user").Return(nil, nil)

        user, err := authService.GetUserProfile("nonexistent-user")

        assert.Error(t, err)
        assert.Nil(t, user)
        assert.Equal(t, "user not found", err.Error())

        mockRepo.AssertExpectations(t)
    })

    t.Run("Profile with Repository Error", func(t *testing.T) {
        mockRepo.On("GetUserByID", "error-user").Return(nil, assert.AnError)

        user, err := authService.GetUserProfile("error-user")

        assert.Error(t, err)
        assert.Nil(t, user)

        mockRepo.AssertExpectations(t)
    })
}

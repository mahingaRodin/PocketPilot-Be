package services

import (
	"errors"
	"pocketpilot/internal/models"
	"pocketpilot/internal/utils"
)

type AuthService struct {
	userRepo UserRepository
	jwtSecret string
}

func NewAuthService(userRepo UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtSecret: jwtSecret,
	}
}

//register new user account
func (s *AuthService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil,err
	}

	if exists {
		return nil, errors.New("email already registered")
	}

	//hassh the pssword
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil,err
	}

	//user creation
	    user := &models.User{
        Email:        req.Email,
        PasswordHash: hashedPassword,
        FirstName:    req.FirstName,
        LastName:     req.LastName,
    }

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil,err
	}

	//generate jwt token
	token, err := utils.GenerateToken(user.ID,user.Email,s.jwtSecret)
	if err != nil {
		return nil,err
	}

	    return &models.AuthResponse{
        Token: token,
        User:  user,
    }, nil

}


func (s *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetUserByEmail(req.Email)

	if err != nil {
		return nil,err
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil , errors.New("invalid password")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User: user,
	}, nil
}


func (s *AuthService) GetUserProfile(userID string) (*models.User, error) {
    user, err := s.userRepo.GetUserByID(userID)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("user not found")
    }

    return user, nil
}
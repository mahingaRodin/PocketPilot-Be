package services

import (
	"errors"
	"pocketpilot-api/internal/models"
    "pocketpilot-api/internal/repository"
    "pocketpilot-api/internal/utils"
)

type AuthService struct {
	userRepo *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtSecret: jwtSecret,
	}
}

//register new user account
func (s *AuthService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	exists, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil,err
	}

	if exists != nil {
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
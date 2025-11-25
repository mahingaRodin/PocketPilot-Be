package services

import "pocketpilot-api/internal/models"

type UserRepository interface {
    GetUserByEmail(email string) (*models.User, error) 
	CreateUser(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	EmailExists(email string) (bool, error)
}
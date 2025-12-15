package services

import "pocketpilot/internal/models"

type UserRepository interface {
    GetUserByEmail(email string) (*models.User, error) 
	CreateUser(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	EmailExists(email string) (bool, error)
}

type ExpenseRepository interface {
    CreateExpense(*models.Expense) error
    GetExpenseByID(string) (*models.Expense, error)
    GetExpensesByID(string) (*models.Expense, error)
    GetExpensesByUser(string, int, int) ([]*models.Expense, error)
    UpdateExpense(*models.Expense) error
    DeleteExpense(string, string) error
    GetExpensesByTeam(string, int, int) ([]*models.Expense, error)
}

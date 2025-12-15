package services

import (
	"pocketpilot/internal/models"
	// "testing"
	// "time"

	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	// "github.com/stretchr/testify/require"
)

type MockExpenseRepository struct {
	mock.Mock
}

func (m *MockExpenseRepository) CreateExpense(expense *models.Expense) error {
	args := m.Called(expense)
	return args.Error(0)
}

func (m *MockExpenseRepository) GetExpensesByID(id string) (*models.Expense, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Expense), args.Error(1)
}

func (m *MockExpenseRepository) GetExpensesByUser(userID string, limit, offset int) ([]*models.Expense, error) {
    args := m.Called(userID, limit, offset)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]*models.Expense), args.Error(1)
}

func (m *MockExpenseRepository) UpdateExpense(expense *models.Expense) error {
    args := m.Called(expense)
    return args.Error(0)
}

func (m *MockExpenseRepository) DeleteExpense(id, userID string) error {
    args := m.Called(id, userID)
    return args.Error(0)
}

func (m *MockExpenseRepository) GetExpensesByTeam(teamID string, limit, offset int) ([]*models.Expense, error) {
    args := m.Called(teamID, limit, offset)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]*models.Expense), args.Error(1)
}

// func TestExpenseService_CreateExpense(t *testing.T) {
// 	mockExpenseRepo := new(MockExpenseRepository)
// 	mockUserRepo := new(MockUserRepository)
// 	// expsenseService := &ExpenseService{expenseRepo: mockExpenseRepo, userRepo: mockUserRepo}
// }



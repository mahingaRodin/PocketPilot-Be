package services

import (
	"errors"
	"pocketpilot/internal/models"
	"pocketpilot/internal/repository"
	"time"
)

type ExpenseService struct {
    expenseRepo ExpenseRepository
    userRepo    *repository.UserRepositoryImpl
}

func NewExpenseService(expenseRepo ExpenseRepository, userRepo *repository.UserRepositoryImpl) *ExpenseService {
    return &ExpenseService{
        expenseRepo: expenseRepo,
        userRepo:    userRepo,
    }
}

// CreateExpense creates a new expense for a user
func (s *ExpenseService) CreateExpense(userID string, req *models.CreateExpenseRequest) (*models.Expense, error) {
    // Validate expense date
    if _, err := time.Parse("2006-01-02", req.ExpenseDate); err != nil {
        return nil, errors.New("invalid expense date format, use YYYY-MM-DD")
    }

    expense := &models.Expense{
        UserID:         userID,
        TeamID:         req.TeamID,
        Amount:         req.Amount,
        Currency:       req.Currency,
        Description:    req.Description,
        Category:       req.Category,
        ExpenseDate:    req.ExpenseDate,
        ReceiptImageURL: req.ReceiptImageURL,
        Status:         "pending",
    }

    err := s.expenseRepo.CreateExpense(expense)
    if err != nil {
        return nil, err
    }

    return expense, nil
}

// GetExpense retrieves an expense by ID with authorization
func (s *ExpenseService) GetExpense(expenseID, userID string) (*models.Expense, error) {
    expense, err := s.expenseRepo.GetExpenseByID(expenseID)
    if err != nil {
        return nil, err
    }
    if expense == nil {
        return nil, errors.New("expense not found")
    }

    // Check if user owns the expense or it's a team expense they have access to
    if expense.UserID != userID {
        return nil, errors.New("access denied")
    }

    return expense, nil
}

// GetUserExpenses retrieves all expenses for a user
func (s *ExpenseService) GetUserExpenses(userID string, page, limit int) ([]*models.Expense, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }
    offset := (page - 1) * limit

    return s.expenseRepo.GetExpensesByUser(userID, limit, offset)
}

// UpdateExpense updates an existing expense
func (s *ExpenseService) UpdateExpense(expenseID, userID string, req *models.UpdateExpenseRequest) (*models.Expense, error) {
    // Get existing expense
    expense, err := s.expenseRepo.GetExpenseByID(expenseID)
    if err != nil {
        return nil, err
    }
    if expense == nil {
        return nil, errors.New("expense not found")
    }

    // Check ownership
    if expense.UserID != userID {
        return nil, errors.New("access denied")
    }

    // Update fields if provided
    if req.Amount != nil {
        expense.Amount = *req.Amount
    }
    if req.Description != nil {
        expense.Description = *req.Description
    }
    if req.Category != nil {
        expense.Category = *req.Category
    }
    if req.ExpenseDate != nil {
        if _, err := time.Parse("2006-01-02", *req.ExpenseDate); err != nil {
            return nil, errors.New("invalid expense date format, use YYYY-MM-DD")
        }
        expense.ExpenseDate = *req.ExpenseDate
    }
    if req.Status != nil {
        expense.Status = *req.Status
    }

    err = s.expenseRepo.UpdateExpense(expense)
    if err != nil {
        return nil, err
    }

    return expense, nil
}

// DeleteExpense deletes an expense
func (s *ExpenseService) DeleteExpense(expenseID, userID string) error {
    return s.expenseRepo.DeleteExpense(expenseID, userID)
}

// GetTeamExpenses retrieves expenses for a team
func (s *ExpenseService) GetTeamExpenses(teamID, userID string, page, limit int) ([]*models.Expense, error) {
    // TODO: Add team membership validation
    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }
    offset := (page - 1) * limit

    return s.expenseRepo.GetExpensesByTeam(teamID, limit, offset)
}
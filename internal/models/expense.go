package models

import (
    "time"
)

type Expense struct {
    ID             string    `json:"id"`
    UserID         string    `json:"user_id"`
    TeamID         *string   `json:"team_id,omitempty"`
    Amount         float64   `json:"amount"`
    Currency       string    `json:"currency"`
    Description    string    `json:"description"`
    Category       string    `json:"category"`
    ExpenseDate    string    `json:"expense_date"` // YYYY-MM-DD
    ReceiptImageURL *string  `json:"receipt_image_url,omitempty"`
    Status         string    `json:"status"` // pending, approved, rejected
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

type CreateExpenseRequest struct {
    Amount         float64 `json:"amount" binding:"required,gt=0"`
    Currency       string  `json:"currency" binding:"required"`
    Description    string  `json:"description" binding:"required"`
    Category       string  `json:"category" binding:"required"`
    ExpenseDate    string  `json:"expense_date" binding:"required"`
    TeamID         *string `json:"team_id,omitempty"`
    ReceiptImageURL *string `json:"receipt_image_url,omitempty"`
}

type UpdateExpenseRequest struct {
    Amount      *float64 `json:"amount,omitempty"`
    Description *string  `json:"description,omitempty"`
    Category    *string  `json:"category,omitempty"`
    ExpenseDate *string  `json:"expense_date,omitempty"`
    Status      *string  `json:"status,omitempty"`
}

type ExpenseResponse struct {
    Expense
    User *User `json:"user,omitempty"`
    // Team *Team `json:"team,omitempty"`
}
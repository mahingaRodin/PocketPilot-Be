package repository

import (
    "database/sql"
    "errors"
    "pocketpilot/internal/models"
    // "fmt"
    "time"
)

type ExpenseRepositoryImpl struct {
    db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepositoryImpl {
    return &ExpenseRepositoryImpl{db: db}
}

// CreateExpense creates a new expense
func (r *ExpenseRepositoryImpl) CreateExpense(expense *models.Expense) error {
    query := `
        INSERT INTO expenses (user_id, team_id, amount, currency, description, category, expense_date, receipt_image_url, status)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, created_at, updated_at
    `
    
    err := r.db.QueryRow(
        query,
        expense.UserID,
        expense.TeamID,
        expense.Amount,
        expense.Currency,
        expense.Description,
        expense.Category,
        expense.ExpenseDate,
        expense.ReceiptImageURL,
        expense.Status,
    ).Scan(&expense.ID, &expense.CreatedAt, &expense.UpdatedAt)
    
    return err
}

// GetExpenseByID retrieves an expense by ID
func (r *ExpenseRepositoryImpl) GetExpenseByID(id string) (*models.Expense, error) {
    query := `
        SELECT id, user_id, team_id, amount, currency, description, category, 
               expense_date, receipt_image_url, status, created_at, updated_at
        FROM expenses 
        WHERE id = $1
    `
    
    expense := &models.Expense{}
    err := r.db.QueryRow(query, id).Scan(
        &expense.ID,
        &expense.UserID,
        &expense.TeamID,
        &expense.Amount,
        &expense.Currency,
        &expense.Description,
        &expense.Category,
        &expense.ExpenseDate,
        &expense.ReceiptImageURL,
        &expense.Status,
        &expense.CreatedAt,
        &expense.UpdatedAt,
    )
    
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    
    return expense, nil
}

// GetExpensesByUser retrieves all expenses for a user
func (r *ExpenseRepositoryImpl) GetExpensesByUser(userID string, limit, offset int) ([]*models.Expense, error) {
    query := `
        SELECT id, user_id, team_id, amount, currency, description, category, 
               expense_date, receipt_image_url, status, created_at, updated_at
        FROM expenses 
        WHERE user_id = $1
        ORDER BY expense_date DESC, created_at DESC
        LIMIT $2 OFFSET $3
    `
    
    rows, err := r.db.Query(query, userID, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var expenses []*models.Expense
    for rows.Next() {
        expense := &models.Expense{}
        err := rows.Scan(
            &expense.ID,
            &expense.UserID,
            &expense.TeamID,
            &expense.Amount,
            &expense.Currency,
            &expense.Description,
            &expense.Category,
            &expense.ExpenseDate,
            &expense.ReceiptImageURL,
            &expense.Status,
            &expense.CreatedAt,
            &expense.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        expenses = append(expenses, expense)
    }
    
    return expenses, nil
}

// UpdateExpense updates an existing expense
func (r *ExpenseRepositoryImpl) UpdateExpense(expense *models.Expense) error {
    query := `
        UPDATE expenses 
        SET amount = $1, currency = $2, description = $3, category = $4, 
            expense_date = $5, receipt_image_url = $6, status = $7, updated_at = $8
        WHERE id = $9 AND user_id = $10
        RETURNING updated_at
    `
    
    expense.UpdatedAt = time.Now()
    err := r.db.QueryRow(
        query,
        expense.Amount,
        expense.Currency,
        expense.Description,
        expense.Category,
        expense.ExpenseDate,
        expense.ReceiptImageURL,
        expense.Status,
        expense.UpdatedAt,
        expense.ID,
        expense.UserID,
    ).Scan(&expense.UpdatedAt)
    
    return err
}

// DeleteExpense deletes an expense
func (r *ExpenseRepositoryImpl) DeleteExpense(id, userID string) error {
    query := `DELETE FROM expenses WHERE id = $1 AND user_id = $2`
    result, err := r.db.Exec(query, id, userID)
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return errors.New("expense not found or access denied")
    }
    
    return nil
}

// GetExpensesByTeam retrieves expenses for a team
func (r *ExpenseRepositoryImpl) GetExpensesByTeam(teamID string, limit, offset int) ([]*models.Expense, error) {
    query := `
        SELECT id, user_id, team_id, amount, currency, description, category, 
               expense_date, receipt_image_url, status, created_at, updated_at
        FROM expenses 
        WHERE team_id = $1
        ORDER BY expense_date DESC, created_at DESC
        LIMIT $2 OFFSET $3
    `
    
    rows, err := r.db.Query(query, teamID, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var expenses []*models.Expense
    for rows.Next() {
        expense := &models.Expense{}
        err := rows.Scan(
            &expense.ID,
            &expense.UserID,
            &expense.TeamID,
            &expense.Amount,
            &expense.Currency,
            &expense.Description,
            &expense.Category,
            &expense.ExpenseDate,
            &expense.ReceiptImageURL,
            &expense.Status,
            &expense.CreatedAt,
            &expense.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        expenses = append(expenses, expense)
    }
    
    return expenses, nil
}
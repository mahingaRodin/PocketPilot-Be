package repository

import (
    "database/sql"
    "errors"
    "pocketpilot-api/internal/models"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *models.User) error {
    query := `
        INSERT INTO users (email, password_hash, first_name, last_name)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at
    `
    
    err := r.db.QueryRow(
        query,
        user.Email,
        user.PasswordHash,
        user.FirstName,
        user.LastName,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
    
    return err
}

// GetUserByEmail finds a user by email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
    query := `
        SELECT id, email, password_hash, first_name, last_name, created_at, updated_at
        FROM users 
        WHERE email = $1
    `
    
    user := &models.User{}
    err := r.db.QueryRow(query, email).Scan(
        &user.ID,
        &user.Email,
        &user.PasswordHash,
        &user.FirstName,
        &user.LastName,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil // User not found
        }
        return nil, err
    }
    
    return user, nil
}

// GetUserByID finds a user by ID
func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
    query := `
        SELECT id, email, password_hash, first_name, last_name, created_at, updated_at
        FROM users 
        WHERE id = $1
    `
    
    user := &models.User{}
    err := r.db.QueryRow(query, id).Scan(
        &user.ID,
        &user.Email,
        &user.PasswordHash,
        &user.FirstName,
        &user.LastName,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    
    return user, nil
}

// EmailExists checks if an email is already registered
func (r *UserRepository) EmailExists(email string) (bool, error) {
    query := `SELECT COUNT(*) FROM users WHERE email = $1`
    
    var count int
    err := r.db.QueryRow(query, email).Scan(&count)
    if err != nil {
        return false, err
    }
    
    return count > 0, nil
}
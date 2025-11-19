package repository

import (
	"backend/internal/database"
	"backend/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db:db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, name, role, phone, avatar_url, is_active, register_date, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	now := time.Now()

	user.RegisterDate = now
	user.UpdatedAt = now

	_, err := r.db.Pool.Exec(ctx, query, user.ID, user.Email, user.PasswordHash, user.Name, user.Role,
        user.Phone, user.AvatarURL, user.IsActive, user.RegisterDate, user.UpdatedAt)

	return err
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `
        SELECT id, email, name, role, phone, avatar_url, is_active, last_login, register_date, updated_at
        FROM users WHERE id = $1
    `

	var user models.User
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
        &user.Phone, &user.AvatarURL, &user.IsActive, &user.LastLogin,
        &user.RegisterDate, &user.UpdatedAt,)

	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("usr not found")
		}
		return  nil, err
	}

	return &user, nil

} 

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    query := `
        SELECT id, email, name, role, phone, avatar_url, is_active, last_login, register_date, updated_at
        FROM users WHERE email = $1
    `
    
    var user models.User
    err := r.db.Pool.QueryRow(ctx, query, email).Scan(
        &user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
        &user.Phone, &user.AvatarURL, &user.IsActive, &user.LastLogin,
        &user.RegisterDate, &user.UpdatedAt,
    )
    
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, err
    }
    
    return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
    query := `
        UPDATE users 
        SET name = $1, phone = $2, avatar_url = $3, updated_at = $4
        WHERE id = $5
    `
    
    user.UpdatedAt = time.Now()
    _, err := r.db.Pool.Exec(ctx, query,
        user.Name, user.Phone, user.AvatarURL, user.UpdatedAt, user.ID)
    
    return err
}

func (r *UserRepository) GetUsers(ctx context.Context) ([]*models.User, error) {
    query := `
        SELECT id, email, name, role, phone, avatar_url, 
               is_active, last_login, register_date, updated_at
        FROM users 
        ORDER BY register_date DESC
    `
    rows, err := r.db.Pool.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []*models.User
    for rows.Next() {
        var user models.User
        err := rows.Scan(
            &user.ID, &user.Email, &user.Name, &user.Role,
            &user.Phone, &user.AvatarURL, &user.IsActive, &user.LastLogin,
            &user.RegisterDate, &user.UpdatedAt,
        )

        if err != nil {
            return nil, err
        }
        users = append(users, &user)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return users, nil
}
package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleTeacher Role = "teacher"
	RoleStudent Role = "student"
	RoleUser Role = "user"
)

type User struct {
	ID 				string 		`json:"id" db:"id"`
	Email 			string 		`json:"email" db:"email"`
	PasswordHash 	string 		`json:"-" db:"password_hash"`
	Name 			string 		`json:"name" db:"name"`
	Role 			Role 		`json:"role" db:"role"`
	Phone 			*string 	`json:"phone,omitempty" db:"phone"`
	AvatarURL 		*string 	`json:"avatar_url,omitempty" db:"avatar_url"`
	IsActive 		bool 		`json:"is_active" db:"is_active"`
	LastLogin 		*time.Time 	`json:"last_login,omitempty" db:"last_login"`
	RegisterDate 	time.Time 	`json:"register_date" db:"register_date"`
	UpdatedAt 		time.Time 	`json:"updated_at" db:"updated_at"`
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	return err == nil
}
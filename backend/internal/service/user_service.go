package service

import (
	"context"
	"errors"
	"time"

	"backend/internal/models"
	"backend/internal/repository"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

type CreateUserRequest struct {
	Email    string 		`json:"email"`
    Password string 		`json:"password"`
    Name     string 		`json:"name"`
    Role     models.Role 	`json:"role"`
    Phone    string 		`json:"phone,omitempty"`
}

func (s *UserService) CreateUser (ctx context.Context, req CreateUserRequest) (*models.User, error) {
	existingUser, _ := s.userRepo.GetUserByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	user := &models.User{
		ID: generateID(),
		Email: req.Email,
		Name: req.Name,
		Role: req.Role,
		Phone: &req.Phone,
		IsActive: true, 
		UpdatedAt: time.Now(),
		RegisterDate: time.Now(),
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil

}

func generateID() string {
    return uuid.New().String()
}

func (s *UserService) GetUser(ctx context.Context, id string) (*models.User, error) {
    return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserService) GetUsers(ctx context.Context) ([]*models.User, error) {
	return s.userRepo.GetUsers(ctx)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    return s.userRepo.GetUserByEmail(ctx, email)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, updates map[string]interface{}) error {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err == nil {
		return err
	}

	if name, ok := updates["name"].(string); ok {
		user.Name = name
	}

	if phone, ok := updates["phone"].(string); ok {
		user.Phone = &phone
	}

	if avatarURL, ok := updates["avatar_url"].(string); ok {
		user.AvatarURL = &avatarURL
	}

	return s.userRepo.UpdateUser(ctx, user)
}




package service

import (
	"boiler/internal/user/model"
	"boiler/internal/user/repository"
)

// UserService handles user-related business logic
type UserService struct {
	UserRepository *repository.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{UserRepository: userRepository}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user *model.User) error {
	return s.UserRepository.CreateUser(user)
}

// GetUser gets a user by id
func (s *UserService) GetUser(id string) (*model.User, error) {
	return s.UserRepository.GetUser(id)
}

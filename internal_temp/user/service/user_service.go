package service

import (
	"github.com/FeisalDy/nogo/internal/common/errors"
	"github.com/FeisalDy/nogo/internal/common/utils"
	"github.com/FeisalDy/nogo/internal/user/dto"
	"github.com/FeisalDy/nogo/internal/user/model"
	"github.com/FeisalDy/nogo/internal/user/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepository,
	}
}

// CreateUser creates a new user without role assignment
// Role assignment should be handled by the application layer
func (s *UserService) CreateUser(registerDTO *dto.RegisterUserDTO) (*model.User, error) {
	existingUser, _ := s.userRepo.GetUserByEmail(registerDTO.Email)
	if existingUser != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(registerDTO.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: &registerDTO.Username,
		Email:    registerDTO.Email,
		Password: &hashedPassword,
	}

	userCreated, err := s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return userCreated, nil
}

func (s *UserService) Login(loginDTO *dto.LoginUserDTO) (*model.User, error) {
	user, err := s.userRepo.GetUserByEmail(loginDTO.Email)
	if err != nil {
		return nil, errors.ErrUserInvalidCredentials
	}

	if user.Password == nil {
		return nil, errors.ErrUserInvalidCredentials
	}

	if !utils.ComparePassword(*user.Password, loginDTO.Password) {
		return nil, errors.ErrUserInvalidCredentials
	}

	return user, nil
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

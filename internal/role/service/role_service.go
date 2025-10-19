package service

import (
	"github.com/FeisalDy/nogo/internal/common/errors"
	"github.com/FeisalDy/nogo/internal/role/dto"
	"github.com/FeisalDy/nogo/internal/role/model"
	"github.com/FeisalDy/nogo/internal/role/repository"
	userRepo "github.com/FeisalDy/nogo/internal/user/repository"
)

type RoleService struct {
	roleRepo *repository.RoleRepository
	userRepo *userRepo.UserRepository
}

func NewRoleService(roleRepo *repository.RoleRepository, userRepo *userRepo.UserRepository) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

func (s *RoleService) CreateRole(req dto.CreateRoleDTO) (*model.Role, error) {
	exists, err := s.roleRepo.ExistsByName(req.Name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.ErrRoleAlreadyExists
	}

	role := &model.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) GetRoleByID(id uint) (*model.Role, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.ErrRoleNotFound
	}

	return role, nil
}

func (s *RoleService) GetRoleByName(name string) (*model.Role, error) {
	role, err := s.roleRepo.GetByName(name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.ErrRoleNotFound
	}

	return role, nil
}

func (s *RoleService) GetAllRoles() ([]model.Role, error) {
	roles, err := s.roleRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *RoleService) UpdateRole(id uint, req dto.UpdateRoleDTO) (*model.Role, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.ErrRoleNotFound
	}

	if req.Name != nil && *req.Name != role.Name {
		exists, err := s.roleRepo.ExistsByName(*req.Name)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.ErrRoleAlreadyExists
		}
		role.Name = *req.Name
	}

	if req.Description != nil {
		role.Description = req.Description
	}

	if err := s.roleRepo.Update(role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) DeleteRole(id uint) error {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.ErrRoleNotFound
	}

	return s.roleRepo.Delete(id)
}

func (s *RoleService) AssignRoleToUser(userID, roleID uint) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.ErrRoleNotFound
	}

	hasRole, err := s.userRepo.HasRoleByID(userID, roleID)
	if err != nil {
		return err
	}
	if hasRole {
		return errors.ErrUserAlreadyHasRole
	}

	return s.userRepo.AssignRoleToUser(userID, roleID)
}

func (s *RoleService) RemoveRoleFromUser(userID, roleID uint) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.ErrRoleNotFound
	}

	hasRole, err := s.userRepo.HasRoleByID(userID, roleID)
	if err != nil {
		return err
	}
	if !hasRole {
		return errors.ErrUserDoesNotHaveRole
	}

	return s.userRepo.RemoveRoleFromUser(userID, roleID)
}

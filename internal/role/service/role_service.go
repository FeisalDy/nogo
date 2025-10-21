package service

import (
	"github.com/FeisalDy/nogo/internal/common/errors"
	"github.com/FeisalDy/nogo/internal/role/dto"
	"github.com/FeisalDy/nogo/internal/role/model"
	"github.com/FeisalDy/nogo/internal/role/repository"
)

type RoleService struct {
	roleRepo *repository.RoleRepository
}

func NewRoleService(roleRepo *repository.RoleRepository) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
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

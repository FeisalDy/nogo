package handler

import (
	"net/http"
	"strconv"

	"github.com/FeisalDy/nogo/internal/common/errors"
	"github.com/FeisalDy/nogo/internal/common/utils"
	"github.com/FeisalDy/nogo/internal/role/dto"
	"github.com/FeisalDy/nogo/internal/role/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RoleHandler struct {
	roleService *service.RoleService
	validator   *validator.Validate
}

func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
		validator:   validator.New(),
	}
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req dto.CreateRoleDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err, errors.ErrCodeRoleValidation)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.RespondValidationError(c, err, errors.ErrCodeRoleValidation)
		return
	}

	role, err := h.roleService.CreateRole(req)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondSuccess(c, http.StatusCreated, role, "Role created successfully")
}

func (h *RoleHandler) GetRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.RespondWithAppError(c, errors.ErrInvalidParam.WithDetails(map[string]any{
			"reason": err.Error(),
		}))
		return
	}

	role, err := h.roleService.GetRoleByID(uint(id))
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, role, "Role retrieved successfully")
}

func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.roleService.GetAllRoles()
	if err != nil {
		utils.RespondWithAppError(c, errors.ErrInternalServer)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, roles, "Roles retrieved successfully")
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.RespondWithAppError(c, errors.ErrInvalidParam.WithDetails(map[string]any{
			"reason": err.Error(),
		}))
		return
	}

	var req dto.UpdateRoleDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err, errors.ErrCodeRoleValidation)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.RespondValidationError(c, err, errors.ErrCodeRoleValidation)
		return
	}

	role, err := h.roleService.UpdateRole(uint(id), req)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, role, "Role updated successfully")
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.RespondWithAppError(c, errors.ErrInvalidParam)
		return
	}

	if err := h.roleService.DeleteRole(uint(id)); err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, gin.H{"message": "Role deleted successfully"}, "Role deleted successfully")
}

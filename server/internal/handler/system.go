package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/pkg/response"
	"github.com/rms-survey/server/internal/service"
)

// SystemHandler handles all /api/system/* endpoints.
type SystemHandler struct {
	svc     *service.SystemService
	userSvc *service.UserService
}

func NewSystemHandler(svc *service.SystemService) *SystemHandler {
	return &SystemHandler{svc: svc}
}

// NewSystemHandlerWithUser creates a SystemHandler with user service for system/user/* routes.
func NewSystemHandlerWithUser(svc *service.SystemService, userSvc *service.UserService) *SystemHandler {
	return &SystemHandler{svc: svc, userSvc: userSvc}
}

// RegisterRoutes wires all system routes onto the given gin.RouterGroup.
// The group should be /api/system.
func (h *SystemHandler) RegisterRoutes(rg *gin.RouterGroup) {
	// Public: GET /api/system
	rg.GET("", h.GetSystemInfo)

	// System info update (admin only)
	rg.POST("/update", h.UpdateSystemInfo)

	// AI setting
	rg.GET("/aiSetting", h.GetAISetting)

	// Role management
	rg.GET("/role/list", h.ListRoles)
	rg.POST("/role/create", h.CreateRole)
	rg.POST("/role/update", h.UpdateRole)
	rg.POST("/role/delete", h.DeleteRole)

	// System user management (mirrors /api/user/* but under /api/system/user/*)
	if h.userSvc != nil {
		userGrp := rg.Group("/user")
		userGrp.GET("/list", h.SystemListUsers)
		userGrp.POST("/create", h.SystemCreateUser)
		userGrp.POST("/update", h.SystemUpdateUser)
		userGrp.POST("/delete", h.SystemDeleteUser)
	}
}

// GetSystemInfo returns current system info (public, no auth required).
func (h *SystemHandler) GetSystemInfo(c *gin.Context) {
	info, err := h.svc.GetSysInfo()
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, info)
}

// UpdateSystemInfo updates system settings (admin only).
func (h *SystemHandler) UpdateSystemInfo(c *gin.Context) {
	var req dto.SysInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.UpdateSysInfo(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// ListRoles returns paginated role list.
func (h *SystemHandler) ListRoles(c *gin.Context) {
	var q dto.RoleQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.ListRoles(q)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

// CreateRole creates a new role.
func (h *SystemHandler) CreateRole(c *gin.Context) {
	var req dto.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.CreateRole(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// UpdateRole updates an existing role.
func (h *SystemHandler) UpdateRole(c *gin.Context) {
	var req dto.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.UpdateRole(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// DeleteRole deletes a role, refusing if it's the last one.
func (h *SystemHandler) DeleteRole(c *gin.Context) {
	var req dto.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if req.ID == "" {
		c.Status(http.StatusOK)
		return
	}
	total, err := h.svc.CountRoles()
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if total <= 1 {
		response.Fail(c, response.CodeError, "must retain at least one role")
		return
	}
	if err := h.svc.DeleteRole(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// GetAISetting returns the AI configuration (token masked).
func (h *SystemHandler) GetAISetting(c *gin.Context) {
	info, err := h.svc.GetSysInfo()
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, info)
}

// SystemListUsers handles GET /api/system/user/list
func (h *SystemHandler) SystemListUsers(c *gin.Context) {
	var req dto.UserQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.userSvc.ListUsers(req)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

// SystemCreateUser handles POST /api/system/user/create
func (h *SystemHandler) SystemCreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.userSvc.CreateUser(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// SystemUpdateUser handles POST /api/system/user/update
func (h *SystemHandler) SystemUpdateUser(c *gin.Context) {
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.userSvc.UpdateUser(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// SystemDeleteUser handles POST /api/system/user/delete
func (h *SystemHandler) SystemDeleteUser(c *gin.Context) {
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if req.ID == "" {
		c.Status(http.StatusOK)
		return
	}
	if err := h.userSvc.DeleteUser(req.ID); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

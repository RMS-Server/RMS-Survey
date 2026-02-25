package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/model"
	"github.com/surveyking/server/internal/pkg/response"
	"github.com/surveyking/server/internal/service"
	"github.com/xuri/excelize/v2"
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

	// Permission list
	rg.Any("/permission/list", h.ListPermissions)
	rg.GET("/permission/diff", h.PermissionDiff)

	// Dept management
	rg.GET("/dept/list", h.ListDepts)
	rg.POST("/dept/create", h.CreateDept)
	rg.POST("/dept/update", h.UpdateDept)
	rg.POST("/dept/delete", h.DeleteDept)
	rg.POST("/dept/sort", h.SortDept)

	// Dict management
	rg.GET("/dict/list", h.ListDicts)
	rg.POST("/dict/create", h.CreateDict)
	rg.POST("/dict/update", h.UpdateDict)
	rg.POST("/dict/delete", h.DeleteDict)

	// DictItem management
	rg.GET("/dictItem/list", h.ListDictItems)
	rg.POST("/dictItem/create", h.CreateDictItem)
	rg.POST("/dictItem/update", h.UpdateDictItem)
	rg.POST("/dictItem/delete", h.DeleteDictItem)
	rg.POST("/dictItem/import", h.ImportDictItem)

	// Position management
	rg.GET("/position/list", h.ListPositions)
	rg.POST("/position/create", h.CreatePosition)
	rg.POST("/position/update", h.UpdatePosition)
	rg.POST("/position/delete", h.DeletePosition)

	// Tag management
	rg.GET("/tag/list", h.ListTags)
	rg.POST("/tag/create", h.CreateTag)
	rg.POST("/tag/delete", h.DeleteTag)

	// Check username
	rg.GET("/checkUsernameExist", h.CheckUsernameExist)

	// RSA key rotation
	rg.POST("/rotateRsaKey", h.RotateRsaKey)

	// System user management (mirrors /api/user/* but under /api/system/user/*)
	if h.userSvc != nil {
		userGrp := rg.Group("/user")
		userGrp.GET("/list", h.SystemListUsers)
		userGrp.POST("/create", h.SystemCreateUser)
		userGrp.POST("/update", h.SystemUpdateUser)
		userGrp.POST("/delete", h.SystemDeleteUser)
		userGrp.POST("/updatePosition", h.SystemUpdateUserPosition)
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

// ListPermissions returns a static permission list (placeholder).
func (h *SystemHandler) ListPermissions(c *gin.Context) {
	perms := []dto.PermissionView{
		{Code: "system:role:list", Name: "Role List", Module: "system"},
		{Code: "system:role:create", Name: "Role Create", Module: "system"},
		{Code: "system:role:update", Name: "Role Update", Module: "system"},
		{Code: "system:role:delete", Name: "Role Delete", Module: "system"},
		{Code: "system:user:list", Name: "User List", Module: "system"},
		{Code: "system:user:create", Name: "User Create", Module: "system"},
		{Code: "system:user:update", Name: "User Update", Module: "system"},
		{Code: "system:user:delete", Name: "User Delete", Module: "system"},
		{Code: "system:dept:list", Name: "Dept List", Module: "system"},
		{Code: "system:dept:create", Name: "Dept Create", Module: "system"},
		{Code: "system:dept:update", Name: "Dept Update", Module: "system"},
		{Code: "system:dept:delete", Name: "Dept Delete", Module: "system"},
		{Code: "system:dict:list", Name: "Dict List", Module: "system"},
		{Code: "system:dict:create", Name: "Dict Create", Module: "system"},
		{Code: "system:dict:update", Name: "Dict Update", Module: "system"},
		{Code: "system:dict:delete", Name: "Dict Delete", Module: "system"},
		{Code: "system:dictItem:list", Name: "DictItem List", Module: "system"},
		{Code: "system:dictItem:create", Name: "DictItem Create", Module: "system"},
		{Code: "system:dictItem:update", Name: "DictItem Update", Module: "system"},
		{Code: "system:dictItem:delete", Name: "DictItem Delete", Module: "system"},
		{Code: "system:position:list", Name: "Position List", Module: "system"},
		{Code: "system:position:create", Name: "Position Create", Module: "system"},
		{Code: "system:position:update", Name: "Position Update", Module: "system"},
		{Code: "system:position:delete", Name: "Position Delete", Module: "system"},
	}
	response.OK(c, perms)
}

// ListDepts returns the department tree.
func (h *SystemHandler) ListDepts(c *gin.Context) {
	depts, err := h.svc.ListDepts()
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, depts)
}

// CreateDept creates a new department.
func (h *SystemHandler) CreateDept(c *gin.Context) {
	var req dto.DeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.CreateDept(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// UpdateDept updates an existing department.
func (h *SystemHandler) UpdateDept(c *gin.Context) {
	var req dto.DeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.UpdateDept(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// DeleteDept deletes a department, refusing if it's the last one.
func (h *SystemHandler) DeleteDept(c *gin.Context) {
	var req dto.DeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if req.ID == "" {
		c.Status(http.StatusOK)
		return
	}
	total, err := h.svc.CountDepts()
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if total <= 1 {
		response.Fail(c, response.CodeError, "must retain at least one department")
		return
	}
	if err := h.svc.DeleteDept(req.ID); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// SortDept handles department sort order updates.
func (h *SystemHandler) SortDept(c *gin.Context) {
	var items []dto.DeptSortRequest
	if err := c.ShouldBindJSON(&items); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.UpdateDeptSortCodes(items); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// ListDicts returns paginated dict list.
func (h *SystemHandler) ListDicts(c *gin.Context) {
	var q dto.CommDictQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.ListDicts(q)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

// CreateDict creates a new data dictionary.
func (h *SystemHandler) CreateDict(c *gin.Context) {
	var req dto.CommDictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.CreateDict(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// UpdateDict updates an existing data dictionary.
func (h *SystemHandler) UpdateDict(c *gin.Context) {
	var req dto.CommDictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.UpdateDict(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// DeleteDict deletes a data dictionary.
func (h *SystemHandler) DeleteDict(c *gin.Context) {
	var req dto.CommDictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.DeleteDict(req.ID); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// ListDictItems returns paginated dict item list.
func (h *SystemHandler) ListDictItems(c *gin.Context) {
	var q dto.CommDictItemQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.ListDictItems(q)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

// CreateDictItem creates a new dict item.
func (h *SystemHandler) CreateDictItem(c *gin.Context) {
	var req dto.CommDictItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.SaveDictItem(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// UpdateDictItem updates an existing dict item.
func (h *SystemHandler) UpdateDictItem(c *gin.Context) {
	var req dto.CommDictItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.SaveDictItem(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// DeleteDictItem deletes a dict item.
func (h *SystemHandler) DeleteDictItem(c *gin.Context) {
	var req dto.CommDictItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.DeleteDictItem(req.ID); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// ListPositions returns paginated position list.
func (h *SystemHandler) ListPositions(c *gin.Context) {
	var q dto.PositionQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.ListPositions(q)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

// CreatePosition creates a new position.
func (h *SystemHandler) CreatePosition(c *gin.Context) {
	var req dto.PositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.CreatePosition(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// UpdatePosition updates an existing position.
func (h *SystemHandler) UpdatePosition(c *gin.Context) {
	var req dto.PositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.UpdatePosition(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// DeletePosition deletes a position.
func (h *SystemHandler) DeletePosition(c *gin.Context) {
	var req dto.PositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.DeletePosition(req.ID); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// ListTags returns paginated tag list.
func (h *SystemHandler) ListTags(c *gin.Context) {
	var q dto.SystemTagQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.ListTags(q)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

// CreateTag creates a new tag.
func (h *SystemHandler) CreateTag(c *gin.Context) {
	var req dto.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.CreateTag(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// DeleteTag deletes a tag.
func (h *SystemHandler) DeleteTag(c *gin.Context) {
	var req dto.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.DeleteTag(req.ID); err != nil {
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

// PermissionDiff is a no-op stub (compares DB vs code permissions).
func (h *SystemHandler) PermissionDiff(c *gin.Context) {
	c.Status(http.StatusOK)
}

// CheckUsernameExist returns true if the username is already taken.
func (h *SystemHandler) CheckUsernameExist(c *gin.Context) {
	username := c.Query("username")
	if h.userSvc == nil {
		response.OK(c, false)
		return
	}
	response.OK(c, h.userSvc.CheckUsernameExist(username))
}

// ImportDictItem handles POST /api/system/dictItem/import (stub).
func (h *SystemHandler) ImportDictItem(c *gin.Context) {
	dictCode := c.Query("dictCode")
	if dictCode == "" {
		response.Fail(c, response.CodeError, "dictCode is required")
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		response.Fail(c, response.CodeError, "failed to read uploaded file: "+err.Error())
		return
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		response.Fail(c, response.CodeError, "failed to parse Excel file: "+err.Error())
		return
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		response.Fail(c, response.CodeError, "Excel file has no sheets")
		return
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		response.Fail(c, response.CodeError, "failed to read sheet rows: "+err.Error())
		return
	}

	// Collect all items first, then insert in one transaction.
	// Skip header row; columns: itemName, itemValue
	items := make([]model.CommDictItem, 0, len(rows))
	for i, row := range rows {
		if i == 0 || len(row) < 2 {
			continue
		}
		items = append(items, model.CommDictItem{
			ID:        uuid.New().String(),
			DictCode:  dictCode,
			ItemName:  row[0],
			ItemValue: row[1],
		})
	}
	if err := h.svc.BatchImportDictItems(items); err != nil {
		response.Fail(c, response.CodeError, "import failed: "+err.Error())
		return
	}
	c.Status(http.StatusOK)
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

// SystemUpdateUserPosition handles POST /api/system/user/updatePosition (stub).
func (h *SystemHandler) SystemUpdateUserPosition(c *gin.Context) {
	var req dto.UpdateUserPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if req.ID == "" {
		response.Fail(c, response.CodeError, "id is required")
		return
	}
	if err := h.userSvc.UpdateUserPosition(req.ID, req.DeptID, req.PositionIDs); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// RotateRsaKey handles POST /api/system/rotateRsaKey — forces RSA key pair regeneration.
func (h *SystemHandler) RotateRsaKey(c *gin.Context) {
	if h.userSvc == nil {
		response.Fail(c, response.CodeError, "user service unavailable")
		return
	}
	pub, err := h.userSvc.RotateRSAKey()
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, pub)
}

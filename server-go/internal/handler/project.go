package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/pkg/response"
	"github.com/surveyking/server/internal/repository"
	"github.com/surveyking/server/internal/service"
	"gorm.io/gorm"
)

// ProjectHandler handles all project HTTP endpoints.
type ProjectHandler struct {
	svc *service.ProjectService
}

// NewProjectHandler creates a ProjectHandler wired to the given DB.
func NewProjectHandler(db *gorm.DB) *ProjectHandler {
	repo := repository.NewProjectRepo(db)
	svc := service.NewProjectService(repo)
	return &ProjectHandler{svc: svc}
}

// RegisterRoutes wires all project routes onto the provided gin group.
func (h *ProjectHandler) RegisterRoutes(project gin.IRouter) {
	project.GET("/list", h.ListProjects)
	project.GET("", h.GetProject)
	project.GET("/setting", h.GetSetting)
	project.POST("/create", h.CreateProject)
	project.POST("/update", h.UpdateProject)
	project.POST("/delete", h.DeleteProject)
	project.GET("/trash", h.GetDeleted)
	project.POST("/destroy", h.DestroyProject)
	project.POST("/restore", h.RestoreProject)
	project.GET("/partner/list", h.ListPartners)
	project.POST("/partner/create", h.AddPartner)
	project.POST("/partner/delete", h.RemovePartner)
	project.GET("/partner/download", h.DownloadPartner)
	project.POST("/partner/import", h.ImportPartner)
	project.POST("/selectUser", h.SelectUser)
	project.POST("/selectDept", h.SelectDept)
	project.POST("/selectRole", h.SelectRole)
	project.POST("/selectPosition", h.SelectPosition)
	project.POST("/selectDict", h.SelectDict)
	project.POST("/selectTemplate", h.SelectTemplate)
	project.POST("/selectRepo", h.SelectRepo)
	project.POST("/selectTag", h.SelectTag)
}

func currentUser(c *gin.Context) *dto.UserInfo {
	if v, ok := c.Get("currentUser"); ok {
		if u, ok := v.(*dto.UserInfo); ok {
			return u
		}
	}
	return &dto.UserInfo{}
}

func (h *ProjectHandler) ListProjects(c *gin.Context) {
	var query dto.ProjectQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.ListProjects(&query, currentUser(c))
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.Fail(c, response.CodeError, "id required")
		return
	}
	result, err := h.svc.GetProject(id)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *ProjectHandler) GetSetting(c *gin.Context) {
	projectID := c.Query("projectId")
	if projectID == "" {
		response.Fail(c, response.CodeError, "projectId required")
		return
	}
	result, err := h.svc.GetSetting(projectID)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req dto.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.CreateProject(&req, currentUser(c))
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	var req dto.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.UpdateProject(&req, currentUser(c)); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	var req dto.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.DeleteProject(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *ProjectHandler) GetDeleted(c *gin.Context) {
	result, err := h.svc.GetDeleted(currentUser(c))
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *ProjectHandler) DestroyProject(c *gin.Context) {
	var req dto.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.DestroyProject(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *ProjectHandler) RestoreProject(c *gin.Context) {
	var req dto.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.RestoreProject(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *ProjectHandler) ListPartners(c *gin.Context) {
	var query dto.ProjectPartnerQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.ListPartners(&query)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *ProjectHandler) AddPartner(c *gin.Context) {
	var req dto.ProjectPartnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.AddPartner(&req, currentUser(c)); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *ProjectHandler) RemovePartner(c *gin.Context) {
	var req dto.ProjectPartnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.RemovePartner(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

// DownloadPartner returns a stub Excel file of project partners.
func (h *ProjectHandler) DownloadPartner(c *gin.Context) {
	c.Header("Content-Disposition", `attachment; filename="partners.xlsx"`)
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", []byte("PK\x03\x04"))
}

// ImportPartner is a stub for importing project partners from a file.
func (h *ProjectHandler) ImportPartner(c *gin.Context) {
	response.OK(c, nil)
}

// SelectUser returns users matching a keyword for the editor.
func (h *ProjectHandler) SelectUser(c *gin.Context) {
	response.OK(c, []interface{}{})
}

// SelectDept returns departments matching a keyword for the editor.
func (h *ProjectHandler) SelectDept(c *gin.Context) {
	response.OK(c, []interface{}{})
}

// SelectRole returns roles matching a keyword for the editor.
func (h *ProjectHandler) SelectRole(c *gin.Context) {
	response.OK(c, []interface{}{})
}

// SelectPosition returns positions matching a keyword for the editor.
func (h *ProjectHandler) SelectPosition(c *gin.Context) {
	response.OK(c, []interface{}{})
}

// SelectDict returns dictionaries for the editor.
func (h *ProjectHandler) SelectDict(c *gin.Context) {
	response.OK(c, []interface{}{})
}

// SelectTemplate returns templates matching a keyword for the editor.
func (h *ProjectHandler) SelectTemplate(c *gin.Context) {
	response.OK(c, map[string]interface{}{})
}

// SelectRepo returns repos matching a keyword for the editor.
func (h *ProjectHandler) SelectRepo(c *gin.Context) {
	response.OK(c, []interface{}{})
}

// SelectTag returns tags matching a keyword for the editor.
func (h *ProjectHandler) SelectTag(c *gin.Context) {
	response.OK(c, []interface{}{})
}

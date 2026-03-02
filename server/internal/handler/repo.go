package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/pkg/response"
	"github.com/rms-survey/server/internal/service"
)

// RepoHandler handles question repo API endpoints.
type RepoHandler struct {
	svc *service.RepoService
}

func NewRepoHandler(svc *service.RepoService) *RepoHandler {
	return &RepoHandler{svc: svc}
}

func (h *RepoHandler) List(c *gin.Context) {
	var query dto.RepoQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.List(query)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *RepoHandler) Create(c *gin.Context) {
	var req dto.RepoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.Create(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *RepoHandler) Update(c *gin.Context) {
	var req dto.RepoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.Update(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *RepoHandler) Delete(c *gin.Context) {
	var req dto.RepoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.Delete(req.ID); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *RepoHandler) BatchCreate(c *gin.Context) {
	var req dto.RepoTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.BatchCreate(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *RepoHandler) Unbind(c *gin.Context) {
	var req dto.RepoTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.Unbind(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *RepoHandler) Pick(c *gin.Context) {
	var conditions []dto.RandomSurveyCondition
	if err := c.ShouldBindJSON(&conditions); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.PickQuestions(conditions)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *RepoHandler) Import(c *gin.Context) {
	repoID := c.PostForm("repoId")
	if repoID == "" {
		response.Fail(c, response.CodeError, "repoId is required")
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, response.CodeError, "file is required")
		return
	}
	f, err := fileHeader.Open()
	if err != nil {
		response.Fail(c, response.CodeError, "failed to open file")
		return
	}
	defer f.Close()

	if err := h.svc.ImportFromExcel(repoID, f); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *RepoHandler) Export(c *gin.Context) {
	var query dto.RepoQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	data, err := h.svc.Export(query)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="repo_export.xlsx"`))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (h *RepoHandler) ListUserBook(c *gin.Context) {
	var query dto.UserBookQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.ListUserBook(query)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *RepoHandler) CreateUserBook(c *gin.Context) {
	var req dto.UserBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.CreateUserBook(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *RepoHandler) UpdateUserBook(c *gin.Context) {
	var req dto.UserBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	view, err := h.svc.UpdateUserBook(req)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, view)
}

func (h *RepoHandler) DeleteUserBook(c *gin.Context) {
	var req dto.UserBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.DeleteUserBook(req.ID); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}


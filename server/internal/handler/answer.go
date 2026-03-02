package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/pkg/response"
	"github.com/rms-survey/server/internal/repository"
	"github.com/rms-survey/server/internal/service"
	"gorm.io/gorm"
)

// AnswerHandler handles all answer HTTP endpoints.
type AnswerHandler struct {
	svc *service.AnswerService
}

// NewAnswerHandler creates an AnswerHandler wired to the given DB.
func NewAnswerHandler(db *gorm.DB) *AnswerHandler {
	answerRepo := repository.NewAnswerRepo(db)
	projectRepo := repository.NewProjectRepo(db)
	svc := service.NewAnswerService(answerRepo, projectRepo)
	return &AnswerHandler{svc: svc}
}

// RegisterRoutes wires all answer routes onto the provided gin group.
func (h *AnswerHandler) RegisterRoutes(answer gin.IRouter) {
	answer.GET("/list", h.ListAnswers)
	answer.GET("/trash", h.ListDeleted)
	answer.GET("", h.GetAnswer)
	answer.POST("/create", h.CreateAnswer)
	answer.POST("/update", h.UpdateAnswer)
	answer.POST("/delete", h.DeleteAnswer)
	answer.POST("/destroy", h.DestroyAnswer)
	answer.POST("/restore", h.RestoreAnswer)
	answer.GET("/download", h.Download)
	answer.POST("/upload", h.Upload)
}

func (h *AnswerHandler) ListAnswers(c *gin.Context) {
	var query dto.AnswerQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.ListAnswers(&query)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *AnswerHandler) ListDeleted(c *gin.Context) {
	var query dto.AnswerQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.ListDeleted(&query)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *AnswerHandler) GetAnswer(c *gin.Context) {
	var query dto.AnswerQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.GetAnswer(&query)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *AnswerHandler) CreateAnswer(c *gin.Context) {
	var req dto.AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.CreateAnswer(&req, currentUser(c)); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *AnswerHandler) UpdateAnswer(c *gin.Context) {
	var req dto.AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.UpdateAnswer(&req, currentUser(c)); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *AnswerHandler) DeleteAnswer(c *gin.Context) {
	var req dto.AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.DeleteAnswer(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *AnswerHandler) DestroyAnswer(c *gin.Context) {
	var req dto.AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.DestroyAnswer(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *AnswerHandler) RestoreAnswer(c *gin.Context) {
	var req dto.AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.RestoreAnswer(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *AnswerHandler) Download(c *gin.Context) {
	var query dto.DownloadQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	data, filename, err := h.svc.ExportAnswers(&query)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

// Upload handles POST /api/answer/upload — imports answers from an Excel file.
func (h *AnswerHandler) Upload(c *gin.Context) {
	projectID := c.PostForm("projectId")
	autoSchema := c.PostForm("autoSchema") == "true"
	parentID := c.PostForm("parentId")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, response.CodeError, "file is required")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Fail(c, response.CodeError, "failed to open file")
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		response.Fail(c, response.CodeError, "failed to read file")
		return
	}

	result, err := h.svc.UploadAnswers(projectID, autoSchema, parentID, fileData, fileHeader.Filename, currentUser(c))
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}

	response.OK(c, result)
}

package handler

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/pkg/response"
	"github.com/rms-survey/server/internal/repository"
	"github.com/rms-survey/server/internal/service"
	"gorm.io/gorm"
)

// SurveyHandler handles public survey and survey management endpoints.
type SurveyHandler struct {
	svc     *service.SurveyService
	fileSvc *service.FileService
}

// NewSurveyHandler creates a SurveyHandler wired to the given DB.
func NewSurveyHandler(db *gorm.DB) *SurveyHandler {
	projectRepo := repository.NewProjectRepo(db)
	answerRepo := repository.NewAnswerRepo(db)
	svc := service.NewSurveyService(projectRepo, answerRepo, db)
	return &SurveyHandler{svc: svc}
}

// NewSurveyHandlerWithFile creates a SurveyHandler with file service for preview.
func NewSurveyHandlerWithFile(db *gorm.DB, fileSvc *service.FileService) *SurveyHandler {
	projectRepo := repository.NewProjectRepo(db)
	answerRepo := repository.NewAnswerRepo(db)
	svc := service.NewSurveyService(projectRepo, answerRepo, db)
	return &SurveyHandler{svc: svc, fileSvc: fileSvc}
}

// RegisterRoutes wires survey routes onto the provided gin groups.
// public -> /api/public
// survey -> /api/survey
func (h *SurveyHandler) RegisterRoutes(public, survey gin.IRouter) {
	public.POST("/loadProject", h.LoadProject)
	public.POST("/validateProject", h.ValidateProject)
	public.POST("/statistics", h.Statistics)
	public.POST("/saveAnswer", h.SaveAnswer)
	public.POST("/tempSaveAnswer", h.TempSaveAnswer)
	public.POST("/upload", h.Upload)
	public.POST("/uploadAttachment", h.UploadAttachment)
	public.GET("/preview/:attachmentId", h.Preview)
	public.POST("/loadQuery", h.LoadQuery)
	public.POST("/getQueryResult", h.GetQueryResult)
	public.POST("/loadDict", h.LoadDict)
	public.POST("/loadExamResult", h.LoadExamResult)
	public.POST("/loadLinkResult", h.LoadLinkResult)

	survey.GET("/setting", h.GetSetting)
	survey.POST("/setting", h.UpdateSetting)
	survey.GET("/logic", h.GetLogic)
	survey.POST("/logic", h.UpdateLogic)
}

func (h *SurveyHandler) LoadProject(c *gin.Context) {
	var req dto.SurveyLoadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.LoadProject(&req)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *SurveyHandler) ValidateProject(c *gin.Context) {
	var req dto.SurveyLoadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.ValidateProject(&req)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *SurveyHandler) Statistics(c *gin.Context) {
	var req dto.SurveyLoadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.StatProject(&req)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *SurveyHandler) SaveAnswer(c *gin.Context) {
	var req dto.AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.SaveAnswer(&req)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *SurveyHandler) TempSaveAnswer(c *gin.Context) {
	var req dto.AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.TempSaveAnswer(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *SurveyHandler) Upload(c *gin.Context) {
	if h.fileSvc == nil {
		response.Fail(c, response.CodeError, "file service not available")
		return
	}
	fh, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, response.CodeError, "missing file: "+err.Error())
		return
	}
	f, err := fh.Open()
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	defer f.Close()
	result, err := h.fileSvc.Upload(fh.Filename, f)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

// UploadAttachment handles question-specific attachment upload with validation.
func (h *SurveyHandler) UploadAttachment(c *gin.Context) {
	if h.fileSvc == nil {
		response.Fail(c, response.CodeError, "file service not available")
		return
	}

	projectID := c.PostForm("projectId")
	questionID := c.PostForm("questionId")
	if projectID == "" || questionID == "" {
		response.Fail(c, response.CodeError, "projectId and questionId required")
		return
	}

	fh, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, response.CodeError, "missing file: "+err.Error())
		return
	}

	// Get project to find question attachment config
	project, err := h.svc.GetProjectByID(projectID)
	if err != nil {
		response.Fail(c, response.CodeError, "project not found")
		return
	}

	// Parse survey schema to find question config
	cfg := h.svc.GetQuestionAttachmentConfig(project, questionID)

	// Validate and upload
	f, err := fh.Open()
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	defer f.Close()

	view, err := h.fileSvc.UploadWithValidation(fh.Filename, fh.Size, f, cfg)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}

	response.OK(c, &dto.AttachmentInfo{
		FileID:   view.ID,
		FileName: view.OriginalName,
		FileSize: fh.Size,
		FileType: filepath.Ext(view.OriginalName),
	})
}

func (h *SurveyHandler) Preview(c *gin.Context) {
	if h.fileSvc == nil {
		response.Fail(c, response.CodeError, "file service not available")
		return
	}
	attachmentID := c.Param("attachmentId")
	data, contentType, err := h.fileSvc.LoadFileBytes(attachmentID)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Header("Cache-Control", "max-age=2592000")
	c.Data(200, contentType, data)
}

func (h *SurveyHandler) LoadQuery(c *gin.Context) {
	var req dto.PublicQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.LoadQuery(&req)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *SurveyHandler) GetQueryResult(c *gin.Context) {
	var req dto.PublicQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.GetQueryResult(&req)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *SurveyHandler) LoadDict(c *gin.Context) {
	var req dto.PublicDictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.LoadDict(&req)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *SurveyHandler) LoadExamResult(c *gin.Context) {
	var req dto.PublicExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.LoadExamResult(&req)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *SurveyHandler) LoadLinkResult(c *gin.Context) {
	var req dto.PublicLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.LoadLinkResult(&req)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *SurveyHandler) GetSetting(c *gin.Context) {
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

func (h *SurveyHandler) UpdateSetting(c *gin.Context) {
	var req dto.SurveySettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.UpdateSetting(&req, currentUser(c)); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *SurveyHandler) GetLogic(c *gin.Context) {
	projectID := c.Query("projectId")
	if projectID == "" {
		response.Fail(c, response.CodeError, "projectId required")
		return
	}
	result, err := h.svc.GetLogic(projectID)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *SurveyHandler) UpdateLogic(c *gin.Context) {
	var req dto.SurveyLogicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.UpdateLogic(&req, currentUser(c)); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

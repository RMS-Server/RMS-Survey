package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/pkg/response"
	"github.com/surveyking/server/internal/repository"
	"github.com/surveyking/server/internal/service"
	"gorm.io/gorm"
)

// SurveyHandler handles public survey and survey management endpoints.
type SurveyHandler struct {
	svc *service.SurveyService
}

// NewSurveyHandler creates a SurveyHandler wired to the given DB.
func NewSurveyHandler(db *gorm.DB) *SurveyHandler {
	projectRepo := repository.NewProjectRepo(db)
	answerRepo := repository.NewAnswerRepo(db)
	svc := service.NewSurveyService(projectRepo, answerRepo, db)
	return &SurveyHandler{svc: svc}
}

// RegisterRoutes wires survey routes onto the provided gin groups.
// public -> /api/public
// survey -> /api/survey
func (h *SurveyHandler) RegisterRoutes(public, survey gin.IRouter) {
	public.POST("/loadProject", h.LoadProject)
	public.POST("/saveAnswer", h.SaveAnswer)
	public.POST("/tempSaveAnswer", h.TempSaveAnswer)

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

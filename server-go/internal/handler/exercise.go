package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/pkg/response"
	"github.com/surveyking/server/internal/service"
)

// ExerciseHandler handles exercise history API endpoints.
type ExerciseHandler struct {
	svc *service.ExerciseService
}

func NewExerciseHandler(svc *service.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{svc: svc}
}

func (h *ExerciseHandler) List(c *gin.Context) {
	var query dto.HistoryExerciseQuery
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

// GetDetail returns the full detail view for a single exercise record.
func (h *ExerciseHandler) GetDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.Fail(c, response.CodeError, "id is required")
		return
	}
	result, err := h.svc.GetDetail(id)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

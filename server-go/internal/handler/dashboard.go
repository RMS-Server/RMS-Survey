package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/pkg/response"
	"github.com/surveyking/server/internal/service"
)

// DashboardHandler handles dashboard API endpoints.
type DashboardHandler struct {
	svc *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) List(c *gin.Context) {
	var query dto.DashboardQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	views, err := h.svc.List(query)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, views)
}

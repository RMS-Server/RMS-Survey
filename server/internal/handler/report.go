package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rms-survey/server/internal/pkg/response"
	"github.com/rms-survey/server/internal/service"
)

// ReportHandler handles report API endpoints.
type ReportHandler struct {
	svc *service.ReportService
}

func NewReportHandler(svc *service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

func (h *ReportHandler) GetData(c *gin.Context) {
	shortID := c.Param("shortId")
	data, err := h.svc.GetData(shortID)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, data)
}

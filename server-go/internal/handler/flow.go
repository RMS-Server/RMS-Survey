package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/pkg/response"
	"github.com/surveyking/server/internal/service"
)

// FlowHandler handles workflow API endpoints.
type FlowHandler struct {
	svc *service.FlowService
}

func NewFlowHandler(svc *service.FlowService) *FlowHandler {
	return &FlowHandler{svc: svc}
}

func (h *FlowHandler) GetFlow(c *gin.Context) {
	projectID := c.Query("projectId")
	view, err := h.svc.GetFlowEntry(projectID)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, view)
}

func (h *FlowHandler) SaveFlow(c *gin.Context) {
	var req dto.FlowEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *FlowHandler) Deploy(c *gin.Context) {
	response.OK(c, nil)
}

func (h *FlowHandler) GetAuditRecord(c *gin.Context) {
	processInstanceID := c.Query("processInstanceId")
	records, err := h.svc.GetAuditRecord(processInstanceID)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, records)
}

func (h *FlowHandler) GetFlowTasks(c *gin.Context) {
	var query dto.FlowTaskQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.GetFlowTasks(query)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *FlowHandler) GetRevertNodes(c *gin.Context) {
	processInstanceID := c.Query("processInstanceId")
	nodes, err := h.svc.GetRevertNodes(processInstanceID)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nodes)
}

func (h *FlowHandler) ApprovalTask(c *gin.Context) {
	var req dto.ApprovalTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.ApprovalTask(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *FlowHandler) Statics(c *gin.Context) {
	view, err := h.svc.Statics()
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, view)
}

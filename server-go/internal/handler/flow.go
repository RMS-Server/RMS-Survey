package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/pkg/response"
	"github.com/rms-survey/server/internal/service"
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
	if err := h.svc.SaveFlow(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *FlowHandler) Deploy(c *gin.Context) {
	projectID := c.Query("projectId")
	if projectID == "" {
		// Also accept from JSON body.
		var body struct {
			ProjectID string `json:"projectId"`
		}
		if err := c.ShouldBindJSON(&body); err == nil {
			projectID = body.ProjectID
		}
	}
	if err := h.svc.Deploy(projectID); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
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
	// Populate operator info from authenticated user context.
	if u, ok := c.Get("currentUser"); ok {
		if userInfo, ok := u.(*dto.UserInfo); ok {
			if req.OperatorID == "" {
				req.OperatorID = userInfo.UserID
			}
			if req.OperatorName == "" {
				req.OperatorName = userInfo.Username
			}
		}
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

// GetTaskInfo returns a single answer/task by taskId query param.
func (h *FlowHandler) GetTaskInfo(c *gin.Context) {
	taskID := c.Query("taskId")
	answer, err := h.svc.GetAnswer(taskID)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, answer)
}

// LoadSchema returns the flow entry schema for a given project/process/task context.
func (h *FlowHandler) LoadSchema(c *gin.Context) {
	var query dto.SchemaQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	view, err := h.svc.GetFlowEntry(query.ProjectID)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, view)
}

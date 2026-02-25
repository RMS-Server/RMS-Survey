package dto

// FlowEntryRequest is used to save a flow definition.
type FlowEntryRequest struct {
	ProjectID  string      `json:"projectId"`
	FlowData   interface{} `json:"flowData"`
	NodeConfig interface{} `json:"nodeConfig"`
}

// FlowEntryView is the response for a flow definition.
type FlowEntryView struct {
	ID         string      `json:"id"`
	ProjectID  string      `json:"projectId"`
	FlowData   interface{} `json:"flowData"`
	NodeConfig interface{} `json:"nodeConfig"`
	Status     string      `json:"status"`
}

// FlowTaskQuery is the query for listing flow tasks.
type FlowTaskQuery struct {
	PageRequest
	Status string `form:"status"`
}

// FlowTaskView represents a single flow task item.
type FlowTaskView struct {
	ID                string `json:"id"`
	ProcessInstanceID string `json:"processInstanceId"`
	ProjectID         string `json:"projectId"`
	ProjectName       string `json:"projectName"`
	Status            string `json:"status"`
	CreateBy          string `json:"createBy"`
	CreatedAt         string `json:"createdAt"`
}

// FlowOperationView represents an audit record entry.
type FlowOperationView struct {
	ID                string `json:"id"`
	ProcessInstanceID string `json:"processInstanceId"`
	OperatorID        string `json:"operatorId"`
	OperatorName      string `json:"operatorName"`
	Action            string `json:"action"`
	Comment           string `json:"comment"`
	CreatedAt         string `json:"createdAt"`
}

// ApprovalTaskRequest is used to approve or reject a task.
type ApprovalTaskRequest struct {
	TaskID       string `json:"taskId"`
	Action       string `json:"action"` // approve | reject | cancel
	Comment      string `json:"comment"`
	RevertTo     string `json:"revertTo"`
	OperatorID   string `json:"operatorId"`
	OperatorName string `json:"operatorName"`
}

// RevokeView represents a node that can be reverted to.
type RevokeView struct {
	NodeID   string `json:"nodeId"`
	NodeName string `json:"nodeName"`
}

// SchemaQuery is used to load a survey schema by permission.
type SchemaQuery struct {
	ProjectID         string `form:"projectId"`
	ProcessInstanceID string `form:"processInstanceId"`
	TaskID            string `form:"taskId"`
}

// FlowStaticsView holds flow statistics.
type FlowStaticsView struct {
	PendingCount  int64 `json:"pendingCount"`
	RunningCount  int64 `json:"runningCount"`
	ApprovedCount int64 `json:"approvedCount"`
	RejectedCount int64 `json:"rejectedCount"`
}

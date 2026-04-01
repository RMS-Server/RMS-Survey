package dto

import "encoding/json"

// AnswerQuery holds query parameters for listing answers.
type AnswerQuery struct {
	PageRequest
	ProjectID string `json:"projectId" form:"projectId"`
	ID        string `json:"id" form:"id"`
	TempSave  *int   `json:"tempSave" form:"tempSave"`
	Deleted   bool   `json:"deleted" form:"deleted"`
}

// AnswerRequest is used for create/update/delete answer operations.
type AnswerRequest struct {
	ID        string          `json:"id"`
	ProjectID string          `json:"projectId"`
	Answer    json.RawMessage `json:"answer"`
	MetaInfo  json.RawMessage `json:"metaInfo"`
	TempSave  *int            `json:"tempSave"`
	// IDs for batch operations
	IDs []string `json:"ids"`
	// IsRead for marking answer as read/unread
	IsRead *bool `json:"isRead"`
}

// AnswerView is the response DTO for an answer.
type AnswerView struct {
	ID         string          `json:"id"`
	ProjectID  string          `json:"projectId"`
	Answer     json.RawMessage `json:"answer"`
	MetaInfo   json.RawMessage `json:"metaInfo"`
	TempSave   *int            `json:"tempSave"`
	ExamScore  *float32        `json:"examScore"`
	CreateBy   string          `json:"createBy"`
	CreatedAt  string          `json:"createAt"`
	UpdatedAt  string          `json:"updateAt"`
	IsRead     bool            `json:"isRead"`
	ReadAt     string          `json:"readAt,omitempty"`
	ReadBy     string          `json:"readBy,omitempty"`
	ReadByName string          `json:"readByName,omitempty"`
	IPAddress  string          `json:"ipAddress"`
}

// DownloadQuery holds parameters for answer export.
type DownloadQuery struct {
	ProjectID string `json:"projectId" form:"projectId"`
	Locale    string `json:"locale" form:"locale"`
}

// AnswerImportResult is the response for answer Excel import.
type AnswerImportResult struct {
	ProjectID string          `json:"projectId,omitempty"`
	Schema    json.RawMessage `json:"schema,omitempty"`
}

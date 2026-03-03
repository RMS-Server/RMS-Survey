package dto

import "github.com/rms-survey/server/internal/model"

// SystemRoleView is the response DTO for a role (full model).
type SystemRoleView struct {
	model.Role
}

// RoleRequest is used for create/update/delete role operations.
type RoleRequest struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Code      string   `json:"code"`
	Remark    string   `json:"remark"`
	Authority string   `json:"authority"`
	Status    *bool    `json:"status"`
	UserIds   []string `json:"userIds"`
}

// RoleQuery holds filter params for listing roles.
type RoleQuery struct {
	PageRequest
	Name string `json:"name" form:"name"`
}

// SysInfoView is the response DTO for system info.
type SysInfoView struct {
	model.SysInfo
	PublicKey string `json:"publicKey,omitempty"`
}

// SysInfoRequest is used for updating system info.
type SysInfoRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Avatar       string `json:"avatar"`
	Locale       string `json:"locale"`
	Setting      string `json:"setting"`
	AISetting    string `json:"aiSetting"`
	RegisterInfo string `json:"registerInfo"`
}

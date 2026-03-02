package dto

import "encoding/json"

// ProjectQuery holds query parameters for listing projects.
type ProjectQuery struct {
	PageRequest
	Name     string `json:"name" form:"name"`
	Status   *int   `json:"status" form:"status"`
	ParentID string `json:"parentId" form:"parentId"`
	Deleted  bool   `json:"deleted" form:"deleted"`
}

// ProjectRequest is used for create/update/delete project operations.
type ProjectRequest struct {
	ID       string          `json:"id"`
	ParentID string          `json:"parentId"`
	Name     string          `json:"name"`
	Survey   json.RawMessage `json:"survey"`
	Setting  json.RawMessage `json:"setting"`
	Status   *int            `json:"status"`
	Mode     string          `json:"mode"`
	Priority *int            `json:"priority"`
	// IDs for batch operations
	IDs []string `json:"ids"`
}

// ProjectView is the response DTO for a project.
type ProjectView struct {
	ID        string          `json:"id"`
	ParentID  string          `json:"parentId"`
	Name      string          `json:"name"`
	Survey    json.RawMessage `json:"survey"`
	Setting   json.RawMessage `json:"setting"`
	Status    int             `json:"status"`
	Mode      string          `json:"mode"`
	Priority  int             `json:"priority"`
	CreateBy  string          `json:"createBy"`
	CreatedAt string          `json:"createAt"`
	UpdatedAt string          `json:"updateAt"`
}

// ProjectPartnerQuery holds query parameters for listing project partners.
type ProjectPartnerQuery struct {
	PageRequest
	ProjectID string `json:"projectId" form:"projectId"`
}

// ProjectPartnerRequest is used for add/remove project partner.
type ProjectPartnerRequest struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	GroupID   string `json:"groupId"`
	Type      *int   `json:"type"`
}

// ProjectPartnerView is the response DTO for a project partner.
type ProjectPartnerView struct {
	ID        string         `json:"id"`
	ProjectID string         `json:"projectId"`
	UserID    string         `json:"userId"`
	UserName  string         `json:"userName"`
	GroupID   string         `json:"groupId"`
	Type      *int           `json:"type"`
	Status    int            `json:"status"`
	User      *SimpleUserView `json:"user,omitempty"`
}

// SelectUserRequest is used by the editor to search for users.
type SelectUserRequest struct {
	Keyword string `json:"keyword"`
}

// SelectDeptRequest is used by the editor to search for departments.
type SelectDeptRequest struct {
	Keyword string `json:"keyword"`
}

// SelectRoleRequest is used by the editor to search for roles.
type SelectRoleRequest struct {
	Keyword string `json:"keyword"`
}

// SelectPositionRequest is used by the editor to search for positions.
type SelectPositionRequest struct {
	Keyword string `json:"keyword"`
}

// SelectTemplateRequest is used by the editor to search for templates.
type SelectTemplateRequest struct {
	Keyword string `json:"keyword"`
}

// SelectRepoRequest is used by the editor to search for repos.
type SelectRepoRequest struct {
	Keyword string `json:"keyword"`
}

// SelectTagRequest is used by the editor to search for tags.
type SelectTagRequest struct {
	Keyword string `json:"keyword"`
}

// WhiteListRequest is used for importing project partners.
type WhiteListRequest struct {
	ProjectID string `json:"projectId" form:"projectId"`
	FileID    string `json:"fileId" form:"fileId"`
}

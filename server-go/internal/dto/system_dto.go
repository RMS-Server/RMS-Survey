package dto

import "github.com/surveyking/server/internal/model"

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

// PermissionView represents a single permission entry.
type PermissionView struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Module string `json:"module"`
}

// DeptView is the response DTO for a department.
type DeptView struct {
	model.Dept
	Children []*DeptView `json:"children,omitempty"`
}

// DeptRequest is used for create/update/delete dept operations.
type DeptRequest struct {
	ID        string `json:"id"`
	ParentID  string `json:"parentId"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Code      string `json:"code"`
	ManagerID string `json:"managerId"`
	SortCode  *int   `json:"sortCode"`
	Status    string `json:"status"`
	Remark    string `json:"remark"`
}

// DeptSortRequest holds sort order updates for a batch of departments.
type DeptSortRequest struct {
	ID       string `json:"id"`
	SortCode int    `json:"sortCode"`
}

// CommDictView is the response DTO for a data dictionary.
type CommDictView struct {
	model.CommDict
}

// CommDictQuery holds filter params for listing dicts.
type CommDictQuery struct {
	PageRequest
	Name string `json:"name" form:"name"`
	Code string `json:"code" form:"code"`
}

// CommDictRequest is used for create/update/delete dict operations.
type CommDictRequest struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Remark   string `json:"remark"`
	DictType *int   `json:"dictType"`
}

// CommDictItemView is the response DTO for a dict item.
type CommDictItemView struct {
	model.CommDictItem
}

// CommDictItemQuery holds filter params for listing dict items.
type CommDictItemQuery struct {
	PageRequest
	DictCode string `json:"dictCode" form:"dictCode"`
}

// CommDictItemRequest is used for create/update/delete dict item operations.
type CommDictItemRequest struct {
	ID              string `json:"id"`
	DictCode        string `json:"dictCode"`
	ItemName        string `json:"itemName"`
	ItemValue       string `json:"itemValue"`
	ItemOrder       *int   `json:"itemOrder"`
	ItemLevel       *int   `json:"itemLevel"`
	ParentItemValue string `json:"parentItemValue"`
}

// SysInfoView is the response DTO for system info.
type SysInfoView struct {
	model.SysInfo
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

// PositionView is the response DTO for a position.
type PositionView struct {
	model.Position
}

// PositionQuery holds filter params for listing positions.
type PositionQuery struct {
	PageRequest
	Name string `json:"name" form:"name"`
}

// PositionRequest is used for create/update/delete position operations.
type PositionRequest struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Code               string `json:"code"`
	IsVirtual          bool   `json:"isVirtual"`
	DataPermissionType string `json:"dataPermissionType"`
}

// TagView is the response DTO for a tag.
type TagView struct {
	model.Tag
}

// SystemTagQuery holds filter params for listing system tags.
type SystemTagQuery struct {
	PageRequest
	Category string `json:"category" form:"category"`
	EntityID string `json:"entityId" form:"entityId"`
}

// TagRequest is used for create/delete tag operations.
type TagRequest struct {
	ID       string `json:"id"`
	EntityID string `json:"entityId"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

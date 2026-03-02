package service

import (
	"github.com/google/uuid"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/model"
	"github.com/rms-survey/server/internal/repository"
)

// SystemService handles business logic for system management.
type SystemService struct {
	repo *repository.SystemRepo
}

func NewSystemService(repo *repository.SystemRepo) *SystemService {
	return &SystemService{repo: repo}
}

func newID() string {
	return uuid.New().String()
}

// --- Role ---

func (s *SystemService) ListRoles(q dto.RoleQuery) (dto.PageResponse[dto.SystemRoleView], error) {
	offset, limit := pageOffset(q.PageIndex, q.PageSize)
	roles, total, err := s.repo.ListRoles(q.Name, offset, limit)
	if err != nil {
		return dto.PageResponse[dto.SystemRoleView]{}, err
	}
	views := make([]dto.SystemRoleView, len(roles))
	for i, r := range roles {
		views[i] = dto.SystemRoleView{Role: r}
	}
	return dto.PageResponse[dto.SystemRoleView]{List: views, Total: total}, nil
}

func (s *SystemService) CreateRole(req dto.RoleRequest) error {
	role := model.Role{
		Name:      req.Name,
		Code:      req.Code,
		Remark:    req.Remark,
		Authority: req.Authority,
		Status:    req.Status,
	}
	role.ID = newID()
	return s.repo.CreateRole(&role)
}

func (s *SystemService) UpdateRole(req dto.RoleRequest) error {
	role, err := s.repo.GetRoleByID(req.ID)
	if err != nil {
		return err
	}
	role.Name = req.Name
	role.Code = req.Code
	role.Remark = req.Remark
	role.Authority = req.Authority
	role.Status = req.Status
	return s.repo.UpdateRole(role)
}

func (s *SystemService) DeleteRole(req dto.RoleRequest) error {
	return s.repo.DeleteRole(req.ID)
}

func (s *SystemService) CountRoles() (int64, error) {
	return s.repo.CountRoles()
}

// --- Dept ---

func (s *SystemService) ListDepts() ([]dto.DeptView, error) {
	depts, err := s.repo.ListDepts()
	if err != nil {
		return nil, err
	}
	return buildDeptTree(depts), nil
}

func (s *SystemService) CreateDept(req dto.DeptRequest) error {
	dept := model.Dept{
		ParentID:  req.ParentID,
		Name:      req.Name,
		ShortName: req.ShortName,
		Code:      req.Code,
		ManagerID: req.ManagerID,
		SortCode:  req.SortCode,
		Status:    req.Status,
		Remark:    req.Remark,
	}
	dept.ID = newID()
	return s.repo.CreateDept(&dept)
}

func (s *SystemService) UpdateDept(req dto.DeptRequest) error {
	dept := model.Dept{
		ParentID:  req.ParentID,
		Name:      req.Name,
		ShortName: req.ShortName,
		Code:      req.Code,
		ManagerID: req.ManagerID,
		SortCode:  req.SortCode,
		Status:    req.Status,
		Remark:    req.Remark,
	}
	dept.ID = req.ID
	return s.repo.UpdateDept(&dept)
}

func (s *SystemService) DeleteDept(id string) error {
	return s.repo.DeleteDept(id)
}

func (s *SystemService) CountDepts() (int64, error) {
	return s.repo.CountDepts()
}

// UpdateDeptSortCodes applies sort_code updates for a batch of departments.
func (s *SystemService) UpdateDeptSortCodes(items []dto.DeptSortRequest) error {
	for _, item := range items {
		if item.ID == "" {
			continue
		}
		if err := s.repo.UpdateDeptSortCode(item.ID, item.SortCode); err != nil {
			return err
		}
	}
	return nil
}

// buildDeptTree converts a flat dept list (pre-sorted by sort_code) into a tree.
// The input slice order is preserved for root nodes and children.
func buildDeptTree(depts []model.Dept) []dto.DeptView {
	nodeMap := make(map[string]*dto.DeptView, len(depts))
	for i := range depts {
		d := depts[i]
		nodeMap[d.ID] = &dto.DeptView{Dept: d}
	}
	var roots []dto.DeptView
	// Iterate in original slice order to preserve sort_code ordering.
	for i := range depts {
		node := nodeMap[depts[i].ID]
		if node.ParentID == "" || node.ParentID == "0" {
			roots = append(roots, *node)
		} else if parent, ok := nodeMap[node.ParentID]; ok {
			parent.Children = append(parent.Children, node)
		} else {
			roots = append(roots, *node)
		}
	}
	return roots
}

// --- Dict ---

func (s *SystemService) ListDicts(q dto.CommDictQuery) (dto.PageResponse[dto.CommDictView], error) {
	offset, limit := pageOffset(q.PageIndex, q.PageSize)
	dicts, total, err := s.repo.ListDicts(q.Name, q.Code, offset, limit)
	if err != nil {
		return dto.PageResponse[dto.CommDictView]{}, err
	}
	views := make([]dto.CommDictView, len(dicts))
	for i, d := range dicts {
		views[i] = dto.CommDictView{CommDict: d}
	}
	return dto.PageResponse[dto.CommDictView]{List: views, Total: total}, nil
}

func (s *SystemService) CreateDict(req dto.CommDictRequest) error {
	d := model.CommDict{
		Code:     req.Code,
		Name:     req.Name,
		Remark:   req.Remark,
		DictType: req.DictType,
	}
	d.ID = newID()
	return s.repo.CreateDict(&d)
}

func (s *SystemService) UpdateDict(req dto.CommDictRequest) error {
	d := model.CommDict{
		Code:     req.Code,
		Name:     req.Name,
		Remark:   req.Remark,
		DictType: req.DictType,
	}
	d.ID = req.ID
	return s.repo.UpdateDict(&d)
}

func (s *SystemService) DeleteDict(id string) error {
	return s.repo.DeleteDict(id)
}

// --- DictItem ---

func (s *SystemService) ListDictItems(q dto.CommDictItemQuery) (dto.PageResponse[dto.CommDictItemView], error) {
	offset, limit := pageOffset(q.PageIndex, q.PageSize)
	items, total, err := s.repo.ListDictItems(q.DictCode, offset, limit)
	if err != nil {
		return dto.PageResponse[dto.CommDictItemView]{}, err
	}
	views := make([]dto.CommDictItemView, len(items))
	for i, item := range items {
		views[i] = dto.CommDictItemView{CommDictItem: item}
	}
	return dto.PageResponse[dto.CommDictItemView]{List: views, Total: total}, nil
}

func (s *SystemService) SaveDictItem(req dto.CommDictItemRequest) error {
	item := model.CommDictItem{
		DictCode:        req.DictCode,
		ItemName:        req.ItemName,
		ItemValue:       req.ItemValue,
		ItemOrder:       req.ItemOrder,
		ItemLevel:       req.ItemLevel,
		ParentItemValue: req.ParentItemValue,
	}
	if req.ID == "" {
		item.ID = newID()
	} else {
		item.ID = req.ID
	}
	return s.repo.SaveDictItem(&item)
}

func (s *SystemService) DeleteDictItem(id string) error {
	return s.repo.DeleteDictItem(id)
}

// BatchImportDictItems inserts multiple dict items atomically.
func (s *SystemService) BatchImportDictItems(items []model.CommDictItem) error {
	return s.repo.BatchCreateDictItems(items)
}

// --- SysInfo ---

func (s *SystemService) GetSysInfo() (*dto.SysInfoView, error) {
	info, err := s.repo.GetSysInfo()
	if err != nil {
		return nil, err
	}
	return &dto.SysInfoView{SysInfo: *info}, nil
}

func (s *SystemService) UpdateSysInfo(req dto.SysInfoRequest) error {
	info, err := s.repo.GetSysInfo()
	if err != nil {
		return err
	}
	if info.ID == "" {
		info.ID = newID()
		t := true
		info.IsDefault = &t
	}
	info.Name = req.Name
	info.Description = req.Description
	info.Avatar = req.Avatar
	info.Locale = req.Locale
	info.Setting = req.Setting
	info.AISetting = req.AISetting
	info.RegisterInfo = req.RegisterInfo
	return s.repo.SaveSysInfo(info)
}

// --- Position ---

func (s *SystemService) ListPositions(q dto.PositionQuery) (dto.PageResponse[dto.PositionView], error) {
	offset, limit := pageOffset(q.PageIndex, q.PageSize)
	positions, total, err := s.repo.ListPositions(q.Name, offset, limit)
	if err != nil {
		return dto.PageResponse[dto.PositionView]{}, err
	}
	views := make([]dto.PositionView, len(positions))
	for i, p := range positions {
		views[i] = dto.PositionView{Position: p}
	}
	return dto.PageResponse[dto.PositionView]{List: views, Total: total}, nil
}

func (s *SystemService) CreatePosition(req dto.PositionRequest) error {
	p := model.Position{
		Name:               req.Name,
		Code:               req.Code,
		IsVirtual:          req.IsVirtual,
		DataPermissionType: req.DataPermissionType,
	}
	p.ID = newID()
	return s.repo.CreatePosition(&p)
}

func (s *SystemService) UpdatePosition(req dto.PositionRequest) error {
	p := model.Position{
		Name:               req.Name,
		Code:               req.Code,
		IsVirtual:          req.IsVirtual,
		DataPermissionType: req.DataPermissionType,
	}
	p.ID = req.ID
	return s.repo.UpdatePosition(&p)
}

func (s *SystemService) DeletePosition(id string) error {
	return s.repo.DeletePosition(id)
}

// --- Tag ---

func (s *SystemService) ListTags(q dto.SystemTagQuery) (dto.PageResponse[dto.TagView], error) {
	offset, limit := pageOffset(q.PageIndex, q.PageSize)
	tags, total, err := s.repo.ListTags(q.Category, q.EntityID, offset, limit)
	if err != nil {
		return dto.PageResponse[dto.TagView]{}, err
	}
	views := make([]dto.TagView, len(tags))
	for i, t := range tags {
		views[i] = dto.TagView{Tag: t}
	}
	return dto.PageResponse[dto.TagView]{List: views, Total: total}, nil
}

func (s *SystemService) CreateTag(req dto.TagRequest) error {
	tag := model.Tag{
		EntityID: req.EntityID,
		Name:     req.Name,
		Category: req.Category,
	}
	tag.ID = newID()
	return s.repo.CreateTag(&tag)
}

func (s *SystemService) DeleteTag(id string) error {
	return s.repo.DeleteTag(id)
}

// pageOffset converts 1-based pageIndex + pageSize to offset/limit.
func pageOffset(pageIndex, pageSize int) (int, int) {
	if pageSize <= 0 {
		return 0, 0
	}
	if pageIndex <= 0 {
		pageIndex = 1
	}
	return (pageIndex - 1) * pageSize, pageSize
}

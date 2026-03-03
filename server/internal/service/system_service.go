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

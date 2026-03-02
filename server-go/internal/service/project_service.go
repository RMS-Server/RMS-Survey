package service

import (
	"encoding/json"
	"fmt"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/model"
	"github.com/rms-survey/server/internal/repository"
)

// ProjectService handles business logic for projects.
type ProjectService struct {
	repo *repository.ProjectRepo
}

// NewProjectService creates a new ProjectService.
func NewProjectService(repo *repository.ProjectRepo) *ProjectService {
	return &ProjectService{repo: repo}
}

// ListProjects returns a paginated list of projects.
func (s *ProjectService) ListProjects(query *dto.ProjectQuery, userInfo *dto.UserInfo) (*dto.PageResponse[dto.ProjectView], error) {
	projects, total, err := s.repo.ListProjects(query, userInfo)
	if err != nil {
		return nil, err
	}
	views := make([]dto.ProjectView, 0, len(projects))
	for _, p := range projects {
		views = append(views, toProjectView(p))
	}
	return &dto.PageResponse[dto.ProjectView]{List: views, Total: total}, nil
}

// GetProject returns a single project by ID.
func (s *ProjectService) GetProject(id string) (*dto.ProjectView, error) {
	p, err := s.repo.GetProject(id)
	if err != nil {
		return nil, err
	}
	v := toProjectView(*p)
	return &v, nil
}

// GetSetting returns the setting JSON for a project.
func (s *ProjectService) GetSetting(projectID string) (json.RawMessage, error) {
	p, err := s.repo.GetProject(projectID)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(p.Setting), nil
}

// CreateProject inserts a new project.
func (s *ProjectService) CreateProject(req *dto.ProjectRequest, userInfo *dto.UserInfo) (*dto.ProjectView, error) {
	id, err := nanoid.New()
	if err != nil {
		return nil, fmt.Errorf("generate id: %w", err)
	}
	priority := 1000
	if req.Priority != nil {
		priority = *req.Priority
	}
	status := 0
	if req.Status != nil {
		status = *req.Status
	}
	p := &model.Project{
		BaseModel: model.BaseModel{
			ID:       id,
			CreateBy: userInfo.UserID,
		},
		ParentID: req.ParentID,
		Name:     req.Name,
		Survey:   string(req.Survey),
		Setting:  string(req.Setting),
		Status:   status,
		Mode:     req.Mode,
		Priority: priority,
	}
	if err := s.repo.CreateProject(p); err != nil {
		return nil, err
	}
	v := toProjectView(*p)
	return &v, nil
}

// UpdateProject updates an existing project.
func (s *ProjectService) UpdateProject(req *dto.ProjectRequest, userInfo *dto.UserInfo) error {
	p, err := s.repo.GetProject(req.ID)
	if err != nil {
		return err
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Survey != nil {
		p.Survey = string(req.Survey)
	}
	if req.Setting != nil {
		p.Setting = string(req.Setting)
	}
	if req.Status != nil {
		p.Status = *req.Status
	}
	if req.Mode != "" {
		p.Mode = req.Mode
	}
	if req.Priority != nil {
		p.Priority = *req.Priority
	}
	p.UpdateBy = userInfo.UserID
	return s.repo.UpdateProject(p)
}

// DeleteProject soft-deletes a project.
func (s *ProjectService) DeleteProject(req *dto.ProjectRequest) error {
	return s.repo.SoftDeleteProject(req.ID)
}

// GetDeleted returns soft-deleted projects.
func (s *ProjectService) GetDeleted(userInfo *dto.UserInfo) ([]dto.ProjectView, error) {
	projects, err := s.repo.ListDeleted(userInfo)
	if err != nil {
		return nil, err
	}
	views := make([]dto.ProjectView, 0, len(projects))
	for _, p := range projects {
		views = append(views, toProjectView(p))
	}
	return views, nil
}

// DestroyProject permanently deletes projects.
func (s *ProjectService) DestroyProject(req *dto.ProjectRequest) error {
	ids := req.IDs
	if len(ids) == 0 && req.ID != "" {
		ids = []string{req.ID}
	}
	for _, id := range ids {
		if err := s.repo.HardDeleteProject(id); err != nil {
			return err
		}
	}
	return nil
}

// RestoreProject recovers a project from trash.
func (s *ProjectService) RestoreProject(req *dto.ProjectRequest) error {
	ids := req.IDs
	if len(ids) == 0 && req.ID != "" {
		ids = []string{req.ID}
	}
	for _, id := range ids {
		if err := s.repo.RestoreProject(id); err != nil {
			return err
		}
	}
	return nil
}

// ListPartners returns partners for a project.
func (s *ProjectService) ListPartners(query *dto.ProjectPartnerQuery) (*dto.PageResponse[dto.ProjectPartnerView], error) {
	partners, total, err := s.repo.ListPartners(query)
	if err != nil {
		return nil, err
	}
	views := make([]dto.ProjectPartnerView, 0, len(partners))
	for _, p := range partners {
		views = append(views, dto.ProjectPartnerView{
			ID:        p.ID,
			ProjectID: p.ProjectID,
			UserID:    p.UserID,
			UserName:  p.UserName,
			GroupID:   p.GroupID,
			Type:      p.Type,
			Status:    p.Status,
		})
	}
	return &dto.PageResponse[dto.ProjectPartnerView]{List: views, Total: total}, nil
}

// AddPartner adds a project partner.
func (s *ProjectService) AddPartner(req *dto.ProjectPartnerRequest, userInfo *dto.UserInfo) error {
	id, err := nanoid.New()
	if err != nil {
		return fmt.Errorf("generate id: %w", err)
	}
	p := &model.ProjectPartner{
		BaseModelNoSoftDelete: model.BaseModelNoSoftDelete{
			ID:       id,
			CreateBy: userInfo.UserID,
		},
		ProjectID: req.ProjectID,
		UserID:    req.UserID,
		UserName:  req.UserName,
		GroupID:   req.GroupID,
		Type:      req.Type,
	}
	return s.repo.CreatePartner(p)
}

// RemovePartner removes a project partner.
func (s *ProjectService) RemovePartner(req *dto.ProjectPartnerRequest) error {
	return s.repo.DeletePartner(req.ID)
}

// ListPartnersAll returns all partners for a project without pagination.
func (s *ProjectService) ListPartnersAll(projectID string) ([]dto.ProjectPartnerView, error) {
	partners, err := s.repo.ListPartnersAll(projectID)
	if err != nil {
		return nil, err
	}
	views := make([]dto.ProjectPartnerView, 0, len(partners))
	for _, p := range partners {
		views = append(views, dto.ProjectPartnerView{
			ID:        p.ID,
			ProjectID: p.ProjectID,
			UserID:    p.UserID,
			UserName:  p.UserName,
			GroupID:   p.GroupID,
			Type:      p.Type,
			Status:    p.Status,
		})
	}
	return views, nil
}

func toProjectView(p model.Project) dto.ProjectView {
	return dto.ProjectView{
		ID:        p.ID,
		ParentID:  p.ParentID,
		Name:      p.Name,
		Survey:    json.RawMessage(p.Survey),
		Setting:   json.RawMessage(p.Setting),
		Status:    p.Status,
		Mode:      p.Mode,
		Priority:  p.Priority,
		CreateBy:  p.CreateBy,
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
	}
}

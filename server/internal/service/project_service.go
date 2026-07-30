package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/model"
	"github.com/rms-survey/server/internal/repository"
)

// ErrInvalidSurveyPayload is returned when a survey update body is present
// but semantically empty (null / empty / {}). Rejecting these loudly keeps
// a misbehaving client from blanking a saved survey.
var ErrInvalidSurveyPayload = errors.New("survey payload is empty or null")

// isAdmin checks if the user has admin role.
func isAdmin(userInfo *dto.UserInfo) bool {
	if userInfo == nil {
		return false
	}
	for _, r := range userInfo.Roles {
		if r == "admin" || r == "ADMIN" {
			return true
		}
	}
	return false
}

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
		views = append(views, toProjectView(p, userInfo))
	}
	return &dto.PageResponse[dto.ProjectView]{List: views, Total: total}, nil
}

// GetProject returns a single project by ID.
func (s *ProjectService) GetProject(id string, userInfo *dto.UserInfo) (*dto.ProjectView, error) {
	p, err := s.repo.HasProjectAccess(id, userInfo)
	if err != nil {
		return nil, err
	}
	v := toProjectView(*p, userInfo)
	return &v, nil
}

// GetSetting returns the setting JSON for a project.
func (s *ProjectService) GetSetting(projectID string, userInfo *dto.UserInfo) (json.RawMessage, error) {
	p, err := s.repo.HasProjectAccess(projectID, userInfo)
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
	setting := req.Setting
	if len(bytes.TrimSpace(setting)) == 0 {
		setting = json.RawMessage(`{}`)
	}
	p := &model.Project{
		BaseModel: model.BaseModel{
			ID:       id,
			CreateBy: userInfo.UserID,
		},
		ParentID: req.ParentID,
		Name:     req.Name,
		Survey:   string(req.Survey),
		Setting:  string(setting),
		Status:   status,
		Mode:     req.Mode,
		Priority: priority,
	}
	if err := s.repo.CreateProject(p); err != nil {
		return nil, err
	}
	v := toProjectView(*p, userInfo)
	return &v, nil
}

// UpdateProjectMeta updates partial meta fields on an existing project.
// Only pointer-non-nil fields are written; everything else is preserved.
// It physically cannot touch survey or setting — that is the whole point.
func (s *ProjectService) UpdateProjectMeta(req *dto.ProjectMetaUpdateRequest, userInfo *dto.UserInfo) error {
	p, err := s.repo.HasProjectAccess(req.ID, userInfo)
	if err != nil {
		return err
	}
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.ParentID != nil {
		p.ParentID = *req.ParentID
	}
	if req.Status != nil {
		p.Status = *req.Status
	}
	if req.Mode != nil {
		p.Mode = *req.Mode
	}
	if req.Priority != nil {
		p.Priority = *req.Priority
	}
	p.UpdateBy = userInfo.UserID
	return s.repo.UpdateProject(p)
}

// UpdateProjectSurvey replaces the survey schema JSON on an existing project.
// Does not touch any other column. Rejects null / empty / {} bodies — gin's
// binding:"required" alone lets those slip through since json.RawMessage
// is just []byte under the hood.
func (s *ProjectService) UpdateProjectSurvey(req *dto.ProjectSurveyUpdateRequest, userInfo *dto.UserInfo) error {
	body := bytes.TrimSpace(req.Survey)
	if len(body) == 0 || bytes.Equal(body, []byte("null")) || bytes.Equal(body, []byte("{}")) {
		return ErrInvalidSurveyPayload
	}
	p, err := s.repo.HasProjectAccess(req.ID, userInfo)
	if err != nil {
		return err
	}
	p.Survey = string(req.Survey)
	p.UpdateBy = userInfo.UserID
	return s.repo.UpdateProject(p)
}

// DeleteProject soft-deletes a project. Only owner can delete.
func (s *ProjectService) DeleteProject(req *dto.ProjectRequest, userInfo *dto.UserInfo) error {
	isOwner, err := s.repo.IsProjectOwner(req.ID, userInfo)
	if err != nil {
		return err
	}
	if !isOwner {
		return ErrAccessDenied
	}
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
		views = append(views, toProjectView(p, userInfo))
	}
	return views, nil
}

// DestroyProject permanently deletes projects. Only owner can destroy.
func (s *ProjectService) DestroyProject(req *dto.ProjectRequest, userInfo *dto.UserInfo) error {
	ids := req.IDs
	if len(ids) == 0 && req.ID != "" {
		ids = []string{req.ID}
	}
	for _, id := range ids {
		isOwner, err := s.repo.IsProjectOwner(id, userInfo)
		if err != nil {
			return err
		}
		if !isOwner {
			return ErrAccessDenied
		}
		if err := s.repo.HardDeleteProject(id); err != nil {
			return err
		}
	}
	return nil
}

// RestoreProject recovers a project from trash. Only owner can restore.
func (s *ProjectService) RestoreProject(req *dto.ProjectRequest, userInfo *dto.UserInfo) error {
	ids := req.IDs
	if len(ids) == 0 && req.ID != "" {
		ids = []string{req.ID}
	}
	for _, id := range ids {
		isOwner, err := s.repo.IsProjectOwner(id, userInfo)
		if err != nil {
			return err
		}
		if !isOwner {
			return ErrAccessDenied
		}
		if err := s.repo.RestoreProject(id); err != nil {
			return err
		}
	}
	return nil
}

// ListPartners returns partners for a project. Requires project access.
func (s *ProjectService) ListPartners(query *dto.ProjectPartnerQuery, userInfo *dto.UserInfo) (*dto.PageResponse[dto.ProjectPartnerView], error) {
	if _, err := s.repo.HasProjectAccess(query.ProjectID, userInfo); err != nil {
		return nil, err
	}
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

// AddPartner adds a project partner. Only owner can add partners.
func (s *ProjectService) AddPartner(req *dto.ProjectPartnerRequest, userInfo *dto.UserInfo) error {
	isOwner, err := s.repo.IsProjectOwner(req.ProjectID, userInfo)
	if err != nil {
		return err
	}
	if !isOwner {
		return ErrAccessDenied
	}
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

// RemovePartner removes a project partner. Only owner can remove partners.
func (s *ProjectService) RemovePartner(req *dto.ProjectPartnerRequest, userInfo *dto.UserInfo) error {
	// First get the partner to find the project ID
	partner, err := s.repo.GetPartnerByID(req.ID)
	if err != nil {
		return err
	}
	isOwner, err := s.repo.IsProjectOwner(partner.ProjectID, userInfo)
	if err != nil {
		return err
	}
	if !isOwner {
		return ErrAccessDenied
	}
	return s.repo.DeletePartner(req.ID)
}

// ListPartnersAll returns all partners for a project without pagination. Requires project access.
func (s *ProjectService) ListPartnersAll(projectID string, userInfo *dto.UserInfo) ([]dto.ProjectPartnerView, error) {
	if _, err := s.repo.HasProjectAccess(projectID, userInfo); err != nil {
		return nil, err
	}
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

func toProjectView(p model.Project, userInfo *dto.UserInfo) dto.ProjectView {
	isOwner := userInfo != nil && (p.CreateBy == userInfo.UserID || isAdmin(userInfo))
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
		IsOwner:   isOwner,
	}
}

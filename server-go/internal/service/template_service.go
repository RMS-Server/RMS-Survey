package service

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/model"
	"github.com/surveyking/server/internal/repository"
)

// TemplateService handles template business logic.
type TemplateService struct {
	repo *repository.TemplateRepository
}

func NewTemplateService(repo *repository.TemplateRepository) *TemplateService {
	return &TemplateService{repo: repo}
}

func (s *TemplateService) List(query dto.TemplateQuery) (*dto.PageResponse[dto.TemplateView], error) {
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	pageIndex := query.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	offset := (pageIndex - 1) * pageSize

	items, total, err := s.repo.List(query.RepoID, query.Name, query.Category, query.Mode, query.QuestionType, offset, pageSize)
	if err != nil {
		return nil, err
	}
	views := make([]dto.TemplateView, 0, len(items))
	for _, t := range items {
		var tdata interface{}
		if t.TemplateData != "" {
			_ = json.Unmarshal([]byte(t.TemplateData), &tdata)
		}
		views = append(views, dto.TemplateView{
			ID:           t.ID,
			RepoID:       t.RepoID,
			SerialNo:     t.SerialNo,
			Name:         t.Name,
			QuestionType: t.QuestionType,
			TemplateData: tdata,
			Mode:         t.Mode,
			Category:     t.Category,
			Tag:          t.Tag,
			Priority:     t.Priority,
			PreviewURL:   t.PreviewURL,
			Shared:       t.Shared,
		})
	}
	return &dto.PageResponse[dto.TemplateView]{List: views, Total: total}, nil
}

func (s *TemplateService) Get(query dto.TemplateQuery) (*dto.TemplateView, error) {
	t, err := s.repo.GetByID(query.ID)
	if err != nil {
		return nil, err
	}
	var tdata interface{}
	if t.TemplateData != "" {
		_ = json.Unmarshal([]byte(t.TemplateData), &tdata)
	}
	return &dto.TemplateView{
		ID:           t.ID,
		RepoID:       t.RepoID,
		SerialNo:     t.SerialNo,
		Name:         t.Name,
		QuestionType: t.QuestionType,
		TemplateData: tdata,
		Mode:         t.Mode,
		Category:     t.Category,
		Tag:          t.Tag,
		Priority:     t.Priority,
		PreviewURL:   t.PreviewURL,
		Shared:       t.Shared,
	}, nil
}

func (s *TemplateService) Create(req dto.TemplateRequest) (string, error) {
	now := time.Now()
	var tdataStr string
	if req.TemplateData != nil {
		b, _ := json.Marshal(req.TemplateData)
		tdataStr = string(b)
	}
	t := &model.Template{
		RepoID:       req.RepoID,
		Name:         req.Name,
		QuestionType: req.QuestionType,
		TemplateData: tdataStr,
		Mode:         req.Mode,
		Category:     req.Category,
		Tag:          req.Tag,
		Priority:     req.Priority,
		PreviewURL:   req.PreviewURL,
		Shared:       req.Shared,
	}
	t.ID = uuid.New().String()
	t.CreatedAt = now
	t.UpdatedAt = now
	return t.ID, s.repo.Create(t)
}

func (s *TemplateService) Update(req dto.TemplateRequest) error {
	t, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if req.TemplateData != nil {
		b, _ := json.Marshal(req.TemplateData)
		t.TemplateData = string(b)
	}
	t.Name = req.Name
	t.QuestionType = req.QuestionType
	t.Mode = req.Mode
	t.Category = req.Category
	t.Tag = req.Tag
	t.Priority = req.Priority
	t.PreviewURL = req.PreviewURL
	t.Shared = req.Shared
	return s.repo.Update(t)
}

func (s *TemplateService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *TemplateService) ListCategories(mode string) ([]string, error) {
	return s.repo.ListCategories(mode)
}

func (s *TemplateService) ListTags(mode, category string) ([]string, error) {
	return s.repo.ListTags(mode, category)
}

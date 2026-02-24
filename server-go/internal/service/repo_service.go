package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/model"
	"github.com/surveyking/server/internal/repository"
)

// RepoService handles question repo business logic.
type RepoService struct {
	repoRepo     *repository.RepoRepository
	userBookRepo *repository.UserBookRepository
}

func NewRepoService(repoRepo *repository.RepoRepository, userBookRepo *repository.UserBookRepository) *RepoService {
	return &RepoService{repoRepo: repoRepo, userBookRepo: userBookRepo}
}

func (s *RepoService) List(query dto.RepoQuery) (*dto.PageResponse[dto.RepoView], error) {
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	pageIndex := query.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	offset := (pageIndex - 1) * pageSize

	items, total, err := s.repoRepo.List(query.Name, query.Category, query.Mode, offset, pageSize)
	if err != nil {
		return nil, err
	}
	views := make([]dto.RepoView, 0, len(items))
	for _, r := range items {
		views = append(views, dto.RepoView{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Category:    r.Category,
			Mode:        r.Mode,
			Shared:      r.Shared,
			Tag:         r.Tag,
			Priority:    r.Priority,
			IsPractice:  r.IsPractice,
		})
	}
	return &dto.PageResponse[dto.RepoView]{List: views, Total: total}, nil
}

func (s *RepoService) Create(req dto.RepoRequest) error {
	now := time.Now()
	r := &model.Repo{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Mode:        req.Mode,
		Shared:      req.Shared,
		Tag:         req.Tag,
		Priority:    req.Priority,
		IsPractice:  req.IsPractice,
	}
	r.ID = uuid.New().String()
	r.CreatedAt = now
	r.UpdatedAt = now
	return s.repoRepo.Create(r)
}

func (s *RepoService) Update(req dto.RepoRequest) error {
	r, err := s.repoRepo.GetByID(req.ID)
	if err != nil {
		return err
	}
	r.Name = req.Name
	r.Description = req.Description
	r.Category = req.Category
	r.Mode = req.Mode
	r.Shared = req.Shared
	r.Tag = req.Tag
	r.Priority = req.Priority
	r.IsPractice = req.IsPractice
	return s.repoRepo.Update(r)
}

func (s *RepoService) Delete(id string) error {
	return s.repoRepo.Delete(id)
}

func (s *RepoService) ListUserBook(query dto.UserBookQuery) (*dto.PageResponse[dto.UserBookView], error) {
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	pageIndex := query.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	offset := (pageIndex - 1) * pageSize

	items, total, err := s.userBookRepo.List(query.RepoID, query.Status, offset, pageSize)
	if err != nil {
		return nil, err
	}
	views := make([]dto.UserBookView, 0, len(items))
	for _, ub := range items {
		views = append(views, dto.UserBookView{
			ID:           ub.ID,
			Name:         ub.Name,
			TemplateID:   ub.TemplateID,
			WrongTimes:   ub.WrongTimes,
			CorrectTimes: ub.CorrectTimes,
			Note:         ub.Note,
			Status:       ub.Status,
			Type:         ub.Type,
			RepoID:       ub.RepoID,
			IsMarked:     ub.IsMarked,
		})
	}
	return &dto.PageResponse[dto.UserBookView]{List: views, Total: total}, nil
}

func (s *RepoService) CreateUserBook(req dto.UserBookRequest) error {
	now := time.Now()
	ub := &model.UserBook{
		TemplateID: req.TemplateID,
		RepoID:     req.RepoID,
		Note:       req.Note,
		Status:     req.Status,
		IsMarked:   req.IsMarked,
	}
	ub.ID = uuid.New().String()
	ub.CreatedAt = now
	ub.UpdatedAt = now
	return s.userBookRepo.Create(ub)
}

func (s *RepoService) UpdateUserBook(req dto.UserBookRequest) (*dto.UserBookView, error) {
	ub, err := s.userBookRepo.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	ub.Note = req.Note
	ub.Status = req.Status
	ub.IsMarked = req.IsMarked
	if err := s.userBookRepo.Update(ub); err != nil {
		return nil, err
	}
	return &dto.UserBookView{
		ID:           ub.ID,
		Name:         ub.Name,
		TemplateID:   ub.TemplateID,
		WrongTimes:   ub.WrongTimes,
		CorrectTimes: ub.CorrectTimes,
		Note:         ub.Note,
		Status:       ub.Status,
		Type:         ub.Type,
		RepoID:       ub.RepoID,
		IsMarked:     ub.IsMarked,
	}, nil
}

func (s *RepoService) DeleteUserBook(id string) error {
	return s.userBookRepo.Delete(id)
}

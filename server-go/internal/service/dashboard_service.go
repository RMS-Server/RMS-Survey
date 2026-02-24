package service

import (
	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/repository"
)

// DashboardService handles dashboard business logic.
type DashboardService struct {
	repo *repository.DashboardRepository
}

func NewDashboardService(repo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) List(query dto.DashboardQuery) ([]dto.DashboardView, error) {
	items, err := s.repo.List(query.ProjectID, query.Key)
	if err != nil {
		return nil, err
	}
	views := make([]dto.DashboardView, 0, len(items))
	for _, d := range items {
		views = append(views, dto.DashboardView{
			ID:        d.ID,
			Key:       d.Key,
			Type:      d.Type,
			ProjectID: d.ProjectID,
			Setting:   d.Setting,
		})
	}
	return views, nil
}

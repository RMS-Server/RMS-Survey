package service

import (
	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/repository"
)

// ReportService handles report data business logic.
type ReportService struct {
	repo *repository.ReportRepository
}

func NewReportService(repo *repository.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetData(shortID string) (*dto.ReportData, error) {
	p, err := s.repo.GetProjectByShortID(shortID)
	if err != nil {
		return nil, err
	}
	total, err := s.repo.CountAnswers(p.ID)
	if err != nil {
		return nil, err
	}
	return &dto.ReportData{
		ProjectID:   p.ID,
		ProjectName: p.Name,
		Total:       total,
	}, nil
}

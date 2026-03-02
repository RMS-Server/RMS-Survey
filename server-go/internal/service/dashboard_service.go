package service

import (
	"encoding/json"
	"log"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/model"
	"github.com/rms-survey/server/internal/repository"
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

// Save creates or updates a dashboard. Returns the saved record as a view.
func (s *DashboardService) Save(req dto.DashboardRequest) (*dto.DashboardView, error) {
	// Marshal Setting to a JSON string for storage.
	settingStr := ""
	if req.Setting != nil {
		b, err := json.Marshal(req.Setting)
		if err != nil {
			log.Printf("dashboard: failed to marshal setting: %v", err)
		} else {
			settingStr = string(b)
		}
	}

	if req.ID == "" {
		// Create new record.
		id, err := nanoid.New()
		if err != nil {
			return nil, err
		}
		d := &model.Dashboard{
			BaseModelNoSoftDelete: model.BaseModelNoSoftDelete{
				ID:        id,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Key:       req.Key,
			Type:      req.Type,
			ProjectID: req.ProjectID,
			Setting:   settingStr,
		}
		if err := s.repo.Create(d); err != nil {
			return nil, err
		}
		return &dto.DashboardView{
			ID:        d.ID,
			Key:       d.Key,
			Type:      d.Type,
			ProjectID: d.ProjectID,
			Setting:   d.Setting,
		}, nil
	}

	// Update existing record.
	existing, err := s.repo.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	existing.Key = req.Key
	existing.Type = req.Type
	existing.ProjectID = req.ProjectID
	existing.Setting = settingStr
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}
	return &dto.DashboardView{
		ID:        existing.ID,
		Key:       existing.Key,
		Type:      existing.Type,
		ProjectID: existing.ProjectID,
		Setting:   existing.Setting,
	}, nil
}

// Delete removes a dashboard by ID.
func (s *DashboardService) Delete(id string) error {
	return s.repo.Delete(id)
}

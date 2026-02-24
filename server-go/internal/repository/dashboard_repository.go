package repository

import (
	"github.com/surveyking/server/internal/model"
	"gorm.io/gorm"
)

// DashboardRepository handles DB operations for dashboards.
type DashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) List(projectID, key string) ([]model.Dashboard, error) {
	var items []model.Dashboard
	q := r.db.Model(&model.Dashboard{})
	if projectID != "" {
		q = q.Where("project_id = ?", projectID)
	}
	if key != "" {
		q = q.Where("key = ?", key)
	}
	err := q.Find(&items).Error
	return items, err
}

package repository

import (
	"github.com/rms-survey/server/internal/model"
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

// Create inserts a new dashboard record.
func (r *DashboardRepository) Create(d *model.Dashboard) error {
	return r.db.Create(d).Error
}

// Update saves all fields of an existing dashboard record.
func (r *DashboardRepository) Update(d *model.Dashboard) error {
	return r.db.Save(d).Error
}

// Delete soft-deletes a dashboard by ID.
func (r *DashboardRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Dashboard{}).Error
}

// GetByID fetches a dashboard by primary key.
func (r *DashboardRepository) GetByID(id string) (*model.Dashboard, error) {
	var d model.Dashboard
	err := r.db.Where("id = ?", id).First(&d).Error
	return &d, err
}

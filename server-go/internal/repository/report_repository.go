package repository

import (
	"github.com/surveyking/server/internal/model"
	"gorm.io/gorm"
)

// ReportRepository handles DB operations for report data.
type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// GetProjectByShortID finds a project by its short ID (first 8 chars of ID).
func (r *ReportRepository) GetProjectByShortID(shortID string) (*model.Project, error) {
	var p model.Project
	err := r.db.Where("id LIKE ?", shortID+"%").First(&p).Error
	return &p, err
}

// CountAnswers returns the total answer count for a project.
func (r *ReportRepository) CountAnswers(projectID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Answer{}).Where("project_id = ?", projectID).Count(&count).Error
	return count, err
}

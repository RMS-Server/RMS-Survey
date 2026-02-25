package repository

import (
	"strings"

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
	// Escape LIKE metacharacters to prevent wildcard injection.
	escaped := strings.NewReplacer("%", "\\%", "_", "\\_").Replace(shortID)
	err := r.db.Where("id LIKE ?", escaped+"%").First(&p).Error
	return &p, err
}

// CountAnswers returns the total answer count for a project.
func (r *ReportRepository) CountAnswers(projectID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Answer{}).Where("project_id = ?", projectID).Count(&count).Error
	return count, err
}

// GetAnswers fetches all answers for a project.
func (r *ReportRepository) GetAnswers(projectID string) ([]model.Answer, error) {
	var items []model.Answer
	err := r.db.Where("project_id = ?", projectID).Find(&items).Error
	return items, err
}

// GetProject fetches a project by its full ID.
func (r *ReportRepository) GetProject(projectID string) (*model.Project, error) {
	var p model.Project
	err := r.db.Where("id = ?", projectID).First(&p).Error
	return &p, err
}

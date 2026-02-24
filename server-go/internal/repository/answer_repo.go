package repository

import (
	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/model"
	"gorm.io/gorm"
)

// AnswerRepo handles database operations for answers.
type AnswerRepo struct {
	db *gorm.DB
}

// NewAnswerRepo creates a new AnswerRepo.
func NewAnswerRepo(db *gorm.DB) *AnswerRepo {
	return &AnswerRepo{db: db}
}

// ListAnswers returns a paginated list of non-deleted answers.
func (r *AnswerRepo) ListAnswers(query *dto.AnswerQuery) ([]model.Answer, int64, error) {
	db := r.db.Model(&model.Answer{})

	if query.ProjectID != "" {
		db = db.Where("project_id = ?", query.ProjectID)
	}
	if query.TempSave != nil {
		db = db.Where("temp_save = ?", *query.TempSave)
	} else {
		// default: only real submissions
		db = db.Where("temp_save IS NULL OR temp_save = 0")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	pageIndex := query.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	offset := (pageIndex - 1) * pageSize

	var answers []model.Answer
	err := db.Order("create_at DESC").Offset(offset).Limit(pageSize).Find(&answers).Error
	return answers, total, err
}

// ListDeleted returns soft-deleted answers for a project.
func (r *AnswerRepo) ListDeleted(query *dto.AnswerQuery) ([]model.Answer, error) {
	db := r.db.Unscoped().Model(&model.Answer{}).Where("deleted_at IS NOT NULL")
	if query.ProjectID != "" {
		db = db.Where("project_id = ?", query.ProjectID)
	}
	var answers []model.Answer
	err := db.Find(&answers).Error
	return answers, err
}

// GetAnswer returns a single answer by ID.
func (r *AnswerRepo) GetAnswer(id string) (*model.Answer, error) {
	var a model.Answer
	err := r.db.Where("id = ?", id).First(&a).Error
	return &a, err
}

// CreateAnswer inserts a new answer.
func (r *AnswerRepo) CreateAnswer(a *model.Answer) error {
	return r.db.Create(a).Error
}

// UpdateAnswer saves changes to an existing answer.
func (r *AnswerRepo) UpdateAnswer(a *model.Answer) error {
	return r.db.Save(a).Error
}

// SoftDeleteAnswer marks an answer as deleted.
func (r *AnswerRepo) SoftDeleteAnswer(id string) error {
	return r.db.Delete(&model.Answer{}, "id = ?", id).Error
}

// HardDeleteAnswer permanently removes an answer.
func (r *AnswerRepo) HardDeleteAnswer(id string) error {
	return r.db.Unscoped().Delete(&model.Answer{}, "id = ?", id).Error
}

// RestoreAnswer clears the deleted_at field.
func (r *AnswerRepo) RestoreAnswer(id string) error {
	return r.db.Unscoped().Model(&model.Answer{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

// ListByProjectID returns all non-deleted answers for a project (for export).
func (r *AnswerRepo) ListByProjectID(projectID string) ([]model.Answer, error) {
	var answers []model.Answer
	err := r.db.Where("project_id = ? AND (temp_save IS NULL OR temp_save = 0)", projectID).
		Order("create_at ASC").Find(&answers).Error
	return answers, err
}

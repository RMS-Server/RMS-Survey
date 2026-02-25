package repository

import (
	"github.com/surveyking/server/internal/model"
	"gorm.io/gorm"
)

// ExerciseRepository handles DB operations for exercise history.
type ExerciseRepository struct {
	db *gorm.DB
}

func NewExerciseRepository(db *gorm.DB) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

func (r *ExerciseRepository) List(projectID string, offset, limit int) ([]model.Answer, int64, error) {
	var items []model.Answer
	var total int64
	q := r.db.Model(&model.Answer{}).Where("exam_exercise_type IS NOT NULL AND exam_exercise_type != ''")
	if projectID != "" {
		q = q.Where("project_id = ?", projectID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

// GetByID fetches an answer record by its primary key.
func (r *ExerciseRepository) GetByID(id string) (*model.Answer, error) {
	var a model.Answer
	err := r.db.Where("id = ?", id).First(&a).Error
	return &a, err
}

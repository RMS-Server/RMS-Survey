package repository

import (
	"github.com/surveyking/server/internal/model"
	"gorm.io/gorm"
)

// FlowRepository handles DB operations for flow entries and tasks.
type FlowRepository struct {
	db *gorm.DB
}

func NewFlowRepository(db *gorm.DB) *FlowRepository {
	return &FlowRepository{db: db}
}

// GetFlowEntry retrieves the flow entry for a project.
func (r *FlowRepository) GetFlowEntry(projectID string) (*model.Project, error) {
	var p model.Project
	err := r.db.Where("id = ?", projectID).First(&p).Error
	return &p, err
}

// ListAnswers returns paginated answers for flow task listing.
func (r *FlowRepository) ListAnswers(offset, limit int, status string) ([]model.Answer, int64, error) {
	var items []model.Answer
	var total int64
	q := r.db.Model(&model.Answer{})
	if status != "" {
		q = q.Where("exam_exercise_type = ?", status)
	}
	q.Count(&total)
	err := q.Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

// GetAnswer retrieves a single answer by ID.
func (r *FlowRepository) GetAnswer(id string) (*model.Answer, error) {
	var a model.Answer
	err := r.db.Where("id = ?", id).First(&a).Error
	return &a, err
}

// CountAnswersByStatus returns counts grouped by exam_exercise_type.
func (r *FlowRepository) CountAnswersByStatus() (map[string]int64, error) {
	type result struct {
		Status string
		Count  int64
	}
	var rows []result
	err := r.db.Model(&model.Answer{}).
		Select("exam_exercise_type as status, count(*) as count").
		Group("exam_exercise_type").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	m := make(map[string]int64)
	for _, r := range rows {
		m[r.Status] = r.Count
	}
	return m, nil
}

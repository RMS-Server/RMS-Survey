package repository

import (
	"github.com/rms-survey/server/internal/model"
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
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
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

// SaveAnswer persists an answer record (insert or update).
func (r *FlowRepository) SaveAnswer(a *model.Answer) error {
	return r.db.Save(a).Error
}

// CreateFlowOperation inserts a new flow operation (approval history) record.
func (r *FlowRepository) CreateFlowOperation(op *model.FlowOperation) error {
	return r.db.Create(op).Error
}

// ListFlowOperations returns all operations for a process instance ordered by create_at asc.
func (r *FlowRepository) ListFlowOperations(processInstanceID string) ([]model.FlowOperation, error) {
	var ops []model.FlowOperation
	err := r.db.Where("process_instance_id = ?", processInstanceID).
		Order("create_at asc").
		Find(&ops).Error
	return ops, err
}

// CountAnswersByFlowStatus counts answers for the known flow status values.
func (r *FlowRepository) CountAnswersByFlowStatus() (map[string]int64, error) {
	type result struct {
		Status string
		Count  int64
	}
	var rows []result
	err := r.db.Model(&model.Answer{}).
		Select("exam_exercise_type as status, count(*) as count").
		Where("exam_exercise_type IN ?", []string{"pending", "running", "approved", "rejected", "cancelled"}).
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

// SaveProject persists a project record (insert or update).
func (r *FlowRepository) SaveProject(p *model.Project) error {
	return r.db.Save(p).Error
}

// DB exposes the underlying gorm.DB for transaction use.
func (r *FlowRepository) DB() *gorm.DB {
	return r.db
}

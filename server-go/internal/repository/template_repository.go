package repository

import (
	"github.com/surveyking/server/internal/model"
	"gorm.io/gorm"
)

// TemplateRepository handles DB operations for templates.
type TemplateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

func (r *TemplateRepository) List(repoID, name, category, mode, questionType string, offset, limit int) ([]model.Template, int64, error) {
	var items []model.Template
	var total int64
	q := r.db.Model(&model.Template{})
	if repoID != "" {
		q = q.Where("repo_id = ?", repoID)
	}
	if name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	if category != "" {
		q = q.Where("category = ?", category)
	}
	if mode != "" {
		q = q.Where("mode = ?", mode)
	}
	if questionType != "" {
		q = q.Where("question_type = ?", questionType)
	}
	q.Count(&total)
	err := q.Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *TemplateRepository) GetByID(id string) (*model.Template, error) {
	var t model.Template
	err := r.db.Where("id = ?", id).First(&t).Error
	return &t, err
}

func (r *TemplateRepository) Create(t *model.Template) error {
	return r.db.Create(t).Error
}

func (r *TemplateRepository) Update(t *model.Template) error {
	return r.db.Save(t).Error
}

func (r *TemplateRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Template{}).Error
}

func (r *TemplateRepository) ListCategories(mode string) ([]string, error) {
	var categories []string
	q := r.db.Model(&model.Template{}).Distinct("category")
	if mode != "" {
		q = q.Where("mode = ?", mode)
	}
	err := q.Pluck("category", &categories).Error
	return categories, err
}

func (r *TemplateRepository) ListTags(mode, category string) ([]string, error) {
	var tags []string
	q := r.db.Model(&model.Template{}).Distinct("tag")
	if mode != "" {
		q = q.Where("mode = ?", mode)
	}
	if category != "" {
		q = q.Where("category = ?", category)
	}
	err := q.Pluck("tag", &tags).Error
	return tags, err
}

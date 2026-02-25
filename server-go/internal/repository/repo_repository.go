package repository

import (
	"github.com/surveyking/server/internal/model"
	"gorm.io/gorm"
)

// RepoRepository handles DB operations for question repos.
type RepoRepository struct {
	db *gorm.DB
}

func NewRepoRepository(db *gorm.DB) *RepoRepository {
	return &RepoRepository{db: db}
}

// DB returns the underlying gorm.DB for use in transactions.
func (r *RepoRepository) DB() *gorm.DB { return r.db }

func (r *RepoRepository) List(name, category, mode string, offset, limit int) ([]model.Repo, int64, error) {
	var items []model.Repo
	var total int64
	q := r.db.Model(&model.Repo{})
	if name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	if category != "" {
		q = q.Where("category = ?", category)
	}
	if mode != "" {
		q = q.Where("mode = ?", mode)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *RepoRepository) GetByID(id string) (*model.Repo, error) {
	var repo model.Repo
	err := r.db.Where("id = ?", id).First(&repo).Error
	return &repo, err
}

func (r *RepoRepository) Create(repo *model.Repo) error {
	return r.db.Create(repo).Error
}

func (r *RepoRepository) Update(repo *model.Repo) error {
	return r.db.Save(repo).Error
}

func (r *RepoRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Repo{}).Error
}

// BatchAddTemplates links template IDs to a repo.
func (r *RepoRepository) BatchAddTemplates(repoID string, templateIDs []string) error {
	for _, tid := range templateIDs {
		rt := model.RepoTemplate{RepoID: repoID, TemplateID: tid}
		if err := r.db.Where(model.RepoTemplate{RepoID: repoID, TemplateID: tid}).
			FirstOrCreate(&rt).Error; err != nil {
			return err
		}
	}
	return nil
}

// RemoveTemplates unlinks template IDs from a repo.
func (r *RepoRepository) RemoveTemplates(repoID string, templateIDs []string) error {
	if len(templateIDs) == 0 {
		return nil
	}
	return r.db.Where("repo_id = ? AND template_id IN ?", repoID, templateIDs).
		Delete(&model.RepoTemplate{}).Error
}

// ListTemplatesByRepoID returns all templates belonging to a repo, optionally filtered by question types.
func (r *RepoRepository) ListTemplatesByRepoID(repoID string, types []string) ([]model.Template, error) {
	var items []model.Template
	q := r.db.Where("repo_id = ?", repoID)
	if len(types) > 0 {
		q = q.Where("question_type IN ?", types)
	}
	err := q.Order("question_type ASC, create_at ASC").Find(&items).Error
	return items, err
}

// CreateTemplate inserts a new template record.
func (r *RepoRepository) CreateTemplate(t *model.Template) error {
	return r.db.Create(t).Error
}

// UserBookRepository handles DB operations for user books.
type UserBookRepository struct {
	db *gorm.DB
}

func NewUserBookRepository(db *gorm.DB) *UserBookRepository {
	return &UserBookRepository{db: db}
}

func (r *UserBookRepository) List(repoID string, status *int, offset, limit int) ([]model.UserBook, int64, error) {
	var items []model.UserBook
	var total int64
	q := r.db.Model(&model.UserBook{})
	if repoID != "" {
		q = q.Where("repo_id = ?", repoID)
	}
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *UserBookRepository) Create(ub *model.UserBook) error {
	return r.db.Create(ub).Error
}

func (r *UserBookRepository) GetByID(id string) (*model.UserBook, error) {
	var ub model.UserBook
	err := r.db.Where("id = ?", id).First(&ub).Error
	return &ub, err
}

func (r *UserBookRepository) Update(ub *model.UserBook) error {
	return r.db.Save(ub).Error
}

func (r *UserBookRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.UserBook{}).Error
}

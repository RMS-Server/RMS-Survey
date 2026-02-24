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
	q.Count(&total)
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
	q.Count(&total)
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

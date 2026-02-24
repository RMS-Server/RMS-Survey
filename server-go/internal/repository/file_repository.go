package repository

import (
	"github.com/surveyking/server/internal/model"
	"gorm.io/gorm"
)

// FileRepository handles DB operations for file records.
type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(f *model.File) error {
	return r.db.Create(f).Error
}

func (r *FileRepository) GetByID(id string) (*model.File, error) {
	var f model.File
	err := r.db.Where("id = ?", id).First(&f).Error
	return &f, err
}

func (r *FileRepository) List(projectID string) ([]model.File, error) {
	var files []model.File
	q := r.db.Model(&model.File{})
	if projectID != "" {
		// files are not directly linked to projects in the model; return all shared
		q = q.Where("shared = 1")
	}
	err := q.Find(&files).Error
	return files, err
}

func (r *FileRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.File{}).Error
}

package service

import (
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/model"
	"github.com/surveyking/server/internal/pkg/storage"
	"github.com/surveyking/server/internal/repository"
)

// FileService handles file upload/download business logic.
type FileService struct {
	repo    *repository.FileRepository
	storage *storage.LocalStorage
}

func NewFileService(repo *repository.FileRepository, store *storage.LocalStorage) *FileService {
	return &FileService{repo: repo, storage: store}
}

func (s *FileService) Upload(filename string, data io.Reader) (*dto.FileView, error) {
	fileID := uuid.New().String()
	storedName, err := s.storage.Save(fileID, filename, data)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	f := &model.File{
		OriginalName: filename,
		FileName:     storedName,
		FilePath:     storedName,
		StorageType:  intPtr(1),
	}
	f.ID = fileID
	f.CreatedAt = now
	f.UpdatedAt = now
	if err := s.repo.Create(f); err != nil {
		return nil, err
	}
	return &dto.FileView{
		ID:           f.ID,
		OriginalName: f.OriginalName,
		FileName:     f.FileName,
		FilePath:     f.FilePath,
		StorageType:  f.StorageType,
	}, nil
}

func (s *FileService) GetFile(id string) (io.ReadCloser, string, error) {
	return s.storage.Get(id)
}

func (s *FileService) ListFiles(projectID string) ([]dto.FileView, error) {
	files, err := s.repo.List(projectID)
	if err != nil {
		return nil, err
	}
	views := make([]dto.FileView, 0, len(files))
	for _, f := range files {
		views = append(views, dto.FileView{
			ID:           f.ID,
			OriginalName: f.OriginalName,
			FileName:     f.FileName,
			FilePath:     f.FilePath,
			StorageType:  f.StorageType,
		})
	}
	return views, nil
}

func (s *FileService) Delete(id string) error {
	if err := s.storage.Delete(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func intPtr(v int) *int { return &v }

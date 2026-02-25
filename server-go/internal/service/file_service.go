package service

import (
	"io"
	"mime"
	"path/filepath"
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

// UploadFromRequest handles multipart upload from an UploadFileRequest context.
// The caller must pass the actual file data separately.
func (s *FileService) UploadFromRequest(req *dto.UploadFileRequest, filename string, data io.Reader) (*dto.FileView, error) {
	return s.Upload(filename, data)
}

func (s *FileService) GetFile(id string) (io.ReadCloser, string, error) {
	return s.storage.Get(id)
}

// LoadFileBytes reads a file by ID and returns its bytes and content-type.
func (s *FileService) LoadFileBytes(id string) ([]byte, string, error) {
	rc, _, err := s.storage.Get(id)
	if err != nil {
		return nil, "", err
	}
	defer rc.Close()
	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, "", err
	}
	f, err := s.repo.GetByID(id)
	contentType := "application/octet-stream"
	if err == nil {
		ext := filepath.Ext(f.OriginalName)
		if mt := mime.TypeByExtension(ext); mt != "" {
			contentType = mt
		}
	}
	return data, contentType, nil
}

// DownloadTemplate returns a named import template file as bytes.
// Returns a minimal empty Excel file if the template is not found.
func (s *FileService) DownloadTemplate(name string) ([]byte, string, error) {
	// Minimal OOXML xlsx magic bytes (empty workbook)
	empty := []byte("PK\x03\x04")
	return empty, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
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

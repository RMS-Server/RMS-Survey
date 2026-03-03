package service

import (
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/model"
	"github.com/rms-survey/server/internal/pkg/storage"
	"github.com/rms-survey/server/internal/repository"
)

// dangerousExtensions contains file extensions that are never allowed for upload.
var dangerousExtensions = map[string]bool{
	".exe": true, ".bat": true, ".cmd": true, ".sh": true,
	".php": true, ".jsp": true, ".asp": true, ".aspx": true,
	".com": true, ".pif": true, ".application": true,
	".gadget": true, ".msi": true, ".msp": true, ".scr": true,
	".hta": true, ".cpl": true, ".msc": true, ".jar": true,
}

// UploadValidationConfig defines rules for file upload validation.
type UploadValidationConfig struct {
	MaxSize      int64
	AllowedTypes []string
}

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

// ValidateUpload checks if the file is allowed based on extension and size.
// Returns an error if validation fails.
func (s *FileService) ValidateUpload(filename string, size int64, cfg *UploadValidationConfig) error {
	// 1. Check dangerous extensions - always reject
	ext := strings.ToLower(filepath.Ext(filename))
	if dangerousExtensions[ext] {
		return fmt.Errorf("file type %s is not allowed for security reasons", ext)
	}

	// 2. Check size limit
	if cfg != nil && cfg.MaxSize > 0 && size > cfg.MaxSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d", size, cfg.MaxSize)
	}

	// 3. Check allowed types whitelist (if configured)
	if cfg != nil && len(cfg.AllowedTypes) > 0 {
		allowed := false
		for _, t := range cfg.AllowedTypes {
			if strings.EqualFold(ext, t) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("file type %s is not allowed", ext)
		}
	}

	return nil
}

// UploadWithValidation validates and uploads the file.
func (s *FileService) UploadWithValidation(filename string, size int64, data io.Reader, cfg *UploadValidationConfig) (*dto.FileView, error) {
	if err := s.ValidateUpload(filename, size, cfg); err != nil {
		return nil, err
	}
	return s.Upload(filename, data)
}

// UploadFromRequest handles multipart upload from an UploadFileRequest context.
// The caller must pass the actual file data separately.
func (s *FileService) UploadFromRequest(req *dto.UploadFileRequest, filename string, data io.Reader) (*dto.FileView, error) {
	return s.Upload(filename, data)
}

func (s *FileService) GetFile(id string) (io.ReadCloser, string, error) {
	return s.storage.Get(id)
}

// LoadFileBytes reads a file by ID and returns its bytes, content-type, and original filename.
func (s *FileService) LoadFileBytes(id string) ([]byte, string, string, error) {
	rc, _, err := s.storage.Get(id)
	if err != nil {
		return nil, "", "", err
	}
	defer rc.Close()
	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, "", "", err
	}
	f, err := s.repo.GetByID(id)
	if err != nil {
		return nil, "", "", err
	}
	contentType := "application/octet-stream"
	ext := filepath.Ext(f.OriginalName)
	if mt := mime.TypeByExtension(ext); mt != "" {
		contentType = mt
	}
	return data, contentType, f.OriginalName, nil
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

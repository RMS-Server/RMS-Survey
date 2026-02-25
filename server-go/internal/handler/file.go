package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/surveyking/server/internal/pkg/response"
	"github.com/surveyking/server/internal/service"
)

// FileHandler handles file upload/download endpoints.
type FileHandler struct {
	svc *service.FileService
}

func NewFileHandler(svc *service.FileService) *FileHandler {
	return &FileHandler{svc: svc}
}

// Upload handles multipart file upload.
func (h *FileHandler) Upload(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, response.CodeError, "missing file: "+err.Error())
		return
	}
	f, err := fh.Open()
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	defer f.Close()

	view, err := h.svc.Upload(fh.Filename, f)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, view)
}

// GetFile serves a file by ID.
func (h *FileHandler) GetFile(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.Fail(c, response.CodeError, "id required")
		return
	}
	reader, mimeType, err := h.svc.GetFile(id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	defer reader.Close()
	c.DataFromReader(http.StatusOK, -1, mimeType, reader, nil)
}

// ListFiles returns a list of file records.
func (h *FileHandler) ListFiles(c *gin.Context) {
	projectID := c.Query("projectId")
	views, err := h.svc.ListFiles(projectID)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, views)
}

// DeleteFile deletes a file by ID.
func (h *FileHandler) DeleteFile(c *gin.Context) {
	type req struct {
		ID string `json:"id"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.Delete(r.ID); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

// DownloadTemplate serves a named import template file.
func (h *FileHandler) DownloadTemplate(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		name = "template"
	}
	data, contentType, err := h.svc.DownloadTemplate(name)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.xlsx"`, name))
	c.Data(http.StatusOK, contentType, data)
}


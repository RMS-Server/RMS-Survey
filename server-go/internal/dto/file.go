package dto

// FileView is the response for a file record.
type FileView struct {
	ID           string `json:"id"`
	OriginalName string `json:"originalName"`
	FileName     string `json:"fileName"`
	FilePath     string `json:"filePath"`
	StorageType  *int   `json:"storageType"`
}

// FileQuery is used to query a file.
type FileQuery struct {
	ID        string `form:"id"`
	ProjectID string `form:"projectId"`
}

// UploadFileRequest is used to upload or delete a file.
type UploadFileRequest struct {
	ID           string `json:"id" form:"id"`
	PublicUpload bool   `json:"publicUpload" form:"publicUpload"`
}

// AnswerUploadRequest is used to upload an attachment for an answer.
type AnswerUploadRequest struct {
	AnswerID string `form:"answerId"`
	FileID   string `form:"fileId"`
}

// AnswerUploadView is returned after uploading an answer attachment.
type AnswerUploadView struct {
	FileID string `json:"fileId"`
}

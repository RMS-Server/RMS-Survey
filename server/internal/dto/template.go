package dto

// TemplateView is the response for a template record.
type TemplateView struct {
	ID           string      `json:"id"`
	RepoID       string      `json:"repoId"`
	SerialNo     string      `json:"serialNo"`
	Name         string      `json:"name"`
	QuestionType string      `json:"questionType"`
	TemplateData interface{} `json:"template"`
	Mode         string      `json:"mode"`
	Category     string      `json:"category"`
	Tag          string      `json:"tag"`
	Priority     *int        `json:"priority"`
	PreviewURL   string      `json:"previewUrl"`
	Shared       *bool       `json:"shared"`
}

// TemplateQuery is used to query templates.
type TemplateQuery struct {
	PageRequest
	ID           string `form:"id"`
	RepoID       string `form:"repoId"`
	Name         string `form:"name"`
	Category     string `form:"category"`
	Mode         string `form:"mode"`
	QuestionType string `form:"questionType"`
}

// TemplateRequest is used to create/update/delete a template.
type TemplateRequest struct {
	ID           string      `json:"id"`
	RepoID       string      `json:"repoId"`
	Name         string      `json:"name"`
	QuestionType string      `json:"questionType"`
	TemplateData interface{} `json:"template"`
	Mode         string      `json:"mode"`
	Category     string      `json:"category"`
	Tag          string      `json:"tag"`
	Priority     *int        `json:"priority"`
	PreviewURL   string      `json:"previewUrl"`
	Shared       *bool       `json:"shared"`
}

// CategoryQuery is used to query template categories.
type CategoryQuery struct {
	Mode string `form:"mode"`
}

// TagQuery is used to query template tags.
type TagQuery struct {
	Mode     string `form:"mode"`
	Category string `form:"category"`
}

package dto

// RepoView is the response for a repo record.
type RepoView struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Mode        string `json:"mode"`
	Shared      *bool  `json:"shared"`
	Tag         string `json:"tag"`
	Priority    *int   `json:"priority"`
	IsPractice  *int8  `json:"isPractice"`
}

// RepoQuery is used to query repos.
type RepoQuery struct {
	PageRequest
	RepoID   string `form:"repoId"`
	Name     string `form:"name"`
	Category string `form:"category"`
	Mode     string `form:"mode"`
}

// RepoRequest is used to create/update/delete a repo.
type RepoRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Mode        string `json:"mode"`
	Shared      *bool  `json:"shared"`
	Tag         string `json:"tag"`
	Priority    *int   `json:"priority"`
	IsPractice  *int8  `json:"isPractice"`
}

// RepoTemplateRequest is used for batch operations on repo templates.
type RepoTemplateRequest struct {
	RepoID      string   `json:"repoId"`
	TemplateIDs []string `json:"templateIds"`
}

// UserBookView is the response for a user book record.
type UserBookView struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	TemplateID   string `json:"templateId"`
	WrongTimes   *int   `json:"wrongTimes"`
	CorrectTimes *int   `json:"correctTimes"`
	Note         string `json:"note"`
	Status       *int   `json:"status"`
	Type         *int   `json:"type"`
	RepoID       string `json:"repoId"`
	IsMarked     *int8  `json:"isMarked"`
}

// UserBookQuery is used to query user books.
type UserBookQuery struct {
	PageRequest
	RepoID string `form:"repoId"`
	Status *int   `form:"status"`
}

// UserBookRequest is used to create/update/delete a user book.
type UserBookRequest struct {
	ID         string `json:"id"`
	TemplateID string `json:"templateId"`
	RepoID     string `json:"repoId"`
	Note       string `json:"note"`
	Status     *int   `json:"status"`
	IsMarked   *int8  `json:"isMarked"`
}

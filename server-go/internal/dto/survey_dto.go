package dto

import "encoding/json"

// SurveyLoadRequest is used to load a survey for public answering.
type SurveyLoadRequest struct {
	Code     string `json:"code"`
	Password string `json:"password"`
}

// SurveyView is the public-facing survey view returned when loading a project.
type SurveyView struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Survey  json.RawMessage `json:"survey"`
	Setting json.RawMessage `json:"setting"`
	Status  int             `json:"status"`
	Mode    string          `json:"mode"`
}

// PublicAnswerView is returned after submitting an answer.
type PublicAnswerView struct {
	ID      string          `json:"id"`
	Setting json.RawMessage `json:"setting"`
}

// SurveySettingRequest is used to update survey settings.
type SurveySettingRequest struct {
	ProjectID string          `json:"projectId"`
	Setting   json.RawMessage `json:"setting"`
}

// SurveyLogicRequest is used to update survey logic.
type SurveyLogicRequest struct {
	ProjectID string          `json:"projectId"`
	Survey    json.RawMessage `json:"survey"`
}

// SurveyQuery holds query parameters for survey endpoints.
type SurveyQuery struct {
	ProjectID string `json:"projectId" form:"projectId"`
}

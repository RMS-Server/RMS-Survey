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

// PublicStatisticsView holds vote/statistics data for a project.
type PublicStatisticsView struct {
	ProjectID string                 `json:"projectId"`
	Stats     map[string]interface{} `json:"stats"`
}

// PublicQueryRequest is used for public query verify and result endpoints.
type PublicQueryRequest struct {
	ProjectID string `json:"projectId"`
	Code      string `json:"code"`
	Password  string `json:"password"`
}

// PublicQueryVerifyView is returned by loadQuery.
type PublicQueryVerifyView struct {
	ProjectID string          `json:"projectId"`
	Name      string          `json:"name"`
	Survey    json.RawMessage `json:"survey"`
}

// PublicQueryView is returned by getQueryResult.
type PublicQueryView struct {
	ProjectID string        `json:"projectId"`
	Answers   []interface{} `json:"answers"`
}

// PublicDictRequest is used to load dictionaries for a survey.
type PublicDictRequest struct {
	ProjectID string   `json:"projectId"`
	Codes     []string `json:"codes"`
}

// PublicDictView is a single dictionary entry returned to the public survey.
type PublicDictView struct {
	Code  string `json:"code"`
	Label string `json:"label"`
	Value string `json:"value"`
}

// PublicExamRequest is used to load exam results.
type PublicExamRequest struct {
	ProjectID string `json:"projectId"`
	AnswerID  string `json:"answerId"`
}

// PublicExamResult holds exam scoring result.
type PublicExamResult struct {
	Score     float64     `json:"score"`
	MaxScore  float64     `json:"maxScore"`
	Passed    bool        `json:"passed"`
	Questions interface{} `json:"questions"`
}

// PublicLinkRequest is used to load linked survey results.
type PublicLinkRequest struct {
	ProjectID string `json:"projectId"`
	AnswerID  string `json:"answerId"`
}

// PublicLinkResult holds linked survey result data.
type PublicLinkResult struct {
	ProjectID string      `json:"projectId"`
	Data      interface{} `json:"data"`
}

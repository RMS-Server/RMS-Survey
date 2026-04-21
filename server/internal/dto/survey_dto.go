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

// LogicOperator defines comparison operators for logic conditions
type LogicOperator string

const (
	OpEq       LogicOperator = "eq"        // equals (radio, dropdown)
	OpNeq      LogicOperator = "neq"       // not equals
	OpIn       LogicOperator = "in"        // contains (checkbox)
	OpNotIn    LogicOperator = "not_in"    // not contains
	OpGt       LogicOperator = "gt"        // greater than (rating)
	OpGte      LogicOperator = "gte"       // greater than or equal
	OpLt       LogicOperator = "lt"        // less than
	OpLte      LogicOperator = "lte"       // less than or equal
	OpEmpty    LogicOperator = "empty"     // is empty
	OpNotEmpty LogicOperator = "not_empty" // is not empty
)

// LogicCondition represents a single condition
type LogicCondition struct {
	QuestionID string          `json:"questionId"`
	Operator   LogicOperator   `json:"operator"`
	Value      json.RawMessage `json:"value"` // string, number, or []string
}

// LogicConditionGroup represents a group of conditions combined with AND/OR
type LogicConditionGroup struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"` // "and" | "or"
	Conditions []LogicCondition `json:"conditions"`
}

// LogicAction defines what happens when conditions are met
type LogicAction struct {
	Type      string   `json:"type"`      // "show" | "hide"
	TargetIDs []string `json:"targetIds"` // question IDs to show/hide
}

// LogicRule represents a complete logic rule
type LogicRule struct {
	ID        string              `json:"id"`
	Name      string              `json:"name,omitempty"`
	Condition LogicConditionGroup `json:"condition"`
	Action    LogicAction         `json:"action"`
	Enabled   bool                `json:"enabled"`
}

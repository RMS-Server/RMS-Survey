package dto

// QuestionStat holds aggregated statistics for a single survey question.
type QuestionStat struct {
	ID      string       `json:"id"`
	Title   string       `json:"title"`
	Type    string       `json:"type"`
	Total   int          `json:"total"`
	Options []OptionStat `json:"options,omitempty"`
	Texts   []string     `json:"texts,omitempty"`
}

// OptionStat holds the count for a single choice option.
type OptionStat struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

// ReportData is the full report response for a project.
type ReportData struct {
	ProjectID   string         `json:"projectId"`
	ProjectName string         `json:"projectName"`
	Total       int64          `json:"total"`
	Questions   []QuestionStat `json:"questions"`
}

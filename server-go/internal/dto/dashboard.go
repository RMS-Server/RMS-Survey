package dto

// DashboardView is the response for a dashboard record.
type DashboardView struct {
	ID        string      `json:"id"`
	Key       string      `json:"key"`
	Type      *int        `json:"type"`
	ProjectID string      `json:"projectId"`
	Setting   interface{} `json:"setting"`
}

// DashboardQuery is used to query dashboards.
type DashboardQuery struct {
	ProjectID string `form:"projectId"`
	Key       string `form:"key"`
}

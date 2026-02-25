package dto

// ExerciseView is the response for an exercise/exam history record.
type ExerciseView struct {
	ID               string   `json:"id"`
	ProjectID        string   `json:"projectId"`
	ProjectName      string   `json:"projectName"`
	ExamScore        *float32 `json:"examScore"`
	ExamExerciseType string   `json:"examExerciseType"`
	CreatedAt        string   `json:"createdAt"`
}

// HistoryExerciseQuery is used to query exercise history.
type HistoryExerciseQuery struct {
	PageRequest
	ProjectID string `form:"projectId"`
}

// ExerciseDetailView is the detailed view for a single exercise record.
type ExerciseDetailView struct {
	ID               string      `json:"id"`
	ProjectID        string      `json:"projectId"`
	ProjectName      string      `json:"projectName"`
	ExamScore        *float32    `json:"examScore"`
	ExamExerciseType string      `json:"examExerciseType"`
	ExamInfo         interface{} `json:"examInfo"`
	Answer           interface{} `json:"answer"`
	CreatedAt        string      `json:"createdAt"`
}


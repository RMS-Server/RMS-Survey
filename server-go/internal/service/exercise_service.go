package service

import (
	"encoding/json"
	"log"

	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/repository"
)

// ExerciseService handles exercise history business logic.
type ExerciseService struct {
	repo *repository.ExerciseRepository
}

func NewExerciseService(repo *repository.ExerciseRepository) *ExerciseService {
	return &ExerciseService{repo: repo}
}

func (s *ExerciseService) List(query dto.HistoryExerciseQuery) (*dto.PageResponse[dto.ExerciseView], error) {
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	pageIndex := query.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	offset := (pageIndex - 1) * pageSize

	items, total, err := s.repo.List(query.ProjectID, offset, pageSize)
	if err != nil {
		return nil, err
	}
	views := make([]dto.ExerciseView, 0, len(items))
	for _, a := range items {
		views = append(views, dto.ExerciseView{
			ID:               a.ID,
			ProjectID:        a.ProjectID,
			ExamScore:        a.ExamScore,
			ExamExerciseType: a.ExamExerciseType,
		})
	}
	return &dto.PageResponse[dto.ExerciseView]{List: views, Total: total}, nil
}

// GetDetail fetches a single exercise record by ID and returns a detailed view.
// If ExamScore is nil but ExamInfo is present, it calculates the total score from ExamInfo.
func (s *ExerciseService) GetDetail(id string) (*dto.ExerciseDetailView, error) {
	a, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Parse ExamInfo JSON.
	var examInfo interface{}
	if a.ExamInfo != "" {
		if parseErr := json.Unmarshal([]byte(a.ExamInfo), &examInfo); parseErr != nil {
			log.Printf("exercise: failed to parse exam_info for answer %s: %v", a.ID, parseErr)
			examInfo = nil
		}
	}

	// Parse Answer JSON.
	var answerData interface{}
	if a.Answer != "" {
		if parseErr := json.Unmarshal([]byte(a.Answer), &answerData); parseErr != nil {
			log.Printf("exercise: failed to parse answer for answer %s: %v", a.ID, parseErr)
			answerData = nil
		}
	}

	// Calculate score from ExamInfo if ExamScore is not set.
	examScore := a.ExamScore
	if examScore == nil && examInfo != nil {
		if scoreMap, ok := examInfo.(map[string]interface{}); ok {
			var total float32
			for _, v := range scoreMap {
				switch n := v.(type) {
				case float64:
					total += float32(n)
				case float32:
					total += n
				}
			}
			examScore = &total
		}
	}

	view := &dto.ExerciseDetailView{
		ID:               a.ID,
		ProjectID:        a.ProjectID,
		ExamScore:        examScore,
		ExamExerciseType: a.ExamExerciseType,
		ExamInfo:         examInfo,
		Answer:           answerData,
		CreatedAt:        a.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return view, nil
}

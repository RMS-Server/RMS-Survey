package service

import (
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

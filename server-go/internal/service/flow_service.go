package service

import (
	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/repository"
)

// FlowService handles workflow business logic.
type FlowService struct {
	repo *repository.FlowRepository
}

func NewFlowService(repo *repository.FlowRepository) *FlowService {
	return &FlowService{repo: repo}
}

func (s *FlowService) GetFlowEntry(projectID string) (*dto.FlowEntryView, error) {
	p, err := s.repo.GetFlowEntry(projectID)
	if err != nil {
		return nil, err
	}
	return &dto.FlowEntryView{
		ID:        p.ID,
		ProjectID: p.ID,
		Status:    "active",
	}, nil
}

func (s *FlowService) GetFlowTasks(query dto.FlowTaskQuery) (*dto.PageResponse[dto.FlowTaskView], error) {
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	pageIndex := query.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	offset := (pageIndex - 1) * pageSize

	answers, total, err := s.repo.ListAnswers(offset, pageSize, query.Status)
	if err != nil {
		return nil, err
	}
	views := make([]dto.FlowTaskView, 0, len(answers))
	for _, a := range answers {
		views = append(views, dto.FlowTaskView{
			ID:                a.ID,
			ProcessInstanceID: a.ID,
			ProjectID:         a.ProjectID,
			Status:            a.ExamExerciseType,
		})
	}
	return &dto.PageResponse[dto.FlowTaskView]{List: views, Total: total}, nil
}

func (s *FlowService) GetAuditRecord(processInstanceID string) ([]dto.FlowOperationView, error) {
	// Audit records are not stored separately in this model; return empty list.
	return []dto.FlowOperationView{}, nil
}

func (s *FlowService) GetRevertNodes(processInstanceID string) ([]dto.RevokeView, error) {
	return []dto.RevokeView{}, nil
}

func (s *FlowService) ApprovalTask(req dto.ApprovalTaskRequest) error {
	a, err := s.repo.GetAnswer(req.TaskID)
	if err != nil {
		return err
	}
	a.ExamExerciseType = req.Action
	_, err = s.repo.GetFlowEntry(a.ProjectID)
	return err
}

func (s *FlowService) Statics() (*dto.FlowStaticsView, error) {
	counts, err := s.repo.CountAnswersByStatus()
	if err != nil {
		return nil, err
	}
	return &dto.FlowStaticsView{
		PendingCount:  counts["pending"],
		RunningCount:  counts["running"],
		ApprovedCount: counts["approved"],
		RejectedCount: counts["rejected"],
	}, nil
}

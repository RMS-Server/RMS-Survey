package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/model"
	"github.com/surveyking/server/internal/repository"
	"gorm.io/gorm"
)

// FlowService handles workflow business logic.
type FlowService struct {
	repo *repository.FlowRepository
}

func NewFlowService(repo *repository.FlowRepository) *FlowService {
	return &FlowService{repo: repo}
}

// GetFlowEntry returns the flow definition view for a project.
// Project.Status: 0=inactive, 1=active, else=inactive.
func (s *FlowService) GetFlowEntry(projectID string) (*dto.FlowEntryView, error) {
	p, err := s.repo.GetFlowEntry(projectID)
	if err != nil {
		return nil, err
	}

	status := "inactive"
	if p.Status == 1 {
		status = "active"
	}

	view := &dto.FlowEntryView{
		ID:        p.ID,
		ProjectID: p.ID,
		Status:    status,
	}

	// Attempt to unmarshal project.Setting into FlowData/NodeConfig.
	if p.Setting != "" {
		var setting map[string]interface{}
		if err := json.Unmarshal([]byte(p.Setting), &setting); err == nil {
			view.FlowData = setting["flowData"]
			view.NodeConfig = setting["nodeConfig"]
		}
	}

	return view, nil
}

// GetFlowTasks returns a paginated list of flow tasks.
// ExamExerciseType IS the flow status field in this schema — it stores the
// workflow state (pending/running/approved/rejected/cancelled) for each answer.
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
			CreateBy:          a.CreateBy,
			CreatedAt:         a.CreatedAt.Format(time.RFC3339),
		})
	}
	return &dto.PageResponse[dto.FlowTaskView]{List: views, Total: total}, nil
}

// GetAuditRecord returns the approval history for a process instance.
func (s *FlowService) GetAuditRecord(processInstanceID string) ([]dto.FlowOperationView, error) {
	ops, err := s.repo.ListFlowOperations(processInstanceID)
	if err != nil {
		return nil, err
	}
	views := make([]dto.FlowOperationView, 0, len(ops))
	for _, op := range ops {
		views = append(views, dto.FlowOperationView{
			ID:                op.ID,
			ProcessInstanceID: op.ProcessInstanceID,
			OperatorID:        op.OperatorID,
			OperatorName:      op.OperatorName,
			Action:            op.Action,
			Comment:           op.Comment,
			CreatedAt:         op.CreatedAt.Format(time.RFC3339),
		})
	}
	return views, nil
}

// GetRevertNodes returns unique action nodes from the operation history as revert targets.
func (s *FlowService) GetRevertNodes(processInstanceID string) ([]dto.RevokeView, error) {
	ops, err := s.repo.ListFlowOperations(processInstanceID)
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{})
	views := make([]dto.RevokeView, 0)
	for _, op := range ops {
		if _, exists := seen[op.Action]; exists {
			continue
		}
		seen[op.Action] = struct{}{}
		views = append(views, dto.RevokeView{
			NodeID:   op.ID,
			NodeName: op.Action + " by " + op.OperatorName,
		})
	}
	return views, nil
}

// ApprovalTask processes an approval action atomically:
// updates the answer status and records the operation in one transaction.
func (s *FlowService) ApprovalTask(req dto.ApprovalTaskRequest) error {
	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewFlowRepository(tx)

		a, err := txRepo.GetAnswer(req.TaskID)
		if err != nil {
			return err
		}

		a.ExamExerciseType = req.Action
		if err := txRepo.SaveAnswer(a); err != nil {
			return err
		}

		op := &model.FlowOperation{
			BaseModelCreateOnly: model.BaseModelCreateOnly{
				ID:        uuid.New().String(),
				CreatedAt: time.Now(),
				CreateBy:  req.OperatorID,
			},
			ProcessInstanceID: a.ID,
			AnswerID:          a.ID,
			OperatorID:        req.OperatorID,
			OperatorName:      req.OperatorName,
			Action:            req.Action,
			Comment:           req.Comment,
		}
		return txRepo.CreateFlowOperation(op)
	})
}

// SaveFlow marshals FlowData+NodeConfig into project.Setting and persists it.
func (s *FlowService) SaveFlow(req dto.FlowEntryRequest) error {
	p, err := s.repo.GetFlowEntry(req.ProjectID)
	if err != nil {
		return err
	}

	setting := map[string]interface{}{
		"flowData":   req.FlowData,
		"nodeConfig": req.NodeConfig,
	}
	raw, err := json.Marshal(setting)
	if err != nil {
		return fmt.Errorf("marshal flow setting: %w", err)
	}
	p.Setting = string(raw)
	return s.repo.SaveProject(p)
}

// Deploy activates a project by setting its status to 1 (active).
func (s *FlowService) Deploy(projectID string) error {
	p, err := s.repo.GetFlowEntry(projectID)
	if err != nil {
		return err
	}
	p.Status = 1
	return s.repo.SaveProject(p)
}

// GetAnswer retrieves a single answer by ID.
func (s *FlowService) GetAnswer(id string) (*model.Answer, error) {
	return s.repo.GetAnswer(id)
}

// Statics returns flow task counts using the filtered CountAnswersByFlowStatus query.
func (s *FlowService) Statics() (*dto.FlowStaticsView, error) {
	counts, err := s.repo.CountAnswersByFlowStatus()
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

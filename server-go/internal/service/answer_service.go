package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/model"
	"github.com/surveyking/server/internal/repository"
	"github.com/xuri/excelize/v2"
)

// AnswerService handles business logic for answers.
type AnswerService struct {
	repo        *repository.AnswerRepo
	projectRepo *repository.ProjectRepo
}

// NewAnswerService creates a new AnswerService.
func NewAnswerService(repo *repository.AnswerRepo, projectRepo *repository.ProjectRepo) *AnswerService {
	return &AnswerService{repo: repo, projectRepo: projectRepo}
}

// ListAnswers returns a paginated list of answers.
func (s *AnswerService) ListAnswers(query *dto.AnswerQuery) (*dto.PageResponse[dto.AnswerView], error) {
	answers, total, err := s.repo.ListAnswers(query)
	if err != nil {
		return nil, err
	}
	views := make([]dto.AnswerView, 0, len(answers))
	for _, a := range answers {
		views = append(views, toAnswerView(a))
	}
	return &dto.PageResponse[dto.AnswerView]{List: views, Total: total}, nil
}

// ListDeleted returns soft-deleted answers.
func (s *AnswerService) ListDeleted(query *dto.AnswerQuery) ([]dto.AnswerView, error) {
	answers, err := s.repo.ListDeleted(query)
	if err != nil {
		return nil, err
	}
	views := make([]dto.AnswerView, 0, len(answers))
	for _, a := range answers {
		views = append(views, toAnswerView(a))
	}
	return views, nil
}

// GetAnswer returns a single answer.
func (s *AnswerService) GetAnswer(query *dto.AnswerQuery) (*dto.AnswerView, error) {
	a, err := s.repo.GetAnswer(query.ID)
	if err != nil {
		return nil, err
	}
	v := toAnswerView(*a)
	return &v, nil
}

// CreateAnswer inserts a new answer.
func (s *AnswerService) CreateAnswer(req *dto.AnswerRequest, userInfo *dto.UserInfo) error {
	id, err := nanoid.New()
	if err != nil {
		return fmt.Errorf("generate id: %w", err)
	}
	a := &model.Answer{
		BaseModel: model.BaseModel{
			ID:       id,
			CreateBy: userInfo.UserID,
		},
		ProjectID: req.ProjectID,
		Answer:    string(req.Answer),
		MetaInfo:  string(req.MetaInfo),
		TempSave:  req.TempSave,
	}
	return s.repo.CreateAnswer(a)
}

// UpdateAnswer updates an existing answer.
func (s *AnswerService) UpdateAnswer(req *dto.AnswerRequest, userInfo *dto.UserInfo) error {
	a, err := s.repo.GetAnswer(req.ID)
	if err != nil {
		return err
	}
	if req.Answer != nil {
		a.Answer = string(req.Answer)
	}
	if req.MetaInfo != nil {
		a.MetaInfo = string(req.MetaInfo)
	}
	a.UpdateBy = userInfo.UserID
	return s.repo.UpdateAnswer(a)
}

// DeleteAnswer soft-deletes an answer.
func (s *AnswerService) DeleteAnswer(req *dto.AnswerRequest) error {
	ids := req.IDs
	if len(ids) == 0 && req.ID != "" {
		ids = []string{req.ID}
	}
	for _, id := range ids {
		if err := s.repo.SoftDeleteAnswer(id); err != nil {
			return err
		}
	}
	return nil
}

// DestroyAnswer permanently deletes answers.
func (s *AnswerService) DestroyAnswer(req *dto.AnswerRequest) error {
	ids := req.IDs
	if len(ids) == 0 && req.ID != "" {
		ids = []string{req.ID}
	}
	for _, id := range ids {
		if err := s.repo.HardDeleteAnswer(id); err != nil {
			return err
		}
	}
	return nil
}

// RestoreAnswer recovers answers from trash.
func (s *AnswerService) RestoreAnswer(req *dto.AnswerRequest) error {
	ids := req.IDs
	if len(ids) == 0 && req.ID != "" {
		ids = []string{req.ID}
	}
	for _, id := range ids {
		if err := s.repo.RestoreAnswer(id); err != nil {
			return err
		}
	}
	return nil
}

// ExportAnswers exports answers to an Excel file with dynamic question headers.
func (s *AnswerService) ExportAnswers(query *dto.DownloadQuery) ([]byte, string, error) {
	project, err := s.projectRepo.GetProject(query.ProjectID)
	if err != nil {
		return nil, "", fmt.Errorf("get project: %w", err)
	}

	// Parse survey JSON to extract question order and titles.
	var schema surveySchema
	if project.Survey != "" {
		if err := json.Unmarshal([]byte(project.Survey), &schema); err != nil {
			return nil, "", fmt.Errorf("parse survey schema: %w", err)
		}
	}

	type questionMeta struct {
		id    string
		title string
	}
	var questions []questionMeta
	for _, page := range schema.Pages {
		for _, el := range page.Elements {
			questions = append(questions, questionMeta{id: el.ID, title: el.Title})
		}
	}

	answers, err := s.repo.ListByProjectID(query.ProjectID)
	if err != nil {
		return nil, "", err
	}

	f := excelize.NewFile()
	sheet := "Sheet1"

	// Build dynamic headers: fixed columns first, then one column per question.
	headers := make([]string, 0, 2+len(questions))
	headers = append(headers, "ID", "Submit Time")
	for _, q := range questions {
		title := q.title
		if title == "" {
			title = q.id
		}
		headers = append(headers, title)
	}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Write one row per answer.
	for rowIdx, a := range answers {
		row := rowIdx + 2

		// Parse answer JSON into a flat map keyed by question ID.
		var answerMap map[string]interface{}
		if a.Answer != "" {
			_ = json.Unmarshal([]byte(a.Answer), &answerMap)
		}

		values := make([]interface{}, 0, 2+len(questions))
		values = append(values, a.ID, a.CreatedAt.Format("2006-01-02 15:04:05"))
		for _, q := range questions {
			if answerMap != nil {
				values = append(values, formatAnswerValue(answerMap[q.id]))
			} else {
				values = append(values, "")
			}
		}

		for colIdx, v := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
			f.SetCellValue(sheet, cell, v)
		}
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", fmt.Errorf("write excel: %w", err)
	}

	filename := fmt.Sprintf("answers_%s_%s.xlsx", query.ProjectID, time.Now().Format("20060102150405"))
	return buf.Bytes(), filename, nil
}

// formatAnswerValue converts an answer value to a string suitable for Excel.
func formatAnswerValue(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case []interface{}:
		// multi-select: join with comma
		parts := make([]string, 0, len(val))
		for _, item := range val {
			parts = append(parts, fmt.Sprintf("%v", item))
		}
		result := ""
		for i, p := range parts {
			if i > 0 {
				result += ","
			}
			result += p
		}
		return result
	default:
		return fmt.Sprintf("%v", val)
	}
}

func toAnswerView(a model.Answer) dto.AnswerView {
	var answerJSON json.RawMessage
	if a.Answer != "" {
		answerJSON = json.RawMessage(a.Answer)
	}
	var metaJSON json.RawMessage
	if a.MetaInfo != "" {
		metaJSON = json.RawMessage(a.MetaInfo)
	}
	return dto.AnswerView{
		ID:        a.ID,
		ProjectID: a.ProjectID,
		Answer:    answerJSON,
		MetaInfo:  metaJSON,
		TempSave:  a.TempSave,
		ExamScore: a.ExamScore,
		CreateBy:  a.CreateBy,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
		UpdatedAt: a.UpdatedAt.Format(time.RFC3339),
	}
}

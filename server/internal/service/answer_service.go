package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/model"
	"github.com/rms-survey/server/internal/repository"
	"github.com/xuri/excelize/v2"
)

// uploadSurveySchema represents the survey structure for Excel import.
// It uses Pages structure for compatibility with export function.
type uploadSurveySchema struct {
	ID    string                   `json:"id,omitempty"`
	Title string                   `json:"title,omitempty"`
	Pages []uploadSurveyPage       `json:"pages,omitempty"`
}

// uploadSurveyElement represents a question element for Excel import.
type uploadSurveyElement struct {
	ID       string                `json:"id"`
	Title    string                `json:"title"`
	Type     string                `json:"type"`
	Children []uploadSurveyElement `json:"children,omitempty"`
}

// uploadSurveyPage represents a page in the survey for Excel import.
type uploadSurveyPage struct {
	Elements []uploadSurveyElement `json:"elements"`
}

// AnswerService handles business logic for answers.
type AnswerService struct {
	repo        *repository.AnswerRepo
	projectRepo *repository.ProjectRepo
	projectSvc  *ProjectService
}

// NewAnswerService creates a new AnswerService.
func NewAnswerService(repo *repository.AnswerRepo, projectRepo *repository.ProjectRepo) *AnswerService {
	return &AnswerService{repo: repo, projectRepo: projectRepo, projectSvc: NewProjectService(projectRepo)}
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

// UploadAnswers parses an Excel file and imports answers.
// If projectId is provided, it matches Excel columns to existing survey questions.
// If autoSchema is true, it creates a new project based on the Excel header.
func (s *AnswerService) UploadAnswers(projectID string, autoSchema bool, parentID string, fileData []byte, filename string, userInfo *dto.UserInfo) (*dto.AnswerImportResult, error) {
	f, err := excelize.OpenReader(bytes.NewReader(fileData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse Excel: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("Excel file has no sheets")
	}
	sheet := sheets[0]
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel must have at least a header row and one data row")
	}

	// Extract name from filename (without extension)
	name := filename
	if idx := len(filename) - 1; idx > 0 {
		for i := idx; i >= 0; i-- {
			if filename[i] == '.' {
				name = filename[:i]
				break
			}
		}
	}

	headerRow := rows[0]
	var schema uploadSurveySchema
	var resultProjectID string

	if projectID != "" {
		// Use existing project - filter schema by Excel columns
		project, err := s.projectRepo.GetProject(projectID)
		if err != nil {
			return nil, fmt.Errorf("get project: %w", err)
		}
		resultProjectID = projectID

		var existingSchema uploadSurveySchema
		if project.Survey != "" {
			if err := json.Unmarshal([]byte(project.Survey), &existingSchema); err != nil {
				return nil, fmt.Errorf("parse survey schema: %w", err)
			}
		}

		// Build question map for quick lookup by title
		questionMap := make(map[string]uploadSurveyElement)
		for _, page := range existingSchema.Pages {
			for _, el := range page.Elements {
				questionMap[el.Title] = el
			}
		}

		// Filter schema: only include questions that match Excel columns
		// Build a single page with matching elements
		schema.ID = existingSchema.ID
		var elements []uploadSurveyElement
		for _, title := range headerRow {
			if el, ok := questionMap[title]; ok {
				elements = append(elements, el)
			} else {
				// Create placeholder element for unmatched column
				id, _ := nanoid.New()
				childID, _ := nanoid.New()
				elements = append(elements, uploadSurveyElement{
					ID:    id,
					Title: title,
					Type:  "fillBlank",
					Children: []uploadSurveyElement{{ID: childID}},
				})
			}
		}
		schema.Pages = []uploadSurveyPage{{Elements: elements}}
	} else if autoSchema {
		// Create new project from Excel header
		schema = createSurveyFromExcelHeader(headerRow)
		schema.Title = name

		setting := `{"mode":"survey","status":1}`
		id, _ := nanoid.New()
		p := &model.Project{
			BaseModel: model.BaseModel{ID: id, CreateBy: userInfo.UserID},
			ParentID:  parentID,
			Name:      name,
			Survey:    mustMarshalJSON(schema),
			Setting:   setting,
			Mode:      "survey",
			Status:    1,
			Priority:  1000,
		}
		if err := s.projectRepo.CreateProject(p); err != nil {
			return nil, fmt.Errorf("create project: %w", err)
		}
		resultProjectID = id
	} else {
		return nil, fmt.Errorf("either projectId or autoSchema must be provided")
	}

	// Parse data rows into answers
	// Get elements from the first page
	var elements []uploadSurveyElement
	if len(schema.Pages) > 0 {
		elements = schema.Pages[0].Elements
	}

	answers := make([]model.Answer, 0, len(rows)-1)
	for rowIdx, row := range rows {
		if rowIdx == 0 {
			continue // skip header
		}

		answerMap := make(map[string]interface{})
		for colIdx, title := range headerRow {
			if colIdx >= len(row) {
				continue
			}
			cellValue := row[colIdx]
			if cellValue == "" {
				continue
			}

			// Find the question element for this column
			var questionID, optionID string
			for _, el := range elements {
				if el.Title == title {
					questionID = el.ID
					if len(el.Children) > 0 {
						optionID = el.Children[0].ID
					}
					break
				}
			}

			if questionID == "" {
				continue
			}

			if optionID == "" {
				optionID = questionID
			}

			// Store answer as map[optionID]value format for compatibility with Java backend
			answerMap[questionID] = map[string]string{optionID: cellValue}
		}

		if len(answerMap) == 0 {
			continue // skip empty rows
		}

		id, _ := nanoid.New()
		answerJSON, _ := json.Marshal(answerMap)
		answers = append(answers, model.Answer{
			BaseModel: model.BaseModel{
				ID:       id,
				CreateBy: userInfo.UserID,
			},
			ProjectID: resultProjectID,
			Answer:    string(answerJSON),
		})
	}

	if len(answers) > 0 {
		if err := s.repo.CreateAnswers(answers); err != nil {
			return nil, fmt.Errorf("save answers: %w", err)
		}
	}

	schemaJSON, _ := json.Marshal(schema)
	return &dto.AnswerImportResult{
		ProjectID: resultProjectID,
		Schema:    schemaJSON,
	}, nil
}

// createSurveyFromExcelHeader creates a survey schema from Excel header row.
// Uses Pages structure for compatibility with export function.
func createSurveyFromExcelHeader(headers []string) uploadSurveySchema {
	var elements []uploadSurveyElement
	usedIDs := make(map[string]bool)

	for _, title := range headers {
		id := generateNanoID(8, usedIDs) // Use 8 chars for better uniqueness
		childID := generateNanoID(8, usedIDs)
		elements = append(elements, uploadSurveyElement{
			ID:       id,
			Title:    title,
			Type:     "fillBlank",
			Children: []uploadSurveyElement{{ID: childID}},
		})
	}

	return uploadSurveySchema{
		Pages: []uploadSurveyPage{{Elements: elements}},
	}
}

// generateNanoID generates a unique nano ID of specified length.
func generateNanoID(length int, used map[string]bool) string {
	for {
		id, _ := nanoid.New()
		if len(id) > length {
			id = id[:length]
		}
		if !used[id] {
			used[id] = true
			return id
		}
	}
}

// mustMarshalJSON marshals to JSON or returns empty object on error.
func mustMarshalJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(b)
}

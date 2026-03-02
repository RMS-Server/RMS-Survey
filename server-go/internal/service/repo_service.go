package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	mathrand "math/rand"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/model"
	"github.com/surveyking/server/internal/repository"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// RepoService handles question repo business logic.
type RepoService struct {
	repoRepo     *repository.RepoRepository
	userBookRepo *repository.UserBookRepository
}

func NewRepoService(repoRepo *repository.RepoRepository, userBookRepo *repository.UserBookRepository) *RepoService {
	return &RepoService{repoRepo: repoRepo, userBookRepo: userBookRepo}
}

func (s *RepoService) List(query dto.RepoQuery) (*dto.PageResponse[dto.RepoView], error) {
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	pageIndex := query.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	offset := (pageIndex - 1) * pageSize

	items, total, err := s.repoRepo.List(query.Name, query.Category, query.Mode, offset, pageSize)
	if err != nil {
		return nil, err
	}
	views := make([]dto.RepoView, 0, len(items))
	for _, r := range items {
		views = append(views, dto.RepoView{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Category:    r.Category,
			Mode:        r.Mode,
			Shared:      r.Shared,
			Tag:         r.Tag,
			Priority:    r.Priority,
			IsPractice:  r.IsPractice,
		})
	}
	return &dto.PageResponse[dto.RepoView]{List: views, Total: total}, nil
}

func (s *RepoService) Create(req dto.RepoRequest) error {
	now := time.Now()
	r := &model.Repo{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Mode:        req.Mode,
		Shared:      req.Shared,
		Tag:         req.Tag,
		Priority:    req.Priority,
		IsPractice:  req.IsPractice,
	}
	r.ID = uuid.New().String()
	r.CreatedAt = now
	r.UpdatedAt = now
	return s.repoRepo.Create(r)
}

func (s *RepoService) Update(req dto.RepoRequest) error {
	r, err := s.repoRepo.GetByID(req.ID)
	if err != nil {
		return err
	}
	r.Name = req.Name
	r.Description = req.Description
	r.Category = req.Category
	r.Mode = req.Mode
	r.Shared = req.Shared
	r.Tag = req.Tag
	r.Priority = req.Priority
	r.IsPractice = req.IsPractice
	return s.repoRepo.Update(r)
}

func (s *RepoService) Delete(id string) error {
	return s.repoRepo.Delete(id)
}

func (s *RepoService) ListUserBook(query dto.UserBookQuery) (*dto.PageResponse[dto.UserBookView], error) {
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	pageIndex := query.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	offset := (pageIndex - 1) * pageSize

	items, total, err := s.userBookRepo.List(query.RepoID, query.Status, offset, pageSize)
	if err != nil {
		return nil, err
	}
	views := make([]dto.UserBookView, 0, len(items))
	for _, ub := range items {
		views = append(views, dto.UserBookView{
			ID:           ub.ID,
			Name:         ub.Name,
			TemplateID:   ub.TemplateID,
			WrongTimes:   ub.WrongTimes,
			CorrectTimes: ub.CorrectTimes,
			Note:         ub.Note,
			Status:       ub.Status,
			Type:         ub.Type,
			RepoID:       ub.RepoID,
			IsMarked:     ub.IsMarked,
		})
	}
	return &dto.PageResponse[dto.UserBookView]{List: views, Total: total}, nil
}

func (s *RepoService) CreateUserBook(req dto.UserBookRequest) error {
	now := time.Now()
	ub := &model.UserBook{
		TemplateID: req.TemplateID,
		RepoID:     req.RepoID,
		Note:       req.Note,
		Status:     req.Status,
		IsMarked:   req.IsMarked,
	}
	ub.ID = uuid.New().String()
	ub.CreatedAt = now
	ub.UpdatedAt = now
	return s.userBookRepo.Create(ub)
}

func (s *RepoService) UpdateUserBook(req dto.UserBookRequest) (*dto.UserBookView, error) {
	ub, err := s.userBookRepo.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	ub.Note = req.Note
	ub.Status = req.Status
	ub.IsMarked = req.IsMarked
	if err := s.userBookRepo.Update(ub); err != nil {
		return nil, err
	}
	return &dto.UserBookView{
		ID:           ub.ID,
		Name:         ub.Name,
		TemplateID:   ub.TemplateID,
		WrongTimes:   ub.WrongTimes,
		CorrectTimes: ub.CorrectTimes,
		Note:         ub.Note,
		Status:       ub.Status,
		Type:         ub.Type,
		RepoID:       ub.RepoID,
		IsMarked:     ub.IsMarked,
	}, nil
}

func (s *RepoService) DeleteUserBook(id string) error {
	return s.userBookRepo.Delete(id)
}

// BatchCreate adds multiple template IDs to a repo.
func (s *RepoService) BatchCreate(req dto.RepoTemplateRequest) error {
	return s.repoRepo.BatchAddTemplates(req.RepoID, req.TemplateIDs)
}

// Unbind removes template IDs from a repo.
func (s *RepoService) Unbind(req dto.RepoTemplateRequest) error {
	return s.repoRepo.RemoveTemplates(req.RepoID, req.TemplateIDs)
}

// PickQuestions randomly selects questions from a repo based on conditions.
func (s *RepoService) PickQuestions(conditions []dto.RandomSurveyCondition) ([]interface{}, error) {
	result := make([]interface{}, 0)
	for _, cond := range conditions {
		if cond.RepoID == "" {
			continue
		}
		items, err := s.repoRepo.ListTemplatesByRepoID(cond.RepoID, cond.Types)
		if err != nil {
			continue
		}
		picked := make([]model.Template, len(items))
		copy(picked, items)
		// Fisher-Yates shuffle
		for i := len(picked) - 1; i > 0; i-- {
			j := mathrand.Intn(i + 1)
			picked[i], picked[j] = picked[j], picked[i]
		}
		if cond.QuestionsNum > 0 && cond.QuestionsNum < len(picked) {
			picked = picked[:cond.QuestionsNum]
		}
		for _, item := range picked {
			result = append(result, item)
		}
	}
	return result, nil
}

// templateSchema is the per-question JSON stored in t_template.template column.
type templateSchema struct {
	ID       string           `json:"id"`
	Title    string           `json:"title"`
	Type     string           `json:"type"`
	Children []templateSchema `json:"children"`
	Tags     []string         `json:"tags"`
	Attribute *templateAttr  `json:"attribute"`
}

type templateAttr struct {
	ExamScore          *float64 `json:"examScore"`
	ExamAnalysis       string   `json:"examAnalysis"`
	ExamCorrectAnswer  string   `json:"examCorrectAnswer"`
	ExamAnswerMode     string   `json:"examAnswerMode"`
}

// Export returns a multi-sheet Excel file with questions grouped by type.
func (s *RepoService) Export(query dto.RepoQuery) ([]byte, error) {
	types := []string{"Radio", "Checkbox", "Judge", "FillBlank", "Textarea"}
	templates, err := s.repoRepo.ListTemplatesByRepoID(query.RepoID, types)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	defer f.Close()

	// Group by question type
	byType := map[string][]model.Template{}
	for _, t := range templates {
		byType[t.QuestionType] = append(byType[t.QuestionType], t)
	}

	sheetDefs := []struct {
		qtype   string
		name    string
		headers []string
	}{
		{"Radio", "单选题", []string{"序号", "题干", "选项A", "选项B", "选项C", "选项D", "选项E", "选项F", "选项G", "选项H", "解析", "分数", "答案", "标签"}},
		{"Checkbox", "多选题", []string{"序号", "题干", "选项A", "选项B", "选项C", "选项D", "选项E", "选项F", "选项G", "选项H", "解析", "分数", "答案", "标签"}},
		{"Judge", "判断题", []string{"序号", "题干", "选项A", "选项B", "解析", "分数", "答案", "标签"}},
		{"FillBlank", "填空题", []string{"序号", "题干", "空1", "空2", "空3", "空4", "空5", "空6", "空7", "空8", "解析", "单空分数", "标签"}},
		{"Textarea", "简答题", []string{"序号", "题干", "答案", "解析", "分数", "标签"}},
	}

	createdSheets := 0
	for _, def := range sheetDefs {
		rows := byType[def.qtype]
		if len(rows) == 0 {
			continue
		}

		var sheetName string
		if createdSheets == 0 {
			// Rename the default "Sheet1"
			_ = f.SetSheetName("Sheet1", def.name)
			sheetName = def.name
		} else {
			_, _ = f.NewSheet(def.name)
			sheetName = def.name
		}
		createdSheets++

		// Write headers
		for col, h := range def.headers {
			cell, _ := excelize.CoordinatesToCellName(col+1, 1)
			_ = f.SetCellValue(sheetName, cell, h)
		}

		// Write data rows
		writtenRow := 0
		for _, t := range rows {
			var schema templateSchema
			if err := json.Unmarshal([]byte(t.TemplateData), &schema); err != nil {
				continue
			}
			writtenRow++
			rowData := buildExcelRow(def.qtype, writtenRow, t.Tag, schema)
			for col, val := range rowData {
				cell, _ := excelize.CoordinatesToCellName(col+1, writtenRow+1)
				_ = f.SetCellValue(sheetName, cell, val)
			}
		}
	}

	// If no sheets were created (no data), keep the default empty sheet
	if createdSheets == 0 {
		_ = f.SetSheetName("Sheet1", "题库")
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// buildExcelRow converts a template schema into a row of cell values for the given question type.
func buildExcelRow(qtype string, idx int, tag string, s templateSchema) []interface{} {
	attr := s.Attribute
	score := ""
	analysis := ""
	if attr != nil {
		if attr.ExamScore != nil {
			score = fmt.Sprintf("%g", *attr.ExamScore)
		}
		analysis = attr.ExamAnalysis
	}

	switch qtype {
	case "Radio", "Checkbox":
		opts := make([]string, 8)
		var answers []string
		for i, child := range s.Children {
			if i >= 8 {
				break
			}
			opts[i] = child.Title
			if child.Attribute != nil && child.Attribute.ExamCorrectAnswer != "" {
				answers = append(answers, string(rune('A'+i)))
			}
		}
		row := []interface{}{idx, s.Title}
		for _, o := range opts {
			row = append(row, o)
		}
		row = append(row, analysis, score, strings.Join(answers, ","), tag)
		return row

	case "Judge":
		optA, optB, answer := "", "", ""
		for i, child := range s.Children {
			if i == 0 {
				optA = child.Title
				if child.Attribute != nil && child.Attribute.ExamCorrectAnswer != "" {
					answer = "A"
				}
			} else if i == 1 {
				optB = child.Title
				if answer == "" && child.Attribute != nil && child.Attribute.ExamCorrectAnswer != "" {
					answer = "B"
				}
			}
		}
		return []interface{}{idx, s.Title, optA, optB, analysis, score, answer, tag}

	case "FillBlank":
		blanks := make([]string, 8)
		for i, child := range s.Children {
			if i >= 8 {
				break
			}
			if child.Attribute != nil {
				blanks[i] = child.Attribute.ExamCorrectAnswer
			}
		}
		row := []interface{}{idx, s.Title}
		for _, b := range blanks {
			row = append(row, b)
		}
		row = append(row, analysis, score, tag)
		return row

	case "Textarea":
		answer := ""
		if len(s.Children) > 0 && s.Children[0].Attribute != nil {
			answer = s.Children[0].Attribute.ExamCorrectAnswer
		}
		return []interface{}{idx, s.Title, answer, analysis, score, tag}
	}
	return nil
}

// ImportFromExcel parses an uploaded Excel file and creates Template records in the repo.
// All inserts are wrapped in a single transaction; on any error the whole import is rolled back.
func (s *RepoService) ImportFromExcel(repoID string, file multipart.File) error {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return fmt.Errorf("failed to open excel: %w", err)
	}
	defer f.Close()

	sheetTypeMap := map[string]string{
		"单选题": "Radio",
		"多选题": "Checkbox",
		"判断题": "Judge",
		"填空题": "FillBlank",
		"简答题": "Textarea",
	}

	// Collect all templates first, then insert in one transaction.
	var toInsert []*model.Template
	now := time.Now()

	for sheetName, qtype := range sheetTypeMap {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			continue
		}
		if len(rows) < 2 {
			continue
		}

		headerIdx := map[string]int{}
		for col, h := range rows[0] {
			headerIdx[h] = col
		}

		serialNo := 0
		for _, row := range rows[1:] {
			t, err := parseExcelRow(qtype, serialNo+1, row, headerIdx)
			if err != nil || t == nil {
				continue
			}
			serialNo++
			t.RepoID = repoID
			t.ID = uuid.New().String()
			t.CreatedAt = now
			t.UpdatedAt = now
			toInsert = append(toInsert, t)
		}
	}

	if len(toInsert) == 0 {
		return nil
	}

	return s.repoRepo.DB().Transaction(func(tx *gorm.DB) error {
		for _, t := range toInsert {
			if err := tx.Create(t).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// parseExcelRow converts a spreadsheet row into a Template model.
func parseExcelRow(qtype string, serialNo int, row []string, hdr map[string]int) (*model.Template, error) {
	get := func(key string) string {
		idx, ok := hdr[key]
		if !ok || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}

	title := get("题干")
	if title == "" {
		return nil, nil
	}

	schema := templateSchema{
		ID:    uuid.New().String(),
		Title: title,
		Type:  qtype,
		Attribute: &templateAttr{
			ExamAnalysis: get("解析"),
		},
	}

	scoreStr := get("分数")
	if scoreStr == "" {
		scoreStr = get("单空分数")
	}
	if scoreStr != "" {
		var v float64
		if _, err := fmt.Sscanf(scoreStr, "%f", &v); err == nil {
			schema.Attribute.ExamScore = &v
		}
	}

	switch qtype {
	case "Radio", "Checkbox":
		schema.Attribute.ExamAnswerMode = "onlyOne"
		if qtype == "Checkbox" {
			schema.Attribute.ExamAnswerMode = "multiple"
		}
		answerStr := get("答案")
		correctSet := map[rune]bool{}
		for _, ch := range strings.ToUpper(answerStr) {
			if ch >= 'A' && ch <= 'H' {
				correctSet[ch] = true
			}
		}
		optLabels := []string{"选项A", "选项B", "选项C", "选项D", "选项E", "选项F", "选项G", "选项H"}
		for i, label := range optLabels {
			text := get(label)
			if text == "" {
				break
			}
			child := templateSchema{
				ID:    uuid.New().String(),
				Title: text,
				Attribute: &templateAttr{},
			}
			if correctSet[rune('A'+i)] {
				child.Attribute.ExamCorrectAnswer = "true"
			}
			schema.Children = append(schema.Children, child)
		}

	case "Judge":
		answerStr := strings.ToUpper(get("答案"))
		for i, label := range []string{"选项A", "选项B"} {
			text := get(label)
			if text == "" {
				break
			}
			child := templateSchema{
				ID:    uuid.New().String(),
				Title: text,
				Attribute: &templateAttr{},
			}
			if (i == 0 && answerStr == "A") || (i == 1 && answerStr == "B") {
				child.Attribute.ExamCorrectAnswer = "true"
			}
			schema.Children = append(schema.Children, child)
		}

	case "FillBlank":
		for i := 1; i <= 8; i++ {
			ans := get(fmt.Sprintf("空%d", i))
			if ans == "" {
				break
			}
			schema.Children = append(schema.Children, templateSchema{
				ID:    uuid.New().String(),
				Attribute: &templateAttr{ExamCorrectAnswer: ans},
			})
		}

	case "Textarea":
		ans := get("答案")
		if ans != "" {
			schema.Children = []templateSchema{{
				ID:        uuid.New().String(),
				Attribute: &templateAttr{ExamCorrectAnswer: ans},
			}}
		}
	}

	data, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	return &model.Template{
		SerialNo:     fmt.Sprintf("%d", serialNo),
		Name:         title,
		QuestionType: qtype,
		TemplateData: string(data),
		Mode:         "exam",
	}, nil
}

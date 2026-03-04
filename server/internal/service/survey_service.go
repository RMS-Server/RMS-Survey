package service

import (
	"encoding/json"

	"github.com/rms-survey/server/internal/config"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/model"
	"github.com/rms-survey/server/internal/repository"
	"gorm.io/gorm"
)

// SurveyService handles public survey loading and answer submission.
type SurveyService struct {
	projectRepo *repository.ProjectRepo
	answerRepo  *repository.AnswerRepo
	systemRepo  *repository.SystemRepo
}

// NewSurveyService creates a new SurveyService.
func NewSurveyService(projectRepo *repository.ProjectRepo, answerRepo *repository.AnswerRepo, db *gorm.DB) *SurveyService {
	return &SurveyService{
		projectRepo: projectRepo,
		answerRepo:  answerRepo,
		systemRepo:  repository.NewSystemRepo(db),
	}
}

// LoadProject loads a survey for public answering by project code (ID).
func (s *SurveyService) LoadProject(req *dto.SurveyLoadRequest) (*dto.SurveyView, error) {
	p, err := s.projectRepo.GetProject(req.Code)
	if err != nil {
		return nil, err
	}
	return &dto.SurveyView{
		ID:      p.ID,
		Name:    p.Name,
		Survey:  json.RawMessage(p.Survey),
		Setting: json.RawMessage(p.Setting),
		Status:  p.Status,
		Mode:    p.Mode,
	}, nil
}

// SaveAnswer submits a public answer.
func (s *SurveyService) SaveAnswer(req *dto.AnswerRequest) (*dto.PublicAnswerView, error) {
	svc := NewAnswerService(s.answerRepo, s.projectRepo)
	if err := svc.CreatePublicAnswer(req); err != nil {
		return nil, err
	}
	// Return setting from project so frontend can show end page
	p, err := s.projectRepo.GetProject(req.ProjectID)
	if err != nil {
		return nil, err
	}
	return &dto.PublicAnswerView{
		ID:      req.ID,
		Setting: json.RawMessage(p.Setting),
	}, nil
}

// TempSaveAnswer temporarily saves an answer (in-progress).
func (s *SurveyService) TempSaveAnswer(req *dto.AnswerRequest) error {
	tempSave := 1
	req.TempSave = &tempSave
	svc := NewAnswerService(s.answerRepo, s.projectRepo)
	return svc.CreatePublicAnswer(req)
}

// GetSetting returns the survey setting for a project.
func (s *SurveyService) GetSetting(projectID string) (json.RawMessage, error) {
	p, err := s.projectRepo.GetProject(projectID)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(p.Setting), nil
}

// UpdateSetting updates the survey setting for a project.
func (s *SurveyService) UpdateSetting(req *dto.SurveySettingRequest, userInfo *dto.UserInfo) error {
	p, err := s.projectRepo.GetProject(req.ProjectID)
	if err != nil {
		return err
	}
	p.Setting = string(req.Setting)
	p.UpdateBy = userInfo.UserID
	return s.projectRepo.UpdateProject(p)
}

// GetLogic returns the survey logic (survey JSON) for a project.
func (s *SurveyService) GetLogic(projectID string) (json.RawMessage, error) {
	p, err := s.projectRepo.GetProject(projectID)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(p.Survey), nil
}

// UpdateLogic updates the survey logic (survey JSON) for a project.
func (s *SurveyService) UpdateLogic(req *dto.SurveyLogicRequest, userInfo *dto.UserInfo) error {
	p, err := s.projectRepo.GetProject(req.ProjectID)
	if err != nil {
		return err
	}
	p.Survey = string(req.Survey)
	p.UpdateBy = userInfo.UserID
	return s.projectRepo.UpdateProject(p)
}

// ValidateProject validates a survey (password check, status check).
func (s *SurveyService) ValidateProject(req *dto.SurveyLoadRequest) (*dto.SurveyView, error) {
	return s.LoadProject(req)
}

// StatProject counts how many times each option was selected across all submitted answers.
// Returns Stats: map[fieldId]map[optionValue]count.
func (s *SurveyService) StatProject(req *dto.SurveyLoadRequest) (*dto.PublicStatisticsView, error) {
	p, err := s.projectRepo.GetProject(req.Code)
	if err != nil {
		return nil, err
	}

	answers, err := s.answerRepo.ListByProjectID(p.ID)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]interface{})
	for _, a := range answers {
		if a.Answer == "" {
			continue
		}
		// Answer JSON is map[fieldId]interface{} where value may be a string or []string
		var answerMap map[string]interface{}
		if err := json.Unmarshal([]byte(a.Answer), &answerMap); err != nil {
			continue
		}
		for fieldID, val := range answerMap {
			existing, ok := stats[fieldID]
			var counts map[string]int
			if ok {
				// Safe assertion: we always store map[string]int for this key.
				counts, ok = existing.(map[string]int)
				if !ok {
					counts = make(map[string]int)
				}
			} else {
				counts = make(map[string]int)
			}
			stats[fieldID] = counts

			switch v := val.(type) {
			case string:
				if v != "" {
					counts[v]++
				}
			case []interface{}:
				for _, item := range v {
					if s, ok := item.(string); ok && s != "" {
						counts[s]++
					}
				}
			}
		}
	}

	return &dto.PublicStatisticsView{ProjectID: p.ID, Stats: stats}, nil
}

// LoadQuery returns the public query verify view including the survey schema.
func (s *SurveyService) LoadQuery(req *dto.PublicQueryRequest) (*dto.PublicQueryVerifyView, error) {
	p, err := s.projectRepo.GetProject(req.ProjectID)
	if err != nil {
		return nil, err
	}
	return &dto.PublicQueryVerifyView{
		ProjectID: p.ID,
		Name:      p.Name,
		Survey:    json.RawMessage(p.Survey),
	}, nil
}

// GetQueryResult returns submitted answers for a project, optionally filtered by answer code.
func (s *SurveyService) GetQueryResult(req *dto.PublicQueryRequest) (*dto.PublicQueryView, error) {
	answers, err := s.answerRepo.ListByProjectID(req.ProjectID)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, 0, len(answers))
	for _, a := range answers {
		result = append(result, map[string]interface{}{
			"id":     a.ID,
			"answer": json.RawMessage(a.Answer),
		})
	}

	return &dto.PublicQueryView{ProjectID: req.ProjectID, Answers: result}, nil
}

// LoadDict returns dictionary items for the requested dict codes.
func (s *SurveyService) LoadDict(req *dto.PublicDictRequest) ([]dto.PublicDictView, error) {
	if len(req.Codes) == 0 {
		return []dto.PublicDictView{}, nil
	}

	items, err := s.systemRepo.ListDictItemsByCodes(req.Codes)
	if err != nil {
		return nil, err
	}

	views := make([]dto.PublicDictView, 0, len(items))
	for _, item := range items {
		views = append(views, dto.PublicDictView{
			Code:  item.DictCode,
			Label: item.ItemName,
			Value: item.ItemValue,
		})
	}
	return views, nil
}

// examInfo is the JSON structure stored in answer.exam_info.
type examInfo struct {
	MaxScore float64 `json:"maxScore"`
	Passed   bool    `json:"passed"`
}

// LoadExamResult reads exam score and pass status from the stored answer.
func (s *SurveyService) LoadExamResult(req *dto.PublicExamRequest) (*dto.PublicExamResult, error) {
	a, err := s.answerRepo.GetAnswer(req.AnswerID)
	if err != nil {
		return nil, err
	}

	var score float64
	if a.ExamScore != nil {
		score = float64(*a.ExamScore)
	}

	var info examInfo
	if a.ExamInfo != "" {
		_ = json.Unmarshal([]byte(a.ExamInfo), &info)
	}

	return &dto.PublicExamResult{
		Score:    score,
		MaxScore: info.MaxScore,
		Passed:   info.Passed,
	}, nil
}

// LoadLinkResult returns the answer data for a linked survey answer.
func (s *SurveyService) LoadLinkResult(req *dto.PublicLinkRequest) (*dto.PublicLinkResult, error) {
	a, err := s.answerRepo.GetAnswer(req.AnswerID)
	if err != nil {
		return nil, err
	}

	var data interface{}
	if a.Answer != "" {
		var raw json.RawMessage = json.RawMessage(a.Answer)
		data = raw
	}

	return &dto.PublicLinkResult{ProjectID: req.ProjectID, Data: data}, nil
}

// GetProjectByID returns a project by ID.
func (s *SurveyService) GetProjectByID(id string) (*model.Project, error) {
	return s.projectRepo.GetProject(id)
}

// GetQuestionAttachmentConfig extracts attachment config for a question.
func (s *SurveyService) GetQuestionAttachmentConfig(project *model.Project, questionID string) *UploadValidationConfig {
	// Default config from global settings
	defaultCfg := &UploadValidationConfig{
		MaxSize:      config.C.Upload.MaxSize,
		AllowedTypes: config.C.Upload.AllowedTypes,
	}

	if project == nil || project.Survey == "" {
		return defaultCfg
	}

	// Parse survey schema
	var schema struct {
		Pages []struct {
			Elements []struct {
				ID         string `json:"id"`
				Attachment *struct {
					Enabled      bool     `json:"enabled"`
					MaxFiles     int      `json:"maxFiles"`
					MaxSize      int64    `json:"maxSize"`
					AllowedTypes []string `json:"allowedTypes"`
				} `json:"attachment"`
			} `json:"elements"`
		} `json:"pages"`
	}

	if err := json.Unmarshal([]byte(project.Survey), &schema); err != nil {
		return defaultCfg
	}

	// Find the question
	for _, page := range schema.Pages {
		for _, elem := range page.Elements {
			if elem.ID == questionID && elem.Attachment != nil && elem.Attachment.Enabled {
				// Use question-specific config, falling back to defaults for zero values
				maxSize := elem.Attachment.MaxSize
				if maxSize <= 0 {
					maxSize = config.C.Upload.MaxSize
				}
				allowedTypes := elem.Attachment.AllowedTypes
				if len(allowedTypes) == 0 {
					allowedTypes = config.C.Upload.AllowedTypes
				}
				return &UploadValidationConfig{
					MaxSize:      maxSize,
					AllowedTypes: allowedTypes,
				}
			}
		}
	}

	return defaultCfg
}

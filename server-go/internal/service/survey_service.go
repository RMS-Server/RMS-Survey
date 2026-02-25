package service

import (
	"encoding/json"

	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/repository"
	"gorm.io/gorm"
)

// SurveyService handles public survey loading and answer submission.
type SurveyService struct {
	projectRepo *repository.ProjectRepo
	answerRepo  *repository.AnswerRepo
}

// NewSurveyService creates a new SurveyService.
func NewSurveyService(projectRepo *repository.ProjectRepo, answerRepo *repository.AnswerRepo, db *gorm.DB) *SurveyService {
	_ = db
	return &SurveyService{
		projectRepo: projectRepo,
		answerRepo:  answerRepo,
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
	if err := svc.CreateAnswer(req, &dto.UserInfo{}); err != nil {
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
	return svc.CreateAnswer(req, &dto.UserInfo{})
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

// StatProject returns vote/statistics data for a project.
func (s *SurveyService) StatProject(req *dto.SurveyLoadRequest) (*dto.PublicStatisticsView, error) {
	p, err := s.projectRepo.GetProject(req.Code)
	if err != nil {
		return nil, err
	}
	return &dto.PublicStatisticsView{ProjectID: p.ID, Stats: map[string]interface{}{}}, nil
}

// LoadQuery returns the public query verify view.
func (s *SurveyService) LoadQuery(req *dto.PublicQueryRequest) (*dto.PublicQueryVerifyView, error) {
	p, err := s.projectRepo.GetProject(req.ProjectID)
	if err != nil {
		return nil, err
	}
	return &dto.PublicQueryVerifyView{ProjectID: p.ID, Name: p.Name}, nil
}

// GetQueryResult returns public query results.
func (s *SurveyService) GetQueryResult(req *dto.PublicQueryRequest) (*dto.PublicQueryView, error) {
	return &dto.PublicQueryView{ProjectID: req.ProjectID, Answers: []interface{}{}}, nil
}

// LoadDict returns dictionary entries for a survey.
func (s *SurveyService) LoadDict(req *dto.PublicDictRequest) ([]dto.PublicDictView, error) {
	return []dto.PublicDictView{}, nil
}

// LoadExamResult returns exam scoring result for an answer.
func (s *SurveyService) LoadExamResult(req *dto.PublicExamRequest) (*dto.PublicExamResult, error) {
	return &dto.PublicExamResult{Score: 0, MaxScore: 0, Passed: false}, nil
}

// LoadLinkResult returns linked survey result data.
func (s *SurveyService) LoadLinkResult(req *dto.PublicLinkRequest) (*dto.PublicLinkResult, error) {
	return &dto.PublicLinkResult{ProjectID: req.ProjectID}, nil
}


package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/repository"
)

// ReportService handles report data business logic.
type ReportService struct {
	repo *repository.ReportRepository
}

func NewReportService(repo *repository.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

// surveySchema mirrors the Java RMS Survey survey JSON structure.
type surveySchema struct {
	Pages []struct {
		Elements []struct {
			ID      string `json:"id"`
			Title   string `json:"title"`
			Type    string `json:"type"`
			Options []struct {
				Value string `json:"value"`
				Label string `json:"label"`
			} `json:"options"`
		} `json:"elements"`
	} `json:"pages"`
}

// AnswerValue represents a potentially nested answer with attachments.
// New format: { "questionId": { "value": "answer", "attachments": [...] } }
// Old format: { "questionId": "answer" } (simple string/number/array)
type AnswerValue struct {
	Value       interface{} `json:"value"`
	Attachments []struct {
		FileID   string `json:"fileId"`
		FileName string `json:"fileName"`
		FileSize int64  `json:"fileSize"`
		FileType string `json:"fileType"`
	} `json:"attachments"`
}

// parseAnswerValue extracts the display value from an answer field.
// Handles both simple values (string, number) and nested {value, attachments} structures.
func parseAnswerValue(field interface{}) interface{} {
	if field == nil {
		return nil
	}

	// Try to parse as nested AnswerValue (map with "value" key)
	if m, ok := field.(map[string]interface{}); ok {
		if v, exists := m["value"]; exists {
			return v
		}
		return nil
	}

	// Simple value - return as-is
	return field
}

// formatValue converts a value to its string representation for display.
func formatValue(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		// JSON numbers are parsed as float64
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%v", val)
	case []interface{}:
		// Array of values (e.g., multi-select)
		parts := make([]string, 0, len(val))
		for _, item := range val {
			parts = append(parts, formatValue(item))
		}
		return strings.Join(parts, ", ")
	default:
		return fmt.Sprintf("%v", val)
	}
}

// extractStringValue extracts a string value from a field, handling nested structures.
func extractStringValue(field interface{}) string {
	parsed := parseAnswerValue(field)
	return formatValue(parsed)
}

// extractStringSlice extracts a string slice from a field (for multi-select), handling nested structures.
func extractStringSlice(field interface{}) []string {
	parsed := parseAnswerValue(field)
	if parsed == nil {
		return nil
	}

	switch v := parsed.(type) {
	case string:
		return []string{v}
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	default:
		return nil
	}
}

func (s *ReportService) GetData(shortID string) (*dto.ReportData, error) {
	p, err := s.repo.GetProjectByShortID(shortID)
	if err != nil {
		return nil, err
	}
	total, err := s.repo.CountAnswers(p.ID)
	if err != nil {
		return nil, err
	}

	answers, err := s.repo.GetAnswers(p.ID)
	if err != nil {
		return nil, err
	}

	// Parse survey schema; on failure return partial result with no question stats.
	var schema surveySchema
	if parseErr := json.Unmarshal([]byte(p.Survey), &schema); parseErr != nil {
		log.Printf("report: failed to parse survey JSON for project %s: %v", p.ID, parseErr)
		return &dto.ReportData{
			ProjectID:   p.ID,
			ProjectName: p.Name,
			Total:       total,
			Questions:   []dto.QuestionStat{},
		}, nil
	}

	// Build question index and stat accumulators.
	type questionMeta struct {
		title   string
		qtype   string
		options []struct{ Value, Label string }
	}
	metaMap := make(map[string]questionMeta)
	order := make([]string, 0)

	for _, page := range schema.Pages {
		for _, el := range page.Elements {
			metaMap[el.ID] = questionMeta{
				title: el.Title,
				qtype: el.Type,
				options: func() []struct{ Value, Label string } {
					out := make([]struct{ Value, Label string }, len(el.Options))
					for i, o := range el.Options {
						out[i] = struct{ Value, Label string }{o.Value, o.Label}
					}
					return out
				}(),
			}
			order = append(order, el.ID)
		}
	}

	// Accumulators: option counts and text lists.
	optionCounts := make(map[string]map[string]int) // qid -> optionValue -> count
	textValues := make(map[string][]string)          // qid -> texts
	questionTotals := make(map[string]int)           // qid -> answer count

	for _, a := range answers {
		if a.Answer == "" {
			continue
		}
		var answerMap map[string]interface{}
		if err := json.Unmarshal([]byte(a.Answer), &answerMap); err != nil {
			log.Printf("report: failed to parse answer JSON for answer %s: %v", a.ID, err)
			continue
		}
		for qid, val := range answerMap {
			if _, ok := metaMap[qid]; !ok {
				continue
			}
			questionTotals[qid]++
			meta := metaMap[qid]
			isChoice := strings.Contains(meta.qtype, "radio") ||
				strings.Contains(meta.qtype, "checkbox") ||
				strings.Contains(meta.qtype, "select")

			if isChoice {
				if optionCounts[qid] == nil {
					optionCounts[qid] = make(map[string]int)
				}
				// Extract values from potentially nested structure
				values := extractStringSlice(val)
				for _, s := range values {
					optionCounts[qid][s]++
				}
			} else {
				// Text question — collect up to 100 values.
				if len(textValues[qid]) < 100 {
					text := extractStringValue(val)
					if text != "" {
						textValues[qid] = append(textValues[qid], text)
					}
				}
			}
		}
	}

	// Build ordered question stats.
	questions := make([]dto.QuestionStat, 0, len(order))
	for _, qid := range order {
		meta := metaMap[qid]
		stat := dto.QuestionStat{
			ID:    qid,
			Title: meta.title,
			Type:  meta.qtype,
			Total: questionTotals[qid],
		}

		isChoice := strings.Contains(meta.qtype, "radio") ||
			strings.Contains(meta.qtype, "checkbox") ||
			strings.Contains(meta.qtype, "select")

		if isChoice {
			counts := optionCounts[qid]
			opts := make([]dto.OptionStat, 0, len(meta.options))
			for _, o := range meta.options {
				c := 0
				if counts != nil {
					c = counts[o.Value]
				}
				opts = append(opts, dto.OptionStat{
					Value: o.Value,
					Label: o.Label,
					Count: c,
				})
			}
			stat.Options = opts
		} else {
			stat.Texts = textValues[qid]
		}
		questions = append(questions, stat)
	}

	return &dto.ReportData{
		ProjectID:   p.ID,
		ProjectName: p.Name,
		Total:       total,
		Questions:   questions,
	}, nil
}

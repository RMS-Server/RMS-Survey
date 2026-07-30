package dto_test

import (
	"encoding/json"
	"testing"

	"github.com/rms-survey/server/internal/dto"
)

func TestProjectViewMarshalsMissingSettingAsEmptyObject(t *testing.T) {
	view := dto.ProjectView{
		Survey:  json.RawMessage(`{"pages":[]}`),
		Setting: json.RawMessage(""),
	}

	body, err := json.Marshal(view)
	if err != nil {
		t.Fatalf("marshal project response: %v", err)
	}

	var response struct {
		Setting map[string]any `json:"setting"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("decode project response: %v", err)
	}
	if response.Setting == nil || len(response.Setting) != 0 {
		t.Fatalf("setting = %#v, want empty object", response.Setting)
	}
}

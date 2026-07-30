package dto

import (
	"bytes"
	"encoding/json"
)

// MarshalJSON keeps legacy projects with a blank setting from breaking the
// entire project response. Blank settings predate the server-side default.
func (v ProjectView) MarshalJSON() ([]byte, error) {
	type projectView ProjectView

	normalized := projectView(v)
	if len(bytes.TrimSpace(normalized.Setting)) == 0 {
		normalized.Setting = json.RawMessage(`{}`)
	}
	return json.Marshal(normalized)
}

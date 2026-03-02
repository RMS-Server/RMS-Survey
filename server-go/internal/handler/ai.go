package handler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rms-survey/server/internal/config"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/pkg/response"
)

// AIHandler handles AI chat SSE endpoints.
type AIHandler struct{}

func NewAIHandler() *AIHandler {
	return &AIHandler{}
}

// GetModels returns available AI model types.
func (h *AIHandler) GetModels(c *gin.Context) {
	models := []dto.AIModelType{
		{ID: config.C.AI.Model, Name: config.C.AI.Model},
	}
	response.OK(c, models)
}

// CreateConversation creates a new conversation session.
func (h *AIHandler) CreateConversation(c *gin.Context) {
	var req dto.AIChatRequest
	_ = c.ShouldBindJSON(&req)
	model := c.Query("model")
	if model == "" {
		model = config.C.AI.Model
	}
	resp := dto.AIConversationResponse{
		ID:    fmt.Sprintf("conv-%d", c.Request.Context().Value("ts")),
		Model: model,
	}
	response.OK(c, resp)
}

// CloseConversation closes a conversation session.
func (h *AIHandler) CloseConversation(c *gin.Context) {
	response.OK(c, nil)
}

// Stream handles SSE streaming chat with the upstream AI API.
func (h *AIHandler) Stream(c *gin.Context) {
	content := c.Query("content")
	model := c.Query("model")
	if model == "" {
		model = config.C.AI.Model
	}
	if content == "" {
		content = "hello"
	}

	messages := []map[string]string{
		{"role": "user", "content": content},
	}
	payload := map[string]interface{}{
		"model":    model,
		"messages": messages,
		"stream":   true,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost,
		config.C.AI.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.C.AI.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	defer resp.Body.Close()

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	scanner := bufio.NewScanner(resp.Body)
	c.Stream(func(w io.Writer) bool {
		if !scanner.Scan() {
			return false
		}
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				c.SSEvent("done", "{}")
				return false
			}
			c.SSEvent("message", data)
		}
		return true
	})
}

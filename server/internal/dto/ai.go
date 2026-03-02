package dto

// AIMessage is a single message in a chat conversation.
type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AIChatRequest is the request body for the AI chat endpoint.
type AIChatRequest struct {
	ConversationID string      `json:"conversationId"`
	Model          string      `json:"model"`
	Messages       []AIMessage `json:"messages"`
	Content        string      `json:"content"`
}

// AIConversationResponse is returned when creating a new conversation.
type AIConversationResponse struct {
	ID    string `json:"id"`
	Model string `json:"model"`
}

// AIModelType represents an available AI model.
type AIModelType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

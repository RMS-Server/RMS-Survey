package dto

// PageRequest holds pagination parameters from query string.
type PageRequest struct {
	PageIndex int `json:"pageIndex" form:"pageIndex"`
	PageSize  int `json:"pageSize" form:"pageSize"`
}

// PageResponse is a generic paginated response.
type PageResponse[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
}

// UserInfo represents the authenticated user stored in JWT claims and gin context.
type UserInfo struct {
	UserID   string   `json:"userId"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

// SimpleUserView is a simplified user view for nested responses.
type SimpleUserView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

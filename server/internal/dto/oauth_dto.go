package dto

// OAuthAuthorizeRequest holds the state and code_challenge for OAuth authorization.
type OAuthAuthorizeRequest struct {
	State         string `json:"state" binding:"required"`
	CodeChallenge string `json:"codeChallenge" binding:"required"`
}

// OAuthAuthorizeResponse returns the authorization URL to redirect to.
type OAuthAuthorizeResponse struct {
	AuthURL string `json:"authUrl"`
}

// OAuthCallbackRequest holds the OAuth callback parameters.
type OAuthCallbackRequest struct {
	Code         string `json:"code" binding:"required"`
	State        string `json:"state" binding:"required"`
	CodeVerifier string `json:"codeVerifier" binding:"required"`
}

// OAuthCallbackResponse returns the local JWT token and user info.
type OAuthCallbackResponse struct {
	Token string   `json:"token"`
	User  UserView `json:"user"`
}

// OAuthRefreshRequest holds the refresh token request.
type OAuthRefreshRequest struct {
	UserID string `json:"userId"`
}

// OAuthRefreshResponse returns a new JWT token.
type OAuthRefreshResponse struct {
	Token string   `json:"token"`
	User  UserView `json:"user"`
}

// SSOUserInfo represents the user info from SSO provider.
type SSOUserInfo struct {
	ID              string `json:"sub"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	PermissionLevel int    `json:"permission_level"`
	GroupID         string `json:"group_id"`
	GroupName       string `json:"group_name"`
	Avatar          string `json:"picture"`
}

// OAuthTokenResponse represents the token response from OAuth provider.
type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

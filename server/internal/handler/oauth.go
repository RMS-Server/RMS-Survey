package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rms-survey/server/internal/config"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/pkg/response"
	"github.com/rms-survey/server/internal/repository"
	"github.com/rms-survey/server/internal/service"
	"gorm.io/gorm"
)

// OAuthHandler handles OAuth authentication endpoints.
type OAuthHandler struct {
	svc *service.OAuthService
	db  *gorm.DB
}

// NewOAuthHandler creates a new OAuthHandler.
func NewOAuthHandler(db *gorm.DB) *OAuthHandler {
	userRepo := repository.NewUserRepo(db)
	return &OAuthHandler{
		svc: service.NewOAuthService(userRepo, db),
		db:  db,
	}
}

// RegisterRoutes registers OAuth routes on the given router.
func (h *OAuthHandler) RegisterRoutes(r gin.IRouter) {
	r.GET("/oauth/authorize", h.Authorize)
	r.GET("/oauth/callback", h.CallbackGet)
	r.POST("/oauth/callback", h.Callback)
	r.POST("/oauth/refresh", h.Refresh)
	r.POST("/oauth/logout", h.Logout)
}

// Authorize handles GET /api/oauth/authorize - returns OAuth config for frontend to redirect.
func (h *OAuthHandler) Authorize(c *gin.Context) {
	if !config.C.OAuth.Enabled {
		response.Fail(c, response.CodeError, "OAuth is not enabled")
		return
	}

	// Return OAuth endpoints configuration for frontend to use
	response.OK(c, gin.H{
		"authUrl":     config.C.OAuth.AuthURL,
		"clientId":    config.C.OAuth.ClientID,
		"redirectUri": config.C.OAuth.RedirectURL,
		"scopes":      config.C.OAuth.Scopes,
	})
}

// CallbackGet handles GET /api/oauth/callback - receives OAuth callback redirect from SSO.
// This returns an HTML page that posts to the POST endpoint with PKCE verifier.
func (h *OAuthHandler) CallbackGet(c *gin.Context) {
	if !config.C.OAuth.Enabled {
		response.Fail(c, response.CodeError, "OAuth is not enabled")
		return
	}

	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		response.Fail(c, response.CodeError, "missing code or state parameter")
		return
	}

	// Return HTML page that extracts code_verifier from sessionStorage and POSTs to callback endpoint
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>OAuth Callback</title>
</head>
<body>
    <div id="status" style="font-family: sans-serif; text-align: center; margin-top: 100px;">
        <p>Processing login...</p>
    </div>
    <script>
    (function() {
        var code = %q;
        var state = %q;
        var codeVerifier = sessionStorage.getItem('oauth_code_verifier');

        if (!codeVerifier) {
            document.getElementById('status').innerHTML =
                '<p style="color: red;">Login session expired. <a href="/login">Click here to try again</a></p>';
            return;
        }

        fetch('/api/oauth/callback', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                code: code,
                state: state,
                codeVerifier: codeVerifier
            })
        })
        .then(function(response) {
            return response.json();
        })
        .then(function(data) {
            if (data.code === 0 && data.data && data.data.token) {
                // Store token in localStorage
                localStorage.setItem('sk-token', data.data.token);
                // Clear OAuth state
                sessionStorage.removeItem('oauth_state');
                sessionStorage.removeItem('oauth_code_verifier');
                // Redirect to project list
                window.location.href = '/project';
            } else {
                throw new Error(data.message || 'Login failed');
            }
        })
        .catch(function(error) {
            document.getElementById('status').innerHTML =
                '<p style="color: red;">Login failed: ' + error.message + '</p>' +
                '<p><a href="/login">Click here to try again</a></p>';
        });
    })();
    </script>
</body>
</html>`, code, state)
}
func (h *OAuthHandler) Callback(c *gin.Context) {
	if !config.C.OAuth.Enabled {
		response.Fail(c, response.CodeError, "OAuth is not enabled")
		return
	}

	var req dto.OAuthCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, "invalid request: "+err.Error())
		return
	}

	if req.Code == "" || req.CodeVerifier == "" {
		response.Fail(c, response.CodeError, "missing required parameters")
		return
	}

	// Exchange code for tokens
	tokenResp, err := h.svc.ExchangeCode(req.Code, req.CodeVerifier)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}

	// Fetch user info from SSO
	ssoUserInfo, err := h.svc.FetchUserInfo(tokenResp.AccessToken)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}

	// Sync user to local database
	userInfo, err := h.svc.SyncUser(ssoUserInfo)
	if err != nil {
		if err == service.ErrOAuthPermissionDenied {
			response.Fail(c, response.CodeError, "permission denied: insufficient permission level")
			return
		}
		response.Fail(c, response.CodeError, err.Error())
		return
	}

	// Store refresh token
	deviceInfo := c.GetHeader("User-Agent")
	if err := h.svc.StoreRefreshToken(userInfo.UserID, tokenResp.RefreshToken, tokenResp.ExpiresIn, deviceInfo); err != nil {
		// Log but don't fail - refresh token storage is optional for functionality
		// The user can still use the system, just won't be able to auto-refresh
	}

	// Generate local JWT
	token, err := h.svc.GenerateLocalJWT(userInfo)
	if err != nil {
		response.Fail(c, response.CodeError, "failed to generate token")
		return
	}

	// Set JWT cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     config.C.JWT.CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	c.Header("Authorization", token)

	// Get user view for response
	userView, err := h.getUserView(userInfo.UserID)
	if err != nil {
		response.Fail(c, response.CodeError, "failed to get user info")
		return
	}

	response.OK(c, dto.OAuthCallbackResponse{
		Token: token,
		User:  *userView,
	})
}

// Refresh handles POST /api/oauth/refresh - refreshes local JWT.
func (h *OAuthHandler) Refresh(c *gin.Context) {
	if !config.C.OAuth.Enabled {
		response.Fail(c, response.CodeError, "OAuth is not enabled")
		return
	}

	// Get user ID from context (requires auth)
	userInfo, ok := c.Get("currentUser")
	if !ok {
		response.Unauthorized(c)
		return
	}
	u := userInfo.(*dto.UserInfo)

	// Refresh access token using stored refresh token
	_, newUserInfo, err := h.svc.RefreshAccessToken(u.UserID)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}

	// Generate new local JWT
	token, err := h.svc.GenerateLocalJWT(newUserInfo)
	if err != nil {
		response.Fail(c, response.CodeError, "failed to generate token")
		return
	}

	// Set JWT cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     config.C.JWT.CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	c.Header("Authorization", token)

	// Get user view for response
	userView, err := h.getUserView(newUserInfo.UserID)
	if err != nil {
		response.Fail(c, response.CodeError, "failed to get user info")
		return
	}

	response.OK(c, dto.OAuthRefreshResponse{
		Token: token,
		User:  *userView,
	})
}

// Logout handles POST /api/oauth/logout - clears OAuth session.
func (h *OAuthHandler) Logout(c *gin.Context) {
	// Clear JWT cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     config.C.JWT.CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	// Get user ID from context if available (optional, allows clearing session)
	if userInfo, ok := c.Get("currentUser"); ok {
		u := userInfo.(*dto.UserInfo)
		_ = h.svc.Logout(u.UserID) // Ignore error, logout always succeeds
	}

	response.OK(c, nil)
}

// getUserView fetches user details for response.
func (h *OAuthHandler) getUserView(userID string) (*dto.UserView, error) {
	var user struct {
		ID     string `gorm:"column:id"`
		Name   string `gorm:"column:name"`
		Phone  string `gorm:"column:phone"`
		Email  string `gorm:"column:email"`
		Avatar string `gorm:"column:avatar"`
		Gender string `gorm:"column:gender"`
		DeptID string `gorm:"column:dept_id"`
		Status int    `gorm:"column:status"`
	}
	if err := h.db.Table("t_user").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	view := &dto.UserView{
		ID:     user.ID,
		Name:   user.Name,
		Phone:  user.Phone,
		Email:  user.Email,
		Avatar: user.Avatar,
		Gender: user.Gender,
		DeptID: user.DeptID,
		Status: user.Status,
	}

	// Get account username
	var account struct {
		AuthAccount string `gorm:"column:auth_account"`
	}
	if err := h.db.Table("t_account").Where("user_id = ?", userID).First(&account).Error; err == nil {
		view.Username = account.AuthAccount
	}

	// Get roles
	var roles []dto.RoleView
	h.db.Table("t_role r").
		Select("r.id, r.name, r.code").
		Joins("JOIN t_user_role ur ON ur.role_id = r.id").
		Where("ur.user_id = ?", userID).
		Find(&roles)
	view.Roles = roles

	return view, nil
}

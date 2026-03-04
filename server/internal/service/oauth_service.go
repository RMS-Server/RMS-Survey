package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rms-survey/server/internal/config"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/model"
	jwtpkg "github.com/rms-survey/server/internal/pkg/jwt"
	"github.com/rms-survey/server/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrOAuthNotEnabled       = errors.New("OAuth is not enabled")
	ErrOAuthPermissionDenied = errors.New("permission denied: permission level insufficient")
	ErrOAuthStateMismatch    = errors.New("OAuth state mismatch")
	ErrOAuthTokenExchange    = errors.New("failed to exchange OAuth token")
	ErrOAuthUserInfo         = errors.New("failed to fetch user info from OAuth provider")
)

// OAuthService handles OAuth2 authentication flow.
type OAuthService struct {
	userRepo *repository.UserRepo
	db       *gorm.DB
}

// NewOAuthService creates a new OAuthService.
func NewOAuthService(userRepo *repository.UserRepo, db *gorm.DB) *OAuthService {
	return &OAuthService{
		userRepo: userRepo,
		db:       db,
	}
}

// GeneratePKCE generates a code_verifier and code_challenge for PKCE.
func (s *OAuthService) GeneratePKCE() (codeVerifier, codeChallenge string, err error) {
	// Generate 32-96 bytes of random data for code_verifier
	verifierBytes := make([]byte, 32)
	if _, err := rand.Read(verifierBytes); err != nil {
		return "", "", err
	}
	codeVerifier = base64URLEncode(verifierBytes)

	// Generate code_challenge using S256 method
	hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge = base64URLEncode(hash[:])

	return codeVerifier, codeChallenge, nil
}

// base64URLEncode encodes bytes using URL-safe base64 without padding.
func base64URLEncode(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

// GenerateState generates a random state string for CSRF protection.
func (s *OAuthService) GenerateState() (string, error) {
	return nanoid.New()
}

// GetAuthURL builds the authorization URL with PKCE parameters.
func (s *OAuthService) GetAuthURL(state, codeChallenge string) string {
	cfg := config.C.OAuth
	params := url.Values{}
	params.Set("client_id", cfg.ClientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", cfg.RedirectURL)
	params.Set("scope", cfg.Scopes)
	params.Set("state", state)
	params.Set("code_challenge", codeChallenge)
	params.Set("code_challenge_method", "S256")

	return fmt.Sprintf("%s?%s", cfg.AuthURL, params.Encode())
}

// ExchangeCode exchanges the authorization code for tokens.
func (s *OAuthService) ExchangeCode(code, codeVerifier string) (*dto.OAuthTokenResponse, error) {
	cfg := config.C.OAuth

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", cfg.ClientID)
	data.Set("code", code)
	data.Set("redirect_uri", cfg.RedirectURL)
	data.Set("code_verifier", codeVerifier)

	req, err := http.NewRequest("POST", cfg.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d, body: %s", ErrOAuthTokenExchange, resp.StatusCode, string(body))
	}

	var tokenResp dto.OAuthTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// FetchUserInfo retrieves user info from the OAuth provider.
func (s *OAuthService) FetchUserInfo(accessToken string) (*dto.SSOUserInfo, error) {
	cfg := config.C.OAuth

	req, err := http.NewRequest("GET", cfg.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d, body: %s", ErrOAuthUserInfo, resp.StatusCode, string(body))
	}

	var userInfo dto.SSOUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	// Handle various possible field names from SSO
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err == nil {
		if userInfo.ID == "" {
			if sub, ok := raw["sub"].(string); ok {
				userInfo.ID = sub
			} else if uid, ok := raw["user_id"].(string); ok {
				userInfo.ID = uid
			} else if uid, ok := raw["id"].(string); ok {
				userInfo.ID = uid
			}
		}
		if userInfo.Name == "" {
			if name, ok := raw["nickname"].(string); ok {
				userInfo.Name = name
			} else if name, ok := raw["preferred_username"].(string); ok {
				userInfo.Name = name
			}
		}
		if userInfo.PermissionLevel == 0 {
			if pl, ok := raw["permission_level"].(float64); ok {
				userInfo.PermissionLevel = int(pl)
			} else if pl, ok := raw["permissionLevel"].(float64); ok {
				userInfo.PermissionLevel = int(pl)
			}
		}
	}

	return &userInfo, nil
}

// SyncUser creates or updates a local user from SSO user info.
func (s *OAuthService) SyncUser(userInfo *dto.SSOUserInfo) (*dto.UserInfo, error) {
	cfg := config.C.OAuth

	// Check permission level
	if userInfo.PermissionLevel < cfg.MinPermissionLevel {
		return nil, ErrOAuthPermissionDenied
	}

	// Find existing user by SSO ID
	var user = s.findUserBySSOID(userInfo.ID)
	now := time.Now()

	if user == nil {
		// Try to find by email
		user = s.findUserByEmail(userInfo.Email)
	}

	if user == nil {
		// Create new user
		userID, _ := nanoid.New()
		user = &model.User{
			Name:               userInfo.Name,
			Email:              userInfo.Email,
			Phone:              userInfo.Phone,
			Avatar:             userInfo.Avatar,
			Status:             1,
		}
		user.ID = userID
		user.SSOID = &userInfo.ID
		user.SSOLastSync = &now
		permLevel := userInfo.PermissionLevel
		user.SSOPermissionLevel = &permLevel
		user.SSOGroupID = &userInfo.GroupID
		user.SSOGroupName = &userInfo.GroupName

		if err := s.createUser(user); err != nil {
			return nil, err
		}
	} else {
		// Update existing user
		user.SSOID = &userInfo.ID
		user.SSOLastSync = &now
		permLevel := userInfo.PermissionLevel
		user.SSOPermissionLevel = &permLevel
		user.SSOGroupID = &userInfo.GroupID
		user.SSOGroupName = &userInfo.GroupName
		if userInfo.Name != "" {
			user.Name = userInfo.Name
		}
		if userInfo.Email != "" {
			user.Email = userInfo.Email
		}
		if userInfo.Phone != "" {
			user.Phone = userInfo.Phone
		}
		if userInfo.Avatar != "" {
			user.Avatar = userInfo.Avatar
		}
		if err := s.updateUser(user); err != nil {
			return nil, err
		}
	}

	// Get user authorities
	authorities, _ := s.getUserAuthorities(user.ID)

	return &dto.UserInfo{
		UserID:   user.ID,
		Username: userInfo.Name,
		Roles:    authorities,
	}, nil
}

// GenerateLocalJWT creates a JWT token for the authenticated user.
func (s *OAuthService) GenerateLocalJWT(userInfo *dto.UserInfo) (string, error) {
	return jwtpkg.GenerateToken(*userInfo)
}

// StoreRefreshToken saves the refresh token to database.
func (s *OAuthService) StoreRefreshToken(userID, refreshToken string, expiresIn int, deviceInfo string) error {
	id, _ := nanoid.New()
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

	// Delete existing sessions for this user (single session per user)
	s.db.Where("user_id = ?", userID).Delete(&model.OAuthSession{})

	session := &model.OAuthSession{
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		DeviceInfo:   deviceInfo,
	}
	session.ID = id

	return s.db.Create(session).Error
}

// RefreshAccessToken uses the stored refresh token to get a new access token.
func (s *OAuthService) RefreshAccessToken(userID string) (*dto.OAuthTokenResponse, *dto.UserInfo, error) {
	// Get stored refresh token
	var session model.OAuthSession
	if err := s.db.Where("user_id = ? AND expires_at > ?", userID, time.Now()).First(&session).Error; err != nil {
		return nil, nil, errors.New("no valid refresh token found")
	}

	cfg := config.C.OAuth

	// Request new token using refresh_token grant
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", cfg.ClientID)
	data.Set("refresh_token", session.RefreshToken)

	req, err := http.NewRequest("POST", cfg.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("token refresh failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	var tokenResp dto.OAuthTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, nil, err
	}

	// Fetch updated user info
	ssoUserInfo, err := s.FetchUserInfo(tokenResp.AccessToken)
	if err != nil {
		return nil, nil, err
	}

	// Sync user and check permission
	userInfo, err := s.SyncUser(ssoUserInfo)
	if err != nil {
		return nil, nil, err
	}

	// Update refresh token in database
	if tokenResp.RefreshToken != "" {
		now := time.Now()
		session.RefreshToken = tokenResp.RefreshToken
		session.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		session.LastRefreshAt = &now
		s.db.Save(&session)
	}

	return &tokenResp, userInfo, nil
}

// CleanExpiredSessions removes expired OAuth sessions.
func (s *OAuthService) CleanExpiredSessions() error {
	return s.db.Where("expires_at < ?", time.Now()).Delete(&model.OAuthSession{}).Error
}

// Logout removes the OAuth session for a user.
func (s *OAuthService) Logout(userID string) error {
	return s.db.Where("user_id = ?", userID).Delete(&model.OAuthSession{}).Error
}

// Helper methods for user operations

func (s *OAuthService) findUserBySSOID(ssoID string) *model.User {
	var user model.User
	if err := s.db.Where("sso_id = ?", ssoID).First(&user).Error; err != nil {
		return nil
	}
	return &user
}

func (s *OAuthService) findUserByEmail(email string) *model.User {
	if email == "" {
		return nil
	}
	var user model.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil
	}
	return &user
}

func (s *OAuthService) createUser(user *model.User) error {
	return s.db.Create(user).Error
}

func (s *OAuthService) updateUser(user *model.User) error {
	return s.db.Save(user).Error
}

func (s *OAuthService) getUserAuthorities(userID string) ([]string, error) {
	var userRoles []struct {
		RoleID string `gorm:"column:role_id"`
	}
	if err := s.db.Table("t_user_role").Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return []string{}, nil
	}

	roleIDs := make([]string, len(userRoles))
	for i, ur := range userRoles {
		roleIDs[i] = ur.RoleID
	}

	var roles []struct {
		Code      string `gorm:"column:code"`
		Authority string `gorm:"column:authority"`
	}
	if err := s.db.Table("t_role").Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return nil, err
	}

	authSet := make(map[string]struct{})
	for _, role := range roles {
		authSet["ROLE_"+role.Code] = struct{}{}
		for _, a := range strings.Split(role.Authority, ",") {
			a = strings.TrimSpace(a)
			if a != "" {
				authSet[a] = struct{}{}
			}
		}
	}

	authorities := make([]string, 0, len(authSet))
	for a := range authSet {
		authorities = append(authorities, a)
	}
	return authorities, nil
}

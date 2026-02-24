package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/surveyking/server/internal/dto"
	jwtpkg "github.com/surveyking/server/internal/pkg/jwt"
	"github.com/surveyking/server/internal/pkg/response"
	"github.com/surveyking/server/internal/repository"
	"github.com/surveyking/server/internal/service"
	"gorm.io/gorm"
)

const cookieName = "sk-token"

// captchaStore holds in-memory captcha answers with expiry.
var captchaStore = struct {
	sync.Mutex
	items map[string]captchaItem
}{items: make(map[string]captchaItem)}

type captchaItem struct {
	answer  string
	expires time.Time
}

// UserHandler handles all user and auth HTTP endpoints.
type UserHandler struct {
	svc *service.UserService
}

// NewUserHandler creates a UserHandler wired to the given DB.
func NewUserHandler(db *gorm.DB) *UserHandler {
	repo := repository.NewUserRepo(db)
	svc := service.NewUserService(repo)
	return &UserHandler{svc: svc}
}

// RegisterRoutes wires all user routes onto the provided gin groups.
// public  -> /api/public
// captcha -> /captcha
// user    -> /api/user
// root    -> r (for /currentUser)
func (h *UserHandler) RegisterRoutes(public, captchaGrp, userGrp, root gin.IRouter) {
	public.POST("/login", h.Login)
	public.POST("/logout", h.Logout)
	public.POST("/register", h.Register)
	public.GET("/rsaPublicKey", h.GetRSAPublicKey)

	captchaGrp.GET("/get", h.GetCaptcha)
	captchaGrp.POST("/check", h.CheckCaptcha)

	root.GET("/currentUser", h.CurrentUser)

	userGrp.PUT("/updatePassword", h.UpdatePassword)
	userGrp.GET("/list", h.ListUsers)
	userGrp.POST("", h.CreateUser)
	userGrp.PUT("", h.UpdateUser)
	userGrp.DELETE("/:id", h.DeleteUser)
	userGrp.GET("/:id", h.GetUser)
	userGrp.POST("/bindRole", h.BindRole)
}

// Login handles POST /api/public/login
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, "invalid request")
		return
	}

	user, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		response.Fail(c, response.CodeError, "username or password error")
		return
	}

	userInfo := dto.UserInfo{
		UserID:   user.ID,
		Username: req.Username,
	}
	if auths, err := h.svc.GetUserAuthorities(user.ID); err == nil {
		userInfo.Roles = auths
	}

	token, err := jwtpkg.GenerateToken(userInfo)
	if err != nil {
		response.Fail(c, response.CodeError, "failed to generate token")
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})
	c.Header("Authorization", token)

	view, _ := h.svc.GetCurrentUser(user.ID)
	response.OK(c, dto.LoginResponse{Token: token, User: *view})
}

// Logout handles POST /api/public/logout
func (h *UserHandler) Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	response.OK(c, nil)
}

// Register handles POST /api/public/register
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, "invalid request")
		return
	}

	createReq := dto.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
		Status:   1,
	}
	if req.Role != "" {
		createReq.Roles = []string{req.Role}
	}

	if err := h.svc.CreateUser(createReq); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

// GetRSAPublicKey handles GET /api/public/rsaPublicKey
func (h *UserHandler) GetRSAPublicKey(c *gin.Context) {
	pub, err := h.svc.GetRSAPublicKey()
	if err != nil {
		response.Fail(c, response.CodeError, "failed to get RSA public key")
		return
	}
	response.OK(c, pub)
}

// CurrentUser handles GET /currentUser
func (h *UserHandler) CurrentUser(c *gin.Context) {
	userInfo, ok := c.Get("currentUser")
	if !ok {
		response.Unauthorized(c)
		return
	}
	u := userInfo.(*dto.UserInfo)
	view, err := h.svc.GetCurrentUser(u.UserID)
	if err != nil {
		response.Fail(c, response.CodeError, "user not found")
		return
	}
	response.OK(c, view)
}

// UpdatePassword handles PUT /api/user/updatePassword
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, "invalid request")
		return
	}

	userInfo, _ := c.Get("currentUser")
	u := userInfo.(*dto.UserInfo)

	if err := h.svc.UpdatePassword(u.UserID, req.OldPassword, req.NewPassword); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

// ListUsers handles GET /api/user/list
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req dto.UserQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, response.CodeError, "invalid request")
		return
	}

	result, err := h.svc.ListUsers(req)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

// CreateUser handles POST /api/user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, "invalid request")
		return
	}

	if err := h.svc.CreateUser(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

// UpdateUser handles PUT /api/user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, "invalid request")
		return
	}

	if err := h.svc.UpdateUser(req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

// DeleteUser handles DELETE /api/user/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteUser(id); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

// GetUser handles GET /api/user/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	view, err := h.svc.GetCurrentUser(id)
	if err != nil {
		response.Fail(c, response.CodeError, "user not found")
		return
	}
	response.OK(c, view)
}

// BindRole handles POST /api/user/bindRole
func (h *UserHandler) BindRole(c *gin.Context) {
	var req dto.BindRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, "invalid request")
		return
	}

	if err := h.svc.BindRoles(req.UserID, req.RoleIDs); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

// GetCaptcha handles GET /captcha/get — returns a simple base64 PNG captcha.
func (h *UserHandler) GetCaptcha(c *gin.Context) {
	code, err := randomDigits(4)
	if err != nil {
		response.Fail(c, response.CodeError, "captcha generation failed")
		return
	}

	id, _ := nanoid.New()
	captchaStore.Lock()
	captchaStore.items[id] = captchaItem{answer: code, expires: time.Now().Add(5 * time.Minute)}
	captchaStore.Unlock()

	img := generateCaptchaImage(code)
	var rawBuf bytes.Buffer
	_ = png.Encode(&rawBuf, img)
	imgBase64 := base64.StdEncoding.EncodeToString(rawBuf.Bytes())

	response.OK(c, gin.H{
		"captchaId":  id,
		"captchaImg": "data:image/png;base64," + imgBase64,
	})
}

// CheckCaptcha handles POST /captcha/check
func (h *UserHandler) CheckCaptcha(c *gin.Context) {
	var req struct {
		CaptchaID   string `json:"captchaId"`
		CaptchaCode string `json:"captchaCode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, "invalid request")
		return
	}

	captchaStore.Lock()
	item, ok := captchaStore.items[req.CaptchaID]
	if ok {
		delete(captchaStore.items, req.CaptchaID)
	}
	captchaStore.Unlock()

	if !ok || time.Now().After(item.expires) || item.answer != req.CaptchaCode {
		response.Fail(c, response.CodeError, "captcha verification failed")
		return
	}
	response.OK(c, true)
}

// randomDigits generates a random n-digit numeric string.
func randomDigits(n int) (string, error) {
	digits := make([]byte, n)
	for i := range digits {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		digits[i] = byte('0' + num.Int64())
	}
	return string(digits), nil
}

// generateCaptchaImage creates a minimal PNG image with the given text.
func generateCaptchaImage(text string) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 120, 40))
	// Fill background white
	for y := 0; y < 40; y++ {
		for x := 0; x < 120; x++ {
			img.Set(x, y, color.White)
		}
	}
	// Draw simple dots for each digit (placeholder — real impl would use font rendering)
	for i, ch := range text {
		x := 15 + i*25
		drawChar(img, x, 15, ch)
	}
	return img
}

// drawChar draws a very simple representation of a digit character.
func drawChar(img *image.RGBA, x, y int, ch rune) {
	c := color.RGBA{R: 50, G: 50, B: 200, A: 255}
	// Draw a 3x5 pixel block as a placeholder character indicator
	for dy := 0; dy < 10; dy++ {
		for dx := 0; dx < 8; dx++ {
			img.Set(x+dx, y+dy, c)
		}
	}
	// Encode digit value as a small pattern
	_ = ch
}


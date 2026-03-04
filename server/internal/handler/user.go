package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"math/big"
	"time"

	"github.com/gin-gonic/gin"
	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/pkg/cache"
	"github.com/rms-survey/server/internal/pkg/response"
	"github.com/rms-survey/server/internal/repository"
	"github.com/rms-survey/server/internal/service"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// captchaStore holds captcha answers with automatic TTL eviction via bigcache.
var captchaStore = cache.New(5 * time.Minute)

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
	// OAuth routes are handled separately by OAuthHandler
	// Local auth routes (login/register/rsaPublicKey) removed

	captchaGrp.GET("/get", h.GetCaptcha)
	captchaGrp.POST("/check", h.CheckCaptcha)

	root.GET("/currentUser", h.CurrentUser)
	root.GET("/userOverview", h.UserOverview)
	root.GET("/listUserTask", h.ListUserTask)
	root.GET("/listHistoryTask", h.ListHistoryTask)
	root.POST("/importUser", h.ImportUser)

	userGrp.PUT("/updatePassword", h.UpdatePassword)
	userGrp.GET("/list", h.ListUsers)
	userGrp.POST("", h.CreateUser)
	userGrp.PUT("", h.UpdateUser)
	userGrp.DELETE("/:id", h.DeleteUser)
	userGrp.GET("/:id", h.GetUser)
	userGrp.POST("/bindRole", h.BindRole)
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

	userInfo, ok := c.Get("currentUser")
	if !ok {
		response.Unauthorized(c)
		return
	}
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

// UserOverview handles GET /userOverview
func (h *UserHandler) UserOverview(c *gin.Context) {
	userInfo, ok := c.Get("currentUser")
	if !ok {
		response.Unauthorized(c)
		return
	}
	u := userInfo.(*dto.UserInfo)
	overview, err := h.svc.GetUserOverview(u.UserID)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, overview)
}

// ImportUser handles POST /importUser — reads an Excel file and bulk-creates users.
func (h *UserHandler) ImportUser(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, response.CodeError, "missing file: "+err.Error())
		return
	}
	f, err := fh.Open()
	if err != nil {
		response.Fail(c, response.CodeError, "open file: "+err.Error())
		return
	}
	defer f.Close()

	xlsx, err := excelize.OpenReader(f)
	if err != nil {
		response.Fail(c, response.CodeError, "parse excel: "+err.Error())
		return
	}
	defer xlsx.Close()

	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		response.Fail(c, response.CodeError, "read rows: "+err.Error())
		return
	}

	cell := func(row []string, idx int) string {
		if idx < len(row) {
			return row[idx]
		}
		return ""
	}

	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}
		name := cell(row, 0)
		if name == "" {
			continue // skip blank rows
		}
		username := cell(row, 1)
		if username == "" {
			username = name
		}
		password := cell(row, 2)
		if password == "" {
			password = "123456"
		}
		req := dto.CreateUserRequest{
			Name:     name,
			Username: username,
			Password: password,
			Phone:    cell(row, 3),
			Email:    cell(row, 4),
			Status:   1,
		}
		_ = h.svc.CreateUser(req) // continue on error
	}

	response.OK(c, nil)
}

// ListUserTask handles GET /listUserTask
func (h *UserHandler) ListUserTask(c *gin.Context) {
	userInfo, ok := c.Get("currentUser")
	if !ok {
		response.Unauthorized(c)
		return
	}
	u := userInfo.(*dto.UserInfo)
	var query dto.MyTaskQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.GetUserTasks(u.UserID, query)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

// ListHistoryTask handles GET /listHistoryTask
func (h *UserHandler) ListHistoryTask(c *gin.Context) {
	userInfo, ok := c.Get("currentUser")
	if !ok {
		response.Unauthorized(c)
		return
	}
	u := userInfo.(*dto.UserInfo)
	var query dto.MyTaskQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.GetHistoryTasks(u.UserID, query)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

// GetCaptcha handles GET /captcha/get — returns a simple base64 PNG captcha.
func (h *UserHandler) GetCaptcha(c *gin.Context) {
	code, err := randomDigits(4)
	if err != nil {
		response.Fail(c, response.CodeError, "captcha generation failed")
		return
	}

	id, _ := nanoid.New()
	_ = captchaStore.Set(id, []byte(code))

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

	stored, err := captchaStore.Get(req.CaptchaID)
	if err != nil || string(stored) != req.CaptchaCode {
		response.Fail(c, response.CodeError, "captcha verification failed")
		return
	}
	_ = captchaStore.Delete(req.CaptchaID)
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

// digitBitmap is a 5x7 pixel bitmap for digits 0-9.
// Each uint8 is a row bitmask (bit 4 = leftmost pixel of 5-wide glyph).
var digitBitmap = [10][7]uint8{
	{0x0E, 0x11, 0x13, 0x15, 0x19, 0x11, 0x0E}, // 0
	{0x04, 0x0C, 0x04, 0x04, 0x04, 0x04, 0x0E}, // 1
	{0x0E, 0x11, 0x01, 0x06, 0x08, 0x10, 0x1F}, // 2
	{0x1F, 0x01, 0x02, 0x06, 0x01, 0x11, 0x0E}, // 3
	{0x02, 0x06, 0x0A, 0x12, 0x1F, 0x02, 0x02}, // 4
	{0x1F, 0x10, 0x1E, 0x01, 0x01, 0x11, 0x0E}, // 5
	{0x06, 0x08, 0x10, 0x1E, 0x11, 0x11, 0x0E}, // 6
	{0x1F, 0x01, 0x02, 0x04, 0x08, 0x08, 0x08}, // 7
	{0x0E, 0x11, 0x11, 0x0E, 0x11, 0x11, 0x0E}, // 8
	{0x0E, 0x11, 0x11, 0x0F, 0x01, 0x02, 0x0C}, // 9
}

// drawChar renders a digit character using a 5x7 bitmap at 2x scale.
func drawChar(img *image.RGBA, x, y int, ch rune) {
	if ch < '0' || ch > '9' {
		return
	}
	c := color.RGBA{R: 50, G: 50, B: 200, A: 255}
	bitmap := digitBitmap[ch-'0']
	for row, mask := range bitmap {
		for col := 0; col < 5; col++ {
			if mask&(1<<uint(4-col)) != 0 {
				// 2x scale
				img.Set(x+col*2, y+row*2, c)
				img.Set(x+col*2+1, y+row*2, c)
				img.Set(x+col*2, y+row*2+1, c)
				img.Set(x+col*2+1, y+row*2+1, c)
			}
		}
	}
}


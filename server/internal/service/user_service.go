package service

import (
	"errors"
	"log"
	"strings"
	"sync"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/model"
	rsapkg "github.com/rms-survey/server/internal/pkg/rsa"
	"github.com/rms-survey/server/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserDisabled       = errors.New("user is disabled")
	ErrUsernameExists     = errors.New("username already exists")
	ErrInvalidPassword    = errors.New("invalid old password")
)

// rsaKeyCache holds the in-memory RSA key pair to avoid repeated DB lookups.
var rsaKeyCache struct {
	sync.RWMutex
	privateKey string
	publicKey  string
	loaded     bool
}

// UserService handles user business logic.
type UserService struct {
	repo *repository.UserRepo
}

// NewUserService creates a new UserService.
func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

// GetRSAPublicKey returns the RSA public key (raw base64, Java-compatible).
func (s *UserService) GetRSAPublicKey() (string, error) {
	rsaKeyCache.RLock()
	if rsaKeyCache.loaded {
		pub := rsaKeyCache.publicKey
		rsaKeyCache.RUnlock()
		return pub, nil
	}
	rsaKeyCache.RUnlock()

	rsaKeyCache.Lock()
	defer rsaKeyCache.Unlock()
	// Double-check after acquiring write lock
	if rsaKeyCache.loaded {
		return rsaKeyCache.publicKey, nil
	}

	// Try loading from DB first
	info, err := s.repo.FindSysInfoByName("rsa_private_key")
	if err == nil {
		// Handle escaped newlines from DB storage
		privKey := strings.ReplaceAll(info.Description, "\\n", "\n")
		pubKeyStored := strings.ReplaceAll(info.Setting, "\\n", "\n")
		// Convert to raw base64 if stored as PEM format
		pubKey := rsapkg.ExtractPublicKeyBase64(pubKeyStored)
		log.Printf("[DEBUG] Loaded RSA keys from DB, privKey len: %d, pubKey len: %d", len(privKey), len(pubKey))
		rsaKeyCache.privateKey = privKey
		rsaKeyCache.publicKey = pubKey
		rsaKeyCache.loaded = true
		return rsaKeyCache.publicKey, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	// Generate new key pair
	priv, pub, err := rsapkg.GenerateKeyPair()
	if err != nil {
		return "", err
	}
	rsaKeyCache.privateKey = priv
	rsaKeyCache.publicKey = pub
	rsaKeyCache.loaded = true

	id, _ := nanoid.New()
	// Store private key as PEM, public key as raw base64 (Java-compatible)
	sysInfo := &model.SysInfo{
		Name:        "rsa_private_key",
		Description: priv,
		Setting:     pub,
	}
	sysInfo.ID = id
	_ = s.repo.SaveSysInfo(sysInfo)

	return pub, nil
}

// getPrivateKey returns the cached RSA private key.
func (s *UserService) getPrivateKey() (string, error) {
	if _, err := s.GetRSAPublicKey(); err != nil {
		return "", err
	}
	rsaKeyCache.RLock()
	priv := rsaKeyCache.privateKey
	rsaKeyCache.RUnlock()
	return priv, nil
}

// Login verifies credentials and returns the user on success.
// encryptedPassword is RSA-encrypted (base64) or plaintext if RSA is unavailable.
func (s *UserService) Login(username, encryptedPassword string) (*model.User, error) {
	log.Printf("[DEBUG] Login attempt - username: %s, encryptedPassword len: %d", username, len(encryptedPassword))

	// Print current public key for debugging
	pubKey, _ := s.GetRSAPublicKey()
	log.Printf("[DEBUG] Current public key (first 100 chars): %s", pubKey[:min(100, len(pubKey))])

	account, err := s.repo.FindByUsername(username)
	if err != nil {
		log.Printf("[DEBUG] User not found: %v", err)
		return nil, ErrInvalidCredentials
	}

	if account.Status != 1 {
		log.Printf("[DEBUG] User disabled - status: %d", account.Status)
		return nil, ErrUserDisabled
	}

	// Attempt RSA decryption; fall back to plaintext for non-encrypted clients.
	password := encryptedPassword
	privKey, privKeyErr := s.getPrivateKey()
	log.Printf("[DEBUG] Private key available: %v, len: %d", privKeyErr == nil, len(privKey))

	if privKeyErr == nil && privKey != "" {
		decrypted, decryptErr := rsapkg.Decrypt(privKey, encryptedPassword)
		if decryptErr == nil {
			password = decrypted
			log.Printf("[DEBUG] RSA decryption successful, password len: %d", len(password))
		} else {
			log.Printf("[DEBUG] RSA decryption failed: %v", decryptErr)
		}
	}

	log.Printf("[DEBUG] Comparing password hash, input password: %q, stored hash len: %d", password, len(account.AuthSecret))
	if err := bcrypt.CompareHashAndPassword([]byte(account.AuthSecret), []byte(password)); err != nil {
		log.Printf("[DEBUG] Password mismatch: %v", err)
		return nil, ErrInvalidCredentials
	}

	user, err := s.repo.FindUserByID(account.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	log.Printf("[DEBUG] Login successful for user: %s", user.Name)
	return user, nil
}

// GetCurrentUser returns a UserView for the given user ID.
func (s *UserService) GetCurrentUser(userID string) (*dto.UserView, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	account, _ := s.repo.FindAccountByUserID(userID)
	roles, _ := s.repo.GetUserRoles(userID)

	view := userToView(user, account, roles, nil)
	return view, nil
}

// ListUsers returns a paginated list of users.
func (s *UserService) ListUsers(req dto.UserQueryRequest) (dto.PageResponse[dto.UserView], error) {
	users, total, err := s.repo.ListUsers(req)
	if err != nil {
		return dto.PageResponse[dto.UserView]{}, err
	}

	views := make([]dto.UserView, 0, len(users))
	for _, u := range users {
		u := u
		account, _ := s.repo.FindAccountByUserID(u.ID)
		roles, _ := s.repo.GetUserRoles(u.ID)
		views = append(views, *userToView(&u, account, roles, nil))
	}

	return dto.PageResponse[dto.UserView]{List: views, Total: total}, nil
}

// CreateUser creates a new user with an associated account.
func (s *UserService) CreateUser(req dto.CreateUserRequest) error {
	// Check username uniqueness
	if _, err := s.repo.FindByUsername(req.Username); err == nil {
		return ErrUsernameExists
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userID, _ := nanoid.New()
	status := req.Status
	if status == 0 {
		status = 1
	}

	user := &model.User{
		Name:   req.Name,
		Phone:  req.Phone,
		Email:  req.Email,
		Gender: req.Gender,
		DeptID: req.DeptID,
		Status: status,
	}
	user.ID = userID

	account := &model.Account{
		AuthAccount: req.Username,
		AuthSecret:  string(hashedPwd),
		AuthType:    "PWD",
		UserType:    "SysUser",
		Status:      status,
	}
	accountID, _ := nanoid.New()
	account.ID = accountID

	if err := s.repo.CreateUser(user, account); err != nil {
		return err
	}

	if len(req.Roles) > 0 {
		_ = s.repo.BindRoles(userID, req.Roles)
	}
	return nil
}

// UpdateUser updates user fields and optionally the account.
func (s *UserService) UpdateUser(req dto.UpdateUserRequest) error {
	user, err := s.repo.FindUserByID(req.ID)
	if err != nil {
		return ErrUserNotFound
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.DeptID != "" {
		user.DeptID = req.DeptID
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := s.repo.UpdateUser(user); err != nil {
		return err
	}

	// Update account if password or status changed
	if req.Password != "" || req.Status != nil || req.Username != "" {
		account, err := s.repo.FindAccountByUserID(req.ID)
		if err != nil {
			return err
		}
		if req.Password != "" {
			if req.OldPassword != "" {
				if err := bcrypt.CompareHashAndPassword([]byte(account.AuthSecret), []byte(req.OldPassword)); err != nil {
					return ErrInvalidPassword
				}
			}
			hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			account.AuthSecret = string(hashed)
		}
		if req.Status != nil {
			account.Status = *req.Status
		}
		if req.Username != "" {
			account.AuthAccount = req.Username
		}
		if err := s.repo.UpdateAccount(account); err != nil {
			return err
		}
	}

	if req.Roles != nil {
		_ = s.repo.BindRoles(req.ID, req.Roles)
	}
	return nil
}

// DeleteUser removes a user and its account.
func (s *UserService) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}

// UpdatePassword changes a user's password after verifying the old one.
func (s *UserService) UpdatePassword(userID, oldPwd, newPwd string) error {
	account, err := s.repo.FindAccountByUserID(userID)
	if err != nil {
		return ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(account.AuthSecret), []byte(oldPwd)); err != nil {
		return ErrInvalidPassword
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	account.AuthSecret = string(hashed)
	return s.repo.UpdateAccount(account)
}

// BindRoles replaces role assignments for a user.
func (s *UserService) BindRoles(userID string, roleIDs []string) error {
	return s.repo.BindRoles(userID, roleIDs)
}

// GetUserOverview returns summary statistics for the current user.
func (s *UserService) GetUserOverview(userID string) (*dto.UserOverview, error) {
	totalProjects, _ := s.repo.CountProjectsByUser(userID)
	totalAnswers, _ := s.repo.CountAnswersByUser(userID)
	todayAnswers, _ := s.repo.CountTodayAnswersByUser(userID)
	return &dto.UserOverview{
		TotalProjects: int(totalProjects),
		TotalAnswers:  int(totalAnswers),
		TodayAnswers:  int(todayAnswers),
	}, nil
}

// GetRegisterRoles returns roles that are available for self-registration.
func (s *UserService) GetRegisterRoles() ([]dto.RegisterRoleView, error) {
	roles, err := s.repo.ListRegisterRoles()
	if err != nil {
		return nil, err
	}
	views := make([]dto.RegisterRoleView, 0, len(roles))
	for _, r := range roles {
		views = append(views, dto.RegisterRoleView{ID: r.ID, Name: r.Name, Code: r.Code})
	}
	return views, nil
}

// FindUserIDByUsername returns the user ID for the given auth_account (username).
func (s *UserService) FindUserIDByUsername(username string) (string, error) {
	account, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", err
	}
	return account.UserID, nil
}

// GetSimpleUserByID returns a simplified user view by ID.
func (s *UserService) GetSimpleUserByID(userID string) *dto.SimpleUserView {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil
	}
	return &dto.SimpleUserView{
		ID:   user.ID,
		Name: user.Name,
	}
}

// GetUserTasks returns paginated pending flow tasks for the current user.
func (s *UserService) GetUserTasks(userID string, query dto.MyTaskQuery) (*dto.PageResponse[dto.MyTaskView], error) {
	pageIndex := query.PageIndex
	if pageIndex < 1 {
		pageIndex = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (pageIndex - 1) * pageSize
	answers, total, err := s.repo.ListPendingAnswersPaged(userID, offset, pageSize)
	if err != nil {
		return nil, err
	}
	views := make([]dto.MyTaskView, 0, len(answers))
	for _, a := range answers {
		views = append(views, dto.MyTaskView{
			ID:        a.ID,
			ProjectID: a.ProjectID,
			Status:    0,
			CreatedAt: a.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &dto.PageResponse[dto.MyTaskView]{List: views, Total: total}, nil
}

// GetHistoryTasks returns paginated completed flow tasks for the current user.
func (s *UserService) GetHistoryTasks(userID string, query dto.MyTaskQuery) (*dto.PageResponse[dto.MyTaskView], error) {
	pageIndex := query.PageIndex
	if pageIndex < 1 {
		pageIndex = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (pageIndex - 1) * pageSize

	answers, total, err := s.repo.ListHistoryAnswers(userID, offset, pageSize)
	if err != nil {
		return nil, err
	}
	views := make([]dto.MyTaskView, 0, len(answers))
	for _, a := range answers {
		views = append(views, dto.MyTaskView{
			ID:        a.ID,
			ProjectID: a.ProjectID,
			Status:    1,
			CreatedAt: a.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &dto.PageResponse[dto.MyTaskView]{List: views, Total: total}, nil
}

// GetUserAuthorities returns all authority strings for a user based on their roles.
func (s *UserService) GetUserAuthorities(userID string) ([]string, error) {
	roles, err := s.repo.GetUserRoles(userID)
	if err != nil {
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
	auths := make([]string, 0, len(authSet))
	for a := range authSet {
		auths = append(auths, a)
	}
	return auths, nil
}

// userToView converts model.User to dto.UserView.
func userToView(user *model.User, account *model.Account, roles []model.Role, deptName *string) *dto.UserView {
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
	if account != nil {
		view.Username = account.AuthAccount
	}
	if deptName != nil {
		view.DeptName = *deptName
	}
	if roles != nil {
		view.Roles = make([]dto.RoleView, 0, len(roles))
		for _, r := range roles {
			view.Roles = append(view.Roles, dto.RoleView{
				ID:   r.ID,
				Name: r.Name,
				Code: r.Code,
			})
		}
	}
	return view
}

package dto

// UserView is the response DTO for user data.
type UserView struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Username string     `json:"username"`
	Phone    string     `json:"phone"`
	Email    string     `json:"email"`
	Avatar   string     `json:"avatar"`
	Gender   string     `json:"gender"`
	DeptID   string     `json:"deptId"`
	DeptName string     `json:"deptName"`
	Status   int        `json:"status"`
	Roles    []RoleView `json:"roles"`
}

// RoleView is a minimal role representation used in UserView.
type RoleView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// LoginRequest holds credentials for the login endpoint.
type LoginRequest struct {
	Username            string `json:"username" binding:"required"`
	Password            string `json:"password" binding:"required"`
	CaptchaVerification string `json:"captchaVerification"`
}

// LoginResponse is returned on successful login.
type LoginResponse struct {
	Token string   `json:"token"`
	User  UserView `json:"user"`
}

// RegisterRequest holds data for self-registration.
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

// CreateUserRequest holds data for admin user creation.
type CreateUserRequest struct {
	Username string   `json:"username" binding:"required"`
	Password string   `json:"password" binding:"required"`
	Name     string   `json:"name" binding:"required"`
	Phone    string   `json:"phone"`
	Email    string   `json:"email"`
	Gender   string   `json:"gender"`
	DeptID   string   `json:"deptId"`
	Status   int      `json:"status"`
	Roles    []string `json:"roles"`
}

// UpdateUserRequest holds data for updating a user.
type UpdateUserRequest struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Phone       string   `json:"phone"`
	Email       string   `json:"email"`
	Gender      string   `json:"gender"`
	Avatar      string   `json:"avatar"`
	DeptID      string   `json:"deptId"`
	Status      *int     `json:"status"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	OldPassword string   `json:"oldPassword"`
	Roles       []string `json:"roles"`
}

// BindRoleRequest holds user-role binding data.
type BindRoleRequest struct {
	UserID  string   `json:"userId" binding:"required"`
	RoleIDs []string `json:"roleIds"`
}

// UserQueryRequest holds query params for listing users.
type UserQueryRequest struct {
	PageRequest
	Name   string `form:"name"`
	RoleID string `form:"roleId"`
	DeptID string `form:"deptId"`
}

// UserOverview holds summary statistics for the current user's dashboard.
type UserOverview struct {
	TotalProjects int `json:"totalProjects"`
	TotalAnswers  int `json:"totalAnswers"`
	TodayAnswers  int `json:"todayAnswers"`
}

// RegisterRoleView is a role available for self-registration.
type RegisterRoleView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// MyTaskQuery holds query params for listing user tasks.
type MyTaskQuery struct {
	PageRequest
	ProjectID string `form:"projectId"`
}

// MyTaskView is a single task entry in the user's task list.
type MyTaskView struct {
	ID          string `json:"id"`
	ProjectID   string `json:"projectId"`
	ProjectName string `json:"projectName"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"createAt"`
}

// UpdateUserPositionRequest holds data for updating a user's position assignments.
type UpdateUserPositionRequest struct {
	ID          string   `json:"id"`
	PositionIDs []string `json:"positionIds"`
	DeptID      string   `json:"deptId"`
}

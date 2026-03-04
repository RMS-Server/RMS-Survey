package model

import (
	"time"
)

// Account maps to t_account.
type Account struct {
	BaseModel
	UserType    string `gorm:"column:user_type;size:100;not null;default:SysUser" json:"userType"`
	UserID      string `gorm:"column:user_id;size:64;not null" json:"userId"`
	AuthType    string `gorm:"column:auth_type;size:20;not null;default:PWD" json:"authType"`
	AuthAccount string `gorm:"column:auth_account;size:100;not null" json:"authAccount"`
	AuthSecret  string `gorm:"column:auth_secret;size:64" json:"-"`
	SecretSalt  string `gorm:"column:secret_salt;size:32" json:"-"`
	Status      int    `gorm:"column:status;not null;default:1" json:"status"`
}

func (Account) TableName() string { return "t_account" }

// Answer maps to t_answer.
type Answer struct {
	BaseModel
	ProjectID        string     `gorm:"column:project_id;size:64;not null" json:"projectId"`
	TempAnswer       string     `gorm:"column:temp_answer;type:text" json:"tempAnswer"`
	Survey           string     `gorm:"column:survey;type:longtext" json:"survey"`
	Answer           string     `gorm:"column:answer;type:text" json:"answer"`
	Attachment       string     `gorm:"column:attachment;size:1024" json:"attachment"`
	MetaInfo         string     `gorm:"column:meta_info;type:text" json:"metaInfo"`
	TempSave         *int       `gorm:"column:temp_save" json:"tempSave"`
	ExamInfo         string     `gorm:"column:exam_info;type:text" json:"examInfo"`
	ExamExerciseType string     `gorm:"column:exam_exercise_type;size:4" json:"examExerciseType"`
	ExamScore        *float32   `gorm:"column:exam_score" json:"examScore"`
	RepoID           string     `gorm:"column:repo_id;size:256" json:"repoId"`
	IsRead           *bool      `gorm:"column:is_read;default:false" json:"isRead"`
	ReadAt           *time.Time `gorm:"column:read_at" json:"readAt"`
	ReadBy           string     `gorm:"column:read_by;size:256" json:"readBy"`
	IPAddress        string     `gorm:"column:ip_address;size:45;index" json:"ipAddress"`
}

func (Answer) TableName() string { return "t_answer" }

// CommDictItem maps to t_comm_dict_item (composite PK: id + item_value).
type CommDictItem struct {
	ID              string    `gorm:"primaryKey;size:64" json:"id"`
	DictCode        string    `gorm:"column:dict_code;size:256" json:"dictCode"`
	ItemName        string    `gorm:"column:item_name;size:256" json:"itemName"`
	ItemValue       string    `gorm:"primaryKey;column:item_value;size:256;not null" json:"itemValue"`
	ItemOrder       *int      `gorm:"column:item_order" json:"itemOrder"`
	ItemLevel       *int      `gorm:"column:item_level" json:"itemLevel"`
	ParentItemValue string    `gorm:"column:parent_item_value;size:64" json:"parentItemValue"`
	CreatedAt       time.Time `gorm:"column:create_at" json:"createAt"`
	CreateBy        string    `gorm:"column:create_by;size:256" json:"createBy"`
	UpdatedAt       time.Time `gorm:"column:update_at" json:"updateAt"`
	UpdateBy        string    `gorm:"column:update_by;size:256" json:"updateBy"`
}

func (CommDictItem) TableName() string { return "t_comm_dict_item" }

// File maps to t_file.
type File struct {
	BaseModel
	OriginalName  string `gorm:"column:original_name;size:256" json:"originalName"`
	FileName      string `gorm:"column:file_name;size:256" json:"fileName"`
	FilePath      string `gorm:"column:file_path;size:512" json:"filePath"`
	ThumbFilePath string `gorm:"column:thumb_file_path;size:512" json:"thumbFilePath"`
	StorageType   *int   `gorm:"column:storage_type" json:"storageType"`
	Shared        int    `gorm:"column:shared;default:0" json:"shared"`
}

func (File) TableName() string { return "t_file" }

// Project maps to t_project.
type Project struct {
	BaseModel
	ParentID string `gorm:"column:parent_id;size:64;default:0" json:"parentId"`
	Name     string `gorm:"column:name;type:text" json:"name"`
	Survey   string `gorm:"column:survey;type:longtext" json:"survey"`
	Setting  string `gorm:"column:setting;type:text" json:"setting"`
	Status   int    `gorm:"column:status;default:0" json:"status"`
	Mode     string `gorm:"column:mode;size:32" json:"mode"`
	Priority int    `gorm:"column:priority;default:1000" json:"priority"`
}

func (Project) TableName() string { return "t_project" }

// ProjectPartner maps to t_project_partner.
type ProjectPartner struct {
	BaseModelNoSoftDelete
	UID            string `gorm:"column:uid;size:64" json:"uid"`
	ProjectID      string `gorm:"column:project_id;size:64" json:"projectId"`
	Type           *int   `gorm:"column:type" json:"type"`
	Status         int    `gorm:"column:status;default:0" json:"status"`
	UserID         string `gorm:"column:user_id;size:64" json:"userId"`
	UserName       string `gorm:"column:user_name;size:256" json:"userName"`
	GroupID        string `gorm:"column:group_id;size:64" json:"groupId"`
	DataPermission string `gorm:"column:data_permission;type:text" json:"dataPermission"`
	InitialValue   string `gorm:"column:initial_value;type:text" json:"initialValue"`
}

func (ProjectPartner) TableName() string { return "t_project_partner" }

// Role maps to t_role.
type Role struct {
	BaseModel
	Name      string `gorm:"column:name;size:50;not null" json:"name"`
	Code      string `gorm:"column:code;size:50;not null" json:"code"`
	Remark    string `gorm:"column:remark;size:100" json:"remark"`
	Authority string `gorm:"column:authority;size:3000" json:"authority"`
	Status    *bool  `gorm:"column:status;default:1" json:"status"`
}

func (Role) TableName() string { return "t_role" }

// SysInfo maps to t_sys_info.
type SysInfo struct {
	BaseModelNoSoftDelete
	Name         string `gorm:"column:name;size:64" json:"name"`
	Description  string `gorm:"column:description;type:text" json:"description"`
	Avatar       string `gorm:"column:avatar;size:64" json:"avatar"`
	Locale       string `gorm:"column:locale;size:64" json:"locale"`
	Version      string `gorm:"column:version;size:64" json:"version"`
	Setting      string `gorm:"column:setting;type:text" json:"setting"`
	AISetting    string `gorm:"column:ai_setting;type:text" json:"aiSetting"`
	RegisterInfo string `gorm:"column:register_info;type:text" json:"registerInfo"`
	IsDefault    *bool  `gorm:"column:is_default" json:"isDefault"`
}

func (SysInfo) TableName() string { return "t_sys_info" }

// Template maps to t_template.
type Template struct {
	BaseModel
	RepoID       string `gorm:"column:repo_id;size:64" json:"repoId"`
	SerialNo     string `gorm:"column:serial_no;size:256" json:"serialNo"`
	Name         string `gorm:"column:name;size:1024" json:"name"`
	QuestionType string `gorm:"column:question_type;size:64" json:"questionType"`
	TemplateData string `gorm:"column:template;type:longtext" json:"template"`
	Mode         string `gorm:"column:mode;size:32" json:"mode"`
	Category     string `gorm:"column:category;size:256" json:"category"`
	Tag          string `gorm:"column:tag;size:512" json:"tag"`
	Priority     *int   `gorm:"column:priority" json:"priority"`
	PreviewURL   string `gorm:"column:preview_url;size:512" json:"previewUrl"`
	Shared       *bool  `gorm:"column:shared;default:0" json:"shared"`
}

func (Template) TableName() string { return "t_template" }

// User maps to t_user.
type User struct {
	BaseModel
	Name               string     `gorm:"column:name;size:50;not null" json:"name"`
	DeptID             string     `gorm:"column:dept_id;size:20" json:"deptId"`
	Gender             string     `gorm:"column:gender;size:10" json:"gender"`
	Birthday           *time.Time `gorm:"column:birthday;type:date" json:"birthday"`
	Phone              string     `gorm:"column:phone;size:20" json:"phone"`
	Email              string     `gorm:"column:email;size:50" json:"email"`
	Avatar             string     `gorm:"column:avatar;size:200" json:"avatar"`
	Status             int        `gorm:"column:status;not null;default:1" json:"status"`
	Profile            string     `gorm:"column:profile;size:255" json:"profile"`
	CorrectTimes       *int       `gorm:"column:correct_times" json:"correctTimes"`
	SSOID              *string    `gorm:"column:sso_id;size:64;uniqueIndex" json:"ssoId"`
	SSOLastSync        *time.Time `gorm:"column:sso_last_sync" json:"ssoLastSync"`
	SSOPermissionLevel *int       `gorm:"column:sso_permission_level" json:"ssoPermissionLevel"`
	SSOGroupID         *string    `gorm:"column:sso_group_id;size:64" json:"ssoGroupId"`
	SSOGroupName       *string    `gorm:"column:sso_group_name;size:128" json:"ssoGroupName"`
}

func (User) TableName() string { return "t_user" }

// UserRole maps to t_user_role.
type UserRole struct {
	BaseModelNoSoftDelete
	UserType string `gorm:"column:user_type;size:100;not null;default:SysUser" json:"userType"`
	UserID   string `gorm:"column:user_id;size:64;not null" json:"userId"`
	RoleID   string `gorm:"column:role_id;size:64;not null" json:"roleId"`
}

func (UserRole) TableName() string { return "t_user_role" }

// OAuthSession stores refresh tokens for automatic token refresh.
type OAuthSession struct {
	BaseModelNoSoftDelete
	UserID        string     `gorm:"column:user_id;size:64;not null;index" json:"userId"`
	RefreshToken  string     `gorm:"column:refresh_token;size:512;not null;uniqueIndex" json:"-"`
	ExpiresAt     time.Time  `gorm:"column:expires_at;not null" json:"expiresAt"`
	DeviceInfo    string     `gorm:"column:device_info;size:256" json:"deviceInfo"`
	LastRefreshAt *time.Time `gorm:"column:last_refresh_at" json:"lastRefreshAt"`
}

func (OAuthSession) TableName() string { return "t_oauth_session" }

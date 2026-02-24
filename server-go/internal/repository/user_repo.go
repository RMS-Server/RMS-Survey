package repository

import (
	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/model"
	"gorm.io/gorm"
)

// UserRepo handles all user-related database operations.
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo creates a new UserRepo.
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// FindByUsername looks up an account by auth_account.
func (r *UserRepo) FindByUsername(username string) (*model.Account, error) {
	var account model.Account
	err := r.db.Where("auth_account = ?", username).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// FindUserByID retrieves a user by primary key.
func (r *UserRepo) FindUserByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAccountByUserID retrieves an account by user_id.
func (r *UserRepo) FindAccountByUserID(userID string) (*model.Account, error) {
	var account model.Account
	err := r.db.Where("user_id = ?", userID).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// ListUsers returns a paginated list of users with optional name filter.
func (r *UserRepo) ListUsers(req dto.UserQueryRequest) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.DeptID != "" {
		query = query.Where("dept_id = ?", req.DeptID)
	}
	if req.RoleID != "" {
		query = query.Where("id IN (SELECT user_id FROM t_user_role WHERE role_id = ?)", req.RoleID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	pageIndex := req.PageIndex
	if pageIndex < 1 {
		pageIndex = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (pageIndex - 1) * pageSize

	err := query.Offset(offset).Limit(pageSize).Find(&users).Error
	return users, total, err
}

// CreateUser inserts a user and its account in a transaction.
func (r *UserRepo) CreateUser(user *model.User, account *model.Account) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		account.UserID = user.ID
		return tx.Create(account).Error
	})
}

// UpdateUser saves changes to a user record.
func (r *UserRepo) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

// UpdateAccount saves changes to an account record.
func (r *UserRepo) UpdateAccount(account *model.Account) error {
	return r.db.Save(account).Error
}

// DeleteUser removes a user and its account in a transaction.
func (r *UserRepo) DeleteUser(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
			return err
		}
		return tx.Where("user_id = ?", id).Delete(&model.Account{}).Error
	})
}

// GetUserRoles returns all roles assigned to a user.
func (r *UserRepo) GetUserRoles(userID string) ([]model.Role, error) {
	var userRoles []model.UserRole
	if err := r.db.Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return []model.Role{}, nil
	}

	roleIDs := make([]string, len(userRoles))
	for i, ur := range userRoles {
		roleIDs[i] = ur.RoleID
	}

	var roles []model.Role
	err := r.db.Where("id IN ?", roleIDs).Find(&roles).Error
	return roles, err
}

// BindRoles replaces all role assignments for a user.
func (r *UserRepo) BindRoles(userID string, roleIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}
		for _, roleID := range roleIDs {
			ur := model.UserRole{
				UserID:   userID,
				RoleID:   roleID,
				UserType: "SysUser",
			}
			if err := tx.Create(&ur).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// FindSysInfo retrieves a sys_info record by key stored in the name field.
func (r *UserRepo) FindSysInfoByName(name string) (*model.SysInfo, error) {
	var info model.SysInfo
	err := r.db.Where("name = ?", name).First(&info).Error
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// SaveSysInfo creates or updates a SysInfo record.
func (r *UserRepo) SaveSysInfo(info *model.SysInfo) error {
	if info.ID == "" {
		return r.db.Create(info).Error
	}
	return r.db.Save(info).Error
}

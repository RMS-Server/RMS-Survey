package repository

import (
	"github.com/rms-survey/server/internal/model"
	"gorm.io/gorm"
)

// SystemRepo handles data access for system management entities.
type SystemRepo struct {
	db *gorm.DB
}

func NewSystemRepo(db *gorm.DB) *SystemRepo {
	return &SystemRepo{db: db}
}

// --- Role ---

func (r *SystemRepo) ListRoles(name string, offset, limit int) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64
	q := r.db.Model(&model.Role{})
	if name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if limit > 0 {
		q = q.Offset(offset).Limit(limit)
	}
	err := q.Find(&roles).Error
	return roles, total, err
}

func (r *SystemRepo) CreateRole(role *model.Role) error {
	return r.db.Create(role).Error
}

func (r *SystemRepo) UpdateRole(role *model.Role) error {
	return r.db.Save(role).Error
}

func (r *SystemRepo) DeleteRole(id string) error {
	return r.db.Delete(&model.Role{}, "id = ?", id).Error
}

func (r *SystemRepo) GetRoleByID(id string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("id = ?", id).First(&role).Error
	return &role, err
}

func (r *SystemRepo) CountRoles() (int64, error) {
	var count int64
	err := r.db.Model(&model.Role{}).Count(&count).Error
	return count, err
}

// --- SysInfo ---

func (r *SystemRepo) GetSysInfo() (*model.SysInfo, error) {
	var info model.SysInfo
	err := r.db.Where("is_default = ?", true).First(&info).Error
	if err == gorm.ErrRecordNotFound {
		// return empty default
		return &model.SysInfo{}, nil
	}
	return &info, err
}

func (r *SystemRepo) SaveSysInfo(info *model.SysInfo) error {
	return r.db.Save(info).Error
}

// ListDictItemsByCodes returns all dict items whose dict_code is in the given list.
func (r *SystemRepo) ListDictItemsByCodes(codes []string) ([]model.CommDictItem, error) {
	var items []model.CommDictItem
	err := r.db.Where("dict_code IN ?", codes).Order("item_order ASC").Find(&items).Error
	return items, err
}

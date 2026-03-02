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

// --- Dept ---

func (r *SystemRepo) ListDepts() ([]model.Dept, error) {
	var depts []model.Dept
	err := r.db.Order("sort_code ASC").Find(&depts).Error
	return depts, err
}

func (r *SystemRepo) CreateDept(dept *model.Dept) error {
	return r.db.Create(dept).Error
}

func (r *SystemRepo) UpdateDept(dept *model.Dept) error {
	return r.db.Save(dept).Error
}

func (r *SystemRepo) DeleteDept(id string) error {
	return r.db.Delete(&model.Dept{}, "id = ?", id).Error
}

func (r *SystemRepo) CountDepts() (int64, error) {
	var count int64
	err := r.db.Model(&model.Dept{}).Count(&count).Error
	return count, err
}

// --- CommDict ---

func (r *SystemRepo) ListDicts(name, code string, offset, limit int) ([]model.CommDict, int64, error) {
	var dicts []model.CommDict
	var total int64
	q := r.db.Model(&model.CommDict{})
	if name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	if code != "" {
		q = q.Where("code LIKE ?", "%"+code+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if limit > 0 {
		q = q.Offset(offset).Limit(limit)
	}
	err := q.Find(&dicts).Error
	return dicts, total, err
}

func (r *SystemRepo) CreateDict(d *model.CommDict) error {
	return r.db.Create(d).Error
}

func (r *SystemRepo) UpdateDict(d *model.CommDict) error {
	return r.db.Save(d).Error
}

func (r *SystemRepo) DeleteDict(id string) error {
	return r.db.Delete(&model.CommDict{}, "id = ?", id).Error
}

// --- CommDictItem ---

func (r *SystemRepo) ListDictItems(dictCode string, offset, limit int) ([]model.CommDictItem, int64, error) {
	var items []model.CommDictItem
	var total int64
	q := r.db.Model(&model.CommDictItem{})
	if dictCode != "" {
		q = q.Where("dict_code = ?", dictCode)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if limit > 0 {
		q = q.Offset(offset).Limit(limit)
	}
	err := q.Order("item_order ASC").Find(&items).Error
	return items, total, err
}

func (r *SystemRepo) GetDictItemByID(id string) (*model.CommDictItem, error) {
	var item model.CommDictItem
	err := r.db.Where("id = ?", id).First(&item).Error
	return &item, err
}

func (r *SystemRepo) SaveDictItem(item *model.CommDictItem) error {
	return r.db.Save(item).Error
}

// BatchCreateDictItems inserts multiple dict items in a single transaction.
func (r *SystemRepo) BatchCreateDictItems(items []model.CommDictItem) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&items).Error
	})
}

func (r *SystemRepo) DeleteDictItem(id string) error {
	return r.db.Delete(&model.CommDictItem{}, "id = ?", id).Error
}

// ListDictItemsByCodes returns all dict items whose dict_code is in the given list.
func (r *SystemRepo) ListDictItemsByCodes(codes []string) ([]model.CommDictItem, error) {
	var items []model.CommDictItem
	err := r.db.Where("dict_code IN ?", codes).Order("item_order ASC").Find(&items).Error
	return items, err
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

// --- Position ---

func (r *SystemRepo) ListPositions(name string, offset, limit int) ([]model.Position, int64, error) {
	var positions []model.Position
	var total int64
	q := r.db.Model(&model.Position{})
	if name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if limit > 0 {
		q = q.Offset(offset).Limit(limit)
	}
	err := q.Find(&positions).Error
	return positions, total, err
}

func (r *SystemRepo) CreatePosition(p *model.Position) error {
	return r.db.Create(p).Error
}

func (r *SystemRepo) UpdatePosition(p *model.Position) error {
	return r.db.Save(p).Error
}

func (r *SystemRepo) DeletePosition(id string) error {
	return r.db.Delete(&model.Position{}, "id = ?", id).Error
}

// --- Tag ---

func (r *SystemRepo) ListTags(category, entityID string, offset, limit int) ([]model.Tag, int64, error) {
	var tags []model.Tag
	var total int64
	q := r.db.Model(&model.Tag{})
	if category != "" {
		q = q.Where("category = ?", category)
	}
	if entityID != "" {
		q = q.Where("entity_id = ?", entityID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if limit > 0 {
		q = q.Offset(offset).Limit(limit)
	}
	err := q.Find(&tags).Error
	return tags, total, err
}

func (r *SystemRepo) CreateTag(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

func (r *SystemRepo) DeleteTag(id string) error {
	return r.db.Delete(&model.Tag{}, "id = ?", id).Error
}

// UpdateDeptSortCode updates the sort_code for a single department by ID.
func (r *SystemRepo) UpdateDeptSortCode(id string, sortCode int) error {
	return r.db.Model(&model.Dept{}).Where("id = ?", id).Update("sort_code", sortCode).Error
}

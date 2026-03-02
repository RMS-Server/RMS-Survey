package repository

import (
	"github.com/surveyking/server/internal/dto"
	"github.com/surveyking/server/internal/model"
	"gorm.io/gorm"
)

// ProjectRepo handles database operations for projects.
type ProjectRepo struct {
	db *gorm.DB
}

// NewProjectRepo creates a new ProjectRepo.
func NewProjectRepo(db *gorm.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

// filterByDataPerm restricts query to projects the user owns or is a partner of.
// This matches the Java DataPermAspect logic:
//   1. Admin users bypass all filters (isAdmin check).
//   2. Project owner: create_by = userId.
//   3. Project partner: userId exists in t_project_partner for the project.
func filterByDataPerm(db *gorm.DB, userInfo *dto.UserInfo) *gorm.DB {
	if isAdmin(userInfo) {
		return db
	}
	// user sees projects they created or are a partner of
	return db.Where(
		"create_by = ? OR id IN (SELECT project_id FROM t_project_partner WHERE user_id = ?)",
		userInfo.UserID, userInfo.UserID,
	)
}

func isAdmin(userInfo *dto.UserInfo) bool {
	for _, r := range userInfo.Roles {
		if r == "admin" || r == "ADMIN" {
			return true
		}
	}
	return false
}

// ListProjects returns a paginated list of non-deleted projects.
func (r *ProjectRepo) ListProjects(query *dto.ProjectQuery, userInfo *dto.UserInfo) ([]model.Project, int64, error) {
	db := r.db.Model(&model.Project{})
	db = filterByDataPerm(db, userInfo)

	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}
	if query.ParentID != "" {
		db = db.Where("parent_id = ?", query.ParentID)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	pageIndex := query.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	offset := (pageIndex - 1) * pageSize

	var projects []model.Project
	err := db.Order("priority DESC, create_at DESC").Offset(offset).Limit(pageSize).Find(&projects).Error
	return projects, total, err
}

// GetProject returns a single project by ID.
func (r *ProjectRepo) GetProject(id string) (*model.Project, error) {
	var p model.Project
	err := r.db.Where("id = ?", id).First(&p).Error
	return &p, err
}

// CreateProject inserts a new project.
func (r *ProjectRepo) CreateProject(p *model.Project) error {
	return r.db.Create(p).Error
}

// UpdateProject saves changes to an existing project.
func (r *ProjectRepo) UpdateProject(p *model.Project) error {
	return r.db.Save(p).Error
}

// SoftDeleteProject marks a project as deleted.
func (r *ProjectRepo) SoftDeleteProject(id string) error {
	return r.db.Delete(&model.Project{}, "id = ?", id).Error
}

// ListDeleted returns soft-deleted projects.
func (r *ProjectRepo) ListDeleted(userInfo *dto.UserInfo) ([]model.Project, error) {
	db := r.db.Unscoped().Model(&model.Project{}).Where("deleted_at IS NOT NULL")
	db = filterByDataPerm(db, userInfo)
	var projects []model.Project
	err := db.Find(&projects).Error
	return projects, err
}

// HardDeleteProject permanently removes a project.
func (r *ProjectRepo) HardDeleteProject(id string) error {
	return r.db.Unscoped().Delete(&model.Project{}, "id = ?", id).Error
}

// RestoreProject clears the deleted_at field.
func (r *ProjectRepo) RestoreProject(id string) error {
	return r.db.Unscoped().Model(&model.Project{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

// ListPartners returns partners for a project.
func (r *ProjectRepo) ListPartners(query *dto.ProjectPartnerQuery) ([]model.ProjectPartner, int64, error) {
	db := r.db.Model(&model.ProjectPartner{}).Where("project_id = ?", query.ProjectID)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	pageIndex := query.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	offset := (pageIndex - 1) * pageSize

	var partners []model.ProjectPartner
	err := db.Offset(offset).Limit(pageSize).Find(&partners).Error
	return partners, total, err
}

// CreatePartner inserts a new project partner.
func (r *ProjectRepo) CreatePartner(p *model.ProjectPartner) error {
	return r.db.Create(p).Error
}

// DeletePartner removes a project partner by ID.
func (r *ProjectRepo) DeletePartner(id string) error {
	return r.db.Delete(&model.ProjectPartner{}, "id = ?", id).Error
}

// ListPartnersAll returns all partners for a project without pagination.
func (r *ProjectRepo) ListPartnersAll(projectID string) ([]model.ProjectPartner, error) {
	var partners []model.ProjectPartner
	err := r.db.Where("project_id = ?", projectID).Order("create_at ASC").Find(&partners).Error
	return partners, err
}

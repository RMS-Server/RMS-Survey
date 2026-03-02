package handler

import (
	"bytes"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/pkg/response"
	"github.com/rms-survey/server/internal/repository"
	"github.com/rms-survey/server/internal/service"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// ProjectHandler handles all project HTTP endpoints.
type ProjectHandler struct {
	svc         *service.ProjectService
	userSvc     *service.UserService
	systemSvc   *service.SystemService
	templateSvc *service.TemplateService
	repoSvc     *service.RepoService
}

// NewProjectHandler creates a ProjectHandler wired to the given DB.
func NewProjectHandler(db *gorm.DB) *ProjectHandler {
	return &ProjectHandler{
		svc:         service.NewProjectService(repository.NewProjectRepo(db)),
		userSvc:     service.NewUserService(repository.NewUserRepo(db)),
		systemSvc:   service.NewSystemService(repository.NewSystemRepo(db)),
		templateSvc: service.NewTemplateService(repository.NewTemplateRepository(db)),
		repoSvc:     service.NewRepoService(repository.NewRepoRepository(db), repository.NewUserBookRepository(db)),
	}
}

// RegisterRoutes wires all project routes onto the provided gin group.
func (h *ProjectHandler) RegisterRoutes(project gin.IRouter) {
	project.GET("/list", h.ListProjects)
	project.GET("", h.GetProject)
	project.GET("/setting", h.GetSetting)
	project.POST("/create", h.CreateProject)
	project.POST("/update", h.UpdateProject)
	project.POST("/delete", h.DeleteProject)
	project.GET("/trash", h.GetDeleted)
	project.POST("/destroy", h.DestroyProject)
	project.POST("/restore", h.RestoreProject)
	project.GET("/partner/list", h.ListPartners)
	project.POST("/partner/create", h.AddPartner)
	project.POST("/partner/delete", h.RemovePartner)
	project.GET("/partner/download", h.DownloadPartner)
	project.POST("/partner/import", h.ImportPartner)
	project.POST("/selectUser", h.SelectUser)
	project.POST("/selectDept", h.SelectDept)
	project.POST("/selectRole", h.SelectRole)
	project.POST("/selectPosition", h.SelectPosition)
	project.POST("/selectDict", h.SelectDict)
	project.POST("/selectTemplate", h.SelectTemplate)
	project.POST("/selectRepo", h.SelectRepo)
	project.POST("/selectTag", h.SelectTag)
}

func currentUser(c *gin.Context) *dto.UserInfo {
	if v, ok := c.Get("currentUser"); ok {
		if u, ok := v.(*dto.UserInfo); ok {
			return u
		}
	}
	return &dto.UserInfo{}
}

func (h *ProjectHandler) ListProjects(c *gin.Context) {
	var query dto.ProjectQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.ListProjects(&query, currentUser(c))
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.Fail(c, response.CodeError, "id required")
		return
	}
	result, err := h.svc.GetProject(id)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *ProjectHandler) GetSetting(c *gin.Context) {
	projectID := c.Query("projectId")
	if projectID == "" {
		response.Fail(c, response.CodeError, "projectId required")
		return
	}
	result, err := h.svc.GetSetting(projectID)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req dto.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.CreateProject(&req, currentUser(c))
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	var req dto.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.UpdateProject(&req, currentUser(c)); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	var req dto.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.DeleteProject(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *ProjectHandler) GetDeleted(c *gin.Context) {
	result, err := h.svc.GetDeleted(currentUser(c))
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *ProjectHandler) DestroyProject(c *gin.Context) {
	var req dto.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.DestroyProject(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *ProjectHandler) RestoreProject(c *gin.Context) {
	var req dto.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.RestoreProject(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *ProjectHandler) ListPartners(c *gin.Context) {
	var query dto.ProjectPartnerQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	result, err := h.svc.ListPartners(&query)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *ProjectHandler) AddPartner(c *gin.Context) {
	var req dto.ProjectPartnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.AddPartner(&req, currentUser(c)); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *ProjectHandler) RemovePartner(c *gin.Context) {
	var req dto.ProjectPartnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	if err := h.svc.RemovePartner(&req); err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, nil)
}

// DownloadPartner exports project partners to an Excel file.
func (h *ProjectHandler) DownloadPartner(c *gin.Context) {
	projectID := c.Query("projectId")
	if projectID == "" {
		response.Fail(c, response.CodeError, "projectId is required")
		return
	}

	partners, err := h.svc.ListPartnersAll(projectID)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}

	f := excelize.NewFile()
	defer f.Close()
	sheet := "Sheet1"

	// Set headers
	f.SetCellValue(sheet, "A1", "名单")
	f.SetCellValue(sheet, "B1", "状态")

	// Status mapping based on Java AppConsts.ProjectPartnerStatus
	statusMap := map[int]string{
		0: "未答题",
		1: "已答题",
		2: "已通过",
		3: "未通过",
	}

	// Write data rows
	for i, p := range partners {
		row := i + 2
		// For RESPONDENT_SYS_USER type (type=1), show user name from user table
		// For RESPONDENT_IMP_USER type (type=2), show userName field directly
		if p.Type != nil && *p.Type == 1 && p.UserID != "" {
			user := h.userSvc.GetSimpleUserByID(p.UserID)
			if user != nil {
				f.SetCellValue(sheet, fmt.Sprintf("A%d", row), user.Name)
			} else {
				f.SetCellValue(sheet, fmt.Sprintf("A%d", row), p.UserName)
			}
		} else {
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), p.UserName)
		}
		statusStr, ok := statusMap[p.Status]
		if !ok {
			statusStr = fmt.Sprintf("%d", p.Status)
		}
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), statusStr)
	}

	// Set column widths
	f.SetColWidth(sheet, "A", "A", 30)
	f.SetColWidth(sheet, "B", "B", 15)

	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		response.Fail(c, response.CodeError, "failed to generate Excel: "+err.Error())
		return
	}

	c.Header("Content-Disposition", `attachment; filename="partners.xlsx"`)
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

// ImportPartner imports project partners from an uploaded Excel file.
// The Excel must have a header row; column 0 is the username.
func (h *ProjectHandler) ImportPartner(c *gin.Context) {
	projectID := c.Query("projectId")
	if projectID == "" {
		response.Fail(c, response.CodeError, "projectId is required")
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, response.CodeError, "file is required")
		return
	}
	f, err := fileHeader.Open()
	if err != nil {
		response.Fail(c, response.CodeError, "failed to open file")
		return
	}
	defer f.Close()

	xlsx, err := excelize.OpenReader(f)
	if err != nil {
		response.Fail(c, response.CodeError, "failed to parse Excel: "+err.Error())
		return
	}
	defer xlsx.Close()

	sheets := xlsx.GetSheetList()
	if len(sheets) == 0 {
		response.OK(c, nil)
		return
	}
	rows, err := xlsx.GetRows(sheets[0])
	if err != nil {
		response.Fail(c, response.CodeError, "failed to read rows: "+err.Error())
		return
	}

	for i, row := range rows {
		if i == 0 || len(row) == 0 || row[0] == "" {
			continue
		}
		username := row[0]
		userID, err := h.userSvc.FindUserIDByUsername(username)
		if err != nil {
			continue
		}
		_ = h.svc.AddPartner(&dto.ProjectPartnerRequest{
			ProjectID: projectID,
			UserID:    userID,
		}, currentUser(c))
	}
	response.OK(c, nil)
}

// selectKeyword extracts the "name" field from a JSON body for keyword filtering.
// The frontend editor sends {"name": "keyword"} for all select endpoints.
func selectKeyword(c *gin.Context) string {
	var body struct {
		Name string `json:"name"`
	}
	_ = c.ShouldBindJSON(&body)
	return body.Name
}

// SelectUser returns users matching a keyword for the editor collaborator picker.
func (h *ProjectHandler) SelectUser(c *gin.Context) {
	name := selectKeyword(c)
	result, err := h.userSvc.ListUsers(dto.UserQueryRequest{
		PageRequest: dto.PageRequest{PageIndex: 1, PageSize: 20},
		Name:        name,
	})
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result.List)
}

// SelectDept returns the full department tree for the editor.
func (h *ProjectHandler) SelectDept(c *gin.Context) {
	result, err := h.systemSvc.ListDepts()
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

// SelectRole returns roles matching a keyword for the editor.
func (h *ProjectHandler) SelectRole(c *gin.Context) {
	name := selectKeyword(c)
	result, err := h.systemSvc.ListRoles(dto.RoleQuery{
		PageRequest: dto.PageRequest{PageIndex: 1, PageSize: 20},
		Name:        name,
	})
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result.List)
}

// SelectPosition returns positions matching a keyword for the editor.
func (h *ProjectHandler) SelectPosition(c *gin.Context) {
	name := selectKeyword(c)
	result, err := h.systemSvc.ListPositions(dto.PositionQuery{
		PageRequest: dto.PageRequest{PageIndex: 1, PageSize: 20},
		Name:        name,
	})
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result.List)
}

// SelectDict returns dictionaries matching a keyword for the editor.
func (h *ProjectHandler) SelectDict(c *gin.Context) {
	name := selectKeyword(c)
	result, err := h.systemSvc.ListDicts(dto.CommDictQuery{
		PageRequest: dto.PageRequest{PageIndex: 1, PageSize: 20},
		Name:        name,
	})
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result.List)
}

// SelectTemplate returns a paginated list of templates matching a keyword for the editor.
func (h *ProjectHandler) SelectTemplate(c *gin.Context) {
	var query dto.TemplateQuery
	_ = c.ShouldBindJSON(&query)
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	if query.PageIndex <= 0 {
		query.PageIndex = 1
	}
	result, err := h.templateSvc.List(query)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result)
}

// SelectRepo returns repos matching a keyword for the editor.
func (h *ProjectHandler) SelectRepo(c *gin.Context) {
	name := selectKeyword(c)
	result, err := h.repoSvc.List(dto.RepoQuery{
		PageRequest: dto.PageRequest{PageIndex: 1, PageSize: 20},
		Name:        name,
	})
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result.List)
}

// SelectTag returns tags matching a category/entityId for the editor.
func (h *ProjectHandler) SelectTag(c *gin.Context) {
	var query dto.SystemTagQuery
	_ = c.ShouldBindJSON(&query)
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	if query.PageIndex <= 0 {
		query.PageIndex = 1
	}
	result, err := h.systemSvc.ListTags(query)
	if err != nil {
		response.Fail(c, response.CodeError, err.Error())
		return
	}
	response.OK(c, result.List)
}

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/surveyking/server/internal/config"
	"github.com/surveyking/server/internal/handler"
	"github.com/surveyking/server/internal/middleware"
	"github.com/surveyking/server/internal/pkg/storage"
	"github.com/surveyking/server/internal/repository"
	"github.com/surveyking/server/internal/service"
	"gorm.io/gorm"
)

// Setup creates the gin engine with all route groups wired to handlers.
func Setup(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(corsMiddleware())
	r.Use(middleware.Auth())

	// --- Repositories ---
	systemRepo := repository.NewSystemRepo(db)
	flowRepo := repository.NewFlowRepository(db)
	fileRepo := repository.NewFileRepository(db)
	repoRepo := repository.NewRepoRepository(db)
	userBookRepo := repository.NewUserBookRepository(db)
	templateRepo := repository.NewTemplateRepository(db)
	exerciseRepo := repository.NewExerciseRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)
	reportRepo := repository.NewReportRepository(db)
	userRepo := repository.NewUserRepo(db)

	// --- Services ---
	systemSvc := service.NewSystemService(systemRepo)
	flowSvc := service.NewFlowService(flowRepo)
	localStorage := &storage.LocalStorage{BasePath: config.C.Storage.LocalPath}
	fileSvc := service.NewFileService(fileRepo, localStorage)
	repoSvc := service.NewRepoService(repoRepo, userBookRepo)
	templateSvc := service.NewTemplateService(templateRepo)
	exerciseSvc := service.NewExerciseService(exerciseRepo)
	dashboardSvc := service.NewDashboardService(dashboardRepo)
	reportSvc := service.NewReportService(reportRepo)
	userSvc := service.NewUserService(userRepo)

	// --- Handlers ---
	userH := handler.NewUserHandler(db)
	projectH := handler.NewProjectHandler(db)
	surveyH := handler.NewSurveyHandlerWithFile(db, fileSvc)
	answerH := handler.NewAnswerHandler(db)
	systemH := handler.NewSystemHandlerWithUser(systemSvc, userSvc)
	flowH := handler.NewFlowHandler(flowSvc)
	aiH := handler.NewAIHandler()
	fileH := handler.NewFileHandler(fileSvc)
	repoH := handler.NewRepoHandler(repoSvc)
	templateH := handler.NewTemplateHandler(templateSvc)
	exerciseH := handler.NewExerciseHandler(exerciseSvc)
	dashboardH := handler.NewDashboardHandler(dashboardSvc)
	reportH := handler.NewReportHandler(reportSvc)

	// --- Route groups ---
	public := r.Group("/api/public")
	captcha := r.Group("/captcha")
	api := r.Group("/api")

	// User routes (public + captcha + protected + root)
	userGrp := api.Group("/user")
	userH.RegisterRoutes(public, captcha, userGrp, api)

	// Survey public routes + protected survey group
	surveyGrp := api.Group("/survey")
	surveyH.RegisterRoutes(public, surveyGrp)

	// Project routes
	projectGrp := api.Group("/project")
	projectH.RegisterRoutes(projectGrp)

	// Answer routes
	answerGrp := api.Group("/answer")
	answerH.RegisterRoutes(answerGrp)

	// System routes
	systemGrp := api.Group("/system")
	systemH.RegisterRoutes(systemGrp)

	// Flow routes — /api/flow (legacy) and /api/workflow (Java-aligned)
	registerFlowRoutes := func(grp *gin.RouterGroup) {
		grp.GET("", flowH.GetFlow)
		grp.POST("/save", flowH.SaveFlow)
		grp.POST("/deploy", flowH.Deploy)
		grp.GET("/auditRecord", flowH.GetAuditRecord)
		grp.GET("/tasks", flowH.GetFlowTasks)
		grp.GET("/revertNodes", flowH.GetRevertNodes)
		grp.POST("/approvalTask", flowH.ApprovalTask)
		grp.GET("/statics", flowH.Statics)
		grp.GET("/getTaskInfo", flowH.GetTaskInfo)
		grp.GET("/loadSchema", flowH.LoadSchema)
	}
	registerFlowRoutes(api.Group("/flow"))
	registerFlowRoutes(api.Group("/workflow"))

	// AI chat routes
	aiGrp := api.Group("/ai/chat")
	{
		aiGrp.GET("/models", aiH.GetModels)
		aiGrp.POST("/conversation", aiH.CreateConversation)
		aiGrp.POST("/create-conversation", aiH.CreateConversation) // Java-aligned alias
		aiGrp.POST("/conversation/close", aiH.CloseConversation)
		aiGrp.POST("/close-conversation", aiH.CloseConversation) // Java-aligned alias
		aiGrp.GET("/stream", aiH.Stream)
	}

	// File routes
	fileGrp := api.Group("/file")
	{
		fileGrp.GET("", fileH.GetFile)
		fileGrp.GET("/list", fileH.ListFiles)
		fileGrp.POST("/upload", fileH.Upload)
		fileGrp.POST("/create", fileH.Upload) // Java uses /create, Go uses /upload
		fileGrp.POST("/delete", fileH.DeleteFile)
		fileGrp.GET("/downloadTemplate", fileH.DownloadTemplate)
	}

	// Repo routes
	repoGrp := api.Group("/repo")
	{
		repoGrp.GET("/list", repoH.List)
		repoGrp.POST("/create", repoH.Create)
		repoGrp.POST("/update", repoH.Update)
		repoGrp.POST("/delete", repoH.Delete)
		repoGrp.POST("/batchCreate", repoH.BatchCreate)
		repoGrp.POST("/unbind", repoH.Unbind)
		repoGrp.POST("/pick", repoH.Pick)
		repoGrp.POST("/import", repoH.Import)
		repoGrp.GET("/export", repoH.Export)
		repoGrp.GET("/userBook/list", repoH.ListUserBook)
		repoGrp.POST("/userBook/create", repoH.CreateUserBook)
		repoGrp.POST("/userBook/update", repoH.UpdateUserBook)
		repoGrp.POST("/userBook/delete", repoH.DeleteUserBook)
		// Java uses /book/* path
		repoGrp.GET("/book/list", repoH.ListUserBook)
		repoGrp.POST("/book/create", repoH.CreateUserBook)
		repoGrp.POST("/book/update", repoH.UpdateUserBook)
		repoGrp.POST("/book/delete", repoH.DeleteUserBook)
	}

	// Template routes
	templateGrp := api.Group("/template")
	{
		templateGrp.GET("/list", templateH.List)
		templateGrp.GET("", templateH.Get)
		templateGrp.POST("/create", templateH.Create)
		templateGrp.POST("/update", templateH.Update)
		templateGrp.POST("/delete", templateH.Delete)
		templateGrp.GET("/category/list", templateH.ListCategory)
		templateGrp.GET("/tag/list", templateH.ListTag)
		// Java uses /listCategory and /listTag
		templateGrp.GET("/listCategory", templateH.ListCategory)
		templateGrp.GET("/listTag", templateH.ListTag)
		templateGrp.GET("/get", templateH.Get)
	}

	// Exercise routes
	exerciseGrp := api.Group("/exercise")
	{
		exerciseGrp.GET("/list", exerciseH.List)
		exerciseGrp.GET("/detail", exerciseH.GetDetail)
	}

	// Dashboard routes
	dashboardGrp := api.Group("/dashboard")
	{
		dashboardGrp.GET("/list", dashboardH.List)
		dashboardGrp.POST("/save", dashboardH.Save)
		dashboardGrp.POST("/delete", dashboardH.Delete)
	}

	// Report routes
	reportGrp := api.Group("/report")
	{
		reportGrp.GET("/:shortId", reportH.GetData)
	}

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

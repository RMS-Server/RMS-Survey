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

	// --- Handlers ---
	userH := handler.NewUserHandler(db)
	projectH := handler.NewProjectHandler(db)
	surveyH := handler.NewSurveyHandler(db)
	answerH := handler.NewAnswerHandler(db)
	systemH := handler.NewSystemHandler(systemSvc)
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

	// Flow routes
	flowGrp := api.Group("/flow")
	{
		flowGrp.GET("", flowH.GetFlow)
		flowGrp.POST("/save", flowH.SaveFlow)
		flowGrp.POST("/deploy", flowH.Deploy)
		flowGrp.GET("/auditRecord", flowH.GetAuditRecord)
		flowGrp.GET("/tasks", flowH.GetFlowTasks)
		flowGrp.GET("/revertNodes", flowH.GetRevertNodes)
		flowGrp.POST("/approvalTask", flowH.ApprovalTask)
		flowGrp.GET("/statics", flowH.Statics)
	}

	// AI chat routes
	aiGrp := api.Group("/ai/chat")
	{
		aiGrp.GET("/models", aiH.GetModels)
		aiGrp.POST("/conversation", aiH.CreateConversation)
		aiGrp.POST("/conversation/close", aiH.CloseConversation)
		aiGrp.GET("/stream", aiH.Stream)
	}

	// File routes
	fileGrp := api.Group("/file")
	{
		fileGrp.GET("", fileH.GetFile)
		fileGrp.GET("/list", fileH.ListFiles)
		fileGrp.POST("/upload", fileH.Upload)
		fileGrp.POST("/delete", fileH.DeleteFile)
	}

	// Repo routes
	repoGrp := api.Group("/repo")
	{
		repoGrp.GET("/list", repoH.List)
		repoGrp.POST("/create", repoH.Create)
		repoGrp.POST("/update", repoH.Update)
		repoGrp.POST("/delete", repoH.Delete)
		repoGrp.GET("/userBook/list", repoH.ListUserBook)
		repoGrp.POST("/userBook/create", repoH.CreateUserBook)
		repoGrp.POST("/userBook/update", repoH.UpdateUserBook)
		repoGrp.POST("/userBook/delete", repoH.DeleteUserBook)
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
	}

	// Exercise routes
	exerciseGrp := api.Group("/exercise")
	{
		exerciseGrp.GET("/list", exerciseH.List)
	}

	// Dashboard routes
	dashboardGrp := api.Group("/dashboard")
	{
		dashboardGrp.GET("/list", dashboardH.List)
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

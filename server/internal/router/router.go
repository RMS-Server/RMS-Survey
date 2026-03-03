package router

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rms-survey/server/internal/config"
	"github.com/rms-survey/server/internal/handler"
	"github.com/rms-survey/server/internal/middleware"
	"github.com/rms-survey/server/internal/pkg/storage"
	"github.com/rms-survey/server/internal/repository"
	"github.com/rms-survey/server/internal/service"
	"gorm.io/gorm"
)

// Setup creates the gin engine with all route groups wired to handlers.
func Setup(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(corsMiddleware())
	r.Use(middleware.Auth())

	// Serve frontend static files
	staticPath := "./static"
	r.Static("/assets", staticPath+"/assets")
	r.StaticFile("/favicon.ico", staticPath+"/favicon.ico")

	// SPA fallback - serve index.html for non-API routes, try static files first
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// Try to serve static file
		if len(path) > 1 {
			filePath := staticPath + path
			if _, err := os.Stat(filePath); err == nil {
				c.File(filePath)
				return
			}
		}
		// Fallback to index.html for SPA
		c.File(staticPath + "/index.html")
	})

	// --- Repositories ---
	systemRepo := repository.NewSystemRepo(db)
	fileRepo := repository.NewFileRepository(db)
	templateRepo := repository.NewTemplateRepository(db)
	userRepo := repository.NewUserRepo(db)

	// --- Services ---
	systemSvc := service.NewSystemService(systemRepo)
	localStorage := &storage.LocalStorage{BasePath: config.C.Storage.LocalPath}
	fileSvc := service.NewFileService(fileRepo, localStorage)
	templateSvc := service.NewTemplateService(templateRepo)
	userSvc := service.NewUserService(userRepo)

	// --- Handlers ---
	userH := handler.NewUserHandler(db)
	projectH := handler.NewProjectHandler(db)
	surveyH := handler.NewSurveyHandlerWithFile(db, fileSvc)
	answerH := handler.NewAnswerHandler(db)
	systemH := handler.NewSystemHandlerWithUser(systemSvc, userSvc)
	fileH := handler.NewFileHandler(fileSvc)
	templateH := handler.NewTemplateHandler(templateSvc)

	// --- Route groups ---
	// Public routes (no auth)
	public := r.Group("/api/public")
	captcha := r.Group("/captcha")

	// Protected API routes
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

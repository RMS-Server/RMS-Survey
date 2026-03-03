package main

import (
	"fmt"
	"log"

	"github.com/rms-survey/server/internal/config"
	"github.com/rms-survey/server/internal/model"
	"github.com/rms-survey/server/internal/router"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto migrate database tables
	if err := db.AutoMigrate(
		&model.User{},
		&model.Account{},
		&model.UserRole{},
		&model.Role{},
		&model.Project{},
		&model.ProjectPartner{},
		&model.Answer{},
		&model.File{},
		&model.Template{},
		&model.SysInfo{},
		&model.CommDictItem{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("database migration completed")

	r := router.Setup(db)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

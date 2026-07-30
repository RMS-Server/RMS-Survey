package service_test

import (
	"encoding/json"
	"testing"

	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/repository"
	"github.com/rms-survey/server/internal/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestCreateProjectDefaultsMissingSettingToEmptyObject(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "test:test@tcp(localhost:3306)/test",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun:                 true,
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("open dry-run database: %v", err)
	}

	svc := service.NewProjectService(repository.NewProjectRepo(db))
	project, err := svc.CreateProject(&dto.ProjectRequest{
		Name:   "Project",
		Survey: json.RawMessage(`{"pages":[]}`),
	}, &dto.UserInfo{UserID: "owner"})
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	if got := string(project.Setting); got != "{}" {
		t.Fatalf("setting = %q, want %q", got, "{}")
	}
}

// Package testutil 提供评分测试共用的内存数据库、测试配置与路由装配。
package testutil

import (
	"io"
	"log/slog"
	"path/filepath"
	"testing"

	"gbevent/internal/config"
	"gbevent/internal/handler"
	"gbevent/internal/model"
	"gbevent/internal/repository"
	"gbevent/internal/router"
	"gbevent/internal/service"
	"gbevent/internal/util"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var lastDB *gorm.DB

// DBOf 返回最近一次 NewDB/NewRouter 使用的数据库。
func DBOf(t *testing.T) *gorm.DB {
	t.Helper()
	if lastDB == nil {
		t.Fatalf("no test db created")
	}
	return lastDB
}

// NewDB 打开文件型 sqlite（WAL + 连接池）并迁移全部表，事务内可再取连接。
func NewDB(t *testing.T) *gorm.DB {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")
	dsn := "file:" + path + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(10000)"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(4)
	lastDB = db
	if err := db.AutoMigrate(
		&model.User{}, &model.Activity{}, &model.Registration{}, &model.CheckInRecord{},
		&model.Comment{}, &model.Favorite{}, &model.Notification{}, &model.AuditLog{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

// TestConfig 返回固定测试配置。
func TestConfig(t *testing.T) *config.Config {
	t.Helper()
	return &config.Config{
		AppEnv:             "test",
		ServerPort:         "18080",
		DBHost:             "127.0.0.1",
		DBPort:             "3306",
		DBName:             "gbevent_db",
		DBUser:             "test",
		DBPassword:         "test",
		JWTSecret:          "test-secret-123456",
		JWTExpireHours:     72,
		RateLimitPerMinute: 10000,
		UploadDir:          t.TempDir(),
		UploadMaxMB:        1,
		CORSOrigins:        []string{"http://localhost:18506"},
	}
}

// NewLogger 返回丢弃输出的测试日志器。
func NewLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// Token 生成测试 JWT。
func Token(t *testing.T, cfg *config.Config, userID uint64, username, role string) string {
	t.Helper()
	tok, err := util.GenerateToken(cfg.JWTSecret, cfg.JWTExpireDuration(), userID, username, role)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	return tok
}

// NewRouter 装配完整路由（全部服务挂载 sqlite 内存库）。
func NewRouter(t *testing.T) *router.Router {
	t.Helper()
	db := NewDB(t)
	cfg := TestConfig(t)
	logger := NewLogger()
	userRepo := repository.NewUserRepository(db)
	activityRepo := repository.NewActivityRepository(db)
	regRepo := repository.NewRegistrationRepository(db)
	checkinRepo := repository.NewCheckInRecordRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	favoriteRepo := repository.NewFavoriteRepository(db)
	notifyRepo := repository.NewNotificationRepository(db)

	userSvc := service.NewUserService(userRepo, logger)
	activitySvc := service.NewActivityService(db, activityRepo, regRepo, commentRepo, notifyRepo, checkinRepo, logger)
	regSvc := service.NewRegistrationService(db, regRepo, activitySvc, notifyRepo, logger)
	checkinSvc := service.NewCheckInRecordService(db, checkinRepo, regRepo, activitySvc, notifyRepo, logger)
	commentSvc := service.NewCommentService(commentRepo, activitySvc, logger)
	favoriteSvc := service.NewFavoriteService(favoriteRepo, activitySvc, logger)
	notifySvc := service.NewNotificationService(notifyRepo, logger)

	return router.New(cfg, db, logger,
		handler.NewUserHandler(userSvc, logger),
		handler.NewActivityHandler(activitySvc, logger),
		handler.NewRegistrationHandler(regSvc, logger),
		handler.NewCheckInRecordHandler(checkinSvc, logger),
		handler.NewCommentHandler(commentSvc, logger),
		handler.NewFavoriteHandler(favoriteSvc, logger),
		handler.NewNotificationHandler(notifySvc, logger),
		handler.NewUploadHandler(cfg, logger),
	)
}

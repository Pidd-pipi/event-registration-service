// Code generated for scoring tests of record 004.
package z4_test

import (
	"gbevent/internal/constants"
	"gbevent/internal/model"
	"gbevent/internal/repository"
	"gbevent/internal/service"
	"gbevent/internal/testutil"
	"gorm.io/gorm"
	"testing"
	"time"
)




type actFixture struct {
	db       *gorm.DB
	actSvc   *service.ActivityService
	activity model.Activity
	other    model.Activity
}

func newActFixture(t *testing.T) *actFixture {
	db := testutil.NewDB(t)
	logger := testutil.NewLogger()
	actRepo := repository.NewActivityRepository(db)
	regRepo := repository.NewRegistrationRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	notifyRepo := repository.NewNotificationRepository(db)
	ckRepo := repository.NewCheckInRecordRepository(db)
	actSvc := service.NewActivityService(db, actRepo, regRepo, commentRepo, notifyRepo, ckRepo, logger)

	now := time.Now()
	act, err := actSvc.Create(2, "要删的活动", "desc", "", constants.ActivityTypeLecture, "某地",
		now.AddDate(0, 0, 1), now.AddDate(0, 0, 1), now.AddDate(0, 0, 2), 10, constants.ActivityStatusPublished)
	if err != nil {
		t.Fatalf("create activity: %v", err)
	}
	other, err := actSvc.Create(2, "保留的活动", "desc", "", constants.ActivityTypeLecture, "某地",
		now.AddDate(0, 0, 3), now.AddDate(0, 0, 3), now.AddDate(0, 0, 4), 10, constants.ActivityStatusPublished)
	if err != nil {
		t.Fatalf("create other activity: %v", err)
	}
	return &actFixture{db: db, actSvc: actSvc, activity: *act, other: *other}
}

func addReg(t *testing.T, db *gorm.DB, actID, userID uint64, voucher string) model.Registration {
	t.Helper()
	reg := model.Registration{
		ActivityID: actID, UserID: userID, Name: "张三", Phone: "13900000001", VoucherNo: voucher,
		Status: constants.RegistrationStatusRegistered, ReviewStatus: constants.ReviewStatusApproved,
	}
	if err := db.Create(&reg).Error; err != nil {
		t.Fatalf("create registration: %v", err)
	}
	return reg
}

func TestDeleteActivityRemovesSignups(t *testing.T) {
	f := newActFixture(t)
	addReg(t, f.db, f.activity.ID, 3, "GB20260816000011")
	addReg(t, f.db, f.activity.ID, 4, "GB20260816000012")
	addReg(t, f.db, f.other.ID, 5, "GB20260816000013")
	if err := f.actSvc.Delete(f.activity.ID, 2, constants.RoleOrganizer); err != nil {
		t.Fatalf("delete activity: %v", err)
	}
	var n int64
	if err := f.db.Model(&model.Registration{}).Where("activity_id = ?", f.activity.ID).Count(&n).Error; err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("expected registrations removed, still %d left", n)
	}
	var kept int64
	if err := f.db.Model(&model.Registration{}).Where("activity_id = ?", f.other.ID).Count(&kept).Error; err != nil {
		t.Fatal(err)
	}
	if kept != 1 {
		t.Fatalf("expected other activity registrations untouched, count=%d", kept)
	}
}

func TestDeleteActivityRemovesComments(t *testing.T) {
	f := newActFixture(t)
	cmtRepo := repository.NewCommentRepository(f.db)
	if err := cmtRepo.Create(&model.Comment{ActivityID: f.activity.ID, UserID: 3, Rating: 5, Content: "很好"}); err != nil {
		t.Fatal(err)
	}
	if err := cmtRepo.Create(&model.Comment{ActivityID: f.activity.ID, UserID: 4, Rating: 4, Content: "不错"}); err != nil {
		t.Fatal(err)
	}
	if err := f.actSvc.Delete(f.activity.ID, 2, constants.RoleOrganizer); err != nil {
		t.Fatalf("delete activity: %v", err)
	}
	var n int64
	if err := f.db.Model(&model.Comment{}).Where("activity_id = ?", f.activity.ID).Count(&n).Error; err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("expected comments removed, still %d left", n)
	}
}

func TestDeleteActivityRemovesCheckIns(t *testing.T) {
	f := newActFixture(t)
	reg := addReg(t, f.db, f.activity.ID, 3, "GB20260816000021")
	ckRepo := repository.NewCheckInRecordRepository(f.db)
	if err := ckRepo.Create(&model.CheckInRecord{RegistrationID: reg.ID, ActivityID: f.activity.ID, CheckInMethod: constants.CheckInMethodVoucher, CheckInTime: time.Now(), OperatorID: 2}); err != nil {
		t.Fatal(err)
	}
	if err := f.actSvc.Delete(f.activity.ID, 2, constants.RoleOrganizer); err != nil {
		t.Fatalf("delete activity: %v", err)
	}
	var n int64
	if err := f.db.Model(&model.CheckInRecord{}).Where("activity_id = ?", f.activity.ID).Count(&n).Error; err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("expected check-in records removed, still %d left", n)
	}
}


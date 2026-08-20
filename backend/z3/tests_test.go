// Code generated for scoring tests of record 003.
package z3_test

import (
	"errors"
	"gbevent/internal/constants"
	"gbevent/internal/model"
	"gbevent/internal/repository"
	"gbevent/internal/service"
	"gbevent/internal/testutil"
	"gbevent/internal/util"
	"gorm.io/gorm"
	"testing"
	"time"
)




type regFixture struct {
	db     *gorm.DB
	regSvc *service.RegistrationService
	ckSvc  *service.CheckInRecordService
	act    model.Activity
	reg    model.Registration
}

func newRegFixture(t *testing.T, actStatus, regStatus, reviewStatus string) *regFixture {
	db := testutil.NewDB(t)
	logger := testutil.NewLogger()
	actRepo := repository.NewActivityRepository(db)
	regRepo := repository.NewRegistrationRepository(db)
	notifyRepo := repository.NewNotificationRepository(db)
	ckRepo := repository.NewCheckInRecordRepository(db)
	actSvc := service.NewActivityService(actRepo, regRepo, notifyRepo, ckRepo, logger)
	regSvc := service.NewRegistrationService(db, regRepo, actSvc, notifyRepo, logger)
	ckSvc := service.NewCheckInRecordService(db, ckRepo, regRepo, actSvc, notifyRepo, logger)

	now := time.Now()
	act, err := actSvc.Create(2, "活动", "desc", "", constants.ActivityTypeLecture, "某地",
		now.AddDate(0, 0, 1), now.AddDate(0, 0, 1), now.AddDate(0, 0, 2), 10, actStatus)
	if err != nil {
		t.Fatalf("create activity: %v", err)
	}
	reg := model.Registration{
		ActivityID: act.ID, UserID: 3, Name: "张三", Phone: "13900000001",
		VoucherNo: "GB20260816000001", Status: regStatus, ReviewStatus: reviewStatus,
	}
	if err := regRepo.Create(&reg); err != nil {
		t.Fatalf("create registration: %v", err)
	}
	return &regFixture{db: db, regSvc: regSvc, ckSvc: ckSvc, act: *act, reg: reg}
}

func assertAppErrorCode(t *testing.T, err error, want int) {
	t.Helper()
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.Code != want {
		t.Fatalf("expected AppError code %d, got: %v", want, err)
	}
}

func TestCancelCheckedInRegistrationFails(t *testing.T) {
	f := newRegFixture(t, constants.ActivityStatusPublished, constants.RegistrationStatusCheckedIn, constants.ReviewStatusApproved)
	_, err := f.regSvc.Cancel(f.reg.ID, 3, constants.RoleUser)
	assertAppErrorCode(t, err, constants.CodeCancelConflict)
}

func TestReviewApprovedRegistrationFails(t *testing.T) {
	f := newRegFixture(t, constants.ActivityStatusPublished, constants.RegistrationStatusRegistered, constants.ReviewStatusApproved)
	_, err := f.regSvc.Review(f.reg.ID, 2, constants.RoleOrganizer, constants.ReviewStatusApproved)
	assertAppErrorCode(t, err, constants.CodeReviewConflict)
}

func TestCheckinCancelledRegistrationFails(t *testing.T) {
	f := newRegFixture(t, constants.ActivityStatusPublished, constants.RegistrationStatusCancelled, constants.ReviewStatusApproved)
	_, err := f.ckSvc.CheckInByVoucher(f.act.ID, 2, f.reg.VoucherNo)
	assertAppErrorCode(t, err, constants.CodeCancelConflict)
}

func TestReviewEndedActivityFails(t *testing.T) {
	f := newRegFixture(t, constants.ActivityStatusEnded, constants.RegistrationStatusRegistered, constants.ReviewStatusPending)
	_, err := f.regSvc.Review(f.reg.ID, 2, constants.RoleOrganizer, constants.ReviewStatusApproved)
	assertAppErrorCode(t, err, constants.CodeConflict)
}

func TestResignupAfterCancelSucceeds(t *testing.T) {
	f := newRegFixture(t, constants.ActivityStatusPublished, constants.RegistrationStatusRegistered, constants.ReviewStatusApproved)
	if _, err := f.regSvc.Cancel(f.reg.ID, 3, constants.RoleUser); err != nil {
		t.Fatalf("cancel should succeed: %v", err)
	}
	reg, err := f.regSvc.Create(f.act.ID, 3, "张三", "13900000001", "再来一次")
	if err != nil {
		t.Fatalf("re-signup after cancel should succeed, got: %v", err)
	}
	if reg == nil || reg.Status != constants.RegistrationStatusRegistered {
		t.Fatalf("unexpected registration: %+v", reg)
	}
}

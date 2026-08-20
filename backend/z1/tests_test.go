// Code generated for scoring tests of record 001.
package z1_test

import (
	"errors"
	"gbevent/internal/constants"
	"gbevent/internal/repository"
	"gbevent/internal/service"
	"gbevent/internal/testutil"
	"gbevent/internal/util"
	"net/http"
	"net/http/httptest"
	"testing"
)




func TestFreshUserSignupSucceeds(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo, testutil.NewLogger())
	u, err := svc.Register("freshman", "123456", "新同学", "", "", constants.RoleUser)
	if err != nil {
		t.Fatalf("register new phone should succeed, got: %v", err)
	}
	if u == nil || u.Username != "freshman" {
		t.Fatalf("unexpected user: %+v", u)
	}
}

func TestLoginUnknownUserReturnsInvalidCredentials(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo, testutil.NewLogger())
	_, _, err := svc.Login("secret", 72, "ghost", "whatever")
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.Code != constants.CodeInvalidCredentials {
		t.Fatalf("expected invalid credentials AppError, got: %v", err)
	}
}

func TestGetUserByIDMissingReturnsNotFound(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo, testutil.NewLogger())
	_, err := svc.GetByID(9999)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got: %v", err)
	}
}




func TestActivityGetMissingReturnsNotFound(t *testing.T) {
	engine := testutil.NewRouter(t).Setup()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/activities/99999", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing activity, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCancelMissingEnrollmentErrors(t *testing.T) {
	cfg := testutil.TestConfig(t)
	engine := testutil.NewRouter(t).Setup()
	tok := testutil.Token(t, cfg, 1, "alice", "user")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/registrations/99999/cancel", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing registration, got %d body=%s", w.Code, w.Body.String())
	}
}

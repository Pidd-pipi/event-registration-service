// Code generated for scoring tests of record 008.
package z8_test

import (
	"bytes"
	"gbevent/internal/constants"
	"gbevent/internal/repository"
	"gbevent/internal/service"
	"gbevent/internal/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
)




func TestAdminRoleSignupRejected(t *testing.T) {
	engine := testutil.NewRouter(t).Setup()
	payload := []byte(`{"username":"badboy","password":"123456","role":"admin"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 422 registering admin role, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestUserCannotListUsers(t *testing.T) {
	cfg := testutil.TestConfig(t)
	engine := testutil.NewRouter(t).Setup()
	tok := testutil.Token(t, cfg, 3, "alice", "user")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for normal user listing users, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestUserCannotAccessActivityMine(t *testing.T) {
	cfg := testutil.TestConfig(t)
	engine := testutil.NewRouter(t).Setup()
	tok := testutil.Token(t, cfg, 3, "alice", "user")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/activities/mine", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for normal user accessing activity mine, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestNormalUserSignupSucceeds(t *testing.T) {
	engine := testutil.NewRouter(t).Setup()
	payload := []byte(`{"username":"orgman","password":"123456","role":"user"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("seed user register should succeed, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCommentMineRequiresAuth(t *testing.T) {
	engine := testutil.NewRouter(t).Setup()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/comments/mine", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token for comment mine, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestFavoriteMineRequiresAuth(t *testing.T) {
	engine := testutil.NewRouter(t).Setup()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/favorites/mine", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token for favorite mine, got %d body=%s", w.Code, w.Body.String())
	}
}




func TestServiceRejectsAdminRoleSignup(t *testing.T) {
	db := testutil.NewDB(t)
	svc := service.NewUserService(repository.NewUserRepository(db), testutil.NewLogger())
	if _, err := svc.Register("adminwannabe", "123456", "", "", "", constants.RoleAdmin); err == nil {
		t.Fatalf("register with admin role should be rejected")
	}
}

// Code generated for scoring tests of record 010.
package za_test

import (
	"gbevent/internal/model"
	"gbevent/internal/testutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)




func auditCount(t *testing.T, where string, args ...any) int64 {
	t.Helper()
	db := testutil.DBOf(t)
	var n int64
	if err := db.Model(&model.AuditLog{}).Where(where, args...).Count(&n).Error; err != nil {
		t.Fatal(err)
	}
	return n
}

func TestGetRequestNotAudited(t *testing.T) {
	engine := testutil.NewRouter(t).Setup()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/activities", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if n := auditCount(t, "LOWER(action) = ?", "get"); n != 0 {
		t.Fatalf("GET request should not be audited, got %d rows", n)
	}
}

func TestFailedRequestNotAudited(t *testing.T) {
	engine := testutil.NewRouter(t).Setup()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code == http.StatusOK {
		t.Fatalf("expected login to fail without body")
	}
	if n := auditCount(t, "LOWER(action) = ?", "post"); n != 0 {
		t.Fatalf("failed request should not be audited, got %d rows", n)
	}
}

func TestAuditActionUppercase(t *testing.T) {
	engine := testutil.NewRouter(t).Setup()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(`{"username":"auditman","password":"123456","role":"user"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("register should succeed, got %d body=%s", w.Code, w.Body.String())
	}
	if n := auditCount(t, "action = ?", "post"); n != 0 {
		t.Fatalf("audit action should be uppercase, found lowercase rows=%d", n)
	}
	if n := auditCount(t, "action = ?", "POST"); n != 1 {
		t.Fatalf("expected one uppercase audit row, got %d", n)
	}
}

func TestAuditEntityTypeFirstSegment(t *testing.T) {
	engine := testutil.NewRouter(t).Setup()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(`{"username":"auditman","password":"123456","role":"user"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("register should succeed, got %d body=%s", w.Code, w.Body.String())
	}
	if n := auditCount(t, "entity_type = ?", "auth/register"); n != 0 {
		t.Fatalf("audit entity_type should be first path segment, found %d rows", n)
	}
	if n := auditCount(t, "entity_type = ?", "auth"); n != 1 {
		t.Fatalf("expected one audit row with entity_type=auth, got %d", n)
	}
}

func TestPaginationTotalCorrect(t *testing.T) {
	engine := testutil.NewRouter(t).Setup()
	db := testutil.DBOf(t)
	now := time.Now()
	for i := 0; i < 2; i++ {
		a := model.Activity{Title: "活动", Description: "", CoverImage: "", ActivityType: "lecture",
			StartTime: now.AddDate(0, 0, 1), EndTime: now.AddDate(0, 0, 1), Location: "某地",
			Capacity: 10, SignupDeadline: now.AddDate(0, 0, 2), Status: "published", OrganizerID: 2}
		if err := db.Create(&a).Error; err != nil {
			t.Fatal(err)
		}
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/activities", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if !containsTotal(w.Body.String(), 2) {
		t.Fatalf("expected total=2, got body=%s", w.Body.String())
	}
}

func containsTotal(body string, want int64) bool {
	for _, cand := range []string{`"total":2`, `"total": 2`} {
		if containsStr(body, cand) {
			return true
		}
	}
	return false
}

func containsStr(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

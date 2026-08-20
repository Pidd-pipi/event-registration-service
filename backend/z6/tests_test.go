// Code generated for scoring tests of record 006.
package z6_test

import (
	"gbevent/internal/model"
	"gbevent/internal/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
)




func TestUnfavoriteMissingActivityErrors(t *testing.T) {
	cfg := testutil.TestConfig(t)
	engine := testutil.NewRouter(t).Setup()
	tok := testutil.Token(t, cfg, 3, "alice", "user")
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/activities/99999/favorite", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 removing non-favorited activity, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestMarkReadNonExistentErrors(t *testing.T) {
	cfg := testutil.TestConfig(t)
	engine := testutil.NewRouter(t).Setup()
	tok := testutil.Token(t, cfg, 3, "alice", "user")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notifications/99999/read", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 marking missing notification, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestMarkReadOthersNotificationErrors(t *testing.T) {
	cfg := testutil.TestConfig(t)
	engine := testutil.NewRouter(t).Setup()
	db := testutil.DBOf(t)
	if err := db.Create(&model.Notification{UserID: 2, NotificationType: "signup_success", Title: "t", Content: "c"}).Error; err != nil {
		t.Fatal(err)
	}
	tok := testutil.Token(t, cfg, 3, "alice", "user")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notifications/1/read", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 marking another user notification, got %d body=%s", w.Code, w.Body.String())
	}
}

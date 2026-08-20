// Code generated for scoring tests of record 009.
package z9_test

import (
	"gbevent/internal/testutil"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)




func makeToken(t *testing.T, secret string, userID uint64, role, issuer string, exp time.Time) string {
	t.Helper()
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": "x",
		"role":     role,
		"iss":      issuer,
		"exp":      exp.Unix(),
	}
	s, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return s
}

func TestExpiredTokenDenied(t *testing.T) {
	cfg := testutil.TestConfig(t)
	engine := testutil.NewRouter(t).Setup()
	tok := makeToken(t, cfg.JWTSecret, 3, "user", "gbevent", time.Now().Add(-time.Hour))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/favorites/mine", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for expired token, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestUnknownRoleTokenDenied(t *testing.T) {
	cfg := testutil.TestConfig(t)
	engine := testutil.NewRouter(t).Setup()
	tok := makeToken(t, cfg.JWTSecret, 3, "hacker", "gbevent", time.Now().Add(time.Hour))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for unknown role token, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestWrongIssuerTokenDenied(t *testing.T) {
	cfg := testutil.TestConfig(t)
	engine := testutil.NewRouter(t).Setup()
	tok := makeToken(t, cfg.JWTSecret, 3, "user", "evil", time.Now().Add(time.Hour))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/favorites/mine", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for wrong issuer token, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestRbacBlocksNormalUserFromAdminList(t *testing.T) {
	cfg := testutil.TestConfig(t)
	engine := testutil.NewRouter(t).Setup()
	tok := testutil.Token(t, cfg, 3, "alice", "user")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for normal user on admin route, got %d body=%s", w.Code, w.Body.String())
	}
}

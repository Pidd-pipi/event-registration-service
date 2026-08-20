// Code generated for scoring tests of record 005.
package z5_test

import (
	"encoding/json"
	"gbevent/internal/middleware"
	"gbevent/internal/testutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)




func doReq(t *testing.T, engine http.Handler, method, path, token string) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("bad json %s: %v", w.Body.String(), err)
	}
	return w, body
}

func TestAuthRejectsGarbageToken(t *testing.T) {
	engine := testutil.NewRouter(t).Setup()
	w, _ := doReq(t, engine, http.MethodGet, "/api/v1/favorites/mine", "garbage-token")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for garbage token, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCommentListReturnsOK(t *testing.T) {
	engine := testutil.NewRouter(t).Setup()
	w, _ := doReq(t, engine, http.MethodGet, "/api/v1/activities/1/comments", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for comment list, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestEmptyInboxListReturnsArray(t *testing.T) {
	cfg := testutil.TestConfig(t)
	engine := testutil.NewRouter(t).Setup()
	tok := testutil.Token(t, cfg, 3, "alice", "user")
	w, body := doReq(t, engine, http.MethodGet, "/api/v1/notifications/mine", tok)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("missing data: %s", w.Body.String())
	}
	if list, ok := data["list"].([]any); !ok || list == nil {
		t.Fatalf("expected empty array list, got: %v", data["list"])
	}
}




func TestErrorHandlerNoErrorsPasses(t *testing.T) {
	panicked := make(chan any, 1)
	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				panicked <- r
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	})
	engine.Use(middleware.ErrorHandler(testutil.NewLogger()))
	engine.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	select {
	case p := <-panicked:
		t.Fatalf("unexpected panic in error handler: %v", p)
	default:
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

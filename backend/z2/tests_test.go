package z2_test

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"gbevent/internal/middleware"

	"github.com/gin-gonic/gin"
)

func callHandler(t *testing.T, h gin.HandlerFunc) {
	t.Helper()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	h(c)
}

func TestRequestIDConcurrentUnique(t *testing.T) {
	var wg sync.WaitGroup
	start := make(chan struct{})
	ids := make(chan string, 8*50)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 50; j++ {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
				middleware.RequestID()(c)
				ids <- w.Header().Get("X-Request-Id")
			}
		}()
	}
	close(start)
	wg.Wait()
	close(ids)
	set := make(map[string]struct{}, 8*50)
	for id := range ids {
		if _, ok := set[id]; ok {
			t.Fatalf("duplicate request id: %s", id)
		}
		set[id] = struct{}{}
	}
}

func TestRateLimiterConcurrentAccess(t *testing.T) {
	rl := middleware.NewRateLimiter(100)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 20; j++ {
				callHandler(t, rl.Limit())
			}
		}()
	}
	close(start)
	wg.Wait()
}

func TestRateLimiterTokensExhausted(t *testing.T) {
	rl := middleware.NewRateLimiter(3)
	blocked := 0
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		rl.Limit()(c)
		if w.Code == http.StatusTooManyRequests {
			blocked++
		}
	}
	if blocked < 1 {
		t.Fatalf("expected rate limiter to exhaust tokens, blocked=%d", blocked)
	}
}

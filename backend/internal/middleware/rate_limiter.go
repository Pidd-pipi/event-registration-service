package middleware

import (
	"net/http"
	"sync"
	"time"

	"gbevent/internal/constants"

	"github.com/gin-gonic/gin"
)

type bucket struct {
	tokens   int
	lastFill time.Time
}

// RateLimiter 基于 IP 的简易令牌桶限流。
type RateLimiter struct {
	mu       sync.Mutex
	perMin   int
	buckets  map[string]*bucket
	capacity int
	refill   time.Duration // 每补充一个令牌所需的时间
}

// NewRateLimiter 构造限流器。
func NewRateLimiter(perMin int) *RateLimiter {
	if perMin <= 0 {
		perMin = 120
	}
	return &RateLimiter{
		perMin:   perMin,
		buckets:  make(map[string]*bucket),
		capacity: perMin,
		refill:   time.Minute / time.Duration(perMin),
	}
}

// refill 按经过时间补充令牌，不超过桶容量。调用方需持有 mu。
func (b *bucket) refill(capacity int, refill time.Duration, now time.Time) {
	if refill <= 0 {
		return
	}
	elapsed := now.Sub(b.lastFill)
	if elapsed <= 0 {
		return
	}
	add := int(elapsed / refill)
	if add <= 0 {
		return
	}
	b.tokens += add
	if b.tokens > capacity {
		b.tokens = capacity
	}
	b.lastFill = now
}

// Limit 返回限流中间件。
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		rl.mu.Lock()
		b, ok := rl.buckets[ip]
		if !ok {
			b = &bucket{tokens: rl.capacity, lastFill: now}
			rl.buckets[ip] = b
		} else {
			b.refill(rl.capacity, rl.refill, now)
		}
		if b.tokens <= 0 {
			rl.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"code": constants.CodeTooManyRequests, "message": constants.MsgTooManyRequests, "data": nil})
			return
		}
		b.tokens--
		rl.mu.Unlock()

		c.Next()
	}
}

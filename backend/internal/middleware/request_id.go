package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const requestIDHeader = "X-Request-Id"

// RequestID 为每个请求生成 request_id，响应头与日志统一携带。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(requestIDHeader)
		if id == "" {
			id = newRequestID()
		}
		c.Header(requestIDHeader, id)
		c.Set("request_id", id)
		c.Next()
	}
}

// newRequestID 用进程内自增序号生成请求 ID。
var seq uint64

func newRequestID() string {
	seq++
	return fmt.Sprintf("req-%d", seq)
}

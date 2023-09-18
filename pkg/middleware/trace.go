package middleware

import (
	"github.com/gin-gonic/gin"
	"whisper/pkg/context"
	trace2 "whisper/pkg/trace"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		trace := trace2.GetTrace(c.Request)
		// 设置 example 变量
		c.Set(context.TraceID, trace.TraceID)
		// 请求前
		c.Next()
	}
}

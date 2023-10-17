package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"io"
	"strconv"
	"time"
	"whisper/pkg/context"
	"whisper/pkg/log"
)

func Proc() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()
		c.Set(context.StartTime, now)
		c.Next()
	}
}

func Params() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		query, _ := json.Marshal(c.Request.URL.Query())
		data := "|method=" + c.Request.Method + "|path=" + path + "|query=" + string(query)

		buf, _ := io.ReadAll(c.Request.Body)
		rdr1 := io.NopCloser(bytes.NewBuffer(buf))
		rdr2 := io.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.
		c.Request.Body = rdr2

		body := readBody(rdr1)
		data += "|req=" + string(body)

		c.Next()

		response := ""
		if _, ok := c.Keys["response"]; ok {
			if data, isStr := c.Keys["response"].(string); isStr {
				response = data
			} else {
				response = cast.ToString(data)
			}
		}

		data += "|http_code=" + strconv.Itoa(c.Writer.Status()) + "|resp=" + response

		ctx := context.Context{
			Context: c,
		}
		log.Logger.Info(&ctx, data)
	}
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(reader)

	s := buf.String()
	return s
}

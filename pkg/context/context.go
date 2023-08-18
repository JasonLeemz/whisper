package context

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"whisper/pkg/errors"
	"whisper/pkg/trace"
)

// Context ...
type Context struct {
	*gin.Context
}

type reply struct {
	TraceID string      `json:"trace_id"`
	ErrMsg  string      `json:"err_msg"`
	ErrNo   int32       `json:"err_no"`
	Data    interface{} `json:"data"`
}

// Reply ...
func (c *Context) Reply(obj interface{}, err *errors.Error) {
	r := &reply{
		TraceID: c.Value(trace.TraceID).(string),
		Data:    obj,
	}
	if err != nil {
		r.ErrNo = err.ErrNo()
		r.ErrMsg = err.Error()
		c.Set("err", err)
	}

	resp, _ := json.Marshal(obj)
	c.Set("response", string(resp))
	c.PureJSON(http.StatusOK, r)
}

// Render ...
func (c *Context) Render(tpl string, data map[string]any) {
	c.HTML(http.StatusOK, tpl, data)
}

type HandlerFunc func(c *Context)

func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		ctx := &Context{
			gCtx,
		}
		h(ctx)
	}
}

// Bind is a shortcut for c.ShouldBindWith(obj, binding.JSON).
func (c *Context) Bind(obj any) error {
	if err := c.Context.ShouldBindJSON(obj); err != nil {
		c.Reply(nil, errors.New(err, errors.ErrNoInvalidInput))
		return err
	}
	return nil
}

func NewContext() *Context {
	ctx := &Context{
		Context: &gin.Context{},
	}

	req := new(http.Request)
	req.Header = make(http.Header)
	tr := trace.GetTrace(req)
	ctx.Set(trace.TraceID, tr.TraceID)

	return ctx
}

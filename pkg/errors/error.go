package errors

import (
	"encoding/json"
	"errors"
	"runtime"
	"strings"
)

type Error struct {
	err   error
	no    int32
	msg   string
	stack []uintptr
}

func callers() []uintptr {
	var pcs [32]uintptr
	l := runtime.Callers(3, pcs[:])
	return pcs[:l]
}

func (e *Error) ErrNo() int32 {
	return e.no
}

func (e *Error) String() string {
	s, _ := json.Marshal(e)
	return string(s)
}

func (e *Error) SetErrNo(code int32) {
	e.no = code
}

func (e *Error) SetError(msg string) {
	e.msg = msg
}

func (e *Error) Error() string {
	errMsg := ""

	if e.err != nil {
		errMsg = e.err.Error()
	}
	if e.msg != "" {
		errMsg = e.msg
	}

	return errMsg
}

type Impl interface {
	Error() string
	ErrNo() int32
	String() string

	SetErrNo(code int32)
	SetError(msg string)
}

func New(errMsg interface{}, ext ...string) *Error {
	var er error
	var emsg = ""
	switch errMsg.(type) {
	case *Error:
		if errMsg.(*Error) == nil {
			return nil
		}
		return &Error{
			err:   errMsg.(*Error).err,
			no:    errMsg.(*Error).ErrNo(),
			msg:   errMsg.(*Error).Error() + "|" + strings.Join(ext, " | "),
			stack: callers(),
		}
	case string:
		er = errors.New(errMsg.(string))
		emsg = errMsg.(string)
	case error:
		if errMsg.(error) == nil {
			return nil
		}
		er = errMsg.(error)
		emsg = errMsg.(error).Error()
	default:
		return nil
	}
	msg := strings.Join(ext, " | ")
	err := &Error{
		err:   er,
		stack: callers(),
	}

	// 错误号
	if code, ok := errorMessages[msg]; ok {
		err.SetErrNo(code)
	} else {
		err.SetErrNo(errorMessages[ErrNoUnclassified])
	}

	// 错误描述
	if msg != "" {
		msg = emsg + "," + msg
	}
	err.SetError(msg)

	return err
}

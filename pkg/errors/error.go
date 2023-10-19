package errors

import (
	"runtime"
)

type Error struct {
	err     error
	no      int
	msg     string
	callers []uintptr // 调用路径,是否需要取决于业务需求
}

func (e *Error) No() int {
	return e.no
}

func (e *Error) Msg() string {
	return e.msg
}

type NewErrOption func(err *Error)

// WithNo 记录错误号
func WithNo(errNo int) NewErrOption {
	return func(err *Error) {
		err.no = errNo
	}
}

// WithMsg 记录错误信息
func WithMsg(errMsg string) NewErrOption {
	return func(err *Error) {
		err.msg += " | " + errMsg
	}
}

// 实例化过程

// New 创建一个error实例
func New(err error, opts ...NewErrOption) *Error {
	e := &Error{
		callers: callers(), // 默认都把堆栈信息记录下来
	}

	if err != nil {
		e.err = err
		e.no = errorMessages[ErrNoUnclassified]
		e.msg = err.Error()
	}

	// 如果传入了可选项，进行赋值操作
	for _, opt := range opts {
		opt(e)
	}

	return e
}

func callers() []uintptr {
	var pcs [32]uintptr
	l := runtime.Callers(3, pcs[:])
	return pcs[:l]
}

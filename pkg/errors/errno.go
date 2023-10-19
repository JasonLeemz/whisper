package errors

type ErrorDesc = string

const (
	ErrNoUnclassified = "unclassified error"

	// 参数错误
	ErrNoInvalidInput  = "invalid input"
	ErrNoMissingField  = "missing field"
	ErrNoOutOfRange    = "out of range"
	ErrNoUnknownParams = "unknown parameter error"

	// 认证和权限错误
	ErrNoUnauthorized = "unauthorized"
)

var errorMessages = map[ErrorDesc]int{
	// 未分类错误
	ErrNoUnclassified: 1000,

	// 参数错误
	ErrNoInvalidInput:  1001,
	ErrNoMissingField:  1002,
	ErrNoOutOfRange:    1003,
	ErrNoUnknownParams: 1004,

	// 认证和权限错误
	ErrNoUnauthorized: 2001,
}

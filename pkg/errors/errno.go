package errors

type ErrorDesc = string

const (
	ErrNoUnclassified = "Unclassified error"

	// 参数错误
	ErrNoInvalidInput  = "Invalid input"
	ErrNoMissingField  = "Missing field"
	ErrNoOutOfRange    = "Out of range"
	ErrNoUnknownParams = "Unknown parameter error"

	// 认证和权限错误
	ErrNoUnauthorized = "Unauthorized"
)

var errorMessages = map[ErrorDesc]int32{
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

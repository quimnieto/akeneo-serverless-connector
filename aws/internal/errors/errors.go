package errors

import "fmt"

type AppError struct {
	Code    string
	Message string
	Details string
}

func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

var (
	ErrInvalidPayload = &AppError{
		Code:    "INVALID_PAYLOAD",
		Message: "Invalid webhook payload",
	}

	ErrMissingField = &AppError{
		Code:    "MISSING_FIELD",
		Message: "Required field is missing",
	}

	ErrSNSPublishFailed = &AppError{
		Code:    "SNS_PUBLISH_FAILED",
		Message: "Failed to publish to SNS",
	}
)

func NewAppError(code, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func (e *AppError) WithDetails(details string) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Details: details,
	}
}


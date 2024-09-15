package errors

type BaseError struct {
	Message    string
	StatusCode int
}

func (e *BaseError) Error() string {
	return e.Message
}

func NewBaseError(message string, statusCode int) *BaseError {
	return &BaseError{
		Message:    message,
		StatusCode: statusCode,
	}
}

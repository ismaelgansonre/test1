package utils

// AppError defines a custom error type for application errors
type AppError struct {
	StatusCode int
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}

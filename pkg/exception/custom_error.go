package exception

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) *AppError {
	return &AppError{Code: 404, Message: message}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{Code: 400, Message: message}
}

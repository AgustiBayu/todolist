package exception

type ErrorHandler struct {
	Code    int
	Message string
}

func (h *ErrorHandler) Error() string {
	return h.Message
}

func ValidationError(message string) *ErrorHandler {
	return &ErrorHandler{
		Code:    400,
		Message: message,
	}
}

func InternalServerError(message string) *ErrorHandler {
	return &ErrorHandler{
		Code:    500,
		Message: message,
	}
}

func NotFoundError(message string) *ErrorHandler {
	return &ErrorHandler{
		Code:    404,
		Message: message,
	}
}

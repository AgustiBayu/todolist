package exception

type ErrorHandler struct {
	Code    int
	Message string
}

func (h *ErrorHandler) Error() string {
	return h.Message
}

func UnauthorizedErr(message string) *ErrorHandler {
	return &ErrorHandler{
		Code:    401,
		Message: message,
	}
}

func ConflictError(message string) *ErrorHandler {
	return &ErrorHandler{
		Code:    409,
		Message: message,
	}
}

func BadRequestError(message string) *ErrorHandler {
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

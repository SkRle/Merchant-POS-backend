package response

// SuccessResponse defines the standard structure for a successful API response
type SuccessResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse defines the standard structure for an error API response
type ErrorResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

// Success returns a standardized success format
func Success(code int, message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Code:    code,
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

// Error returns a standardized error format
func Error(code int, message string, errors interface{}) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Status:  "error",
		Message: message,
		Errors:  errors,
	}
}

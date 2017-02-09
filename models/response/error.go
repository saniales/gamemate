package response

//ErrorResponse represents a response from the system regarding an error during an operation.
type ErrorResponse struct {
	Type      string `json:"Type" xml:"Type" form:"Type"`
	ErrorCode int    `json:"ErrorCode" xml:"ErrorCode" form:"ErrorCode"`
}

//ErrorDetail is supposed to be used only during development to debug errors.
type ErrorDetail struct {
	Error        ErrorResponse `json:"Error" xml:"Error" form:"Error"`
	ErrorMessage string        `json:"ErrorMessage" xml:"ErrorMessage" form:"ErrorMessage"`
}

//FromErrorCode fills the struct with information from the error.
func (receiver *ErrorResponse) FromErrorCode(errorCode int) {
	receiver.Type = "Error"
	receiver.ErrorCode = errorCode
}

//FromError generates a detailed error report (meant to be used during development to debug).
func (receiver *ErrorDetail) FromError(err error, errorCode int) {
	receiver.Error = ErrorResponse{}
	receiver.Error.FromErrorCode(errorCode)
	receiver.ErrorMessage = err.Error()
}

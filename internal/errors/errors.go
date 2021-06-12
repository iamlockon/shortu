package errors

func New(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

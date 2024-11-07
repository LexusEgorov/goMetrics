package dohSimpson

import "fmt"

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("D'oh! #%d, Message: %s", e.Code, e.Message)
}

func NewDoh(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

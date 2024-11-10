package dohsimpson

import "fmt"

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("D'oh! #%d, Message: %s", e.Code, e.Message)
}

func NewDoh(code int, message string) *Error {
	doh := &Error{
		Code:    code,
		Message: message,
	}

	fmt.Println(doh.Error())

	return doh
}

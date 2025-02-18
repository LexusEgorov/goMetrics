// Пакет с кастомными ошибками.
package dohsimpson

import "fmt"

//Содержит http-код ошибки и сообщение.
type Error struct {
	Code    int
	Message string
}

//Конструктор ошибки с выводом в консоль.
func NewDoh(code int, message string) *Error {
	doh := &Error{
		Code:    code,
		Message: message,
	}

	fmt.Println(doh.Error())

	return doh
}

//Имплементирует интерфейс error. Ругается как Гомер Симсон:))
func (e Error) Error() string {
	return fmt.Sprintf("D'oh! #%d, Message: %s", e.Code, e.Message)
}

package server

import "fmt"

type ServerError interface {
	error
	Code() int
	Message() string
}

type ErrorHTTP struct {
	code    int
	message string
}

func (e ErrorHTTP) Code() int {
	return e.code
}

func (e ErrorHTTP) Message() string {
	return e.message
}

func (e ErrorHTTP) Error() string {
	return fmt.Sprintf("%d: %s", e.code, e.message)
}

package errors

import (
	"fmt"
	"time"
)

var (
	ErrorTokenExpired = Error{code: 1000, what: "token expired"}
)

type Error struct {
	when time.Time // 错误发生时间
	code int32     // 错误码
	what string    // 错误原因
}

func (e *Error) SetWhen(when time.Time) {
	e.when = when
}

func (e *Error) When() time.Time {
	return e.when
}

func (e *Error) WhenString() string {
	return e.when.Format("2016-01-02 15:04:05")
}

func (e *Error) What() string {
	return e.what
}

func (e *Error) Code() int32 {
	return e.code
}

func New(code int32, what string) *Error {
	return &Error{
		when: time.Now(),
		code: code,
		what: what,
	}
}

func NewError(err Error) *Error {
	err.SetWhen(time.Now())
	return &err
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s [%d]", e.what, e.code)
}

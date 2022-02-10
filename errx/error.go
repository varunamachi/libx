package errx

import (
	"fmt"
	"runtime"
)

type Error struct {
	Err  error  `json:"err"`
	Msg  string `json:"msg"`
	File string `json:"file"`
	Line int    `json:"line"`
}

func (fxErr *Error) Error() string {
	return fxErr.Err.Error()
}

func (fxErr *Error) Unwrap() error {
	return fxErr.Err
}

func (cfx *Error) String() string {
	return cfx.Err.Error() + ": " + cfx.Msg
}

func Errf(inner error, msg string, args ...interface{}) *Error {
	_, file, line, _ := runtime.Caller(1)

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return &Error{Err: inner, Msg: msg, File: file, Line: line}
}

func Str(err error) string {
	ex, ok := err.(*Error)
	if !ok {
		return err.Error()
	}
	return fmt.Sprintf("%s:%d - %s", ex.File, ex.Line, ex.Msg)
}

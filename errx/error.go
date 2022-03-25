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

func (ex *Error) Error() string {
	if ex.Err == nil {
		return ex.Msg
	}
	return ex.Err.Error()
}

func (ex *Error) Unwrap() error {
	return ex.Err
}

func (ex *Error) String() string {
	if ex.Err == nil {
		return ex.Msg
	}
	return ex.Err.Error() + ": " + ex.Msg
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

func PrintSomeStack(err error) {
	stackPrinter(err, "")
}

func stackPrinter(err error, indent string) {
	if err == nil {
		return
	}
	ex, ok := err.(*Error)
	if ok {
		fmt.Printf("%s> %s:%d - %s\n", indent, ex.File, ex.Line, ex.Msg)
		stackPrinter(ex.Err, indent+"--")
		return
	}
	fmt.Printf("%s> %v\n", indent, err)
}

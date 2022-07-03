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

func StackArray(err error) []string {
	errs := make([]string, 0, 10)
	errStackArray(err, 10, 0, errs)
	return errs
}

func errStackArray(in error, maxDepth, curDepth int, out []string) {
	if in == nil {
		return
	}

	out = append(out, Str(in))
	ex, ok := in.(*Error)
	if ok && curDepth < maxDepth {
		curDepth++
		errStackArray(ex.Err, maxDepth, curDepth, out)
		return
	}
}

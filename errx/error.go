package errx

import (
	"errors"
	"fmt"
	"runtime"
)

type Error struct {
	Err  error  `json:"err"`
	Code string `json:"code"`
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

func Fmt(msg string, args ...interface{}) *Error {
	_, file, line, _ := runtime.Caller(1)
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return &Error{
		Err:  nil,
		Code: "",
		Msg:  msg,
		File: file,
		Line: line,
	}
}

func Errf(inner error, msg string, args ...interface{}) *Error {
	_, file, line, _ := runtime.Caller(1)
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return &Error{
		Err:  inner,
		Code: inner.Error(),
		Msg:  msg,
		File: file,
		Line: line,
	}
}

func Errfx(inner error, code, msg string, args ...interface{}) *Error {
	_, file, line, _ := runtime.Caller(1)
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return &Error{
		Err:  inner,
		Code: code,
		Msg:  msg,
		File: file,
		Line: line,
	}
}

func New(code, msg string) *Error {
	_, file, line, _ := runtime.Caller(1)
	return &Error{
		Err:  errors.New(code),
		Code: code,
		Msg:  msg,
		File: file,
		Line: line,
	}
}

func Str(err error) string {
	if err == nil {
		return "N/A"
	}
	ex, ok := err.(*Error)
	if !ok {
		return err.Error()
	}
	// if ex.Code != ex.Err.Error() {

	// }
	return fmt.Sprintf("%s:%d - %s", ex.File, ex.Line, ex.Msg)
}

func PrintSomeStack(err error) {
	fmt.Println()
	stackPrinter(err, "")
	fmt.Println()
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
	const maxDepth = 5
	errs := make([]string, 0, maxDepth)
	for depth := 0; depth < maxDepth; depth++ {
		ex, ok := err.(*Error)
		if !ok {
			break
		}
		errs = append(errs, Str(ex))
		err = ex.Err
	}
	return errs
}

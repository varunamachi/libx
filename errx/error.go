package errx

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

var MaxStackPrintDepth = 10

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

func Fmt(msg string, args ...interface{}) error {
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

func Errf(inner error, msg string, args ...interface{}) error {
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

func Wrap(inner error) error {
	if inner == nil {
		return nil
	}

	msg, code := "", ""
	ex, ok := inner.(*Error)
	// using if-else on purpose
	if !ok || ex == nil {
		msg = inner.Error()
		code = reflect.TypeOf(inner).String()
	} else {
		msg = ex.Msg
		code = ex.Code
	}

	// errName := reflect.TypeOf(inner).String()
	_, file, line, _ := runtime.Caller(1)

	return &Error{
		Err:  inner,
		Code: code,
		Msg:  msg,
		File: file,
		Line: line,
	}
}

func Todo(emsg string) error {
	_, file, line, _ := runtime.Caller(1)
	return &Error{
		Err:  nil,
		Code: "TODO",
		Msg:  "TODO: " + emsg,
		File: file,
		Line: line,
	}
}

func Errfx(inner error, code, msg string, args ...interface{}) error {
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

func New(code, msg string) error {
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
	if err != nil {
		fmt.Println(err, err != nil)
		ex, ok := err.(*Error)
		if !ok {
			return err.Error()
		}
		return fmt.Sprintf("%s:%d \u2B9E %s", ex.File, ex.Line, ex.Msg)
	}
	return "N/A"
}

func PrintSomeStack(err error) {
	fmt.Println()
	stackPrinter(err, true, 0)
	fmt.Println()
}

func stackPrinter(err error, first bool, idx int) {
	if err == nil {
		return
	}
	ex, ok := err.(*Error)
	if ok {

		sym := ""
		indent := strings.Repeat("\u2500", idx)
		if !first {
			sym = "\u251c"
			// indent += "\u2B9E "
			indent += "\u2BC8 "
		}

		fmt.Printf("%s%s%s:%d \u2B9E %s\n",
			sym, indent, ex.File, ex.Line, ex.Msg)
		if idx < MaxStackPrintDepth {
			stackPrinter(ex.Err, false, idx+1)
		}
		return
	}
	indent := strings.Repeat("\u2500", idx)
	fmt.Printf("\u2514%s\u2BC8 %v\n", indent, err.Error())
}

func StackArray(err error) []string {
	const maxDepth = 5
	errs := make([]string, 0, maxDepth)
	for depth := 0; depth < maxDepth; depth++ {
		ex, ok := err.(*Error)
		if !ok || ex == nil {
			break
		}
		errs = append(errs, ex.Error())
		err = ex.Err
	}
	return errs
}

package errx

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func BadReqf(e error, msg string, args ...interface{}) *echo.HTTPError {
	if len(args) != 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return &echo.HTTPError{
		Code:     http.StatusBadRequest,
		Message:  msg,
		Internal: e,
	}
}

func BadReq(e error) *echo.HTTPError {
	return &echo.HTTPError{
		Code:     http.StatusBadRequest,
		Message:  Str(e),
		Internal: e,
	}
}

func InternalServerErr(e error) *echo.HTTPError {
	return &echo.HTTPError{
		Code:     http.StatusInternalServerError,
		Internal: e,
	}
}

func InternalServerErrf(
	e error, msg string, args ...interface{}) *echo.HTTPError {
	if len(args) != 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return &echo.HTTPError{
		Code:     http.StatusBadRequest,
		Message:  e.Error(),
		Internal: e,
	}
}

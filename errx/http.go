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
		Internal: Errf(e, "bad request"),
	}
}

func IntSrvErrf(
	e error, msg string, args ...interface{}) *echo.HTTPError {
	if len(args) != 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return &echo.HTTPError{
		Code:     http.StatusInternalServerError,
		Message:  msg,
		Internal: Errf(e, "internal server error"),
	}
}

func IntSrvErr(e error) *echo.HTTPError {
	return &echo.HTTPError{
		Code:     http.StatusInternalServerError,
		Message:  e.Error(),
		Internal: Errf(e, "internal server error"),
	}
}

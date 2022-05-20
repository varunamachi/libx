package httpx

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/rt"
)

const (
	EnvPrintAllAccess = "VLIBX_HTTP_PRINT_ALL_ACCESS"
)

//getToken - gets token from context or from header
func getToken(ctx echo.Context) (token *jwt.Token, err error) {
	itk := ctx.Get("token")
	if itk != nil {
		var ok bool
		if token, ok = itk.(*jwt.Token); !ok {
			err = fmt.Errorf("invalid token found in context")
		}
	} else {
		header := ctx.Request().Header.Get("Authorization")
		authSchemeLen := len("Bearer")
		if len(header) > authSchemeLen {
			tokStr := header[authSchemeLen+1:]
			keyFunc := func(t *jwt.Token) (interface{}, error) {
				return auth.GetJWTKey(), nil
			}
			token, err = jwt.Parse(tokStr, keyFunc)
		} else {
			err = fmt.Errorf("unexpected auth scheme used to JWT")
		}
	}
	return token, err
}

//RetrieveSessionInfo - retrieves session information from JWT token
func retrieveUserId(ctx echo.Context) (string, error) {
	token, err := getToken(ctx)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims in JWT")
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		return "", fmt.Errorf("couldnt find userId in token")
	}

	return userId, nil
}

func getAuthzMiddleware(ep *Endpoint, server *Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(etx echo.Context) error {
			userId, err := retrieveUserId(etx)
			if err != nil {
				return &echo.HTTPError{
					Code:     http.StatusForbidden,
					Message:  "invalid JWT information",
					Internal: err,
				}
			}

			// ep, ok := etx.Get("endpoint").(Endpoint)
			// if !ok {
			// 	return &echo.HTTPError{
			// 		Code:    http.StatusInternalServerError,
			// 		Message: "could not find endpoint information",
			// 	}
			// }

			user, err := server.userRetriever(etx.Request().Context(), userId)
			if err != nil {
				return err
			}

			if !auth.HasPerms(user, ep.Permission) || !auth.HasRole(user, ep.Role) {

				return &echo.HTTPError{
					Code:    http.StatusUnauthorized,
					Message: "permission to access resource is denied",
				}
			}

			etx.Set("endpoint", ep)
			etx.Set("user", user)
			return next(etx)
		}
	}
}

func getAccessMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(etx echo.Context) error {
			err := next(etx)
			if err == nil {
				return nil
			}
			printIfInternal := func(err error) bool {
				irr, ok := err.(*errx.Error)
				if !ok {
					return false
				}
				log.Error().
					Int("statusCode", http.StatusInternalServerError).
					Str("file", irr.File).
					Int("line", irr.Line).
					Str("user", GetUserId(etx)).
					Str("method", etx.Request().Method).
					Str("path", etx.Request().URL.Path).
					Msg(irr.Msg)
				errx.PrintSomeStack(irr)
				return true
			}

			if err == nil && rt.EnvBool(EnvPrintAllAccess, false) {
				status := etx.Response().Status
				log.Debug().
					Int("statusCode", status).
					Str("user", GetUserId(etx)).
					Str("method", etx.Request().Method).
					Str("path", etx.Request().URL.Path).
					Msg(http.StatusText(status))
				return nil
			}

			if printIfInternal(err) {
				return err
			}

			httpErr, ok := err.(*echo.HTTPError)
			if ok && printIfInternal(httpErr.Internal) {
				return err
			}
			if ok {
				log.Error().
					Int("statusCode", httpErr.Code).
					Str("user", GetUserId(etx)).
					Str("method", etx.Request().Method).
					Str("path", etx.Request().URL.Path).
					Msg(StrMsg(httpErr))
				return httpErr
			}

			log.Error().
				Int("statusCode", http.StatusInternalServerError).
				Str("user", GetUserId(etx)).
				Str("method", etx.Request().Method).
				Str("path", etx.Request().URL.Path).
				Msg(err.Error())
			return err
		}
	}
}

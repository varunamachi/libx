package httpx

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
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
			err = errx.New("jwt.tokenInvalid",
				"invalid token found in context")
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
			err = errx.New("jwt.invalidScheme",
				"unexpected auth scheme used to JWT")
		}
	}
	return token, err
}

//RetrieveSessionInfo - retrieves session information from JWT token
func retrieveUserId(ctx echo.Context) (string, string, error) {
	token, err := getToken(ctx)
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errx.New("jwt.invalidClaims", "invalid claims in JWT")
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		return "", "", errx.New("jwt.invalidUserId",
			"couldnt find userId in token")
	}

	userType, _ := claims["type"].(string)
	return userId, userType, nil
}

func getAuthzMiddleware(ep *Endpoint, server *Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(etx echo.Context) error {
			userId, userType, err := retrieveUserId(etx)
			if err != nil {
				return &echo.HTTPError{
					// Code:     http.StatusForbidden,
					Code:     http.StatusUnauthorized,
					Message:  "invalid JWT information",
					Internal: err,
				}
			}

			// This check is for non-DB users
			if userType != "" && userType != "user" {
				etx.Set("endpoint", ep)
				etx.Set("userId", userId)
				return next(etx)
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

			hasAccess := auth.HasPerms(user, ep.Permission) &&
				auth.HasRole(user, ep.Role)
			if !hasAccess {
				return &echo.HTTPError{
					// Code:    http.StatusUnauthorized,
					Code:    http.StatusForbidden,
					Message: "permission to access resource is denied",
				}
			}

			etx.Set("endpoint", ep)
			etx.Set("user", user)
			etx.Set("userId", userId)
			return next(etx)
		}
	}
}

func getAccessMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(etx echo.Context) error {
			err := next(etx)
			if err == nil {
				log.Info().
					Int("statusCode", etx.Response().Status).
					Str("user", GetUserId(etx)).
					Str("method", etx.Request().Method).
					Str("path", etx.Request().URL.Path).
					Msg("-- OK --")
				return nil
			}
			printIfInternal := func(status int, err error) bool {
				irr, ok := err.(*errx.Error)
				if !ok {
					return false
				}
				log.Error().
					Int("statusCode", status).
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

			if printIfInternal(http.StatusInternalServerError, err) {
				return err
			}

			httpErr, ok := err.(*echo.HTTPError)
			if ok && printIfInternal(httpErr.Code, httpErr.Internal) {
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

func accessMiddleware(printErrors bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(etx echo.Context) error {
			err := next(etx)
			if err == nil {
				log.Info().
					Int("statusCode", etx.Response().Status).
					Str("user", GetUserId(etx)).
					Str("method", etx.Request().Method).
					Str("path", etx.Request().URL.Path).
					Msg("-- OK --")
				return nil
			}
			if err != nil && printErrors {
				log.Error().
					Int("statusCode", http.StatusInternalServerError).
					Str("user", GetUserId(etx)).
					Str("method", etx.Request().Method).
					Str("path", etx.Request().URL.Path).
					Msg(err.Error())
			}
			return nil
		}
	}
}

func errorHandlerFunc(err error, etx echo.Context) {
	asJson := func(status int, code, msg string) {
		hm := data.M{
			"status":    http.StatusText(status),
			"errorCode": code,
			"msg":       msg,
		}
		if err := etx.JSON(status, hm); err != nil {
			log.Error().Err(err).Msg("failed to send error to client")
		}
	}

	printIfInternal := func(status int, err error) bool {
		irr, ok := err.(*errx.Error)
		if !ok {
			return false
		}
		log.Error().
			Int("statusCode", status).
			Str("file", irr.File).
			Int("line", irr.Line).
			Str("user", GetUserId(etx)).
			Str("method", etx.Request().Method).
			Str("path", etx.Request().URL.Path).
			Msg(irr.Msg)
		errx.PrintSomeStack(irr)
		asJson(status, irr.Code, irr.Msg)
		return true
	}

	if printIfInternal(http.StatusInternalServerError, err) {
		return
	}

	httpErr, ok := err.(*echo.HTTPError)
	if ok && printIfInternal(httpErr.Code, httpErr.Internal) {
		return
	}
	if ok {
		msg := StrMsg(httpErr)
		log.Error().
			Int("statusCode", httpErr.Code).
			Str("user", GetUserId(etx)).
			Str("method", etx.Request().Method).
			Str("path", etx.Request().URL.Path).
			Msg(msg)
		asJson(httpErr.Code, strconv.Itoa(httpErr.Code), msg)
		return
	}

	log.Error().
		Int("statusCode", http.StatusInternalServerError).
		Str("user", GetUserId(etx)).
		Str("method", etx.Request().Method).
		Str("path", etx.Request().URL.Path).
		Msg(err.Error())
	asJson(http.StatusInternalServerError, "500", err.Error())
}

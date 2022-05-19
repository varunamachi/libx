package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/varunamachi/libx/errx"
)

type AuthData map[string]interface{}

var (
	ErrAuthentication = errors.New("auth.user.authenticationError")
	ErrUserRetrieval  = errors.New("auth.user.retrievalError")
	ErrToken          = errors.New("auth.user.authTokenError")
)

type Authenticator interface {
	Authenticate(gtx context.Context, authData AuthData) error
}

type UserGetter interface {
	GetUser(gtx context.Context, authData AuthData) (User, error)
}

type UserAuthenticator interface {
	Authenticator
	UserGetter
}

// Login - authenticates the user, get's the user information and generates a
// JWT token. The user and the token are then returned. And in case of error
// the error is returned
func Login(
	gtx context.Context,
	authr UserAuthenticator,
	data AuthData) (User, string, error) {
	if err := authr.Authenticate(gtx, data); err != nil {
		return nil, "", errx.Errf(ErrAuthentication, err.Error())
	}

	user, err := authr.GetUser(gtx, data)
	if err != nil {
		return nil, "", errx.Errf(ErrUserRetrieval, err.Error())
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = user.Id()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	signed, err := token.SignedString(GetJWTKey())
	if err != nil {
		return nil, "", errx.Errf(ErrToken, err.Error())
	}

	return user, signed, nil

}

//GetJWTKey - gives a unique JWT key
func GetJWTKey() []byte {
	jwtKey := os.Getenv("VLIBX_JWT_KEY")
	if len(jwtKey) == 0 {
		jwtKey = uuid.NewString()
		// TODO - may be need to do something better
		os.Setenv("VLIBX_JWT_KEY", jwtKey)
	}
	return []byte(jwtKey)
}

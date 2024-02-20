package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/varunamachi/libx/errx"
)

const (
	UserSessionTimeout = time.Hour * 24 * 10 // 10 days
	// UserSessionTimeout = time.Minute // 1 minute - for testing
)

type (
	AuthData        map[string]interface{}
	userAndPassword struct {
		UserId   string `json:"userId"`
		Password string `json:"password"`
	}
)

func (ad AuthData) Decode(out any) error {
	if err := mapstructure.Decode(&ad, &out); err != nil {
		return errx.Errf(err, "failed to decode auth data")
	}
	return nil
}

func (ad AuthData) ToUserAndPassword() (
	userId string, password string, err error) {

	var up userAndPassword
	if err := ad.Decode(&up); err != nil {
		return "", "", err
	}
	return up.UserId, up.Password, nil
}

var (
	ErrAuthentication         = errors.New("auth.user.authenticationError")
	ErrUserRetrieval          = errors.New("auth.user.retrievalError")
	ErrToken                  = errors.New("auth.user.authTokenError")
	ErrInsufficientPrivileges = errors.New("auth.user.insufficient.privs")
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

// TODO - remove once idx is operational
// Login - authenticates the user, get's the user information and generates a
// JWT token. The user and the token are then returned. And in case of error
// the error is returned
func Login(
	gtx context.Context,
	authr UserAuthenticator,
	data AuthData) (User, string, error) {
	if err := authr.Authenticate(gtx, data); err != nil {
		return nil, "", errx.Errf(err, "failed to authenticate user")
	}

	user, err := authr.GetUser(gtx, data)
	if err != nil {
		return nil, "", errx.Errf(err, "failed to retrieve user")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = user.Id()
	claims["exp"] = time.Now().Add(UserSessionTimeout).Unix()
	claims["type"] = "user"

	signed, err := token.SignedString(GetJWTKey())
	if err != nil {
		return nil, "", errx.Errf(err, "failed to generate session token")
	}

	return user, signed, nil

}

// GetJWTKey - gives a unique JWT key
func GetJWTKey() []byte {
	jwtKey := os.Getenv("VLIBX_JWT_KEY")
	if len(jwtKey) == 0 {
		jwtKey = uuid.NewString()
		// TODO - may be need to do something better
		os.Setenv("VLIBX_JWT_KEY", jwtKey)
	}
	return []byte(jwtKey)
}

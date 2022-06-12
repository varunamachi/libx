package httpx

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/iox"
)

var (
	ErrNotFound            = errors.New("client.http.notFound")
	ErrUnauthorized        = errors.New("client.http.unauthorized")
	ErrForbidden           = errors.New("client.http.forbidden")
	ErrInternalServerError = errors.New("client.http.internalServerError")
	ErrOtherStatus         = errors.New("client.http.otherStatus")

	ErrInvalidResponse = errors.New("client.http.invalidResponse")
	ErrClientError     = errors.New("client.http.clientError")
)

type ApiResult struct {
	resp   *http.Response
	err    error
	target string
	code   int
}

func newApiResult(req *http.Request, resp *http.Response) *ApiResult {

	target := "[" + req.Method + " " + req.URL.Path + "]"
	res := &ApiResult{
		resp:   resp,
		target: target,
		code:   resp.StatusCode,
	}

	var err *errx.Error

	switch resp.StatusCode {
	case http.StatusNotFound:
		err = errx.Errf(ErrNotFound, "not found: %s", target)
	case http.StatusUnauthorized:
		err = errx.Errf(ErrUnauthorized, "unauthorized: %s", target)
	case http.StatusForbidden:
		err = errx.Errf(ErrUnauthorized, "forbidden: %s", target)
	case http.StatusInternalServerError:
		err = errx.Errf(ErrUnauthorized, "internal Server Error: %s", target)
	default:
		if resp.StatusCode > 400 {
			err = errx.Errf(
				ErrOtherStatus, "http-error: %d - %s", resp.StatusCode, target)
		}
	}
	if err != nil {
		log.Error().Err(err)
		res.err = err
	}
	return res
}

func newErrorResult(req *http.Request, err error, msg string) *ApiResult {
	target := ""
	if req != nil {
		target = req.Method + " " + req.URL.Path
		msg = msg + " - [" + target + "]"
	}

	return &ApiResult{
		err:    errx.Errf(err, msg),
		target: target,
	}
}

func (ar *ApiResult) LoadClose(out interface{}) error {
	defer func() {
		if ar.resp != nil && ar.resp.Body != nil {
			ar.resp.Body.Close()
		}
	}()

	if ar.Error() != nil {
		return ar.Error()
	}

	ar.err = json.NewDecoder(ar.resp.Body).Decode(out)
	return ar.err
}

func (ar *ApiResult) Error() error {
	if ar.err != nil {
		return ar.err
	}

	if ar.resp == nil || ar.resp.Body == nil {
		ar.err = errx.Errf(ErrInvalidResponse,
			"No valid http response received")
	}
	return ar.err
}

type Client struct {
	*http.Client
	address     string
	contextRoot string
	token       string
	user        auth.User
}

func DefaultTransport() *http.Transport {
	return &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 20 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 20 * time.Second,
	}
}

func New(address, contextRoot string) *Client {
	return &Client{
		address:     address,
		contextRoot: contextRoot,
		Client: &http.Client{
			Timeout:   time.Second * 20,
			Transport: DefaultTransport(),
		},
	}
}

func NewCustom(
	address, contextRoot string,
	transport *http.Transport,
	timeout time.Duration) *Client {
	return &Client{
		address:     address,
		contextRoot: contextRoot,
		Client: &http.Client{
			Transport: transport,
			Timeout:   timeout,
		},
	}
}

func (client *Client) createUrl(args ...string) string {
	var buffer bytes.Buffer
	if _, err := buffer.WriteString(client.address); err != nil {
		log.Fatal().Err(err)
	}
	if _, err := buffer.WriteString(client.contextRoot); err != nil {
		log.Fatal().Err(err)
	}
	if !strings.HasSuffix(client.contextRoot, "/") {
		buffer.WriteString("/")
	}
	for i := 0; i < len(args); i++ {
		if _, err := buffer.WriteString(args[i]); err != nil {
			log.Fatal().Err(err)
		}
		if i < len(args)-1 {
			if _, err := buffer.WriteString("/"); err != nil {
				log.Fatal().Err(err)
			}
		}
	}
	return buffer.String()
}

func (client *Client) putOrPost(
	gtx context.Context,
	method string,
	content interface{},
	urlArgs ...string) *ApiResult {

	url := client.createUrl(urlArgs...)
	data, err := json.Marshal(content)
	if err != nil {
		return newErrorResult(nil, err, "failed to marshal data")
	}

	req, err := http.NewRequestWithContext(
		gtx, method, url, bytes.NewBuffer(data))
	if err != nil {
		return newErrorResult(req, err, "failed to create http request")
	}

	// We assume JSON
	req.Header.Set("Content-Type", "application/json")
	if client.token != "" {
		authHeader := fmt.Sprintf("Bearer %s", client.token)
		req.Header.Add("Authorization", authHeader)
	}

	resp, err := client.Do(req)
	if err != nil {
		return newErrorResult(req, err, "failed to perform http request")
	}
	return newApiResult(req, resp)
}

func (client *Client) Get(gtx context.Context, urlArgs ...string) *ApiResult {

	apiURL := client.createUrl(urlArgs...)
	req, err := http.NewRequestWithContext(gtx, "GET", apiURL, nil)
	if err != nil {
		newErrorResult(req, err, "Failed to create http request")
	}

	if client.token != "" {
		authHeader := fmt.Sprintf("Bearer %s", client.token)
		req.Header.Add("Authorization", authHeader)
	}

	resp, err := client.Do(req)

	if err != nil {
		return newErrorResult(req, err, "Failed to perform http request")
	}

	return newApiResult(req, resp)
}

func (client *Client) Post(
	gtx context.Context,
	content interface{},
	urlArgs ...string) *ApiResult {
	return client.putOrPost(gtx, echo.POST, content, urlArgs...)
}

//Put - performs a put request
func (client *Client) Put(
	gtx context.Context,
	content interface{},
	urlArgs ...string) *ApiResult {
	return client.putOrPost(gtx, echo.PUT, content, urlArgs...)
}

//Delete - performs a delete request
func (client *Client) Delete(
	gtx context.Context,
	urlArgs ...string) *ApiResult {
	apiURL := client.createUrl(urlArgs...)
	req, err := http.NewRequestWithContext(gtx, echo.DELETE, apiURL, nil)
	if err != nil {
		newErrorResult(req, err, "Failed to create http request")
	}

	if client.token != "" {
		authHeader := fmt.Sprintf("Bearer %s", client.token)
		req.Header.Add("Authorization", authHeader)
	}

	resp, err := client.Do(req)
	if err != nil {
		return newErrorResult(req, err, "Failed to perform http request")
	}

	return newApiResult(req, resp)
}

func (client *Client) User() auth.User {
	return client.user
}

type AuthData map[string]interface{}

type LoginConfig struct {
	LoginURL string
	UserOut  auth.User
}

func (client *Client) Login(
	gtx context.Context,
	lc *LoginConfig,
	authData AuthData) error {

	if authData == nil {
		return nil
	}

	loginResult := struct {
		Token   string    `json:"token"`
		UserOut auth.User `json:"user"`
	}{
		"",
		lc.UserOut,
	}

	rr := client.Post(gtx, authData, lc.LoginURL)
	if err := rr.LoadClose(&loginResult); err != nil {
		return err
	}
	client.token = loginResult.Token
	client.user = loginResult.UserOut
	return nil
}

func CreateClient(lc *LoginConfig, ctx *cli.Context) (
	*Client, error) {
	host := ctx.String("host")
	ignCertErrs := ctx.Bool("ignore-cert-errors")
	timeOut := ctx.Int("timeout-secs")

	tp := DefaultTransport()
	if ignCertErrs {
		tp.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: ignCertErrs,
		}
	}
	client := NewCustom(host, "", tp, time.Duration(timeOut)*time.Second)

	if lc == nil {
		return client, nil
	}

	userId := ctx.String("user-id")
	password := ctx.String("password")
	if password == "" {
		password = iox.AskPassword("Password")
	}

	err := client.Login(ctx.Context, lc, AuthData{
		"userId":   userId,
		"password": password,
	})
	if err != nil {
		return nil, err
	}

	return client, nil

}

func WithClientFlags(
	withAuth bool,
	envPrefix string,
	flags ...cli.Flag) []cli.Flag {
	flags = append(flags,
		&cli.StringFlag{
			Name: "host",
			Usage: "Full address of the host with URL scheme, host name/IP " +
				"and port",
			EnvVars: []string{
				envPrefix + "_CLIENT_REMOTE_HOST",
				"LIBX_CLIENT_REMOTE_HOST",
			},
			Required: true,
		},
		&cli.BoolFlag{
			Name: "ignore-cert-errors",
			Usage: "Ignore certificate errors while connecting to a HTTPS " +
				"service",
			Value: false,
			EnvVars: []string{
				envPrefix + "_CLIENT_IGNORE_CERT_ERR",
				"LIBX_CLIENT_IGNORE_CERT_ERR",
			},
		},
		&cli.IntFlag{
			Name:  "timeout-secs",
			Usage: "Time out in seconds",
			Value: 20,
			EnvVars: []string{
				envPrefix + "_CLIENT_TIMEOUT_SECS",
				"LIBX_CLIENT_TIMEOUT_SECS",
			},
		},
	)
	if withAuth {
		flags = append(flags,
			&cli.StringFlag{
				Name:     "user-id",
				Usage:    "User present in the remote service",
				Required: false,
				EnvVars: []string{
					envPrefix + "_CLIENT_USER_ID",
					"LIBX_CLIENT_USER_ID",
				},
			},
			&cli.StringFlag{
				Name: "password",
				Usage: "Password for the user, only use for development " +
					"purposes",
				Required: false,
				Hidden:   true,
				EnvVars: []string{
					envPrefix + "_CLIENT_PASSWORD",
					"LIBX_CLIENT_PASSWORD",
				},
			},
		)
	}

	return flags
}

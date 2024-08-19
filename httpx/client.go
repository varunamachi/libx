package httpx

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
)

var (
	ErrNotFound            = errors.New("client.http.notFound")
	ErrUnauthorized        = errors.New("client.http.unauthorized")
	ErrForbidden           = errors.New("client.http.forbidden")
	ErrInternalServerError = errors.New("client.http.internalServerError")
	ErrBadRequest          = errors.New("client.http.badRequest")
	ErrOtherStatus         = errors.New("client.http.otherStatus")
	ErrInvalidResponse     = errors.New("client.http.invalidResponse")
	ErrClientError         = errors.New("client.http.clientError")
)

type AuthData map[string]interface{}

type ApiResult struct {
	resp           *http.Response
	err            error
	reqBuildErrors []error
	target         string
	code           int
}

func newApiResult(req *http.Request, resp *http.Response) *ApiResult {

	target := "[" + req.Method + " " + req.URL.Path + "]"
	res := &ApiResult{
		resp:   resp,
		target: target,
		code:   resp.StatusCode,
	}

	var err error

	switch resp.StatusCode {
	case http.StatusBadRequest:
		err = ErrBadRequest
	case http.StatusNotFound:
		err = ErrNotFound
	case http.StatusUnauthorized:
		err = ErrUnauthorized
	case http.StatusForbidden:
		err = ErrForbidden
	case http.StatusInternalServerError:
		err = ErrInternalServerError
	default:
		if resp.StatusCode >= 400 {
			err = ErrOtherStatus
		}
	}

	if err == nil {
		return res
	}
	if resp.Body == nil {
		res.err = errx.Errf(err, "%s - %s - No body", resp.Status, target)
		return res
	}

	defer resp.Body.Close()
	msg := ""

	bbytes, err := io.ReadAll(resp.Body)
	if err != nil {
		msg = "failed to get error message: " + err.Error()
	} else {
		errMap := map[string]string{}
		reader := bytes.NewReader(bbytes)
		if err := json.NewDecoder(reader).Decode(&errMap); err != nil {
			msg = "unknown error: " + string(bbytes)
		} else {
			msg = errMap["msg"]
			msg = data.Qop(msg == "", string(bbytes), msg)
		}
	}
	// err =
	res.err = errx.Errf(ErrClientError, "%s - %s", resp.Status, msg)
	return res

	// First we check if this is in the form of echo.HttpError, if so we try to
	// get the internal error. If not we try to read the entire body as message
	// var he echo.HTTPError
	// if json.NewDecoder(resp.Body).Decode(&he) != nil {
	// 	data, err := io.ReadAll(resp.Body)
	// 	if err == nil {
	// 		msg = string(data)
	// 	}
	// } else {
	// 	msg = he.Error()
	// }
	// err = errx.Errf(err, "%s - %s - %s", resp.Status, target, msg)
	// res.err = err
	// return res
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

	var reader io.ReadCloser = ar.resp.Body
	var err error
	switch ar.resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(ar.resp.Body)
		if err != nil {
			return errx.Errf(err, "failed to create gzip reader for response")
		}
		defer reader.Close()
	}

	ar.err = json.NewDecoder(reader).Decode(out)
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

func (ar *ApiResult) Close() error {
	defer func() {
		if ar.resp != nil && ar.resp.Body != nil {
			ar.resp.Body.Close()
		}
	}()

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

func NewClient(address, contextRoot string) *Client {
	return &Client{
		address:     address,
		contextRoot: contextRoot,
		Client: &http.Client{
			Timeout:   time.Second * 20,
			Transport: DefaultTransport(),
		},
	}
}

func NewCustomClient(
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

func (client *Client) SetUser(user auth.User) *Client {
	client.user = user
	return client
}

func (client *Client) SetToken(token string) *Client {
	client.token = token
	return client
}

func (client *Client) User() auth.User {
	return client.user
}

func (client *Client) RemoteHost() (string, error) {
	url, err := url.Parse(client.address)
	if err != nil {
		return "", errx.Errf(err, "client: invalid remote host address")
	}
	return url.Host, nil
}

func (client *Client) createUrl(args ...string) string {
	var buffer bytes.Buffer
	// if _, err := buffer.WriteString(client.address); err != nil {
	// 	log.Fatal().Err(err)
	// }
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
	return client.address + path.Clean(buffer.String())
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
	r := newApiResult(req, resp)
	if r.err != nil {
		r.err = errx.Wrap(r.err)
	}
	return r
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

	r := newApiResult(req, resp)
	if r.err != nil {
		r.err = errx.Wrap(r.err)
	}
	return r
}

func (client *Client) Post(
	gtx context.Context,
	content interface{},
	urlArgs ...string) *ApiResult {
	return client.putOrPost(gtx, echo.POST, content, urlArgs...)
}

// Put - performs a put request
func (client *Client) Put(
	gtx context.Context,
	content interface{},
	urlArgs ...string) *ApiResult {
	return client.putOrPost(gtx, echo.PUT, content, urlArgs...)
}

// Delete - performs a delete request
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

	r := newApiResult(req, resp)
	if r.err != nil {
		r.err = errx.Wrap(r.err)
	}
	return r
}

// type RequestBuilder struct {

// }

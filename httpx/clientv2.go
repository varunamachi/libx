package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

var ErrHttpRequestBuildFailed = errors.New("http request build failed")

type RequestBuilder struct {
	method      string
	body        any
	client      *Client
	headers     http.Header
	queryParams map[string]string
	path        string
	withAuth    bool
	err         []error

	// TODO - now only json is supported, when others are to be supported, we
	// need a way to specify encoders
	// contentType string
}

func newRequestBuilder(client *Client, method string) *RequestBuilder {
	return &RequestBuilder{
		method:      method,
		headers:     map[string][]string{},
		queryParams: map[string]string{},
		path:        "",
		withAuth:    true,
		err:         make([]error, 0, 10),
	}
}

func newRequestBuilderWithBody(
	client *Client, method string, body any) *RequestBuilder {
	return &RequestBuilder{
		method:      method,
		body:        body,
		headers:     map[string][]string{},
		queryParams: map[string]string{},
		path:        "",
		withAuth:    true,
		err:         make([]error, 0, 10),
	}
}

func (rb *RequestBuilder) HdrStr(name, value string) *RequestBuilder {
	rb.headers.Add(name, value)
	return rb
}

func (rb *RequestBuilder) HdrInt(name string, value int64) *RequestBuilder {
	rb.headers.Add(name, strconv.FormatInt(value, 10))
	return rb
}

func (rb *RequestBuilder) HdrUint(name string, value uint64) *RequestBuilder {
	rb.headers.Add(name, strconv.FormatUint(value, 10))
	return rb
}

func (rb *RequestBuilder) HdrBool(name string, value bool) *RequestBuilder {
	rb.headers.Add(name, strconv.FormatBool(value))
	return rb
}

func (rb *RequestBuilder) HdrJson(name string, value any) *RequestBuilder {
	j, err := encodeJsonUrl(value)
	if err != nil {
		rb.err = append(rb.err, err)
		return rb
	}

	rb.headers.Add(name, j)
	return rb
}

func (rb *RequestBuilder) QStr(name, value string) *RequestBuilder {
	rb.queryParams[name] = value
	return rb
}

func (rb *RequestBuilder) QInt(name string, value int64) *RequestBuilder {
	rb.queryParams[name] = strconv.FormatInt(value, 10)
	return rb
}

func (rb *RequestBuilder) QUint(name string, value uint64) *RequestBuilder {
	rb.queryParams[name] = strconv.FormatUint(value, 10)
	return rb
}

func (rb *RequestBuilder) QBool(name string, value bool) *RequestBuilder {
	rb.queryParams[name] = strconv.FormatBool(value)
	return rb
}

func (rb *RequestBuilder) QJson(name string, value any) *RequestBuilder {
	j, err := encodeJsonUrl(value)
	if err != nil {
		rb.err = append(rb.err, err)
		return rb
	}

	rb.queryParams[name] = j
	return rb
}

func (rb *RequestBuilder) Path(params ...any) *RequestBuilder {
	var sb strings.Builder
	var err error
	for _, param := range params {
		sb.WriteString("/")
		switch p := param.(type) {
		case int8, int16, int32, int64, int:
			_, err = sb.WriteString(strconv.FormatInt(p.(int64), 10))
		case uint8, uint16, uint32, uint64, uint:
			_, err = sb.WriteString(strconv.FormatUint(p.(uint64), 10))
		case string:
			_, err = sb.WriteString(p)
		case float32, float64:
			_, err = sb.WriteString(
				strconv.FormatFloat(p.(float64), 'f', -1, 64))
		case bool:
			_, err = sb.WriteString(strconv.FormatBool(p))
		default:
			rb.err = append(rb.err,
				fmt.Errorf("invalid param type for path: '%T'", p))
			return rb
		}
		if err != nil {
			rb.err = append(rb.err, err)
			return rb
		}
	}
	rb.path = sb.String()
	return rb
}

func (rb *RequestBuilder) WithAuth(useAuth bool) *RequestBuilder {
	rb.withAuth = true
	return rb
}

func (rb *RequestBuilder) Exec(
	gtx context.Context) (*ApiResult, *RequestBuilder) {
	if len(rb.err) != 0 {
		return &ApiResult{
			err:    ErrHttpRequestBuildFailed,
			target: rb.method + " " + rb.path,
		}, rb
	}

	var err error
	var bodyBytes []byte = nil
	if rb.body != nil {
		bodyBytes, err = json.Marshal(rb.body)
		if err != nil {
			return newErrorResult(nil, err, "failed to marshal data"), rb
		}
	}
	// if !strings.HasSuffix(rb.client.contextRoot, "/") {
	// 	rb.path = "/" + rb.path
	// }
	fullUrl := rb.client.address + path.Clean(rb.client.contextRoot+"/"+rb.path)

	req, err := http.NewRequestWithContext(
		gtx, rb.method, fullUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return newErrorResult(req, err, "failed to create http request"), rb
	}

	req.Header = rb.headers
	if rb.withAuth {
		authHeader := fmt.Sprintf("Bearer %s", rb.client.token)
		req.Header.Add("Authorization", authHeader)
	}

	resp, err := rb.client.Do(req)
	if err != nil {
		return newErrorResult(req, err, "Failed to perform http request"), rb
	}

	return newApiResult(req, resp), rb
}

func encodeJsonUrl(obj any) (string, error) {
	j, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	return url.PathEscape(string(j)), nil
}

func (client *Client) BuildPost(body any) *RequestBuilder {
	return newRequestBuilderWithBody(client, echo.POST, body)
}

func (client *Client) BuildPut(body any) *RequestBuilder {
	return newRequestBuilderWithBody(client, echo.PUT, body)
}

func (client *Client) BuildPatch(body any) *RequestBuilder {
	return newRequestBuilderWithBody(client, echo.PATCH, body)
}

func (client *Client) BuildGet() *RequestBuilder {
	return newRequestBuilder(client, echo.GET)
}

func (client *Client) BuildDelete() *RequestBuilder {
	return newRequestBuilder(client, echo.DELETE)
}

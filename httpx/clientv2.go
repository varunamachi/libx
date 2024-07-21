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
	"time"

	"github.com/labstack/echo/v4"
	"github.com/varunamachi/libx/data"
)

var ErrHttpRequestBuildFailed = errors.New("http request build failed")

type RequestBuilder struct {
	client      *Client
	headers     http.Header
	queryParams map[string]string
	path        string
	withAuth    bool
	errs        []error
	timeout     time.Duration

	// TODO - now only json is supported, when others are to be supported, we
	// need a way to specify encoders
	// contentType string
}

func newRequestBuilder(client *Client) *RequestBuilder {
	return &RequestBuilder{
		client:      client,
		headers:     map[string][]string{},
		queryParams: map[string]string{},
		path:        "",
		withAuth:    true,
		errs:        make([]error, 0, 10),
	}
}

// func newRequestBuilderWithBody(
// 	client *Client, method string, body any) *RequestBuilder {
// 	return &RequestBuilder{
// 		client:      client,
// 		method:      method,
// 		body:        body,
// 		headers:     map[string][]string{},
// 		queryParams: map[string]string{},
// 		path:        "",
// 		withAuth:    true,
// 		errs:        make([]error, 0, 10),
// 	}
// }

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
		rb.errs = append(rb.errs, err)
		return rb
	}

	rb.headers.Add(name, j)
	return rb
}

func (rb *RequestBuilder) CmnParam(cp *data.CommonParams) *RequestBuilder {
	return rb.
		QInt("page", cp.Page).
		QInt("pageSize", cp.PageSize).
		QStr("sort", cp.Sort).
		QBool("sortDesc", cp.SortDescending).
		Filter(cp.Filter)
}

func (rb *RequestBuilder) Filter(f *data.Filter) *RequestBuilder {
	return rb.QJson("filter", f)
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
		rb.errs = append(rb.errs, err)
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
			rb.errs = append(rb.errs,
				fmt.Errorf("invalid param type for path: '%T'", p))
			return rb
		}
		if err != nil {
			rb.errs = append(rb.errs, err)
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

func (rb *RequestBuilder) WithTimeout(duration time.Duration) *RequestBuilder {
	rb.timeout = duration
	return rb
}

func (rb *RequestBuilder) Exec(
	gtx context.Context, method string, body any) *ApiResult {
	if len(rb.errs) != 0 {
		return &ApiResult{
			err:            ErrHttpRequestBuildFailed,
			target:         method + " " + rb.path,
			reqBuildErrors: rb.errs,
		}
	}

	var err error
	var bodyBytes []byte = nil
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return newErrorResult(nil, err, "failed to marshal data")
		}
	}
	// if !strings.HasSuffix(rb.client.contextRoot, "/") {
	// 	rb.path = "/" + rb.path
	// }
	fullUrl := rb.client.address + path.Clean(rb.client.contextRoot+"/"+rb.path)

	if rb.timeout.Seconds() != 0 {
		var cancel context.CancelFunc
		gtx, cancel = context.WithTimeout(gtx, rb.timeout)
		defer cancel()
	}
	req, err := http.NewRequestWithContext(
		gtx, method, fullUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return newErrorResult(req, err, "failed to create http request")
	}

	req.Header = rb.headers
	req.Header.Add("Content-Type", "application/json")
	if rb.withAuth {
		authHeader := fmt.Sprintf("Bearer %s", rb.client.token)
		req.Header.Add("Authorization", authHeader)
	}

	resp, err := rb.client.Do(req)
	if err != nil {
		return newErrorResult(req, err, "Failed to perform http request")
	}

	return newApiResult(req, resp)
}

func (rb *RequestBuilder) Post(gtx context.Context, body any) *ApiResult {
	return rb.Exec(gtx, echo.POST, body)
}

func (rb *RequestBuilder) Put(gtx context.Context, body any) *ApiResult {
	return rb.Exec(gtx, echo.PUT, body)
}

func (rb *RequestBuilder) Patch(gtx context.Context, body any) *ApiResult {
	return rb.Exec(gtx, echo.PATCH, body)
}

func (rb *RequestBuilder) Get(gtx context.Context) *ApiResult {
	return rb.Exec(gtx, echo.GET, nil)
}

func (rb *RequestBuilder) Delete(gtx context.Context) *ApiResult {
	return rb.Exec(gtx, echo.DELETE, nil)
}

func encodeJsonUrl(obj any) (string, error) {
	j, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	return url.PathEscape(string(j)), nil
}

func (client *Client) Build() *RequestBuilder {
	return newRequestBuilder(client)
}

// func (client *Client) BuildPost(body any) *RequestBuilder {
// 	return newRequestBuilderWithBody(client, echo.POST, body)
// }

// func (client *Client) BuildPut(body any) *RequestBuilder {
// 	return newRequestBuilderWithBody(client, echo.PUT, body)
// }

// func (client *Client) BuildPatch(body any) *RequestBuilder {
// 	return newRequestBuilderWithBody(client, echo.PATCH, body)
// }

// func (client *Client) BuildGet() *RequestBuilder {
// 	return newRequestBuilder(client, echo.GET)
// }

// func (client *Client) BuildDelete() *RequestBuilder {
// 	return newRequestBuilder(client, echo.DELETE)
// }

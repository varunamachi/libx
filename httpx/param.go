package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/str"
)

var (
	ErrHttpParam = errors.New("error.http.param")
)

type ParamGetter struct {
	etx  echo.Context
	errs map[string]error
}

func NewParamGetter(etx echo.Context) *ParamGetter {
	return &ParamGetter{
		etx:  etx,
		errs: make(map[string]error),
	}
}

func (pm *ParamGetter) Str(name string) string {
	return pm.etx.Param(name)
}

func (pm *ParamGetter) Int(name string) int {
	param := pm.etx.Param(name)
	val, err := strconv.Atoi(param)
	if err != nil {
		pm.errs[name] = err
	}
	return val
}

func (pm *ParamGetter) Int64(name string) int64 {
	param := pm.etx.Param(name)
	val, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		pm.errs[name] = err
	}
	return val
}

func (pm *ParamGetter) UInt(name string) uint {
	param := pm.etx.Param(name)
	val, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		pm.errs[name] = err
	}
	return uint(val)
}

func (pm *ParamGetter) UInt64(name string) uint64 {
	param := pm.etx.Param(name)
	val, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		pm.errs[name] = err
	}
	return val
}

func (pm *ParamGetter) Float64(name string) float64 {
	param := pm.etx.Param(name)
	val, err := strconv.ParseFloat(param, 64)
	if err != nil {
		pm.errs[name] = err
	}
	return val
}

func (pm *ParamGetter) Bool(name string) bool {
	param := pm.etx.Param(name)
	if str.EqFold(param, "true", "on") {
		return true
	} else if str.EqFold(param, "false", "off") {
		return false
	}
	pm.errs[name] = errors.New("invalid string for bool param")
	return false
}

func (pm *ParamGetter) QueryStr(name string) string {
	if !pm.etx.QueryParams().Has(name) {
		pm.errs[name] = errors.New("query param not found")
	}
	return pm.etx.QueryParam(name)
}

func (pm *ParamGetter) QueryInt(name string) int {
	param := pm.etx.QueryParam(name)
	val, err := strconv.Atoi(param)
	if err != nil {
		pm.errs[name] = err
	}
	return val
}

func (pm *ParamGetter) QueryInt64(name string) int64 {
	param := pm.etx.QueryParam(name)
	val, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		pm.errs[name] = err
	}
	return val
}

func (pm *ParamGetter) QueryUInt(name string) uint {
	param := pm.etx.QueryParam(name)
	val, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		pm.errs[name] = err
	}
	return uint(val)
}

func (pm *ParamGetter) QueryUInt64Param(name string) uint64 {
	param := pm.etx.QueryParam(name)
	val, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		pm.errs[name] = err
	}
	return val
}

func (pm *ParamGetter) QueryFloat64(name string) float64 {
	param := pm.etx.QueryParam(name)
	val, err := strconv.ParseFloat(param, 64)
	if err != nil {
		pm.errs[name] = err
	}
	return val
}

func (pm *ParamGetter) QueryBool(name string) bool {
	param := pm.etx.QueryParam(name)
	if str.EqFold(param, "true", "on") {
		return true
	} else if str.EqFold(param, "false", "off") {
		return false
	}
	pm.errs[name] = errors.New("invalid string for bool param")
	return false
}

func (pm *ParamGetter) QueryStrOr(name, def string) string {
	if !pm.etx.QueryParams().Has(name) {
		return def
	}
	return pm.etx.QueryParam(name)
}

func (pm *ParamGetter) QueryIntOr(name string, def int) int {
	param := pm.etx.QueryParam(name)
	val, err := strconv.Atoi(param)
	if err != nil {
		return def
	}
	return val
}

func (pm *ParamGetter) QueryInt64Or(name string, def int64) int64 {
	param := pm.etx.QueryParam(name)
	val, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return def
	}
	return val
}

func (pm *ParamGetter) QueryUIntOr(name string, def uint) uint {
	param := pm.etx.QueryParam(name)
	val, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return def
	}
	return uint(val)
}

func (pm *ParamGetter) QueryUInt64Or(name string, def uint64) uint64 {
	param := pm.etx.QueryParam(name)
	val, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return def
	}
	return val
}

func (pm *ParamGetter) QueryFloat64Or(name string, def float64) float64 {
	param := pm.etx.QueryParam(name)
	val, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return def
	}
	return val
}

func (pm *ParamGetter) QueryBoolOr(name string, def bool) bool {
	param := pm.etx.QueryParam(name)
	if str.EqFold(param, "true", "on") {
		return true
	} else if str.EqFold(param, "false", "off") {
		return false
	}
	return def
}

func (pm *ParamGetter) QueryJSON(name string, out interface{}) *ParamGetter {
	val := pm.etx.QueryParam(name)
	if len(val) == 0 {
		pm.errs[name] = errors.New("could not find json param")
		return pm
	}
	decoded, err := url.PathUnescape(val)
	if err != nil {
		pm.errs[name] = err
		return pm
	}
	if err = json.Unmarshal([]byte(decoded), out); err != nil {
		pm.errs[name] = err
		return pm
	}
	return pm
}

func (pm *ParamGetter) HasError() bool {
	return len(pm.errs) != 0
}

func (pm *ParamGetter) Error() error {
	if len(pm.errs) == 0 {
		return nil
	}
	buf := bytes.NewBufferString("http parameter error: ")
	index := 0
	for k := range pm.errs {
		buf.WriteString(k)
		if index < len(pm.errs)-1 {
			buf.WriteString(", ")
		}
		index++
	}
	return errx.Errf(ErrHttpParam, buf.String())
}

func (pm *ParamGetter) BadReqError() error {
	if len(pm.errs) == 0 {
		return nil
	}
	buf := bytes.NewBufferString("http parameter error: ")
	index := 0
	for k := range pm.errs {
		buf.WriteString(k)
		if index < len(pm.errs)-1 {
			buf.WriteString(", ")
		}
		index++
	}
	// return errx.Errf(ErrHttpParam, buf.String())
	return errx.BadReq(buf.String())
}

func (pm *ParamGetter) WriteDetailedError(w io.Writer) {
	for p, e := range pm.errs {
		if len(p) > 15 {
			p = p[:15]
		}
		fmt.Fprintf(w, "%-15s  %v", p, e)
	}
}

func MustGetEndpoint(etx echo.Context) *Endpoint {
	obj := etx.Get("endpoint")
	ep, ok := obj.(*Endpoint)
	if !ok {
		panic("failed to get endpoint info from echo.Context")
	}
	return ep
}

func StrMsg(err *echo.HTTPError) string {
	msg, ok := err.Message.(string)
	if !ok {
		return ""
	}
	return msg
}

func MustGetUser(etx echo.Context) auth.User {
	obj := etx.Get("user")
	user, ok := obj.(auth.User)
	if !ok {
		panic("failed to get user info from echo.Context")
	}
	return user
}

func GetUserId(etx echo.Context) string {
	ido := etx.Get("userId")
	if id, ok := ido.(string); ok {
		return id
	}

	obj := etx.Get("user")
	user, ok := obj.(auth.User)
	if !ok {
		return ""
	}
	return user.Id()
}

func GetUser[T auth.User](gtx context.Context) T {
	val := gtx.Value(UserKey)
	if val == nil {
		log.Warn().Msg("user not in context")
		var t T
		return t
	}

	user, ok := val.(T)
	if !ok {
		log.Warn().Msg("invalid context user")
		return user
	}

	return user
}

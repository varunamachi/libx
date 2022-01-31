package env

import (
	"os"
	"strconv"

	"github.com/varunamachi/libx/str"
)

var config *Config

type HttpxConfig struct {
	PrintAllAccess bool
}

type Config struct {
	HttpxConfig
}

func init() {
	config = &Config{
		HttpxConfig: HttpxConfig{
			PrintAllAccess: Bool("SOUSE_PRINT_ALL_HTTP_ACCESS"),
		},
	}
}

func GetConfig() *Config {
	return config
}

func Bool(name string) bool {
	return str.EqFold(os.Getenv(name), "true", "on")
}

func Int64(name string, def int64) int64 {
	ev := os.Getenv(name)
	val, err := strconv.ParseInt(ev, 10, 64)
	if err != nil {
		return def
	}
	return val
}

func UInt64(name string, def uint64) uint64 {
	ev := os.Getenv(name)
	val, err := strconv.ParseUint(ev, 10, 64)
	if err != nil {
		return def
	}
	return val
}

func Float64(name string, def float64) float64 {
	ev := os.Getenv(name)
	val, err := strconv.ParseFloat(ev, 64)
	if err != nil {
		return def
	}
	return val
}

func String(name, def string) string {
	val := os.Getenv(name)
	if val == "" {
		return def
	}
	return val
}

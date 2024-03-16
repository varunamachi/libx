package rt

import (
	"os"
	"strconv"

	"github.com/varunamachi/libx/str"
)

// var config *Config

// type HttpxConfig struct {
// 	PrintAllAccess bool
// }

// type Config struct {
// 	HttpxConfig
// }

// func init() {
// 	config = &Config{
// 		HttpxConfig: HttpxConfig{
// 			PrintAllAccess: EnvBool("SOUSE_PRINT_ALL_HTTP_ACCESS"),
// 		},
// 	}
// }

// func Env() *Config {
// 	return config
// }

func EnvBool(name string, def bool) bool {
	return str.EqFold(os.Getenv(name), "true", "on")
}

func EnvInt64(name string, def int64) int64 {
	ev := os.Getenv(name)
	val, err := strconv.ParseInt(ev, 10, 64)
	if err != nil {
		return def
	}
	return val
}

func EnvUInt64(name string, def uint64) uint64 {
	ev := os.Getenv(name)
	val, err := strconv.ParseUint(ev, 10, 64)
	if err != nil {
		return def
	}
	return val
}

func EnvInt(name string, def int) int {
	ev := os.Getenv(name)
	val, err := strconv.Atoi(ev)
	if err != nil {
		return def
	}
	return val
}

func EnvUInt(name string, def uint) uint {
	ev := os.Getenv(name)
	val, err := strconv.Atoi(ev)
	if err != nil {
		return def
	}
	return uint(val)
}

func EnvFloat64(name string, def float64) float64 {
	ev := os.Getenv(name)
	val, err := strconv.ParseFloat(ev, 64)
	if err != nil {
		return def
	}
	return val
}

func EnvString(name, def string) string {
	val := os.Getenv(name)
	if val == "" {
		return def
	}
	return val
}

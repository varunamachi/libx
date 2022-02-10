package iox

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx/errx"
	"golang.org/x/term"
)

//PrintJSON - dumps JSON representation of given data to stdout
func PrintJSON(o interface{}) {
	b, err := json.MarshalIndent(o, "", "    ")
	if err != nil {
		log.Error().Err(err).Msg("failed to write formatted JSON to console")
		return
	}

	fmt.Println(string(b))
}

func LoadJsonFile(path string, out interface{}) error {
	reader, err := os.Open(path)
	if err != nil {
		log.Error().Err(err).Str("path", path).
			Msg("failed to open JSON file")
		return errx.Errf(err, "failed to open JSON file at %s", path)
	}
	if err = LoadJson(reader, out); err != nil {
		return err
	}
	return nil
}

func LoadJson(reader io.Reader, out interface{}) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		const msg = "Failed to read from reader"
		log.Error().Err(err).Msg(msg)
		return errx.Errf(err, msg)
	}

	if err = json.Unmarshal(data, out); err != nil {
		const msg = "Failed to decode JSON data"
		log.Error().Err(err).Msg(msg)
		return errx.Errf(err, "Failed to decode JSON data")
	}
	return nil
}

//ExistsAsFile - checks if a regular file exists at given path. If a error
//occurs while stating whatever exists at given location, false is returned
func ExistsAsFile(path string) (yes bool) {
	stat, err := os.Stat(path)
	if err == nil && !stat.IsDir() {
		yes = true
	}
	return yes
}

//ExistsAsDir - checks if a directory exists at given path. If a error
//occurs while stating whatever exists at given location, false is returned
func ExistsAsDir(path string) (yes bool) {
	stat, err := os.Stat(path)
	if err == nil && stat.IsDir() {
		yes = true
	}
	return yes
}

func MustGetUserHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get user home")
	}
	return home
}

//FormattedJSON - converts given data to JSON and returns as pretty printed
func FormattedJSON(o interface{}) (string, error) {
	b, err := json.MarshalIndent(o, "", "    ")
	if err != nil {
		return "", errx.Errf(err, "failed to generate formatted JSON")
	}
	return string(b), nil
}

//askSecret - asks password from user, does not echo charectors
func askSecret() (secret string, err error) {
	var pbyte []byte
	pbyte, err = term.ReadPassword(int(syscall.Stdin))
	if err == nil {
		secret = string(pbyte)
		fmt.Println()
	}
	return secret, err
}

//AskPassword - asks password, prints the given name before asking
func AskPassword(name string) (secret string) {
	fmt.Print(name + ": ")
	secret, _ = askSecret()
	return secret
}

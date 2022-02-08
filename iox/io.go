package iox

import (
	"encoding/json"
	"fmt"
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

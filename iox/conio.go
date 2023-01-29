package iox

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

//askSecret - asks password from user, does not echo charectors
func askSecret() (string, error) {
	pbyte, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	secret := string(pbyte)
	fmt.Println()
	return secret, nil
}

//AskPassword - asks password, prints the given name before asking
func AskPassword(name string) string {
	fmt.Print(name + ": ")
	secret, _ := askSecret()
	return strings.TrimSpace(secret)
}

func AskDangerous(question string, def bool) bool {
	uir := NewUserInputReader(os.Stdin, os.Stdout)
	return uir.BoolOr(question, def)
}

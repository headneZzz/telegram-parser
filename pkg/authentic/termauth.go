package authentic

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gotd/td/tg"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
)

// TermAuth implements authentication via terminal.
type TermAuth struct {
	NoSignUp

	UserPhone string
}

func (a TermAuth) Phone(_ context.Context) (string, error) {
	return a.UserPhone, nil
}

func (a TermAuth) Password(_ context.Context) (string, error) {
	fmt.Print("Enter 2FA password: ")
	bytePwd, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytePwd)), nil
}

func (a TermAuth) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {
	fmt.Print("Enter code: ")
	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(code), nil
}

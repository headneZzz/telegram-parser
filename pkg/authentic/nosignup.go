package authentic

import (
	"context"
	"errors"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

// NoSignUp can be embedded to prevent signing up.
type NoSignUp struct{}

func (c NoSignUp) SignUp(ctx context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("not implemented")
}

func (c NoSignUp) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

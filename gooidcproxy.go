package gooidcproxy

import (
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Tokens struct {
	Oauth2Token *oauth2.Token
	IdToken     *oidc.IDToken
	IdTokenRaw  string
}

func (t Tokens) IsAccesTokenExpired(now time.Time) bool {
	return now.After(t.Oauth2Token.Expiry)
}

// TODO: test
func (t Tokens) IsRefreshTokenExpired(now time.Time) bool {
	return now.After(t.IdToken.Expiry)
}

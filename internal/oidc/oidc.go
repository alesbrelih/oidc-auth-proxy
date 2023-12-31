package oidc

import (
	"context"
	"errors"
	"fmt"

	gooidcproxy "github.com/alesbrelih/oidc-auth-proxy"
	"github.com/alesbrelih/oidc-auth-proxy/internal/config"
	"github.com/alesbrelih/oidc-auth-proxy/internal/packageerrors"
	"github.com/alesbrelih/oidc-auth-proxy/internal/strings"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type SignInRequest struct {
	Location string
	State    string
	Nonce    string
}

type OIDC interface {
	SignIn() SignInRequest
	RefreshAccessToken(ctx context.Context, token gooidcproxy.Tokens) (gooidcproxy.Tokens, error)
	Exchange(ctx context.Context, code string) (gooidcproxy.Tokens, error)
}

func New(ctx context.Context, cfg config.Config) (OIDC, error) {
	provider, err := oidc.NewProvider(ctx, cfg.Issuer)
	if err != nil {
		return nil, err
	}

	oauth2Config := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  cfg.RedirectURL,
		Scopes:       []string{"openid"},
	}

	oidcConfig := &oidc.Config{
		ClientID: cfg.ClientID,
	}

	verifier := provider.Verifier(oidcConfig)

	return &service{
		oauth2Config: &oauth2Config,
		verifier:     verifier,
	}, nil
}

type service struct {
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
}

func (s *service) SignIn() SignInRequest {
	state := strings.Random()
	nonce := strings.Random()
	authURL := s.oauth2Config.AuthCodeURL(state, oidc.Nonce(nonce))

	return SignInRequest{
		Location: authURL,
		State:    state,
		Nonce:    nonce,
	}
}

func (s *service) RefreshAccessToken(ctx context.Context, token gooidcproxy.Tokens) (gooidcproxy.Tokens, error) {
	tokenSource := s.oauth2Config.TokenSource(ctx, &oauth2.Token{
		RefreshToken: token.Oauth2Token.RefreshToken,
	})

	oauth2Token, err := tokenSource.Token()
	if err != nil {
		return gooidcproxy.Tokens{}, err
	}

	token.Oauth2Token = oauth2Token

	return token, nil
}
func (s *service) Exchange(ctx context.Context, code string) (gooidcproxy.Tokens, error) {
	oauth2Token, err := s.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return gooidcproxy.Tokens{},
			packageerrors.ErrInternal.
				WithErr(fmt.Errorf("can't exchange the code secret: %w", err))
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return gooidcproxy.Tokens{},
			packageerrors.ErrInternal.
				WithErr(errors.New("no ID token"))
	}

	idToken, err := s.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return gooidcproxy.Tokens{},
			packageerrors.ErrInternal.
				WithErr(fmt.Errorf("can't verify id token: %s", err))
	}

	return gooidcproxy.Tokens{
		Oauth2Token: oauth2Token,
		IdToken:     idToken,
		IdTokenRaw:  rawIDToken,
	}, nil
}

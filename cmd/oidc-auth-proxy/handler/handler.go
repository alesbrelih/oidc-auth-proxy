package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	gooidcproxy "github.com/alesbrelih/oidc-auth-proxy"
	"github.com/alesbrelih/oidc-auth-proxy/internal/generated/oidc/api"
	oidcPkg "github.com/alesbrelih/oidc-auth-proxy/internal/oidc"
	"github.com/alesbrelih/oidc-auth-proxy/internal/packageerrors"
	stringsPkg "github.com/alesbrelih/oidc-auth-proxy/internal/strings"
	"github.com/alesbrelih/oidc-auth-proxy/internal/transform"
)

const (
	sessionCookieName = "_go_oidc_auth_proxy"
	stateCookieName   = "_go_oidc_auth_proxy_state"
	nonceCookieName   = "_go_oidc_auth_proxy_nounce"
)

func New(
	oidcSvc oidcPkg.OIDC,
	transformerSvc transform.Transformer,
) api.Handler {
	return &handler{
		oidcSvc:        oidcSvc,
		transformerSvc: transformerSvc,
		sessions:       map[string]gooidcproxy.Tokens{},
		now:            time.Now,
	}
}

type handler struct {
	oidcSvc        oidcPkg.OIDC
	transformerSvc transform.Transformer
	sessions       map[string]gooidcproxy.Tokens
	now            func() time.Time
}

func (h *handler) OidcSignInGet(ctx context.Context) (*api.OidcSignInGetFound, error) {
	signInRedirectionURL := h.oidcSvc.SignIn()

	stateCookie := http.Cookie{
		Name:  stateCookieName,
		Value: signInRedirectionURL.State,
	}

	nonceCookie := http.Cookie{
		Name:  nonceCookieName,
		Value: signInRedirectionURL.Nonce,
	}

	return &api.OidcSignInGetFound{
		Location: api.NewOptString(signInRedirectionURL.Location),
		SetCookie: []string{
			stateCookie.String(),
			nonceCookie.String(),
		},
	}, nil
}

func (h *handler) OidcAuthGet(ctx context.Context, params api.OidcAuthGetParams) (api.OidcAuthGetRes, error) {
	tokens, ok := h.sessions[params.GoOidcAuthProxy.Value]
	if !ok {
		return &api.OidcAuthGetUnauthorized{}, nil
	}

	if tokens.IsRefreshTokenExpired(h.now()) {
		return &api.OidcAuthGetUnauthorized{}, nil
	}

	if tokens.IsAccesTokenExpired(h.now()) {
		token, err := h.oidcSvc.RefreshAccessToken(ctx, tokens)
		if err != nil {
			return &api.OidcAuthGetUnauthorized{}, nil
		}

		h.sessions[params.GoOidcAuthProxy.Value] = token
	}

	headerValue, err := h.transformerSvc.ClaimsHeader("keycloak", h.sessions[params.GoOidcAuthProxy.Value].IdTokenRaw)
	if err != nil {
		slog.Error("error creating header value: %s", err)

		return nil, err
	}

	return &api.OidcAuthGetAccepted{
		XClaims: api.NewOptString(headerValue),
	}, nil
}

func (h *handler) OidcCallbackGet(ctx context.Context, params api.OidcCallbackGetParams) (*api.OidcCallbackGetFound, error) {
	if params.State.Value != params.GoOidcAuthProxyState.Value {
		return nil, &api.ErrRespStatusCode{
			StatusCode: http.StatusBadRequest,
		}
	}

	tokens, err := h.oidcSvc.Exchange(ctx, params.Code.Value)
	if err != nil {
		slog.Error("error exchanging token: %s", err)

		return nil, err
	}

	if tokens.IdToken.Nonce != params.GoOidcAuthProxyNounce.Value {
		return nil, &api.ErrRespStatusCode{
			StatusCode: http.StatusBadRequest,
		}
	}

	tokenCookieValue := stringsPkg.Random()
	cookie := http.Cookie{
		Name:  sessionCookieName,
		Value: tokenCookieValue,
		Path:  "/",
	}

	h.sessions[tokenCookieValue] = tokens

	return &api.OidcCallbackGetFound{
		SetCookie: api.NewOptString(cookie.String()),
		Location:  api.NewOptString("http://localhost"),
	}, nil

}

func (*handler) NewError(ctx context.Context, err error) *api.ErrRespStatusCode {
	if cast, ok := err.(*api.ErrRespStatusCode); ok {
		return cast
	}

	if cast, ok := err.(*packageerrors.Error); ok {
		return &api.ErrRespStatusCode{
			StatusCode: cast.Code,
			Response: api.ErrResp{
				Error: api.NewOptString(cast.Message),
			},
		}
	}

	return &api.ErrRespStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: api.ErrResp{
			Error: api.NewOptString(http.StatusText(http.StatusInternalServerError)),
		},
	}
}

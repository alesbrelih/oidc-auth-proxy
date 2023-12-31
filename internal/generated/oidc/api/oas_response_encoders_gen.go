// Code generated by ogen, DO NOT EDIT.

package api

import (
	"net/http"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/ogen-go/ogen/conv"
	"github.com/ogen-go/ogen/uri"
)

func encodeOidcAuthGetResponse(response OidcAuthGetRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *OidcAuthGetAccepted:
		// Encoding response headers.
		{
			h := uri.NewHeaderEncoder(w.Header())
			// Encode "X-Claims" header.
			{
				cfg := uri.HeaderParameterEncodingConfig{
					Name:    "X-Claims",
					Explode: false,
				}
				if err := h.EncodeParam(cfg, func(e uri.Encoder) error {
					if val, ok := response.XClaims.Get(); ok {
						return e.EncodeValue(conv.StringToString(val))
					}
					return nil
				}); err != nil {
					return errors.Wrap(err, "encode X-Claims header")
				}
			}
		}
		w.WriteHeader(202)
		span.SetStatus(codes.Ok, http.StatusText(202))

		return nil

	case *OidcAuthGetUnauthorized:
		w.WriteHeader(401)
		span.SetStatus(codes.Error, http.StatusText(401))

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeOidcCallbackGetResponse(response *OidcCallbackGetFound, w http.ResponseWriter, span trace.Span) error {
	// Encoding response headers.
	{
		h := uri.NewHeaderEncoder(w.Header())
		// Encode "Location" header.
		{
			cfg := uri.HeaderParameterEncodingConfig{
				Name:    "Location",
				Explode: false,
			}
			if err := h.EncodeParam(cfg, func(e uri.Encoder) error {
				if val, ok := response.Location.Get(); ok {
					return e.EncodeValue(conv.StringToString(val))
				}
				return nil
			}); err != nil {
				return errors.Wrap(err, "encode Location header")
			}
		}
		// Encode "Set-Cookie" header.
		{
			cfg := uri.HeaderParameterEncodingConfig{
				Name:    "Set-Cookie",
				Explode: false,
			}
			if err := h.EncodeParam(cfg, func(e uri.Encoder) error {
				if val, ok := response.SetCookie.Get(); ok {
					return e.EncodeValue(conv.StringToString(val))
				}
				return nil
			}); err != nil {
				return errors.Wrap(err, "encode Set-Cookie header")
			}
		}
	}
	w.WriteHeader(302)
	span.SetStatus(codes.Ok, http.StatusText(302))

	return nil
}

func encodeOidcSignInGetResponse(response *OidcSignInGetFound, w http.ResponseWriter, span trace.Span) error {
	// Encoding response headers.
	{
		h := uri.NewHeaderEncoder(w.Header())
		// Encode "Location" header.
		{
			cfg := uri.HeaderParameterEncodingConfig{
				Name:    "Location",
				Explode: false,
			}
			if err := h.EncodeParam(cfg, func(e uri.Encoder) error {
				if val, ok := response.Location.Get(); ok {
					return e.EncodeValue(conv.StringToString(val))
				}
				return nil
			}); err != nil {
				return errors.Wrap(err, "encode Location header")
			}
		}
		// Encode "Set-Cookie" header.
		{
			cfg := uri.HeaderParameterEncodingConfig{
				Name:    "Set-Cookie",
				Explode: false,
			}
			if err := h.EncodeParam(cfg, func(e uri.Encoder) error {
				return e.EncodeArray(func(e uri.Encoder) error {
					for i, item := range response.SetCookie {
						if err := func() error {
							return e.EncodeValue(conv.StringToString(item))
						}(); err != nil {
							return errors.Wrapf(err, "[%d]", i)
						}
					}
					return nil
				})
			}); err != nil {
				return errors.Wrap(err, "encode Set-Cookie header")
			}
		}
	}
	w.WriteHeader(302)
	span.SetStatus(codes.Ok, http.StatusText(302))

	return nil
}

func encodeErrorResponse(response *ErrRespStatusCode, w http.ResponseWriter, span trace.Span) error {
	w.Header().Set("Content-Type", "application/json")
	code := response.StatusCode
	if code == 0 {
		// Set default status code.
		code = http.StatusOK
	}
	w.WriteHeader(code)
	st := http.StatusText(code)
	if code >= http.StatusBadRequest {
		span.SetStatus(codes.Error, st)
	} else {
		span.SetStatus(codes.Ok, st)
	}

	e := jx.GetEncoder()
	response.Response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}
	return nil

}

package transform

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"unicode"

	"github.com/alesbrelih/oidc-auth-proxy/internal/config"
	"github.com/alesbrelih/oidc-auth-proxy/internal/packageerrors"
	jwt "github.com/golang-jwt/jwt/v5"
)

//go:embed default.tmpl
var DefaultTemplate string

type TemplateValues struct {
	AuthType string
	Claims   map[string]interface{}
}

func New(cfg config.Config) (Transformer, error) {
	tmpl, err := getTemplate(cfg)
	if err != nil {
		return nil, packageerrors.ErrInternal.
			WithErr(fmt.Errorf("error getting template: %w", err))
	}

	t, err := template.
		New("authTemplate").
		Funcs(template.FuncMap{
			"json": func(v interface{}) string {
				js, _ := json.Marshal(v)
				return string(js)
			},
		}).
		Parse(tmpl)

	if err != nil {
		return nil, packageerrors.ErrInternal.
			WithErr(fmt.Errorf("error initializing template: %w", err))
	}

	return &service{
		template: t,
		parser:   jwt.NewParser(),
	}, nil
}

type Transformer interface {
	ClaimsHeader(authType string, token string) (string, error)
}

type service struct {
	template *template.Template
	parser   *jwt.Parser
}

// ClaimsHeader creates custom claims header in string format.
func (s *service) ClaimsHeader(authType string, token string) (string, error) {
	mapClaims := jwt.MapClaims{}
	_, _, err := s.parser.ParseUnverified(token, mapClaims)
	if err != nil {
		return "", packageerrors.ErrInternal.
			WithErr(fmt.Errorf("could not verify token: %w", err))
	}

	templateValues := TemplateValues{
		AuthType: authType,
		Claims:   mapClaims,
	}

	var result bytes.Buffer
	err = s.template.Execute(&result, templateValues)
	if err != nil {
		return "", packageerrors.ErrInternal.
			WithErr(fmt.Errorf("could not execute template: %w", err))
	}

	removeWhitespaces := func(str string) string {
		return strings.Join(strings.FieldsFunc(str, func(r rune) bool {
			return unicode.IsSpace(r)
		}), "")
	}

	headerValue := removeWhitespaces(result.String())

	return headerValue, nil
}

func getTemplate(cfg config.Config) (string, error) {
	headerValueTemplate := DefaultTemplate
	if cfg.CustomTemplatePath != "" {
		file, err := os.Open(cfg.CustomTemplatePath)
		if err != nil {
			return "", packageerrors.ErrInternal.
				WithErr(fmt.Errorf("custom template was provided but couldnt read it: %w", err))
		}

		defer file.Close()

		templateBytes, err := io.ReadAll(file)
		if err != nil {
			return "", packageerrors.ErrInternal.
				WithErr(fmt.Errorf("error reading custom template: %w", err))
		}

		headerValueTemplate = string(templateBytes)
	}

	return headerValueTemplate, nil
}

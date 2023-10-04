package transform

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"strings"
	"text/template"
	"unicode"

	jwt "github.com/golang-jwt/jwt/v5"
)

//go:embed default.tmpl
var DefaultTemplate string

type TemplateValues struct {
	AuthType string
	Claims   map[string]interface{}
}

func New(tmpl string) (Transformer, error) {
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
		return nil, err
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
		return "", err
	}

	templateValues := TemplateValues{
		AuthType: authType,
		Claims:   mapClaims,
	}

	var result bytes.Buffer
	err = s.template.Execute(&result, templateValues)
	if err != nil {
		return "", err
	}

	headerValue := strings.Join(strings.FieldsFunc(result.String(), func(r rune) bool {
		return unicode.IsSpace(r)
	}), "")

	return headerValue, nil
}

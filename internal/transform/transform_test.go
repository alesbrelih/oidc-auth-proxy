package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	s, err := New(DefaultTemplate)
	assert.NoError(t, err)

	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	got, err := s.ClaimsHeader("keycloak", jwt)
	assert.NoError(t, err)

	want := `{"auth_typ":"keycloak","claims":[{"typ":"iat","val":"1516239022"},{"typ":"name","val":""JohnDoe""},{"typ":"sub","val":""1234567890""}]}`

	assert.Equal(t, want, got)
}

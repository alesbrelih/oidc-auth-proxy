package strings

import (
	"github.com/google/uuid"
)

func Random() string {
	return uuid.New().String()
}

package githubcomment

import (
	"strings"
	"unicode"

	"github.com/google/uuid"
)

type ID string

func (i ID) GetID() string {
	var sb strings.Builder

	runes := []rune(i)

	if len(runes) <= 0 {
		runes = []rune(uuid.New().String())
	}

	for i := 0; i < len(runes); i++ {
		if !unicode.IsNumber(runes[i]) && !unicode.IsLetter(runes[i]) && runes[i] != '-' {
			continue
		}
		sb.WriteRune(runes[i])
	}

	return sb.String()
}

package githubcomment

import (
	"fmt"
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

const magic = "github-info-id"

func makeMagicMarker(id ID) string {
	return fmt.Sprintf("<!---%s-%s--->", magic, id.GetID())
}

func makeBody(body string, id ID) *string {
	var sb strings.Builder
	sb.WriteString(body)
	sb.WriteRune('\n')
	sb.WriteString(makeMagicMarker(id))
	t := sb.String()
	return &t
}

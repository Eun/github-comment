package githubcomment

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const magic = "github-info-id"

func makeMagicMarker(id ID) string {
	return fmt.Sprintf("<!---%s-%s--->", magic, id.GetID())
}

var regexID *regexp.Regexp
var regexMeta *regexp.Regexp

func init() {
	regexID = regexp.MustCompile(fmt.Sprintf(`<!---%s-([0-9a-zA-Z-]+)--->`, magic))
	regexMeta = regexp.MustCompile(`^<!---(.*)--->$`)
}

type Info struct {
	ID   ID
	Body string
	Meta interface{}
}

// ParseInfo parses the body of a comment
func ParseInfo(raw string) (*Info, error) {
	end := strings.IndexRune(raw, '\n')
	if end == -1 {
		return nil, errors.New("no marker found (invalid header)")
	}
	var info Info
	// limit the search
	info.Body = raw[end+1:]
	raw = strings.TrimSpace(raw[:end])
	matches := regexID.FindStringSubmatch(raw)
	if len(matches) != 2 {
		return nil, errors.New("no marker found (regex failure)")
	}
	// we found the marker
	info.ID = ID(matches[1])
	if info.ID.GetID() != matches[1] {
		// ID mismatch
		return nil, errors.New("id mismatch")
	}
	// jump over the marker
	raw = raw[len(matches[0]):]
	if len(raw) <= 0 {
		return &info, nil
	}
	matches = regexMeta.FindStringSubmatch(raw)
	if len(matches) != 2 {
		return nil, errors.New("no meta found (regex failure)")
	}
	// we found the meta
	if err := json.Unmarshal([]byte(matches[1]), &info.Meta); err != nil {
		return nil, err
	}
	return &info, nil
}

// Build builds a info
func (i *Info) Build() (string, error) {
	var sb strings.Builder
	sb.WriteString(makeMagicMarker(i.ID))
	if i.Meta != nil {
		sb.WriteString("<!---")
		bytes, err := json.Marshal(i.Meta)
		if err != nil {
			return "", err
		}
		sb.Write(bytes)
		sb.WriteString("--->")
	}
	sb.WriteRune('\n')
	sb.WriteString(i.Body)
	return sb.String(), nil
}

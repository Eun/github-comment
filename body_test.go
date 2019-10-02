package githubcomment

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseInfo(t *testing.T) {
	tests := []struct {
		Input string
		Info  *Info
		Error string
	}{
		{fmt.Sprintf("%s\n", makeMagicMarker(ID("123"))), &Info{ID: ID("123")}, ""},
		{fmt.Sprintf("%s\nHello World!", makeMagicMarker(ID("123"))), &Info{ID: ID("123"), Body: "Hello World!"}, ""},
		{fmt.Sprintf("%s<!---[1,2,3]--->\nHello World!", makeMagicMarker(ID("123"))), &Info{ID: ID("123"), Meta: []interface{}{float64(1), float64(2), float64(3)}, Body: "Hello World!"}, ""},

		{"Hello World", nil, "no marker found (invalid header)"},
		{"<!---github-info-id-ÖÄL--->\n", nil, "no marker found (regex failure)"},
		{"<!---github-info-id-123---><!Hello World>\n", nil, "no meta found (regex failure)"},
		{"<!---github-info-id-123---><!---Hello World--->\n", nil, "invalid character 'H' looking for beginning of value"},
	}

	for _, test := range tests {
		b, err := ParseInfo(test.Input)
		if test.Error != "" {
			require.EqualError(t, err, test.Error)
		} else {
			require.NoError(t, err)
		}

		require.Equal(t, test.Info, b)
	}
}

func TestBuildInfo(t *testing.T) {
	tests := []struct {
		Info   *Info
		Output string
	}{
		{&Info{ID: ID("123"), Meta: []interface{}{float64(1), float64(2), float64(3)}, Body: "Hello World!"}, fmt.Sprintf("%s<!---[1,2,3]--->\nHello World!", makeMagicMarker(ID("123")))},
	}

	for _, test := range tests {
		raw, err := test.Info.Build()
		require.NoError(t, err)
		require.Equal(t, test.Output, raw)
	}
}

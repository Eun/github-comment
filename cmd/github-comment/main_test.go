package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOwnerAndRepo(t *testing.T) {
	tests := []struct {
		Input string
		Owner string
		Repo  string
		Error error
	}{
		{"", "", "", errors.New("unable to parse repository")},
		{"bob", "", "", errors.New("unable to parse repository")},
		{"bob/repo", "bob", "repo", nil},
		{"bob/repo1/repo2", "bob", "repo1/repo2", nil},
	}

	for _, test := range tests {
		owner, repo, err := parseOwnerAndRepo(test.Input)
		require.Equal(t, test.Owner, owner)
		require.Equal(t, test.Repo, repo)
		require.Equal(t, test.Error, err)
	}
}

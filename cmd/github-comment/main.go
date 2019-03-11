package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	githubcomment "github.com/Eun/github-comment"
	"github.com/alecthomas/kingpin"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	idFlag         = kingpin.Flag("id", "id for this comment").String()
	repositoryFlag = kingpin.Flag("repo", "repository").PlaceHolder("owner/repo").Required().String()
	issueFlag      = kingpin.Flag("issue", "issue id").PlaceHolder("1234").Int()
	prFlag         = kingpin.Flag("pr", "pull request id").PlaceHolder("1234").Int()
	textFlag       = kingpin.Arg("text", "text to post").Required().String()
)

var version string
var commit string
var date string

func main() {
	kingpin.Version(fmt.Sprintf("%s %s %s", version, commit, date))
	kingpin.Parse()

	sanitizeFlags()

	owner, repo, err := parseOwnerAndRepo(*repositoryFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid repository `%s': %v\n", *repositoryFlag, err.Error())
		os.Exit(1)
	}

	if *issueFlag == 0 && *prFlag == 0 {
		fmt.Fprint(os.Stderr, "either --issue or --pr must be specified\n")
		os.Exit(1)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Fprint(os.Stderr, "environment GITHUB_TOKEN is not set\n")
		os.Exit(1)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	ctx := context.Background()

	var id int
	if *issueFlag > 0 {
		id = *issueFlag
	} else {
		id = *prFlag
	}
	if err = githubcomment.PostOrUpdateIssueComment(client, ctx, owner, repo, id, githubcomment.ID(*idFlag), *textFlag); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func sanitizeFlags() {
	if idFlag == nil {
		var nullString string
		idFlag = &nullString
	}

	if repositoryFlag == nil {
		var nullString string
		repositoryFlag = &nullString
	}

	if issueFlag == nil {
		var zero int
		issueFlag = &zero
	}
	if prFlag == nil {
		var zero int
		prFlag = &zero
	}
}

func parseOwnerAndRepo(s string) (owner, repo string, err error) {
	p := strings.SplitN(s, "/", 2)
	if len(p) == 2 {
		return p[0], p[1], nil
	}
	return "", "", errors.New("unable to parse repository")
}

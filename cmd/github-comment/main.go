package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	githubcomment "github.com/Eun/github-comment"
	"github.com/alecthomas/kingpin"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	yaml "gopkg.in/yaml.v2"
)

var (
	idFlag         = kingpin.Flag("id", "id for this comment").String()
	repositoryFlag = kingpin.Flag("repo", "repository").PlaceHolder("owner/repo").Required().String()
	issueFlag      = kingpin.Flag("issue", "issue id").PlaceHolder("1234").Int()
	prFlag         = kingpin.Flag("pr", "pull request id").PlaceHolder("1234").Int()

	getCmd = kingpin.Command("get", "get the text of a posted comment")

	getMetaCmd    = kingpin.Command("get-meta", "get the meta of a posted comment")
	getMetaFormat = getMetaCmd.Flag("meta-format", "format for the meta").PlaceHolder("json|yml").Default("json").String()

	postOrUpdateCmd = kingpin.Command("post", "post or update a new comment").Default()
	setMetaFormat   = postOrUpdateCmd.Flag("meta-format", "format for the meta").PlaceHolder("json|yml").Default("json").String()
	setMetaFlag     = postOrUpdateCmd.Flag("meta", "meta to set").String()
	setTextFlag     = postOrUpdateCmd.Arg("text", "text to post").String()
)

var version string
var commit string
var date string

var comment githubcomment.GithubComment

func main() {
	kingpin.Version(fmt.Sprintf("%s %s %s", version, commit, date))
	cmd := kingpin.Parse()
	sanitizeFlags()
	initComments()
	switch cmd {
	case getCmd.FullCommand():
		getText()
	case getMetaCmd.FullCommand():
		getMeta()
	case postOrUpdateCmd.FullCommand():
		postOrUpdate()
	}
}

func sanitizeFlags() {
	// general
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

	// get meta command
	if getMetaFormat == nil {
		var nullString string
		getMetaFormat = &nullString
	}

	// post command
	if setMetaFormat == nil {
		var nullString string
		setMetaFormat = &nullString
	}

	if setMetaFlag == nil {
		var nullString string
		setMetaFlag = &nullString
	}

	if setTextFlag == nil {
		var nullString string
		setTextFlag = &nullString
	}
}

func initComments() {
	var err error
	comment.Owner, comment.Repository, err = parseOwnerAndRepo(*repositoryFlag)
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

	comment.Client = github.NewClient(tc)
	comment.Context = context.Background()
}

func parseOwnerAndRepo(s string) (owner, repo string, err error) {
	p := strings.SplitN(s, "/", 2)
	if len(p) == 2 {
		return p[0], p[1], nil
	}
	return "", "", errors.New("unable to parse repository")
}

func postOrUpdate() {
	var id int
	if *issueFlag > 0 {
		id = *issueFlag
	} else {
		id = *prFlag
	}

	if *setTextFlag == "" {
		var sb strings.Builder
		_, err := io.Copy(&sb, os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read from stdin: %v\n", err.Error())
			os.Exit(1)
		}
		t := sb.String()
		setTextFlag = &t
	}
	meta, err := readMetaFromFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}

	if err = comment.PostOrUpdateIssueComment(id, githubcomment.ID(*idFlag), *setTextFlag, meta); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func get() *githubcomment.Info {
	var id int
	if *issueFlag > 0 {
		id = *issueFlag
	} else {
		id = *prFlag
	}

	info, err := comment.GetIssueComment(id, githubcomment.ID(*idFlag))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}
	return info
}

func getText() {
	fmt.Fprint(os.Stdout, get().Body)
	os.Exit(0)
}

func getMeta() {
	switch strings.ToLower(*getMetaFormat) {
	case "yml", "yaml":
		yaml.NewEncoder(os.Stdout).Encode(get().Meta)
	default:
		json.NewEncoder(os.Stdout).Encode(get().Meta)
	}

	os.Exit(0)
}

func readMetaFromFlags() (v interface{}, err error) {
	switch strings.ToLower(*setMetaFormat) {
	case "yml", "yaml":
		err = yaml.Unmarshal([]byte(*setMetaFlag), &v)
	default:
		err = json.Unmarshal([]byte(*setMetaFlag), &v)
	}
	return v, err
}

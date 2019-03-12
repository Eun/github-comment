package githubcomment

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
)

type GithubComment struct {
	Client     *github.Client
	Context    context.Context
	Owner      string
	Repository string
}

type IDMustBeSpecifiedError struct{}

func (e IDMustBeSpecifiedError) Error() string {
	return fmt.Sprintf("id cannot be empty")
}

type IssueCommentNotFoundError struct {
	ID ID
}

func (e IssueCommentNotFoundError) Error() string {
	return fmt.Sprintf("comment with the id `%s' not found", e.ID.GetID())
}

// FindIssueComment finds a issue comment and returns it
func (gc *GithubComment) FindIssueComment(issueID int, id ID) (*github.IssueComment, error) {
	if id == "" {
		return nil, IDMustBeSpecifiedError{}
	}
	comments, _, err := gc.Client.Issues.ListComments(gc.Context, gc.Owner, gc.Repository, issueID, nil)
	if err != nil {
		return nil, err
	}

	magicMarker := makeMagicMarker(id)
	for _, comment := range comments {
		if comment.Body == nil || comment.ID == nil {
			continue
		}
		if strings.Contains(*comment.Body, magicMarker) {
			return comment, nil
		}
	}
	return nil, IssueCommentNotFoundError{ID: id}
}

// PostIssueComment posts a new comment with the specified id
func (gc *GithubComment) PostIssueComment(issueID int, id ID, text string, meta interface{}) error {
	info := Info{
		ID:   id,
		Body: text,
		Meta: meta,
	}
	bodyText, err := info.Build()
	if err != nil {
		return err
	}
	_, _, err = gc.Client.Issues.CreateComment(gc.Context, gc.Owner, gc.Repository, issueID, &github.IssueComment{
		Body: &bodyText,
	})
	return err
}

// UpdateIssueComment updates an existing comment
func (gc *GithubComment) UpdateIssueComment(issueID int, id ID, text string, meta interface{}) error {
	comment, err := gc.FindIssueComment(issueID, id)
	if err != nil {
		if _, ok := err.(IssueCommentNotFoundError); !ok {
			return err
		}
		return gc.PostIssueComment(issueID, id, text, meta)
	}
	info := Info{
		ID:   id,
		Body: text,
		Meta: meta,
	}
	bodyText, err := info.Build()
	if err != nil {
		return err
	}
	_, _, err = gc.Client.Issues.EditComment(gc.Context, gc.Owner, gc.Repository, comment.GetID(), &github.IssueComment{
		Body: &bodyText,
	})
	return err
}

// PostOrUpdateIssueComment  posts an new comment if it was not able to update the existing comment,
// if you omit the ID it will always post a new comment
func (gc *GithubComment) PostOrUpdateIssueComment(issueID int, id ID, text string, meta interface{}) error {
	// if id is not specificed
	if id == "" {
		return gc.PostIssueComment(issueID, id, text, meta)
	}
	return gc.UpdateIssueComment(issueID, id, text, meta)
}

// GetIssueComment retruns the info for a comment
func (gc *GithubComment) GetIssueComment(issueID int, id ID) (*Info, error) {
	comment, err := gc.FindIssueComment(issueID, id)
	if err != nil {
		return nil, err
	}

	return ParseInfo(comment.GetBody())
}

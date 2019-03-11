package githubcomment

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
)

func PostIssueComment(client *github.Client, ctx context.Context, owner, repo string, issue int, id ID, text string) error {
	_, _, err := client.Issues.CreateComment(ctx, owner, repo, issue, &github.IssueComment{
		Body: makeBody(text, id),
	})
	return err
}

func UpdateIssueComment(client *github.Client, ctx context.Context, owner, repo string, issue int, id ID, text string) error {
	comments, _, err := client.Issues.ListComments(ctx, owner, repo, issue, nil)
	if err != nil {
		return err
	}
	if len(comments) == 0 {
		// no comments present
		return PostIssueComment(client, ctx, owner, repo, issue, id, text)
	}

	magicMarker := makeMagicMarker(id)

	for _, comment := range comments {
		if comment.Body == nil || comment.ID == nil {
			continue
		}
		if strings.Contains(*comment.Body, magicMarker) {
			_, _, err := client.Issues.EditComment(ctx, owner, repo, *comment.ID, &github.IssueComment{
				Body: makeBody(text, id),
			})
			return err
		}
	}

	// no comment matches our id
	return PostIssueComment(client, ctx, owner, repo, issue, id, text)
}

func PostOrUpdateIssueComment(client *github.Client, ctx context.Context, owner, repo string, issue int, id ID, text string) error {
	// if id is not specificed
	if id == "" {
		return PostIssueComment(client, ctx, owner, repo, issue, id, text)
	}
	return UpdateIssueComment(client, ctx, owner, repo, issue, id, text)
}

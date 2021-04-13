// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package webhook

import (
	"fmt"
	"strings"

	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/git"
	api "code.gitea.io/gitea/modules/structs"

	jsoniter "github.com/json-iterator/go"
	dingtalk "github.com/lunny/dingtalk_webhook"
)

type (
	// DingtalkPayload represents
	DingtalkPayload dingtalk.Payload
)

var (
	_ PayloadConvertor = &DingtalkPayload{}
)

// SetSecret sets the dingtalk secret
func (d *DingtalkPayload) SetSecret(_ string) {}

// JSONPayload Marshals the DingtalkPayload to json
func (d *DingtalkPayload) JSONPayload() ([]byte, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

// Create implements PayloadConvertor Create method
func (d *DingtalkPayload) Create(p *api.CreatePayload) (api.Payloader, error) {
	// created tag/branch
	refName := git.RefEndName(p.Ref)
	title := fmt.Sprintf("[%s] %s %s created", p.Repo.FullName, p.RefType, refName)

	return &DingtalkPayload{
		MsgType: "actionCard",
		ActionCard: dingtalk.ActionCard{
			Text:        title,
			Title:       title,
			HideAvatar:  "0",
			SingleTitle: fmt.Sprintf("view ref %s", refName),
			SingleURL:   p.Repo.HTMLURL + "/src/" + refName,
		},
	}, nil
}

// Delete implements PayloadConvertor Delete method
func (d *DingtalkPayload) Delete(p *api.DeletePayload) (api.Payloader, error) {
	// created tag/branch
	refName := git.RefEndName(p.Ref)
	title := fmt.Sprintf("[%s] %s %s deleted", p.Repo.FullName, p.RefType, refName)

	return &DingtalkPayload{
		MsgType: "actionCard",
		ActionCard: dingtalk.ActionCard{
			Text:        title,
			Title:       title,
			HideAvatar:  "0",
			SingleTitle: fmt.Sprintf("view ref %s", refName),
			SingleURL:   p.Repo.HTMLURL + "/src/" + refName,
		},
	}, nil
}

// Fork implements PayloadConvertor Fork method
func (d *DingtalkPayload) Fork(p *api.ForkPayload) (api.Payloader, error) {
	title := fmt.Sprintf("%s is forked to %s", p.Forkee.FullName, p.Repo.FullName)

	return &DingtalkPayload{
		MsgType: "actionCard",
		ActionCard: dingtalk.ActionCard{
			Text:        title,
			Title:       title,
			HideAvatar:  "0",
			SingleTitle: fmt.Sprintf("view forked repo %s", p.Repo.FullName),
			SingleURL:   p.Repo.HTMLURL,
		},
	}, nil
}

// Push implements PayloadConvertor Push method
func (d *DingtalkPayload) Push(p *api.PushPayload) (api.Payloader, error) {
	var (
		branchName = git.RefEndName(p.Ref)
		commitDesc string
	)

	var titleLink, linkText string
	if len(p.Commits) == 1 {
		commitDesc = "1 new commit"
		titleLink = p.Commits[0].URL
		linkText = fmt.Sprintf("view commit %s", p.Commits[0].ID[:7])
	} else {
		commitDesc = fmt.Sprintf("%d new commits", len(p.Commits))
		titleLink = p.CompareURL
		linkText = fmt.Sprintf("view commit %s...%s", p.Commits[0].ID[:7], p.Commits[len(p.Commits)-1].ID[:7])
	}
	if titleLink == "" {
		titleLink = p.Repo.HTMLURL + "/src/" + branchName
	}

	title := fmt.Sprintf("[%s:%s] %s", p.Repo.FullName, branchName, commitDesc)

	var text string
	// for each commit, generate attachment text
	for i, commit := range p.Commits {
		var authorName string
		if commit.Author != nil {
			authorName = " - " + commit.Author.Name
		}
		text += fmt.Sprintf("[%s](%s) %s", commit.ID[:7], commit.URL,
			strings.TrimRight(commit.Message, "\r\n")) + authorName
		// add linebreak to each commit but the last
		if i < len(p.Commits)-1 {
			text += "\n"
		}
	}

	return &DingtalkPayload{
		MsgType: "actionCard",
		ActionCard: dingtalk.ActionCard{
			Text:        text,
			Title:       title,
			HideAvatar:  "0",
			SingleTitle: linkText,
			SingleURL:   titleLink,
		},
	}, nil
}

// Issue implements PayloadConvertor Issue method
func (d *DingtalkPayload) Issue(p *api.IssuePayload) (api.Payloader, error) {
	text, issueTitle, attachmentText, _ := getIssuesPayloadInfo(p, noneLinkFormatter, true)

	return &DingtalkPayload{
		MsgType: "actionCard",
		ActionCard: dingtalk.ActionCard{
			Text: text + "\r\n\r\n" + attachmentText,
			//Markdown:    "# " + title + "\n" + text,
			Title:       issueTitle,
			HideAvatar:  "0",
			SingleTitle: "view issue",
			SingleURL:   p.Issue.HTMLURL,
		},
	}, nil
}

// IssueComment implements PayloadConvertor IssueComment method
func (d *DingtalkPayload) IssueComment(p *api.IssueCommentPayload) (api.Payloader, error) {
	text, issueTitle, _ := getIssueCommentPayloadInfo(p, noneLinkFormatter, true)

	return &DingtalkPayload{
		MsgType: "actionCard",
		ActionCard: dingtalk.ActionCard{
			Text:        text + "\r\n\r\n" + p.Comment.Body,
			Title:       issueTitle,
			HideAvatar:  "0",
			SingleTitle: "view issue comment",
			SingleURL:   p.Comment.HTMLURL,
		},
	}, nil
}

// PullRequest implements PayloadConvertor PullRequest method
func (d *DingtalkPayload) PullRequest(p *api.PullRequestPayload) (api.Payloader, error) {
	text, issueTitle, attachmentText, _ := getPullRequestPayloadInfo(p, noneLinkFormatter, true)

	return &DingtalkPayload{
		MsgType: "actionCard",
		ActionCard: dingtalk.ActionCard{
			Text: text + "\r\n\r\n" + attachmentText,
			//Markdown:    "# " + title + "\n" + text,
			Title:       issueTitle,
			HideAvatar:  "0",
			SingleTitle: "view pull request",
			SingleURL:   p.PullRequest.HTMLURL,
		},
	}, nil
}

// Review implements PayloadConvertor Review method
func (d *DingtalkPayload) Review(p *api.PullRequestPayload, event models.HookEventType) (api.Payloader, error) {
	var text, title string
	switch p.Action {
	case api.HookIssueReviewed:
		action, err := parseHookPullRequestEventType(event)
		if err != nil {
			return nil, err
		}

		title = fmt.Sprintf("[%s] Pull request review %s : #%d %s", p.Repository.FullName, action, p.Index, p.PullRequest.Title)
		text = p.Review.Content

	}

	return &DingtalkPayload{
		MsgType: "actionCard",
		ActionCard: dingtalk.ActionCard{
			Text:        title + "\r\n\r\n" + text,
			Title:       title,
			HideAvatar:  "0",
			SingleTitle: "view pull request",
			SingleURL:   p.PullRequest.HTMLURL,
		},
	}, nil
}

// Repository implements PayloadConvertor Repository method
func (d *DingtalkPayload) Repository(p *api.RepositoryPayload) (api.Payloader, error) {
	var title, url string
	switch p.Action {
	case api.HookRepoCreated:
		title = fmt.Sprintf("[%s] Repository created", p.Repository.FullName)
		url = p.Repository.HTMLURL
		return &DingtalkPayload{
			MsgType: "actionCard",
			ActionCard: dingtalk.ActionCard{
				Text:        title,
				Title:       title,
				HideAvatar:  "0",
				SingleTitle: "view repository",
				SingleURL:   url,
			},
		}, nil
	case api.HookRepoDeleted:
		title = fmt.Sprintf("[%s] Repository deleted", p.Repository.FullName)
		return &DingtalkPayload{
			MsgType: "text",
			Text: struct {
				Content string `json:"content"`
			}{
				Content: title,
			},
		}, nil
	}

	return nil, nil
}

// Release implements PayloadConvertor Release method
func (d *DingtalkPayload) Release(p *api.ReleasePayload) (api.Payloader, error) {
	text, _ := getReleasePayloadInfo(p, noneLinkFormatter, true)

	return &DingtalkPayload{
		MsgType: "actionCard",
		ActionCard: dingtalk.ActionCard{
			Text:        text,
			Title:       text,
			HideAvatar:  "0",
			SingleTitle: "view release",
			SingleURL:   p.Release.URL,
		},
	}, nil
}

// GetDingtalkPayload converts a ding talk webhook into a DingtalkPayload
func GetDingtalkPayload(p api.Payloader, event models.HookEventType, meta string) (api.Payloader, error) {
	return convertPayloader(new(DingtalkPayload), p, event)
}

package handlers

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"net/http"
	"strings"
	"time"
)

type commitData struct {
	Hash    string
	Message string
	URL     string
	When    time.Time
}

func (c *commitData) ShortHash() string {
	return c.Hash[:7]
}

func (c *commitData) Date() string {
	return c.When.Format("2006-01-02")
}

func newCommitData(ctx *Context, commit *object.Commit) *commitData {
	if commit == nil {
		return nil
	}
	return &commitData{
		Message: strings.Trim(commit.Message, "\n"),
		Hash:    commit.Hash.String(),
		When:    commit.Author.When,
		URL:     ctx.GetCommitURL(commit),
	}
}

type commitsViewData struct {
	Commits []commitData
	*Context
}

func gitLog(ctx *Context) error {
	r := ctx.request
	g := ctx.repo
	ref := ctx.Ref

	query := r.URL.Query()
	next := query.Get("next")
	nextRef := plumbing.ZeroHash

	if next != "" {
		nextRef = plumbing.NewHash(next)
		if nextRef.IsZero() {
			return StatusError{http.StatusBadRequest, errInvalidHash}
		}
	}

	// ... retrieves the commit history
	cIter, err := g.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return err
	}

	data := commitsViewData{Context: ctx}

	// ... just iterates over the commits
	for i := 0; i < 20; i++ {
		c, err := cIter.Next()
		if err != nil {
			break
		}
		if !nextRef.IsZero() {
			if nextRef != c.Hash {
				continue
			}
			nextRef = plumbing.ZeroHash
		}
		data.Commits = append(data.Commits, *newCommitData(ctx, c))
		if len(data.Commits) >= 20 {
			break
		}
	}

	if !nextRef.IsZero() {
		return StatusError{http.StatusNotFound, err}
	}

	return ctx.RenderTemplate("git-commits.html", data)
}

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

func newCommitData(commit *object.Commit) *commitData {
	return &commitData{
		Message: strings.Trim(commit.Message, "\n"),
		Hash:    commit.Hash.String(),
		When:    commit.Author.When,
	}
}

type commitsViewData struct {
	Commits []*commitData
	*RepoConfig
}

func gitLog(env *Env, w http.ResponseWriter, r *http.Request) error {
	rc, err := env.getRepoConfig(r)
	if err != nil {
		return StatusError{http.StatusNotFound, err}
	}

	g, err := rc.open()
	if err != nil {
		return StatusError{http.StatusNotFound, err}
	}

	ref, err := getRef(r, g)
	if err != nil {
		return StatusError{http.StatusNotFound, err}
	}

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
	cIter, err := g.Log(&git.LogOptions{From: ref})
	if err != nil {
		return err
	}

	data := commitsViewData{RepoConfig: rc}

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
		cd := newCommitData(c)
		cd.URL = env.getCommitURL(rc, c)
		data.Commits = append(data.Commits, cd)
		if len(data.Commits) >= 20 {
			break
		}
	}

	if !nextRef.IsZero() {
		return StatusError{http.StatusNotFound, err}
	}

	template, err := env.Template.GetTemplate("git-commits.html")
	if err != nil {
		return err
	}

	return template.ExecuteTemplate(w, "layout", data)
}

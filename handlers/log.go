package handlers

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"net/http"
	"strings"
)

type commitData struct {
	Hash    string
	Message string
	URL     string
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
		data.Commits = append(data.Commits, &commitData{
			Message: strings.Trim(c.Message, "\n"),
			Hash:    c.Hash.String(),
			URL:     env.getCommitURL(rc, c),
		})
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

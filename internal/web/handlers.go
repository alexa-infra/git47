package web

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gorilla/mux"
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

func newCommitData(ctx *requestContext, commit *object.Commit) *commitData {
	if commit == nil {
		return nil
	}
	url, _ := getCommitURL(ctx, commit)
	return &commitData{
		Message: strings.Trim(commit.Message, "\n"),
		Hash:    commit.Hash.String(),
		When:    commit.Author.When,
		URL:     url,
	}
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	err := fmt.Errorf("Not implemented (%s) %s", mux.CurrentRoute(r).GetName(), r.URL.Path)
	http.Error(w, err.Error(), http.StatusNotImplemented)
}

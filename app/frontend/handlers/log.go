package handlers

import (
	"github.com/alexa-infra/git47/app/frontend/server"
	"github.com/alexa-infra/git47/app/frontend/middleware"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"net/http"
	"strings"
	"time"
	"log"
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

func newCommitData(ctx *middleware.RequestContext, commit *object.Commit) *commitData {
	if commit == nil {
		return nil
	}
	url, _ := GetCommitURL(ctx, commit)
	return &commitData{
		Message: strings.Trim(commit.Message, "\n"),
		Hash:    commit.Hash.String(),
		When:    commit.Author.When,
		URL:     url,
	}
}

type commitsViewData struct {
	Commits []commitData
	*middleware.RequestContext
}

func GitLog(env *server.Env) http.HandlerFunc {
	template, err := env.GetTemplate("git-commits.html")
	if err != nil {
		log.Fatal(err)
	}
	return func (w http.ResponseWriter, r *http.Request) {
		ctx, _ := middleware.GetRequestContext(r)
		g := ctx.Repo
		ref := ctx.Ref

		query := r.URL.Query()
		next := query.Get("next")
		nextRef := plumbing.ZeroHash

		if next != "" {
			nextRef = plumbing.NewHash(next)
			if nextRef.IsZero() {
				http.Error(w, errInvalidHash.Error(), http.StatusBadRequest)
				return
			}
		}

		// ... retrieves the commit history
		cIter, err := g.Log(&git.LogOptions{From: ref.Hash()})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := commitsViewData{RequestContext: ctx}

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
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		err = env.RenderTemplate(w, template, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetLogURL(rc *middleware.RequestContext) (string, error) {
	router := rc.Env.Router
	route := router.Get("commits")
	url, err := route.URLPath("repo", rc.Config.Name, "ref", rc.Ref.Name)
	if err != nil {
		return "", err
	}
	return url.Path, nil
}

package handlers

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type summaryViewData struct {
	NumCommits      int
	NumBranches     int
	NumTags         int
	NumFiles        int
	NumContributors int
	*Context
}

func gitSummary(ctx *Context) error {
	g := ctx.repo
	ref := ctx.Ref

	data := summaryViewData{Context: ctx}

	cIter, err := g.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return err
	}
	uniqUsers := make(map[string]bool)
	cIter.ForEach(func(c *object.Commit) error {
		data.NumCommits++
		uniqUsers[c.Author.Email] = true
		return nil
	})
	data.NumContributors = len(uniqUsers)

	refs, err := g.References()
	if err != nil {
		return err
	}
	refs.ForEach(func(ref *plumbing.Reference) error {
		refName := ref.Name()
		if refName.IsBranch() {
			data.NumBranches++
		}
		if refName.IsTag() {
			data.NumTags++
		}
		return nil
	})

	tree, err := ref.Commit.Tree()
	if err != nil {
		return err
	}
	tree.Files().ForEach(func(f *object.File) error {
		data.NumFiles++
		return nil
	})

	return ctx.RenderTemplate("git-summary.html", data)
}

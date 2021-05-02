package core

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var (
	ErrRefNotFound = errors.New("Ref not found")
)

type NamedReference struct {
	Repository *git.Repository
	Name       string
	Kind       string
	Commit     *object.Commit
}

func (ref NamedReference) Hash() plumbing.Hash {
	if ref.Commit == nil {
		return plumbing.ZeroHash
	}
	return ref.Commit.Hash
}

func (ref NamedReference) Invalid() bool {
	return ref.Commit == nil
}

func GetNamedRef(g *git.Repository, ref string) (NamedReference, error) {
	branch, err := GetBranchRef(g, ref)
	if err == nil {
		return branch, nil
	}
	tag, err := GetTagRef(g, ref)
	if err == nil {
		return tag, nil
	}
	return GetCommitRef(g, ref)
}

func GetBranchRef(g *git.Repository, ref string) (NamedReference, error) {
	branch, err := g.Reference(plumbing.NewBranchReferenceName(ref), false)
	if err != nil {
		return NamedReference{}, ErrRefNotFound
	}
	hash := branch.Hash()

	commit, err := g.CommitObject(hash)
	if err != nil {
		return NamedReference{}, err
	}

	return NamedReference{
		Repository: g,
		Name:       ref,
		Kind:       "branch",
		Commit:     commit,
	}, nil
}

func GetTagRef(g *git.Repository, ref string) (NamedReference, error) {
	tag, err := g.Reference(plumbing.NewTagReferenceName(ref), false)
	if err != nil {
		return NamedReference{}, ErrRefNotFound
	}
	hash := tag.Hash()

	commit, err := g.CommitObject(hash)
	if err != nil {
		return NamedReference{}, err
	}

	return NamedReference{
		Repository: g,
		Name:       ref,
		Kind:       "tag",
		Commit:     commit,
	}, nil
}

func GetCommitRef(g *git.Repository, ref string) (NamedReference, error) {
	hash := plumbing.NewHash(ref)
	if hash.IsZero() {
		return NamedReference{}, ErrRefNotFound
	}

	commit, err := g.CommitObject(hash)
	if err != nil {
		return NamedReference{}, ErrRefNotFound
	}
	return NamedReference{
		Repository: g,
		Name:       ref,
		Kind:       "commit",
		Commit:     commit,
	}, nil
}

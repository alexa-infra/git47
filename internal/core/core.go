package core

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"io"
	"path"
	"strings"
)

var (
	ErrFileNotFound = errors.New("File not found")
)

func GetDefaultBranch(g *git.Repository) (NamedReference, error) {
	master, err := GetBranchRef(g, "master")
	if err == nil {
		return master, nil
	}
	return GetBranchRef(g, "main")
}

func GetBlob(ref NamedReference, path string) ([]byte, error) {
	if ref.Invalid() {
		return nil, ErrRefNotFound
	}
	commit := ref.Commit
	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}
	file, err := tree.File(path)
	if err != nil {
		return nil, ErrFileNotFound
	}

	reader, err := file.Reader()
	if err != nil {
		return nil, err
	}

	return io.ReadAll(reader)
}

func GetLastCommit(ref NamedReference, paths ...string) (*object.Commit, error) {
	if ref.Invalid() {
		return nil, ErrRefNotFound
	}
	g := ref.Repository
	path := path.Join(paths...)
	cIter, err := g.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}
	defer cIter.Close()

	var treeEntry *object.TreeEntry = nil
	var lastCommit *object.Commit = nil
	for {
		commit, err := cIter.Next()
		if err != nil {
			if lastCommit != nil {
				return lastCommit, nil
			}
			return nil, err
		}
		tree, err := commit.Tree()
		if err != nil {
			return nil, err
		}
		entry, err := tree.FindEntry(path)
		if err != nil {
			if lastCommit != nil {
				return lastCommit, nil
			}
			return nil, err
		}
		if treeEntry == nil {
			treeEntry = entry
		} else if *treeEntry != *entry {
			return lastCommit, nil
		}
		lastCommit = commit
	}
}

type TreeData struct {
	Files []string
	Dirs  []string
}

func GetTree(ref NamedReference, path string) (TreeData, error) {
	data := TreeData{}
	if ref.Invalid() {
		return data, ErrRefNotFound
	}

	commit := ref.Commit
	tree, err := commit.Tree()
	if err != nil {
		return data, err
	}

	if path != "" {
		tree, err = tree.Tree(path)
		if err != nil {
			return data, ErrFileNotFound
		}
		data.Dirs = append(data.Dirs, "..")
	}

	uniqDirs := make(map[string]bool)
	err = tree.Files().ForEach(func(f *object.File) error {
		if strings.Index(f.Name, "/") > 0 {
			components := strings.Split(f.Name, "/")
			folderName := components[0]
			_, ok := uniqDirs[folderName]
			if ok {
				return nil
			}
			uniqDirs[folderName] = true
			data.Dirs = append(data.Dirs, folderName)
		} else {
			data.Files = append(data.Files, f.Name)
		}
		return nil
	})
	return data, err
}

func GetLog(ref NamedReference, nextRef NamedReference) ([]*object.Commit, error) {
	if ref.Invalid() {
		return nil, ErrRefNotFound
	}
	g := ref.Repository
	cIter, err := g.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, ErrRefNotFound
	}
	data := []*object.Commit{}
	startHash := nextRef.Hash()
	for {
		c, err := cIter.Next()
		if err != nil {
			break
		}
		if !startHash.IsZero() {
			if startHash != c.Hash {
				continue
			}
			startHash = plumbing.ZeroHash
		}
		data = append(data, c)
		if len(data) >= 20 {
			break
		}
	}
	if !startHash.IsZero() {
		return nil, ErrRefNotFound
	}
	return data, nil
}

type SummaryData struct {
	NumCommits      int
	NumBranches     int
	NumTags         int
	NumFiles        int
	NumContributors int
}

func GetSummary(ref NamedReference) (SummaryData, error) {
	data := SummaryData{}
	if ref.Invalid() {
		return data, ErrRefNotFound
	}
	g := ref.Repository
	cIter, err := g.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return data, err
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
		return data, err
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
		return data, err
	}
	tree.Files().ForEach(func(f *object.File) error {
		data.NumFiles++
		return nil
	})
	return data, nil
}

func GetDiff(ref NamedReference) (object.FileStats, error) {
	if ref.Invalid() {
		return nil, ErrRefNotFound
	}
	commit := ref.Commit
	stats, err := commit.Stats()
	if err != nil {
		return nil, err
	}
	return stats, nil
}

package utils

import (
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetLatestCommit(repository *git.Repository, branchName string) (*object.Commit, error) {
	b, err := repository.Branch(branchName)
	if err != nil {
		return nil, err
	}

	commitHash, err := repository.ResolveRevision(plumbing.Revision(b.Name))
	if err != nil {
		return nil, err
	}

	commit, err := repository.CommitObject(*commitHash)
	if err != nil {
		return nil, err
	}

	return commit, nil
}

func GetContentsOfFile(repository *git.Repository, hash plumbing.Hash) ([]byte, error) {
	b, err := repository.BlobObject(hash)
	if err != nil {
		return nil, err
	}

	r, err := b.Reader()
	if err != nil {
		return nil, err
	}

	contents, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

package repository

import (
	"errors"

	"github.com/go-git/go-git/v5"
)

func (r *Repository) PlainOpen(path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return errors.New("repository does not exist")
	}
	r.repo = repo

	return nil
}

package repository

import "github.com/go-git/go-git/v5"

type Repository struct {
	repo          *git.Repository
	submodulePath string
}

// New creates a new Repository
func New() *Repository {
	return &Repository{
		repo:          nil,
		submodulePath: "versionctrls-integration",
	}
}

package repository

import (
	"errors"

	"github.com/go-git/go-git/v5/config"
)

// GetGitUserInfo retrieves the name and email of the Git user
func (r Repository) GetGitUserInfo() (string, string, error) {
	if r.repo == nil {
		return "", "", errors.New("no repository opened")
	}

	cfg, err := r.repo.ConfigScoped(config.GlobalScope)
	if err != nil {
		return "", "", err
	}

	name := cfg.User.Name
	email := cfg.User.Email

	return name, email, nil
}

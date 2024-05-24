package gitmanagement

import (
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// Clone a repository using go-git
func cloneRepo(url, path, username, password string) {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: username,
			Password: password,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Repository cloned successfully using go-git.")
}

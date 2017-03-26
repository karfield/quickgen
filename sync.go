package main

import (
	"os"

	git "github.com/karfield/go-git"
	"github.com/karfield/go-git/plumbing"
)

func SyncTemplates() error {
	templateDir := templateDir()
	if !fileExist(templateDir) {
		if err := os.Mkdir(templateDir, os.FileMode(0755)); err != nil {
			return err
		}
	}

	repo, err := git.PlainOpen(templateDir)
	if err != nil {
		repo, err = git.PlainClone(templateDir, false, &git.CloneOptions{
			URL:           "https://github.com/karfield/quickgen.git",
			ReferenceName: plumbing.ReferenceName("refs/heads/templates"),
			SingleBranch:  true,
			Depth:         1,
		})
	} else {
		err = repo.Pull(&git.PullOptions{
			ReferenceName: plumbing.ReferenceName("refs/heads/templates"),
			SingleBranch:  true,
		})
		if err == git.NoErrAlreadyUpToDate {
			err = nil
		}
	}

	//fmt.Printf("templates has been updated\n")
	return err
}

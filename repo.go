package adot

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"time"
)

func (a *ADot) openRepo() (*git.Repository, error) {
	r, err := git.PlainOpen(a.GitPath)
	if err != nil {
		return nil, errors.Wrapf(err, "pull %v", a.GitPath)
	}

	return r, nil
}

func (a *ADot) worktree() (*git.Repository, *git.Worktree, error) {
	r, err := a.openRepo()
	if err != nil {
		return nil, nil, errors.Wrap(err, "open repo")
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, nil, errors.Wrap(err, "worktree")
	}

	return r, w, nil
}

func (a *ADot) clone() error {
	_, err := git.PlainClone(a.GitPath, false, &git.CloneOptions{
		URL: a.GitURL,
	})
	if err != nil {
		return errors.Wrapf(err, "cloning %v to %v", a.GitURL, a.GitPath)
	}

	return nil
}

func (a *ADot) commit(p string) error {
	hash, err := a.work.Add(p)
	if err != nil {
		return errors.Wrapf(err, "could not add to worktree")
	}
	fmt.Println(hash)

	hash, err = a.work.Commit("Added " + p, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "ADot",
			Email: "adot@slowchop.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return errors.Wrapf(err, "could not commit worktree")
	}
	fmt.Println(hash)

	return nil
}

func (a *ADot) track() error {
	repo, work, err := a.worktree()
	if err != nil {
		return err
	}
	a.repo = repo
	a.work = work

	return nil
}

func (a *ADot) pull() error {
	return a.work.Pull(&git.PullOptions{
		RemoteName:    a.GitRemote,
		ReferenceName: plumbing.NewBranchReferenceName(a.GitBranch),
	})
}


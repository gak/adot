package adot

import (
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
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
	panic("no commit")
	return nil
}


package adot

import (
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
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


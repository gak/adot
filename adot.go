package adot

import (
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"path/filepath"
)

type ADot struct {
	WorkPath   string
	ConfigPath string

	// The URL of the git repo that stores the tracked dotfiles.
	GitURL string

	// The repo remote to push/pull.
	GitRemote string

	// The branch of the repo to track.
	GitBranch string

	// The check out path. Prepare to ~/.adot-git
	GitPath string

	repo *git.Repository
	work *git.Worktree

	files []string
}

func (a *ADot) Service() error {
	return nil
}

func (a *ADot) InitNew(url string) error {
	if err := a.Prepare(); err != nil {
		return errors.Wrap(err, "could not prepare configuration")
	}

	a.GitURL = url
	repo, err := git.PlainInit(a.GitPath, false)
	if err != nil {
		return errors.Wrap(err, "could not initiate repo")
	}
	a.repo = repo

	return nil
}

func (a *ADot) InitExisting(url string) error {
	if err := a.Prepare(); err != nil {
		return errors.Wrap(err, "could not prepare configuration")
	}

	a.GitURL = url

	if err := a.clone(); err != nil {
		return errors.Wrap(err, "clone")
	}
	if err := a.track(); err != nil {
		return errors.Wrap(err, "could not open repo")
	}
	if err := a.down(true); err != nil {
		return errors.Wrap(err, "load")
	}

	return nil
}

func (a *ADot) Prepare() error {
	var err error
	a.WorkPath, err = os.Getwd()
	if err != nil {
		return errors.Wrap(err, "could not determine cwd")
	}

	if a.ConfigPath == "" {
		a.ConfigPath = filepath.Join(a.WorkPath, ".adot")
	}
	a.ConfigPath = expand(a.ConfigPath)

	if a.GitPath == "" {
		a.GitPath = filepath.Join(filepath.Dir(a.ConfigPath), ".adot-git")
	}

	if a.GitBranch == "" {
		a.GitBranch = "master"
	}

	return nil
}

func (a *ADot) Add(p string) error {
	if err := a.upFile(p); err != nil {
		return errors.Wrapf(err, "could not up file %s", p)
	}
	if err := a.commit(p); err != nil {
		return errors.Wrapf(err, "could not commit %s", p)
	}
	if err := a.Push(); err != nil {
		return errors.Wrapf(err, "could not push %s", p)
	}
	return nil
}

func (a *ADot) Remove(p string) error {
	panic("no remove")
}

func (a *ADot) MonitorFile(path string) error {
	a.files = append(a.files, path)

	return nil
}

func (a *ADot) Push() error {
	panic("no push")
	return nil
}

func (a *ADot) Pull() error {
	panic("no pul")
	return nil
}

func (a *ADot) monitorFile(path string) error {
	a.files = append(a.files, path)

	return nil
}

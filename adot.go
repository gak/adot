package adot

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	"path/filepath"
	"strings"
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
	return nil
}

func (a *ADot) Pull() error {
	return a.work.Pull(&git.PullOptions{
		RemoteName:    a.GitRemote,
		ReferenceName: plumbing.NewBranchReferenceName(a.GitBranch),
	})
}

func (a *ADot) Push() error {
	panic("no push")
	return nil
}

// Save copies files from the home directory to the repository.
func (a *ADot) Save() error {
	return nil
}

func (a *ADot) Monitor() error {
	fmt.Printf("Monitoring %d files in %v based on %v", len(a.files), a.WorkPath, a.ConfigPath)

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

func (a *ADot) iterate(fun func(string) error) error {
	fp, err := os.Open(a.ConfigPath)
	if err != nil {
		return errors.Wrapf(err, "Could not open %v", a.ConfigPath)
	}

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		file := scanner.Text()
		file = strings.TrimSpace(file)
		if file == "" {
			continue
		}
		err := fun(file)
		if err != nil {
			return errors.Wrap(err, file)
		}
	}

	if err := scanner.Err(); err != nil {
		return errors.Wrapf(err, "scanner error")
	}

	return nil
}

func (a *ADot) monitorFile(path string) error {
	a.files = append(a.files, path)

	return nil
}

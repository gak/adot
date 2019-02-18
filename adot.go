package adot

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"io"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

type ADot struct {
	ScanPath   string
	ConfigPath string

	// The URL of the git repo that stores the tracked dotfiles.
	GitURL string

	// The repo remote to push/pull.
	GitRemote string

	// The branch of the repo to track. Prepare to master.
	GitBranch string

	// The check out path. Prepare to ~/.adot-git
	GitPath string

	repo *git.Repository
	work *git.Worktree

	files []string
}

func (a *ADot) Service() error {

	/*

	adot init

	 - clone
	 - if file exists, rename
	 - copy file down

	adot push

	 - check repo state to be clean
	 - iterate
	   - if file is different to repo, copy to repo
	 - if there are changes, commit, push
	 - if there's a conflict do a pull rebase
	 - otherwise let the user handle the conflict

	adot pull

	 - check repo state to be clean
	 - git pull
	 - iterate
	   - if any files are different, backup/copy.

	adot service

	 - every 5 minutes
	 - adot push
	 - adot pull

	*/

	//a.Monitor()
	//a.EnsureClone()

	return nil
}

func (a *ADot) Init(url string) error {
	a.GitURL = url

	if err := a.Clone(); err != nil {
		return errors.Wrap(err, "clone")
	}

	if err := a.Track(); err != nil {
		return errors.Wrap(err, "could not open repo")
	}

	if err := a.Load(true); err != nil {
		return errors.Wrap(err, "pull")
	}

	return nil
}

func (a *ADot) Prepare() error {
	if a.ScanPath == "" {
		a.ScanPath = "~"
	}
	a.ScanPath = expand(a.ScanPath)

	if a.ConfigPath == "" {
		a.ConfigPath = filepath.Join(a.ScanPath, ".adot")
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

func (a *ADot) Track() error {
	repo, work, err := a.Worktree()
	if err != nil {
		return err
	}
	a.repo = repo
	a.work = work

	return nil
}

func (a *ADot) Iterate(fun func(string) error) error {
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

func (a *ADot) MonitorFile(path string) error {
	a.files = append(a.files, path)

	return nil
}

// https://stackoverflow.com/a/12518877
func fileExists(s string) (bool, error) {
	if _, err := os.Stat(s); err == nil {
		return true, nil

	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		return false, nil

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		return false, err
	}
}

func dirExists(s string) (bool, error) {
	if stat, err := os.Stat(s); err == nil {
		if stat.IsDir() {
			return true, nil
		} else {
			return false, errors.New(fmt.Sprintf("path is not a directory %v", s))
		}
	}

	return false, nil
}

func copy(src, dst string) error {
	fmt.Println("copy", src, dst)
	return nil

	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}

	return nil
}

func (a *ADot) OpenRepo() (*git.Repository, error) {
	r, err := git.PlainOpen(a.GitPath)
	if err != nil {
		return nil, errors.Wrapf(err, "pull %v", a.GitPath)
	}

	return r, nil
}

func (a *ADot) Worktree() (*git.Repository, *git.Worktree, error) {
	r, err := a.OpenRepo()
	if err != nil {
		return nil, nil, errors.Wrap(err, "open repo")
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, nil, errors.Wrap(err, "worktree")
	}

	return r, w, nil
}

func (a *ADot) Clone() error {
	_, err := git.PlainClone(a.GitPath, false, &git.CloneOptions{
		URL: a.GitURL,
	})
	if err != nil {
		return errors.Wrapf(err, "cloning %v to %v", a.GitURL, a.GitPath)
	}

	return nil
}

func (a *ADot) Pull() error {
	return a.work.Pull(&git.PullOptions{
		RemoteName:    a.GitRemote,
		ReferenceName: plumbing.NewBranchReferenceName(a.GitBranch),
	})
}

func (a *ADot) Commit() error {
	return nil
}

func (a *ADot) Push() error {
	return nil
}

// Load copies files from the repository to the home directory.
func (a *ADot) Load(backup bool) error {
	// Bootstrap the inclusive files to know what to look for.
	if err := a.LoadFile(".adot", backup); err != nil {
		return err
	}

	err := a.Iterate(func(s string) error {
		return a.LoadFile(s, backup)
	})
	if err != nil {
		return errors.Wrap(err, "iterate")
	}

	return nil
}

func (a *ADot) LoadFile(p string, backup bool) error {
	src := path.Join(a.GitPath, p)
	dst := path.Join(expand("~"), p)

	hasSrc, err := fileExists(src)
	if err != nil {
		return errors.Wrap(err, "file exists check")
	}

	if !hasSrc {
		fmt.Printf("%v doesn't exist in repo, skipping.\n", p)
		return nil
	}

	hasDst, err := fileExists(dst)
	if hasDst && backup {
		head, err := a.repo.Head()
		if err != nil {
			return errors.Wrapf(err, "head failed")
		}
		hash := head.Hash().String()

		err = copy(dst, filepath.Join(dst+".adot."+hash))
		if err != nil {
			return errors.Wrapf(err, "creating backup of %v", dst)
		}
	}

	return copy(src, dst)
}

// Save copies files from the home directory to the repository.
func (a *ADot) Save() error {
	return nil
}

func (a *ADot) Monitor() error {
	fmt.Printf("Monitoring %d files in %v based on %v", len(a.files), a.ScanPath, a.ConfigPath)

	return nil
}

// https://stackoverflow.com/a/17617721
func expand(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	if path == "~" {
		// In case of "~", which won't be caught by the "else if"
		path = dir
	} else if strings.HasPrefix(path, "~/") {
		// Use strings.HasPrefix so we don't match paths like
		// "/something/~/something/"
		path = filepath.Join(dir, path[2:])
	}

	return path
}

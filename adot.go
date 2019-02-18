package adot

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
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
	GitRepo    string

	// The branch of the repo to track. Defaults to master.
	GitBranch  string

	// The check out path. Defaults to ~/.adot-git
	GitPath    string

	files []string
}

func (a *ADot) Run() error {
	if err := a.Init(); err != nil {
		return err
	}

	a.Monitor()
	a.EnsureClone()

	return nil
}

func (a *ADot) Init() error {
	if a.ScanPath == "" {
		a.ScanPath = "~"
	}
	a.ScanPath = expand(a.ScanPath)

	if a.ConfigPath == "" {
		a.ConfigPath = path.Join(a.ScanPath, ".adot")
	}
	a.ConfigPath = expand(a.ConfigPath)

	if a.GitPath == "" {
		a.GitPath = path.Join(filepath.Base(a.ConfigPath), ".adot-git")
	}

	if a.GitBranch == "" {
		a.GitBranch = "master"
	}

	fp, err := os.Open(a.ConfigPath)
	if err != nil {
		return errors.Wrapf(err, "Could not open %v", a.ConfigPath)
	}

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		err := a.MonitorFile(scanner.Text())
		if err != nil {
			return errors.Wrapf(err, "adding file %v", err)
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

func (a *ADot) EnsureClone() error {

}

func (a *ADot) Clone() error {

}


func (a *ADot) Pull() error {
	a.git("pull")
}

func (a *ADot) Commit() error {

}

func (a *ADot) Push() error {

}

// Load copies files from the repository to the home directory.
func (a *ADot) Load() error {

}

// Save copies files from the home directory to the repository.
func (a *ADot) Save() error {

}

func (a *ADot) Monitor() error {
	fmt.Printf("Monitoring %d files in %v based on %v", len(a.files), a.ScanPath, a.ConfigPath)
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

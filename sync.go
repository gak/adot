package adot

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

func (a *ADot) configAppend(path string) error {
	fp, err := os.OpenFile(a.ConfigPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrapf(err, "could not open config file")
	}

	_, err = fmt.Fprintf(fp, path+"\n")
	if err != nil {
		return errors.Wrapf(err, "could not write to file %v", a.ConfigPath)
	}

	return nil
}

func (a *ADot) upFile(p string) error {
	src := filepath.Join(a.WorkPath, p)
	dst := filepath.Join(a.GitPath, p)

	if err := copy(src, dst); err != nil {
		return errors.Wrapf(err, "could not copy %s to %s", src, dst)
	}

	return nil
}

func (a *ADot) downFile(p string, backup bool) error {
	src := filepath.Join(a.GitPath, p)
	dst := filepath.Join(a.WorkPath, p)

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

// load copies files from the repository to the home directory.
func (a *ADot) down(backup bool) error {
	// Bootstrap the inclusive files to know what to look for.
	if err := a.downFile(".adot", backup); err != nil {
		return err
	}

	err := a.iterate(func(s string) error {
		return a.downFile(s, backup)
	})
	if err != nil {
		return errors.Wrap(err, "iterate")
	}

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

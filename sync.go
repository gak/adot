package adot

import (
	"fmt"
	"github.com/pkg/errors"
	"path"
	"path/filepath"
)

func (a *ADot) upFile(p string) error {
	src := path.Join(a.WorkPath, p)
	dst := path.Join(a.GitPath, p)

	if err := copy(src, dst); err != nil {
		return errors.Wrapf(err, "could not copy %s to %s", src, dst)
	}

	return nil
}

func (a *ADot) downFile(p string, backup bool) error {
	src := path.Join(a.GitPath, p)
	dst := path.Join(a.WorkPath, p)

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


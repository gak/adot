package adot

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

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

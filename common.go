package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

func templateDir() string {
	u, _ := user.Current()
	return filepath.Join(u.HomeDir, ".quickgen")
}

func fileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func isDirEmpty(dir string) bool {
	stat, err := os.Stat(dir)
	if os.ErrNotExist == err {
		return true
	} else if err != nil {
		return false
	}
	if stat.IsDir() {
		dirs, err := ioutil.ReadDir(dir)
		if err != nil {
			return false
		}
		return len(dirs) == 0
	}
	return false
}

func isWorkDirEmpty() bool {
	cwd, err := os.Getwd()
	if err != nil {
		return false
	}
	return isDirEmpty(cwd)
}

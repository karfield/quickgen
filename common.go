package main

import (
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

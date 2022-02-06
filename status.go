package main

import (
	"os"
	"strings"
	"regexp"
)

func home() string {
	home, _ := os.UserHomeDir()
	return home + "/"
}

func minifyDir(name string) string {
	r, _ := regexp.Compile("^(\\W*\\w)")
	if match := r.FindString(name); match != "" {
		return match
	}
	return name
}

func minifyPath(path string) string {
	if path == "" {
		return path
	}
	if home := home(); strings.HasPrefix(path, home) {
		path = "~/" + path[len(home):len(path)]
	}
	dirs := strings.Split(path, "/")
	for i, d := range dirs[:len(dirs)-1] {
		dirs[i] = minifyDir(d)
	}
	return "\033[94m" + strings.Join(dirs, "/") + "\033[m"
}

func applyVCS(path string, vcs VCS) string {
	root := vcs.RootDir()
	common := path[0:len(root)]
	remainder := path[len(root):len(path)]
	return minifyPath(common) + vcs.Stats() + minifyPath(remainder)
}

func Statusline() string {
	path, _ := os.Getwd()
	return applyVCS(path, git{})
}

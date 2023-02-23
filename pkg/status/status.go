package status

import (
	"os"
	"strings"
	"regexp"
	"github.com/kevna/statusline/pkg/git"
	"io/ioutil"
	"errors"
)

func home() string {
	home, _ := os.UserHomeDir()
	return home
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func minifyPath(path string) string {
	if path == "" {
		return path
	}
	dirs := strings.Split(path, "/")
	home := home()
	keep := 1
	has_git := false
	parent := path
	if _, err := os.Stat(path+"/.git"); !errors.Is(err, os.ErrNotExist) {
		has_git = true
	}
	for i := len(dirs)-1; i > 0; i-- {
		path = parent
		if path == home {
			dirs = dirs[i:]
			dirs[0] = "~"
			break
		}
		if has_git {
			has_git = false
			if repo, err := git.RepoBuilder(path); err == nil {
				keep = 0
				if strings.HasSuffix("/"+dirs[i], "/"+repo.Branch()) {
					keep++
				}
				dirs[i] += repo.Stats() + "\033[94m"
				continue
			}
		}
		shorten := 1
		parent = strings.Join(dirs[:i], "/")
		dir := dirs[i]
		dirL := len(dir)
		list, _ := ioutil.ReadDir(parent)
		for _, f := range list {
			name := f.Name()
			if name == ".git" {
				has_git = true
			}
			if dir == name {
				continue
			}
			limit := min(dirL, len(name))
			for j := shorten; j <= limit; j++ {
				if dir[:j] != name[:j] {
					if j > shorten {
						shorten = j
					}
					break
				}
			}
		}
		if keep > 0 {
			keep--
			continue
		}
		dirs[i] = dir[:shorten]
	}
	return "\033[94m" + strings.Join(dirs, "/") + "\033[m"
}

func Statusline() string {
	path, _ := os.Getwd()
	return minifyPath(path)
}

package git

import (
	"strings"
	"fmt"
	"strconv"
	"os"
	"errors"
)

const icon = "\033[38;5;202m\uE0A0\033[m"

type ab struct {
	set bool
	ahead int
	behind int
}

func (ab ab) String() string {
	if !ab.set {
		return AnsiColour("91", "1") + "↯" + AnsiColour("0")
	}
	result := ""
	if ab.ahead > 0 {
		result += AnsiColour("32") + fmt.Sprintf("↑%d", ab.ahead)
	}
	if ab.behind > 0 {
		result += AnsiColour("31") + fmt.Sprintf("↓%d", ab.behind)
	}
	if result != "" {
		result += AnsiColour("0")
	}
	return result
}

type status struct {
	unmerged int
	staged int
	unstaged int
	untracked int
}

func (s status) Bool() bool {
	return s.unmerged+s.staged+s.unstaged+s.untracked > 0
}

func (s status) String() string {
	var result []string
	if s.unmerged > 0 {
		result = append(result, AnsiColour("91;1"), fmt.Sprintf("%d", s.unmerged))
	}
	if s.staged > 0 {
		result = append(result, AnsiColour("32"), fmt.Sprintf("%d", s.staged))
	}
	if s.unstaged > 0 {
		result = append(result, AnsiColour("31"), fmt.Sprintf("%d", s.unstaged))
	}
	if s.untracked > 0 {
		result = append(result, AnsiColour("90"), fmt.Sprintf("%d", s.untracked))
	}
	if len(result) > 0 {
		result = append(result, AnsiColour("0"))
	}
	return strings.Join(result, "")
}

type VCS interface {
	RootDir(path string) string
	Branch() string
	Stats() string
}

type repo struct {
	branch string
	ab ab
	status status
	stashes int
}

func (r repo) Bool() bool {
	return r.branch != ""
}

func (r repo) RootDir(path string) string {
	// str, _ := runCommand("rev-parse", "--show-toplevel")
	// return str
	dirs := strings.Split(path, "/")
	for {
		path = strings.Join(dirs, "/") + "/.git"
		if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
			return strings.Join(dirs, "/")
		}
		length := len(dirs)-1
		dirs = dirs[:length]
		if length <= 2 { // Note: this cutoff saves us checking top-level directories, which are unlikely to be repos
			return ""
		}
	}
}

func (r repo) Branch() string {
	return r.branch
}

func (r repo) Stats() string {
	result := []string {icon, r.branch}
	result = append(result, fmt.Sprintf("%s", r.ab))
	if r.status.Bool() {
		result = append(result, fmt.Sprintf("(%s)", r.status))
	}
	if r.stashes > 0 {
		result = append(result, fmt.Sprintf("{%d}", r.stashes))
	}
	return strings.Join(result, "")
}

func repoStringBuilder(str string) repo {
	branch := ""
	ab := ab{}
	status := status{}
	stashes := 0
	for _, line := range strings.Split(str, "\n") {
		if line == "" {
			continue;
		}
		switch line[0] {
			case '#':
			// if strings.HasPrefix(line, "?") {
			// }
			fields := strings.Split(line, " ")
			switch fields[1] {
				case "branch.head":
				branch = fields[2]

				case "branch.ab":
				ab.set = true
				ab.ahead, _ = strconv.Atoi(fields[2])
				ab.behind, _ = strconv.Atoi(fields[3][1:])

				case "stash":
				stashes, _ = strconv.Atoi(fields[2])
			}

			case 'u':
			status.unmerged++

			case '1', '2':
			if line[2] != '.' {
				status.staged++
			}
			if line[3] != '.' {
				status.unstaged++
			}

			case '?':
			status.untracked++
		}
	}
	return repo{branch, ab, status, stashes}
}

func RepoBuilder() (repo, error) {
	str, err := runCommand("status", "--porcelain=v2", "--branch", "--show-stash");
	if err != nil {
		return repo{}, err
	}
	return repoStringBuilder(str), nil
}

package git

import (
	"strings"
	"fmt"
)

const icon = "\033[38;5;202m\uE0A0\033[m"

type ab struct {
	ahead int
	behind int
}

func (ab ab) String() string {
	ahead := ab.ahead > 0
	behind := ab.ahead > 0
	if ahead && behind {
		return fmt.Sprintf("\033[30;41m↕%d\033[m", ab.ahead+ab.behind)
	}
	if ahead {
		return fmt.Sprintf("↑%d", ab.ahead)
	}
	if behind {
		return fmt.Sprintf("↓%d", ab.behind)
	}
	return ""
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
		result = append(result, fmt.Sprintf("\033[91;1m%d", s.unmerged))
	}
	if s.staged > 0 {
		result = append(result, fmt.Sprintf("\033[32m%d", s.staged))
	}
	if s.unstaged > 0 {
		result = append(result, fmt.Sprintf("\033[31m%d", s.unstaged))
	}
	if s.untracked > 0 {
		result = append(result, fmt.Sprintf("\033[90m%d", s.untracked))
	}
	if len(result) > 0 {
		result = append(result, "\033[m")
	}
	return strings.Join(result, "")
}

type VCS interface {
	RootDir() string
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

func (r repo) RootDir() string {
	str, _ := runCommand("rev-parse", "--show-toplevel")
	return str
}

func (r repo) Branch() string {
	return r.branch
}

func (r repo) Stats() string {
	result := []string {icon, r.branch}
	// if ab, err := g.AheadBehind(); err != nil {
	// 	result = append(result, "\033[91m↯\033[m")
	// } else {
	// 	result = append(result, fmt.Sprintf("%s", ab))
	// }
	result = append(result, fmt.Sprintf("%s", r.ab))
	if r.status.Bool() {
		result = append(result, fmt.Sprintf("(%s)", r.status))
	}
	if r.stashes > 0 {
		result = append(result, fmt.Sprintf("{%d}", r.stashes))
	}
	return strings.Join(result, "")
}

func RepoBuilder() repo {
	str, _ := runCommand("status", "--porcelain=v2", "--branch", "--show-stash");
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

				// case "branch.ab":
				// ab.ahead = fields[2]
				// ab.behind = fields[3]

				// case "stash":
				// stashes = fields[2]
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


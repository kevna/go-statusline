package main

import (
	"strings"
	"fmt"
)

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
	staged int
	unstaged int
	untracked int
}

func (s status) Bool() bool {
	return s.staged+s.unstaged+s.untracked > 0
}

func (s status) String() string {
	var result []string
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

const icon = "\033[38;5;202m\uE0A0\033[m"

type VCS interface {
	RootDir() string
	Branch() string
	Stats() string
}

type git struct {}

func (g git) RootDir() string {
	str, _ := runCommand("rev-parse", "--show-toplevel")
	return str
}

func (g git) Branch() string {
	str, _ := runCommand("rev-parse", "--symbolic-full-name", "--abbrev-ref", "HEAD")
	return str
}

func (g git) AheadBehind() (ab, error) {
	ahead, err := count("rev-list", "@{push}..HEAD")
	if err != nil {
		return ab{}, err
	}
	behind, err := count("rev-list", "HEAD..@{upstream}")
	if err != nil {
		return ab{}, err
	}
	return ab{
		ahead: ahead,
		behind: behind,
	}, nil
}

func (g git) Status() status {
	str, _ := runCommand("status", "--porcelain")
	result := status{}
	for _, line := range strings.Split(str, "\n") {
		if line == "" {
			continue;
		}
		if strings.HasPrefix(line, "??") {
			result.untracked++
		} else {
			if line[0] != ' ' {
				result.staged++
			}
			if line[1] != ' ' {
				result.unstaged++
			}
		}
	}
	return result
}

func (g git) Stashes() int {
	count, _ := count("stash", "list")
	return count
}

func (g git) Stats() string {
	result := []string {icon}
	if branch := g.Branch(); !strings.HasSuffix(g.RootDir(), branch) {
		result = append(result, branch)
	}
	if ab, err := g.AheadBehind(); err != nil {
		result = append(result, "\033[91m↯\033[m")
	} else {
		result = append(result, fmt.Sprintf("%s", ab))
	}
	if status := g.Status(); status.Bool() {
		result = append(result, fmt.Sprintf("(%s)", status))
	}
	if stashes := g.Stashes(); stashes > 0 {
		result = append(result, fmt.Sprintf("{%d}", stashes))
	}
	return strings.Join(result, "")
}

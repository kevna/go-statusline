package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMinifyDir(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{"~", "~"},
		{"~root", "~r"},
		{"private_dot_config", "p"},
		{"._shares", "._"},
	}
	for _, test := range tests {
		actual := minifyDir(test.input)
		assert.Equal(t, test.expected, actual)
	}
}

func TestMinifyPath(t* testing.T) {
	tests := []struct {
		path string
		keep int
		expected string
	}{
		{"~", 1, "\033[94m~\033[m"},
		{"/etc/X11/xorg.conf.d", 1, "\033[94m/e/X/xorg.conf.d\033[m"},
		{"~/.local/share/chezmoi/private_dot_config/i3", 1, "\033[94m~/.l/s/c/p/i3\033[m"},
		{"~/.local/share/chezmoi/private_dot_config/i3", 2, "\033[94m~/.l/s/c/private_dot_config/i3\033[m"},
	}
	for _, test := range tests {
		actual := minifyPath(test.path, test.keep)
		assert.Equal(t, test.expected, actual)
	}
}

type MockVCS struct {
	root, branch, stats string
}
func (m MockVCS) RootDir() string {
	return m.root
}
func (m MockVCS) Branch() string {
	return m.branch
}
func (m MockVCS) Stats() string {
	return m.stats
}

func TestApplyVCS(t* testing.T) {
	tests := []struct {
		root, branch, stats, input, expected string
	}{
		{
			"~/.local/share/chezmoi",
			"master",
			"\uE0A0master",
			"~/.local/share/chezmoi/private_dot_config/i3",
			"\033[94m~/.l/s/chezmoi\033[m\uE0A0master\033[94m/p/i3\033[m",
		},
		{
			"~/Documents/python/statusline/master",
			"master",
			"\uE0A0",
			"~/Documents/python/statusline/master/statusline",
			"\033[94m~/D/p/statusline/master\033[m\uE0A0\033[94m/statusline\033[m",
		},
		{
			"~/Documents/python/statusline-master",
			"master",
			"\uE0A0",
			"~/Documents/python/statusline-master/statusline",
			"\033[94m~/D/p/statusline-master\033[m\uE0A0\033[94m/statusline\033[m",
		},
		{
			"~/Documents/python/statusline/feature/newfeature",
			"feature/newfeature",
			"\uE0A0",
			"~/Documents/python/statusline/feature/newfeature/statusline",
			"\033[94m~/D/p/s/feature/newfeature\033[m\uE0A0\033[94m/statusline\033[m",
		},
	}
	for _, test := range tests {
		mock := MockVCS{
			root: test.root,
			branch: test.branch,
			stats: test.stats,
		}
		actual := applyVCS(test.input, mock)
		assert.Equal(t, test.expected, actual)
	}
}

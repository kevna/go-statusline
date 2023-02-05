package git

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAbString(t *testing.T) {
	tests := []struct {
		ab ab
		expected string
	}{
		{ab{}, "\033[91;1m↯\033[m"},
		{ab{true, 0, 0}, ""},
		{ab{true, 1, 0}, "\033[32m↑1\033[m"},
		{ab{true, 0, 1}, "\033[31m↓1\033[m"},
		// {ab{true, 1, 10}, "\033[30;41m↕11\033[m"},
		{ab{true, 1, 10}, "\033[32m↑1\033[31m↓10\033[m"},
	}
	for _, test := range tests {
		actual := test.ab.String()
		assert.Equal(t, test.expected, actual)
	}
}

func TestStatusString(t *testing.T) {
	tests := []struct {
		status status
		expected string
	}{
		{status{}, ""},
		{status{1, 0, 0, 0}, "\033[91;1m1\033[m"},
		{status{0, 1, 0, 0}, "\033[32m1\033[m"},
		{status{0, 0, 1, 0}, "\033[31m1\033[m"},
		{status{0, 0, 0, 1}, "\033[90m1\033[m"},
		{status{1, 2, 4, 5}, "\033[91;1m1\033[32m2\033[31m4\033[90m5\033[m"},
	}
	for _, test := range tests {
		actual := test.status.String()
		assert.Equal(t, test.expected, actual)
	}
}

func TestRepoStats(t *testing.T) {
	tests := []struct {
		repo repo
		expected string
	}{
		{repo{}, icon+"\033[91;1m↯\033[m"},
		{repo{"master", ab{}, status{}, 0}, icon+"master\033[91;1m↯\033[m"},
		{repo{"master", ab{true, 0, 0}, status{}, 0}, icon+"master"},
		{
			repo{"master", ab{}, status{1, 1, 1, 1}, 0},
			icon+"master\x1b[91;1m↯\x1b[m(\x1b[91;1m1\x1b[32m1\x1b[31m1\x1b[90m1\x1b[m)",
		},
		{
			repo{"master", ab{true, 1, 10}, status{1, 1, 4, 5}, 3},
			icon+"master\x1b[32m↑1\x1b[31m↓10\x1b[m(\x1b[91;1m1\x1b[32m1\x1b[31m4\x1b[90m5\x1b[m){3}",
		},
	}
	for _, test := range tests {
		actual := test.repo.Stats()
		assert.Equal(t, test.expected, actual)
	}
}


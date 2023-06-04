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
		{ab{}, "\001\033[91;1m\002↯\001\033[0m\002"},
		{ab{true, 0, 0}, ""},
		{ab{true, 1, 0}, "\001\033[32m\002↑1\001\033[0m\002"},
		{ab{true, 0, 1}, "\001\033[31m\002↓1\001\033[0m\002"},
		// {ab{true, 1, 10}, "\001\033[30;41m\002↕11\001\033[0m\002"},
		{ab{true, 1, 10}, "\001\033[32m\002↑1\001\033[31m\002↓10\001\033[0m\002"},
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
		{status{1, 0, 0, 0}, "\001\033[91;1m\0021\001\033[0m\002"},
		{status{0, 1, 0, 0}, "\001\033[32m\0021\001\033[0m\002"},
		{status{0, 0, 1, 0}, "\001\033[31m\0021\001\033[0m\002"},
		{status{0, 0, 0, 1}, "\001\033[90m\0021\001\033[0m\002"},
		{status{1, 2, 4, 5}, "\001\033[91;1m\0021\001\033[32m\0022\001\033[31m\0024\001\033[90m\0025\001\033[0m\002"},
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
		{repo{}, icon+"\001\033[91;1m\002↯\001\033[0m\002"},
		{repo{"master", ab{}, status{}, 0}, icon+"master\001\033[91;1m\002↯\001\033[0m\002"},
		{repo{"master", ab{true, 0, 0}, status{}, 0}, icon+"master"},
		{
			repo{"master", ab{}, status{1, 1, 1, 1}, 0},
			icon+"master\x01\x1b[91;1m\x02↯\x01\x1b[0m\x02(\x01\x1b[91;1m\x021\x01\x1b[32m\x021\x01\x1b[31m\x021\x01\x1b[90m\x021\x01\x1b[0m\x02)",
		},
		{
			repo{"master", ab{true, 1, 10}, status{1, 1, 4, 5}, 3},
			icon+"master\x01\x1b[32m\x02↑1\x01\x1b[31m\x02↓10\x01\x1b[0m\x02(\x01\x1b[91;1m\x021\x01\x1b[32m\x021\x01\x1b[31m\x024\x01\x1b[90m\x025\x01\x1b[0m\x02){3}",
		},
	}
	for _, test := range tests {
		actual := test.repo.Stats()
		assert.Equal(t, test.expected, actual)
	}
}

func TestRepoStringBuilder(t *testing.T) {
	tests := []struct {
		status string
		expected repo
	}{
		{"", repo{}},
		{
			`
# branch.oid (initial)
# branch.head (detached)
1 MM N... 100644 100644 100644 3e2ceb914cf9be46bf235432781840f4145363fd 3e2ceb914cf9be46bf235432781840f4145363fd README.md
			`,
			repo{"(detached)", ab{}, status{0, 1, 1, 0}, 0},
		},
		{
			`
# branch.oid 51c9c58e2175b768137c1e38865f394c76a7d49d
# branch.head master
# branch.upstream origin/master
# branch.ab +1 -10
# stash 3
1 .M N... 100644 100644 100644 3e2ceb914cf9be46bf235432781840f4145363fd 3e2ceb914cf9be46bf235432781840f4145363fd Gopkg.lock
1 .M N... 100644 100644 100644 cecb683e6e626bcba909ddd36d3357d49f0cfd09 cecb683e6e626bcba909ddd36d3357d49f0cfd09 Gopkg.toml
1 .M N... 100644 100644 100644 aea984b7df090ce3a5826a854f3e5364cd8f2ccd aea984b7df090ce3a5826a854f3e5364cd8f2ccd porcelain.go
1 .D N... 100644 100644 000000 6d9532ba55b84ec4faf214f9cdb9ce70ec8f4f5b 6d9532ba55b84ec4faf214f9cdb9ce70ec8f4f5b porcelain_test.go
2 R. N... 100644 100644 100644 44d0a25072ee3706a8015bef72bdd2c4ab6da76d 44d0a25072ee3706a8015bef72bdd2c4ab6da76d R100 hm.rb     hw.rb
u UU N... 100644 100644 100644 100644 ac51efdc3df4f4fd328d1a02ad05331d8e2c9111 36c06c8752c78d2aff89571132f3bf7841a7b5c3 e85207e04dfdd5eb0a1e9febbc67fd837c44a1cd hw.rb
? _porcelain_test.go
? git.go
? git_test.go
? goreleaser.yml
? vendor/
			`,
			repo{"master", ab{true, 1, 10}, status{1, 1, 4, 5}, 3},
		},
	}
	for _, test := range tests {
		actual := repoStringBuilder(test.status)
		assert.Equal(t, test.expected, actual)
	}
}

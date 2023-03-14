package git

import (
	"bytes"
	"os/exec"
	"strings"
)

func runCommand(arg ...string) (string, error) {
	cmd := exec.Command("git", arg...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func count(arg ...string) (int, error) {
	str, err := runCommand(arg...)
	if err != nil {
		return 0, err
	}
	rows := strings.Split(str, "\n")
	return len(rows)-1, nil
}

func AnsiColour(escapes ...string) string {
	return "\033[" + strings.Join(escapes, ";") + "m"
}

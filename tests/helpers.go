package tests

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func RunCommand(command string) (string, error) {
	commandArgs := strings.Split(command, " ")
	cmd := exec.Command("wampproto", commandArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("command execution failed: %w, stderr: %s", err, stderr.String())
	}

	output := strings.TrimSpace(stdout.String())

	return output, nil
}

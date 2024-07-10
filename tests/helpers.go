package tests

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
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

func RunCommandAndDeserialize(t *testing.T, command string, serializer serializers.Serializer) messages.Message {
	output, err := RunCommand(command)
	require.NoError(t, err)
	outputBytes, err := hex.DecodeString(output)
	require.NoError(t, err)

	msg, err := serializer.Deserialize(outputBytes)
	require.NoError(t, err)

	return msg
}

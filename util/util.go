package util

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
)

func GetCommit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error getting commit: %w", err)
	}

	return strings.TrimSuffix(out.String(), "\n"), nil
}

func GetBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error getting branch: %w", err)
	}

	return strings.TrimSuffix(out.String(), "\n"), nil
}

func GetLatency(message twitch.PrivateMessage) string {
	latency := time.Since(message.Time)

	return latency.String()
}

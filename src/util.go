package main

import (
	"bytes"
	"os/exec"
	"strings"
)

func GetCommit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(out.String(), "\n"), nil
}

func GetBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(out.String(), "\n"), nil
}

package commands

import (
	"context"
	"fmt"
	"log"
	"os/exec"
)

func Restart(ctx context.Context) (string, error) {
	err := executeCmd("docker", "compose", "down")
	if err != nil {
		return "", fmt.Errorf("error taking compose down%w", err)
	}

	err = executeCmd("docker", "compose", "up", "-d")
	if err != nil {
		return "", fmt.Errorf("error bringing compose up: %w", err)
	}

	return "Restarting...", nil
}

type LogWriter struct {
	logger *log.Logger
}

func NewLogWriter(logger *log.Logger) *LogWriter {
	lw := &LogWriter{}
	lw.logger = logger

	return lw
}

func (lw LogWriter) Write(p []byte) (int, error) {
	lw.logger.Print(string(p))

	return len(p), nil
}

func executeCmd(command string, args ...string) error {
	cmd := exec.Command(command, args...)

	cmd.Stdout = NewLogWriter(log.Default())
	cmd.Stderr = NewLogWriter(log.Default())

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("error startin cmd: %w", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("error executing cmd: %w", err)
	}

	return nil
}

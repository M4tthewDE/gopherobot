package cmd_test

import (
	"testing"

	"de.com.fdm/gopherobot/cmd"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/stretchr/testify/assert"
)

func TestHTTPStatusCommand(t *testing.T) {
	t.Parallel()

	codes := map[string]string{
		100, "Continue",
		"101", "Switching Protocols",
		"102", "Processing",
		"200", "OK",
		"201", "Created",
	}

	cmdHandler := cmd.CommandHandler{}
	message := twitch.PrivateMessage{
		Message: ";httpstatus 100",
	}

	result := cmdHandler.HTTPStatusCommand(message)

	assert.Equal(t, result, "Continue", "should be equals")

}

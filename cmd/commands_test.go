package cmd

import (
	"fmt"
	"testing"

	"de.com.fdm/gopherobot/provider"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/stretchr/testify/assert"
)

func TestEchoCommand(t *testing.T) {
	t.Parallel()

	cmdHandler := CommandHandler{}
	message := twitch.PrivateMessage{
		Message: ";echo test",
	}
	result := cmdHandler.EchoCommand(message)
	assert.Equal(t, "test", result)

	message.Message = ";echo "
	result = cmdHandler.EchoCommand(message)
	assert.Equal(t, "", result)

	message.Message = ";echo"
	result = cmdHandler.EchoCommand(message)
	assert.Equal(t, "", result)
}

func TestUserIDCommand(t *testing.T) {
	t.Parallel()

	cmdHandler := CommandHandler{
		twitchProvider: &provider.TestTwitchProvider{},
	}

	message := twitch.PrivateMessage{
		Message: ";id test",
	}

	result := cmdHandler.UserIDCommand(message)
	assert.Equal(t, "User-ID of test is 1337", result)

	message.Message = ";id"
	result = cmdHandler.UserIDCommand(message)
	assert.Equal(t, "No user provided", result)
}

func TestUserCommand(t *testing.T) {
	t.Parallel()

	cmdHandler := CommandHandler{
		twitchProvider: &provider.TestTwitchProvider{},
	}

	message := twitch.PrivateMessage{
		Message: ";user 1337",
	}

	result := cmdHandler.UserCommand(message)
	assert.Equal(t, "Username for 1337 is user", result)

	message.Message = ";user"
	result = cmdHandler.UserCommand(message)
	assert.Equal(t, "No ID provided", result)
}

func TestHTTPStatusCommand(t *testing.T) {
	t.Parallel()

	codes := map[string]string{
		"100": "Continue",
		"101": "Switching Protocols",
		"102": "Processing",
		"200": "OK",
		"201": "Created",
		"202": "Accepted",
		"203": "Non-Authoritative Information",
		"204": "No Content",
		"206": "Partial Content",
		"207": "Multi-Status",
		"300": "Multiple Choices",
		"301": "Moved Permanently",
		"302": "Found",
		"303": "See Other",
		"304": "Not Modified",
		"305": "Use Proxy",
		"307": "Temporary Redirect",
		"308": "Permanent Redirect",
		"400": "Bad Request",
		"401": "Unauthorized",
		"402": "Payment Required",
		"403": "Forbidden",
		"404": "Not Found",
		"405": "Method Not Allowed",
		"406": "Not Acceptable",
		"407": "Proxy Authentication Required",
		"408": "Request Timeout",
		"409": "Conflict",
		"410": "Gone",
		"411": "Length Required",
		"412": "Precondition Failed",
		"413": "Request Entity Too Large",
		"414": "Request URI Too Long",
		"415": "Unsupported Media Type",
		"416": "Requested Range Not Satisfiable",
		"417": "Expectation Failed",
		"418": "I'm a teapot",
		"421": "Misdirected Request",
		"422": "Unprocessable Entity",
		"423": "Locked",
		"424": "Failed Dependency",
		"425": "Too Early",
		"426": "Upgrade Required",
		"428": "Precondition Required",
		"429": "Too Many Requests",
		"431": "Request Header Fields Too Large",
		"451": "Unavailable For Legal Reasons",
		"500": "Internal Server Error",
		"501": "Not Implemented",
		"502": "Bad Gateway",
		"503": "Service Unavailable",
		"504": "Gateway Timeout",
		"505": "HTTP Version Not Supported",
		"506": "Variant Also Negotiates",
		"507": "Insufficient Storage",
		"508": "Loop Detected",
		"510": "Not Extended",
		"511": "Network Authentication Required",
	}

	cmdHandler := CommandHandler{}
	message := twitch.PrivateMessage{}

	for status, expected := range codes {
		message.Message = fmt.Sprintf(";httpstatus %s", status)
		result := cmdHandler.HTTPStatusCommand(message)
		assert.Equal(t, expected, result, "should be equals")
	}

	message.Message = ";httpstatus"
	result := cmdHandler.HTTPStatusCommand(message)
	assert.Equal(t, "No code provided", result)

	message.Message = ";httpstatus "
	result = cmdHandler.HTTPStatusCommand(message)
	assert.Equal(t, "No valid code provided", result)

	message.Message = ";httpstatus test"
	result = cmdHandler.HTTPStatusCommand(message)
	assert.Equal(t, "No valid code provided", result)
}

func TestNextLaunchCommand(t *testing.T) {
	t.Parallel()

	cmdHandler := CommandHandler{
		launchProvider: provider.TestLaunchProvider{},
	}
	message := twitch.PrivateMessage{
		Message: ";nextlaunch",
	}

	result := cmdHandler.NextLaunchCommand(message)
	assert.Equal(t, "2021-01-01 12:00:00 +0000 UTC (UTC) | Name: Test-Launch | Details: Test-Details", result)
}

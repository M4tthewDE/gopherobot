package commands

import (
	"testing"

	"de.com.fdm/gopherobot/providers"
)

func TestGetTargetEmoteCode(t *testing.T) {
	t.Parallel()

	emoteCode, err := getTargetEmoteCode(",improveemote OMEGALUL")
	if err != nil {
		t.Fatal(err)
	}

	if emoteCode != "OMEGALUL" {
		t.Fatalf("Expected OMEGALUL, got %s", emoteCode)
	}

	emoteCode, err = getTargetEmoteCode(",improveemote ")
	if err == nil {
		t.Fatalf("Expected error, got %s", emoteCode)
	}

	emoteCode, err = getTargetEmoteCode(",improveemote")
	if err == nil {
		t.Fatalf("Expected error, got %s", emoteCode)
	}
}

func TestImproveBttvEmote(t *testing.T) {
	t.Parallel()

	emoteCode := "60c8d8bef8b3f62601c3e32b"

	emoteBuffer, err := providers.GetBttvEmote(emoteCode)
	if err != nil {
		t.Fatal(err)
	}

	_, err = modifyEmote(emoteBuffer)
	if err != nil {
		t.Fatal(err)
	}
}

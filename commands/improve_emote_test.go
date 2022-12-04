package commands

import (
	"testing"
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

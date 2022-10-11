package commands

import "testing"

func TestGetTargetEmoteCode(t *testing.T) {
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
	emoteCode := "60c8d8bef8b3f62601c3e32b"

	_, err := improveBttvEmote(emoteCode)
	if err != nil {
		t.Fatal(err)
	}
}

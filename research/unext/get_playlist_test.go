package unext

import (
	"net/http"
	"testing"
)

func TestSaveMPDURL(t *testing.T) {
	tokens, err := LoadTokens("tokens.json")
	if err != nil {
		t.Skipf("no saved tokens found: %v", err)
	}

	t.Logf("loaded access_token: %s...", tokens.AccessToken[:50])

	client := &http.Client{}

	playlist, err := GetPlaylistURL(client, tokens.AccessToken)
	if err != nil {
		t.Fatalf("GetPlaylistURL: %v", err)
	}

	if err := SaveMPDURL(playlist); err != nil {
		t.Fatalf("SaveMPDURL: %v", err)
	}

	t.Logf("MPD URL saved to %s", mpdURLFile)
}

func TestLoadAndPrintMPD(t *testing.T) {
	client := &http.Client{}

	if err := LoadAndPrintMPD(client); err != nil {
		t.Fatalf("LoadAndPrintMPD: %v", err)
	}
}

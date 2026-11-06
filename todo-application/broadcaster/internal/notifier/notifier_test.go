package notifier

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"broadcaster/internal/models"
)

func capture(t *testing.T, format string, status int) (map[string]any, error) {
	t.Helper()
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(status)
	}))
	defer srv.Close()

	n, err := New(format, srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	err = n.Send(context.Background(), models.StatusMessage{User: "bot", Message: "A todo was created"})
	return got, err
}

func TestGenericPayload(t *testing.T) {
	got, err := capture(t, "generic", http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	if got["user"] != "bot" || got["message"] != "A todo was created" {
		t.Fatalf("unexpected payload: %v", got)
	}
}

func TestDiscordPayload(t *testing.T) {
	got, err := capture(t, "discord", http.StatusNoContent) // Discord answers 204
	if err != nil {
		t.Fatal(err)
	}
	if got["username"] != "bot" || got["content"] != "A todo was created" {
		t.Fatalf("unexpected payload: %v", got)
	}
}

func TestNon2xxIsError(t *testing.T) {
	if _, err := capture(t, "generic", http.StatusInternalServerError); err == nil {
		t.Fatal("expected an error")
	}
}

func TestInvalidConfig(t *testing.T) {
	if _, err := New("slack", "http://x"); err == nil {
		t.Fatal("expected error for unknown format")
	}
	if _, err := New("generic", ""); err == nil {
		t.Fatal("expected error for empty URL")
	}
}

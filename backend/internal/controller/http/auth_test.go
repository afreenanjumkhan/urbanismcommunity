package http

import (
	"io"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAnonymousSessionCreationDenied(t *testing.T) {
	r := NewRouter(nil, slog.New(slog.NewTextHandler(io.Discard, nil)), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("POST", "/api/v1/sessions", strings.NewReader(`{}`)))
	if w.Code != 401 {
		t.Fatalf("got %d", w.Code)
	}
}
func TestAuthDisabledStatus(t *testing.T) {
	r := NewRouter(nil, slog.New(slog.NewTextHandler(io.Discard, nil)), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/api/v1/auth/status", nil))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "false") {
		t.Fatal(w.Body.String())
	}
}

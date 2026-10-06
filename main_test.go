package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerStatusAndGreeting(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Hello, DevOps!") {
		t.Errorf("body does not contain greeting: %q", rec.Body.String())
	}
}

func TestHandlerShowsVersion(t *testing.T) {
	version = "9.9.9"

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if !strings.Contains(rec.Body.String(), "version=9.9.9") {
		t.Errorf("body does not show injected version: %q", rec.Body.String())
	}
}

package traefik_inline_response_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/tuxgal/traefik_inline_response"
)

func TestHeaderAccess(t *testing.T) {
	t.Parallel()

	config := buildConfig(`
matchers:
  - path:
      abs: /check
    statusCode: 200
    response:
      template: '{{ .Header.Get "X-Forwarded-Host" }}'
`)

	ctx := context.Background()
	next := newNextHandler()

	handler, err := traefik_inline_response.New(ctx, next.handlerFunc(), config, "inline-response")
	if err != nil {
		t.Fatalf("failed to initialize handler: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/check", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("X-Forwarded-Host", "example.com")

	rec := newResponseRecorder()
	handler.ServeHTTP(rec, req)

	result := rec.Result()
	if result == nil {
		t.Fatal("expected a response but got none")
	}
	if result.StatusCode != http.StatusOK {
		t.Errorf("got status %d, want %d", result.StatusCode, http.StatusOK)
	}

	body, err := readBody(result.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	if body != "example.com" {
		t.Errorf("got body %q, want %q", body, "example.com")
	}
}

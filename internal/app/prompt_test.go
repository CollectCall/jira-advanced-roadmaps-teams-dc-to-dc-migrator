package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPromptForAuthPreservesProvidedCredentials(t *testing.T) {
	username := "target-user"
	password := "target-pass"

	if err := promptForAuth(nil, "target", &username, &password); err != nil {
		t.Fatalf("promptForAuth returned error: %v", err)
	}

	if username != "target-user" {
		t.Fatalf("username was changed to %q", username)
	}
	if password != "target-pass" {
		t.Fatalf("password was changed to %q", password)
	}
}

func TestVerifyJiraCredentialsUsesCurrentUserEndpoint(t *testing.T) {
	var sawAuth bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/rest/api/2/myself" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		username, password, ok := r.BasicAuth()
		sawAuth = ok && username == "admin" && password == "secret"
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"admin","key":"admin-key","displayName":"Admin User","active":true}`))
	}))
	defer server.Close()

	user, err := verifyJiraCredentials(server.URL, "admin", "secret")
	if err != nil {
		t.Fatalf("verifyJiraCredentials returned error: %v", err)
	}
	if !sawAuth {
		t.Fatal("expected basic auth credentials on current-user request")
	}
	if user.DisplayName != "Admin User" {
		t.Fatalf("unexpected verified user: %#v", user)
	}
}

func TestVerifyJiraCredentialsReturnsAuthFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/2/myself" {
			t.Fatalf("unexpected request path %s", r.URL.Path)
		}
		http.Error(w, "bad credentials", http.StatusUnauthorized)
	}))
	defer server.Close()

	_, err := verifyJiraCredentials(server.URL, "admin", "wrong")
	if err == nil {
		t.Fatal("expected verification error")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Fatalf("expected 401 in verification error, got %v", err)
	}
}

func TestVerifyJiraCredentialsReturnsDecodeFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer server.Close()

	_, err := verifyJiraCredentials(server.URL, "admin", "secret")
	if err == nil {
		t.Fatal("expected decode error")
	}
}

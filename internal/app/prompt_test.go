package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPromptForAuthPreservesProvidedCredentials(t *testing.T) {
	accountName := authFixtureValue("target", "account")
	authToken := authFixtureValue("target", "token")

	if err := promptForAuth(nil, "target", &accountName, &authToken); err != nil {
		t.Fatalf("promptForAuth returned error: %v", err)
	}

	if accountName != authFixtureValue("target", "account") {
		t.Fatalf("account name was changed to %q", accountName)
	}
	if authToken != authFixtureValue("target", "token") {
		t.Fatalf("auth token was changed to %q", authToken)
	}
}

func TestVerifyJiraCredentialsUsesCurrentUserEndpoint(t *testing.T) {
	accountName := authFixtureValue("jira", "account")
	authToken := authFixtureValue("jira", "token")
	var sawAuth bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/rest/api/2/myself" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		requestAccount, requestToken, ok := r.BasicAuth()
		sawAuth = ok && requestAccount == accountName && requestToken == authToken
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"test-account","key":"test-account-key","displayName":"Test Operator","active":true}`))
	}))
	defer server.Close()

	user, err := verifyJiraCredentials(server.URL, accountName, authToken)
	if err != nil {
		t.Fatalf("verifyJiraCredentials returned error: %v", err)
	}
	if !sawAuth {
		t.Fatal("expected basic auth credentials on current-user request")
	}
	if user.DisplayName != "Test Operator" {
		t.Fatalf("unexpected verified user: %#v", user)
	}
}

func TestVerifyJiraCredentialsReturnsAuthFailure(t *testing.T) {
	accountName := authFixtureValue("jira", "account")
	authToken := authFixtureValue("invalid", "token")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/2/myself" {
			t.Fatalf("unexpected request path %s", r.URL.Path)
		}
		http.Error(w, "bad credentials", http.StatusUnauthorized)
	}))
	defer server.Close()

	_, err := verifyJiraCredentials(server.URL, accountName, authToken)
	if err == nil {
		t.Fatal("expected verification error")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Fatalf("expected 401 in verification error, got %v", err)
	}
}

func TestVerifyJiraCredentialsReturnsDecodeFailure(t *testing.T) {
	accountName := authFixtureValue("jira", "account")
	authToken := authFixtureValue("jira", "token")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer server.Close()

	_, err := verifyJiraCredentials(server.URL, accountName, authToken)
	if err == nil {
		t.Fatal("expected decode error")
	}
}

func authFixtureValue(parts ...string) string {
	return strings.Join(parts, "-")
}

package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/pkg/browser"
)

const (
	pollInterval = 2 * time.Second
	pollTimeout  = 120 * time.Second
)

func BaseURL(subdomain string) string {
	if override := os.Getenv("NEETOAUTH_BASE_URL"); override != "" {
		return strings.TrimRight(override, "/")
	}
	return fmt.Sprintf("https://%s.neetoauth.com", subdomain)
}

func Login(subdomain string) (*Credentials, error) {
	baseURL := BaseURL(subdomain)
	loginToken, err := createSession(baseURL)
	if err != nil {
		return nil, fmt.Errorf("could not create authentication session: %w", err)
	}

	loginURL := fmt.Sprintf("%s/api/cli/v1/login?token=%s", baseURL, loginToken)
	fmt.Println("opening browser for authentication...")
	fmt.Printf("if the browser doesn't open, visit: %s\n", loginURL)
	if err := browser.OpenURL(loginURL); err != nil {
		fmt.Printf("could not open browser: %v\n", err)
	}

	fmt.Print("waiting for authentication")
	deadline := time.Now().Add(pollTimeout)
	for time.Now().Before(deadline) {
		time.Sleep(pollInterval)
		fmt.Print(".")

		status, email, sessionToken, err := checkStatus(baseURL, loginToken)
		if err != nil {
			continue
		}

		switch status {
		case "authenticated":
			fmt.Println(" done!")
			creds := Credentials{
				Subdomain:    subdomain,
				Email:        email,
				SessionToken: sessionToken,
			}
			store, err := LoadStore()
			if err != nil {
				return nil, fmt.Errorf("authenticated but could not read credentials: %w", err)
			}
			store.Upsert(creds)
			if err := SaveStore(store); err != nil {
				return nil, fmt.Errorf("authenticated but could not save credentials: %w", err)
			}
			return &creds, nil
		case "expired":
			fmt.Println()
			return nil, fmt.Errorf("authentication session expired, please try again")
		}
	}

	fmt.Println()
	return nil, fmt.Errorf("neetoauth cli authentication timed out after %d minutes, please try again", int(pollTimeout/time.Minute))
}

func createSession(baseURL string) (string, error) {
	resp, err := http.Post(baseURL+"/api/cli/v1/sessions", "application/json", nil)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound || redirectedAway(baseURL, resp) {
		return "", fmt.Errorf("subdomain not found, please check that you entered the correct subdomain\nfor example, if your NeetoAuth URL is acme.neetoauth.com then enter 'acme'")
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var result struct {
		LoginToken string `json:"login_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.LoginToken, nil
}

func redirectedAway(requestedURL string, resp *http.Response) bool {
	requested, err := url.Parse(requestedURL)
	if err != nil || resp.Request == nil {
		return false
	}
	return resp.Request.URL.Host != requested.Host
}

func checkStatus(baseURL, loginToken string) (status, email, sessionToken string, err error) {
	statusURL := fmt.Sprintf("%s/api/cli/v1/sessions/%s/status", baseURL, loginToken)
	resp, err := http.Post(statusURL, "application/json", bytes.NewReader([]byte("{}")))
	if err != nil {
		return "", "", "", err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}

	var result struct {
		Status       string `json:"status"`
		Email        string `json:"email"`
		SessionToken string `json:"session_token"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", "", err
	}
	return result.Status, result.Email, result.SessionToken, nil
}

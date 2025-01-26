package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストのパスに応じてレスポンスを返す
		switch r.URL.Path {
		case "/success":
			fmt.Fprint(w, "success")
		case "/fail":
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "error")
		case "/retry":
			// リトライ回数に応じて挙動を変える
			if r.URL.Query().Get("count") == "0" {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "error")
			} else {
				fmt.Fprint(w, "success after retry")
			}
		case "/retry_edge":
			// リトライ回数のエッジケース
			if r.URL.Query().Get("count") == "2" { // RetryLimit - 1
				fmt.Fprint(w, "success at retry limit - 1")
			} else if r.URL.Query().Get("count") == "3" { // RetryLimit
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "error at retry limit")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "error")
			}
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "not found")
		}
	}))
	defer ts.Close()

	tests := []struct {
		name       string
		urlPath    string
		expectErr  bool
		expectBody string
	}{
		{
			name:       "Success",
			urlPath:    "/success",
			expectErr:  false,
			expectBody: "success",
		},
		{
			name:       "Server Failure",
			urlPath:    "/fail",
			expectErr:  true,
			expectBody: "",
		},
		{
			name:       "Retry Success",
			urlPath:    "/retry?count=1",
			expectErr:  false,
			expectBody: "success after retry",
		},
		{
			name:       "Retry Edge Case",
			urlPath:    "/retry_edge?count=2",
			expectErr:  false,
			expectBody: "success at retry limit - 1",
		},
		{
			name:       "Retry Exceeded",
			urlPath:    "/retry_edge?count=3",
			expectErr:  true,
			expectBody: "",
		},
		{
			name:       "Not Found",
			urlPath:    "/unknown",
			expectErr:  false,
			expectBody: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testURL := ts.URL + tt.urlPath
			result := download(testURL)

			if (result.Err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, result.Err)
			}

			if result.Body != nil {
				defer result.Body.Close()
				body, err := io.ReadAll(result.Body)
				if err != nil {
					t.Fatalf("failed to read body: %v", err)
				}
				if string(body) != tt.expectBody {
					t.Errorf("expected body: '%s', got: '%s'", tt.expectBody, string(body))
				}
			} else if tt.expectBody != "" {
				t.Errorf("expected body: '%s', got nil", tt.expectBody)
			}
		})
	}
}

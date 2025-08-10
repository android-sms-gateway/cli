package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"

	"e2e/testutils"

	"github.com/stretchr/testify/assert"
)

func TestWebhookRegisterValid(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		event      string
		id         string
		expectJSON bool
	}{
		{
			name:       "register webhook with `sms:received` event",
			url:        "https://example.com/webhook",
			event:      "sms:received",
			id:         "",
			expectJSON: true,
		},
		{
			name:       "register webhook with `sms:delivered` event",
			url:        "https://example.com/delivery",
			event:      "sms:delivered",
			id:         "test-id",
			expectJSON: true,
		},
		{
			name:       "register webhook with `sms:sent` event",
			url:        "https://example.com/status",
			event:      "sms:sent",
			id:         "",
			expectJSON: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := testutils.CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/webhooks", r.URL.Path)

				var webhook struct {
					ID    string `json:"id"`
					URL   string `json:"url"`
					Event string `json:"event"`
				}
				err := json.NewDecoder(r.Body).Decode(&webhook)
				assert.NoError(t, err)
				assert.Equal(t, tt.url, webhook.URL)
				assert.Equal(t, tt.event, webhook.Event)
				if tt.id != "" {
					assert.Equal(t, tt.id, webhook.ID)
				}

				response := map[string]interface{}{
					"id":    "wh-test-123",
					"url":   webhook.URL,
					"event": webhook.Event,
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(response)
			})
			defer mockServer.Close()

			var stdout, stderr bytes.Buffer
			args := []string{
				"--format", "json",
				"webhooks", "register",
				"--event", tt.event,
			}
			if tt.id != "" {
				args = append(args, "--id", tt.id)
			}

			args = append(args, tt.url)

			cmd := exec.Command("./smsgate", args...)
			cmd.Env = append([]string{}, os.Environ()...)
			cmd.Env = append(cmd.Env, fmt.Sprintf("ASG_ENDPOINT=%s", mockServer.URL))
			cmd.Env = append(cmd.Env, "ASG_USERNAME=testuser", "ASG_PASSWORD=testpass")
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			assert.NoError(t, err, "stderr: %s", stderr.String())

			if tt.expectJSON {
				var response map[string]any
				err := json.Unmarshal(stdout.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "wh-test-123", response["id"])
			} else {
				assert.Contains(t, stdout.String(), "Webhook registered successfully")
			}
		})
	}
}

func TestWebhookRegisterInvalid(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		event    string
		id       string
		expected string
	}{
		{
			name:     "missing url",
			url:      "",
			event:    "sms:received",
			id:       "",
			expected: "URL is empty",
		},
		{
			name:     "invalid url",
			url:      "not-a-url",
			event:    "sms:received",
			id:       "",
			expected: "invalid URL",
		},
		{
			name:     "missing event",
			url:      "https://example.com/webhook",
			event:    "",
			id:       "",
			expected: "Required flag \"event\" not set",
		},
		{
			name:     "invalid event",
			url:      "https://example.com/webhook",
			event:    "invalid-event",
			id:       "",
			expected: "Invalid event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := testutils.CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})
			defer mockServer.Close()

			var stdout, stderr bytes.Buffer
			args := []string{
				"webhooks", "register",
			}
			if tt.event != "" {
				args = append(args, "--event", tt.event)
			}
			if tt.id != "" {
				args = append(args, "--id", tt.id)
			}
			if tt.url != "" {
				args = append(args, tt.url)
			}

			cmd := exec.Command("./smsgate", args...)
			cmd.Env = append([]string{}, os.Environ()...)
			cmd.Env = append(cmd.Env, fmt.Sprintf("ASG_ENDPOINT=%s", mockServer.URL))
			cmd.Env = append(cmd.Env, "ASG_USERNAME=testuser", "ASG_PASSWORD=testpass")
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			assert.Error(t, err)
			assert.Contains(t, stderr.String(), tt.expected)
		})
	}
}

func TestWebhookList(t *testing.T) {
	tests := []struct {
		name          string
		setupResponse func() []byte
		expectedCount int
		expectJSON    bool
	}{
		{
			name: "list empty webhooks",
			setupResponse: func() []byte {
				return []byte("[]")
			},
			expectedCount: 0,
			expectJSON:    true,
		},
		{
			name: "list multiple webhooks",
			setupResponse: func() []byte {
				return []byte(`[
					{
						"id": "wh-001",
						"url": "https://example.com/message",
						"event": "message"
					},
					{
						"id": "wh-002", 
						"url": "https://example.com/delivery",
						"event": "delivery"
					}
				]`)
			},
			expectedCount: 2,
			expectJSON:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := testutils.CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/webhooks", r.URL.Path)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(tt.setupResponse())
			})
			defer mockServer.Close()

			var stdout, stderr bytes.Buffer
			cmd := exec.Command(
				"./smsgate",
				"--format", "json",
				"webhooks", "list",
			)
			cmd.Env = append([]string{}, os.Environ()...)
			cmd.Env = append(cmd.Env, fmt.Sprintf("ASG_ENDPOINT=%s", mockServer.URL))
			cmd.Env = append(cmd.Env, "ASG_USERNAME=testuser", "ASG_PASSWORD=testpass")
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			assert.NoError(t, err, "stderr: %s", stderr.String())

			if tt.expectJSON {
				var webhooks []map[string]interface{}
				err := json.Unmarshal(stdout.Bytes(), &webhooks)
				assert.NoError(t, err)
				assert.Len(t, webhooks, tt.expectedCount)
			}
		})
	}
}

func TestWebhookDelete(t *testing.T) {
	tests := []struct {
		name        string
		webhookID   string
		setupStatus int
		expectError bool
		expectedMsg string
	}{
		{
			name:        "delete existing webhook",
			webhookID:   "wh-test-123",
			setupStatus: http.StatusOK,
			expectError: false,
			expectedMsg: "Success",
		},
		{
			name:        "delete non-existing webhook",
			webhookID:   "wh-nonexistent",
			setupStatus: http.StatusNotFound,
			expectError: true,
			expectedMsg: "not found",
		},
		{
			name:        "delete webhook with empty id",
			webhookID:   "",
			setupStatus: http.StatusBadRequest,
			expectError: true,
			expectedMsg: "ID is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := testutils.CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodDelete, r.Method)
				if tt.webhookID != "" {
					assert.Equal(t, "/webhooks/"+tt.webhookID, r.URL.Path)
				}

				switch tt.setupStatus {
				case http.StatusNotFound:
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"error": "Webhook not found"}`))
				case http.StatusBadRequest:
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error": "Invalid webhook ID"}`))
				default:
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"success": true}`))
				}
			})
			defer mockServer.Close()

			var stdout, stderr bytes.Buffer
			args := []string{"webhooks", "delete"}
			if tt.webhookID != "" {
				args = append(args, tt.webhookID)
			}

			cmd := exec.Command("./smsgate", args...)
			cmd.Env = append([]string{}, os.Environ()...)
			cmd.Env = append(cmd.Env, fmt.Sprintf("ASG_ENDPOINT=%s", mockServer.URL))
			cmd.Env = append(cmd.Env, "ASG_USERNAME=testuser", "ASG_PASSWORD=testpass")
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, stderr.String(), tt.expectedMsg)
			} else {
				assert.NoError(t, err, "stderr: %s", stderr.String())
				assert.Contains(t, stdout.String(), tt.expectedMsg)
			}
		})
	}
}

func TestWebhookCommandHelp(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectHelp  bool
		expectUsage bool
	}{
		{
			name:        "webhooks command help",
			args:        []string{"webhooks", "--help"},
			expectHelp:  true,
			expectUsage: true,
		},
		{
			name:        "register command help",
			args:        []string{"webhooks", "register", "--help"},
			expectHelp:  true,
			expectUsage: true,
		},
		{
			name:        "list command help",
			args:        []string{"webhooks", "list", "--help"},
			expectHelp:  true,
			expectUsage: true,
		},
		{
			name:        "delete command help",
			args:        []string{"webhooks", "delete", "--help"},
			expectHelp:  true,
			expectUsage: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			cmd := exec.Command("./smsgate", tt.args...)
			cmd.Env = append([]string{}, os.Environ()...)
			cmd.Env = append(cmd.Env, "ASG_USERNAME=testuser", "ASG_PASSWORD=testpass")
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			assert.NoError(t, err, "stderr: %s", stderr.String())

			output := stdout.String()
			if tt.expectHelp {
				assert.Contains(t, output, "help")
			}
			if tt.expectUsage {
				assert.Contains(t, output, "USAGE:")
			}
		})
	}
}

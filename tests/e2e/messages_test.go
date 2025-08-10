package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"e2e/testutils"

	"github.com/stretchr/testify/assert"
)

func TestMessageSendValid(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		phones         []string
		deviceID       string
		simNumber      int
		priority       int
		dataMessage    bool
		dataPort       uint
		expectJSON     bool
		expectedStatus int
	}{
		{
			name:           "send SMS with single recipient",
			message:        "Hello, this is a test message",
			phones:         []string{"+12025550123"},
			deviceID:       "",
			simNumber:      0,
			priority:       0,
			dataMessage:    false,
			dataPort:       0,
			expectJSON:     false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "send SMS with multiple recipients",
			message:        "Hello everyone!",
			phones:         []string{"+12025550123", "+12025550124", "+12025550125"},
			deviceID:       "",
			simNumber:      0,
			priority:       0,
			dataMessage:    false,
			dataPort:       0,
			expectJSON:     false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "send SMS with device ID",
			message:        "Device specific message",
			phones:         []string{"+12025550123"},
			deviceID:       "device-123",
			simNumber:      0,
			priority:       0,
			dataMessage:    false,
			dataPort:       0,
			expectJSON:     false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "send SMS with SIM number",
			message:        "SIM specific message",
			phones:         []string{"+12025550123"},
			deviceID:       "",
			simNumber:      1,
			priority:       0,
			dataMessage:    false,
			dataPort:       0,
			expectJSON:     false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "send SMS with high priority",
			message:        "High priority message",
			phones:         []string{"+12025550123"},
			deviceID:       "",
			simNumber:      0,
			priority:       100,
			dataMessage:    false,
			dataPort:       0,
			expectJSON:     false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "send data message",
			message:        "SGVsbG8gV29ybGQh", // "Hello World!" in base64
			phones:         []string{"+12025550123"},
			deviceID:       "",
			simNumber:      0,
			priority:       0,
			dataMessage:    true,
			dataPort:       8080,
			expectJSON:     false,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := testutils.CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/messages", r.URL.Path)

				var req map[string]any
				err := json.NewDecoder(r.Body).Decode(&req)
				assert.NoError(t, err)

				// Validate request structure - convert interface{} to []string
				phoneNumbers := req["phoneNumbers"].([]any)
				expectedPhones := make([]any, len(tt.phones))
				for i, phone := range tt.phones {
					expectedPhones[i] = phone
				}
				assert.Equal(t, expectedPhones, phoneNumbers)
				if tt.deviceID != "" {
					assert.Equal(t, tt.deviceID, req["deviceId"])
				}
				if tt.simNumber > 0 {
					assert.Equal(t, float64(tt.simNumber), req["simNumber"])
				}
				if tt.priority != 0 {
					assert.Equal(t, float64(tt.priority), req["priority"])
				}
				if tt.dataMessage {
					assert.NotNil(t, req["dataMessage"])
					dataMsg := req["dataMessage"].(map[string]interface{})
					assert.Equal(t, tt.message, dataMsg["data"])
					assert.Equal(t, float64(tt.dataPort), dataMsg["port"])
				} else {
					assert.NotNil(t, req["textMessage"])
					textMsg := req["textMessage"].(map[string]interface{})
					assert.Equal(t, tt.message, textMsg["text"])
				}

				response := map[string]interface{}{
					"id":        fmt.Sprintf("msg-%d", time.Now().Unix()),
					"state":     "Pending",
					"createdAt": time.Now().Format(time.RFC3339),
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.expectedStatus)
				json.NewEncoder(w).Encode(response)
			})
			defer mockServer.Close()

			var stdout, stderr bytes.Buffer
			args := []string{
				"send",
				"--phones", strings.Join(tt.phones, ","),
			}

			if tt.deviceID != "" {
				args = append(args, "--device-id", tt.deviceID)
			}
			if tt.simNumber > 0 {
				args = append(args, "--sim-number", fmt.Sprintf("%d", tt.simNumber))
			}
			if tt.priority != 0 {
				args = append(args, "--priority", fmt.Sprintf("%d", tt.priority))
			}
			if tt.dataMessage {
				args = append(args, "--data", "--data-port", fmt.Sprintf("%d", tt.dataPort))
			}

			args = append(args, tt.message)

			cmd := exec.Command("./smsgate", args...)
			cmd.Env = append([]string{}, os.Environ()...)
			cmd.Env = append(cmd.Env, fmt.Sprintf("ASG_ENDPOINT=%s", mockServer.URL))
			cmd.Env = append(cmd.Env, "ASG_USERNAME=testuser", "ASG_PASSWORD=testpass")
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			assert.NoError(t, err, "stderr: %s", stderr.String())

			if tt.expectJSON {
				var response map[string]interface{}
				err := json.Unmarshal(stdout.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["id"], "msg-")
			}
		})
	}
}

func TestMessageSendInvalid(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		phones      []string
		extraArgs   []string
		expectedErr string
		setupStatus int
	}{
		{
			name:        "missing message content",
			message:     "",
			phones:      []string{"+12025550123"},
			expectedErr: "Message is empty",
			setupStatus: http.StatusBadRequest,
		},
		{
			name:        "missing phone numbers",
			message:     "Hello, world!",
			phones:      []string{},
			expectedErr: "Required flag \"phones\" not set",
			setupStatus: http.StatusBadRequest,
		},
		{
			name:        "invalid phone format",
			message:     "Hello, world!",
			phones:      []string{"invalid-phone"},
			expectedErr: "validation failed",
			setupStatus: http.StatusBadRequest,
		},
		{
			name:        "invalid SIM number",
			message:     "Hello, world!",
			phones:      []string{"+12025550123"},
			extraArgs:   []string{"--sim-number", "bad-sim"},
			expectedErr: "invalid value \"bad-sim\" for flag -sim-number",
			setupStatus: http.StatusBadRequest,
		},
		{
			name:        "invalid priority",
			message:     "Hello, world!",
			phones:      []string{"+12025550123"},
			extraArgs:   []string{"--priority", "9999"},
			expectedErr: "Priority must be between -128 and 127",
			setupStatus: http.StatusBadRequest,
		},
		{
			name:        "invalid base64 data",
			message:     "not-base64",
			phones:      []string{"+12025550123"},
			extraArgs:   []string{"--data"},
			expectedErr: "Invalid base64 data",
			setupStatus: http.StatusBadRequest,
		},
		{
			name:        "invalid data port",
			message:     "SGVsbG8=", // valid base64
			phones:      []string{"+12025550123"},
			extraArgs:   []string{"--data", "--data-port", "99999"},
			expectedErr: "Data port must be between 1 and 65535",
			setupStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := testutils.CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.setupStatus)
			})
			defer mockServer.Close()

			var stdout, stderr bytes.Buffer
			args := []string{"send"}

			if len(tt.phones) > 0 {
				args = append(args, "--phones", strings.Join(tt.phones, ","))
			}

			if len(tt.extraArgs) > 0 {
				args = append(args, tt.extraArgs...)
			}

			args = append(args, tt.message)

			cmd := exec.Command("./smsgate", args...)
			cmd.Env = append([]string{}, os.Environ()...)
			cmd.Env = append(cmd.Env, fmt.Sprintf("ASG_ENDPOINT=%s", mockServer.URL))
			cmd.Env = append(cmd.Env, "ASG_USERNAME=testuser", "ASG_PASSWORD=testpass")
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			assert.Error(t, err)
			assert.Contains(t, stderr.String(), tt.expectedErr)
		})
	}
}

func TestMessageSendAuthenticationError(t *testing.T) {
	mockServer := testutils.CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Invalid credentials"}`))
	})
	defer mockServer.Close()

	var stdout, stderr bytes.Buffer
	args := []string{
		"send",
		"--phones", "+12025550123",
		"Hello, world!",
	}

	cmd := exec.Command("./smsgate", args...)
	cmd.Env = append([]string{}, os.Environ()...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("ASG_ENDPOINT=%s", mockServer.URL))
	cmd.Env = append(cmd.Env, "ASG_USERNAME=wronguser", "ASG_PASSWORD=wrongpass")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	assert.Error(t, err)
	assert.Contains(t, stderr.String(), "Invalid credentials")
}

func TestMessageStatusValid(t *testing.T) {
	tests := []struct {
		name        string
		messageID   string
		setupStatus int
		setupBody   string
		expectJSON  bool
	}{
		{
			name:        "get existing message status",
			messageID:   "msg-12345",
			setupStatus: http.StatusOK,
			setupBody:   `{"id": "msg-12345", "state": "Delivered", "createdAt": "2023-01-01T00:00:00Z"}`,
			expectJSON:  false,
		},
		{
			name:        "get non-existing message status",
			messageID:   "msg-nonexistent",
			setupStatus: http.StatusNotFound,
			setupBody:   `{"error": "Message not found"}`,
			expectJSON:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := testutils.CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/messages/"+tt.messageID, r.URL.Path)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.setupStatus)
				w.Write([]byte(tt.setupBody))
			})
			defer mockServer.Close()

			var stdout, stderr bytes.Buffer
			args := []string{
				"status",
				tt.messageID,
			}

			cmd := exec.Command("./smsgate", args...)
			cmd.Env = append([]string{}, os.Environ()...)
			cmd.Env = append(cmd.Env, fmt.Sprintf("ASG_ENDPOINT=%s", mockServer.URL))
			cmd.Env = append(cmd.Env, "ASG_USERNAME=testuser", "ASG_PASSWORD=testpass")
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if tt.setupStatus == http.StatusNotFound {
				assert.Error(t, err)
				assert.Contains(t, stderr.String(), "Message not found")
			} else {
				assert.NoError(t, err, "stderr: %s", stderr.String())
				assert.Contains(t, stdout.String(), "msg-12345")
			}
		})
	}
}

func TestMessageStatusInvalid(t *testing.T) {
	tests := []struct {
		name        string
		messageID   string
		expectedErr string
		setupStatus int
	}{
		{
			name:        "empty message ID",
			messageID:   "",
			expectedErr: "Message ID is empty",
			setupStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := testutils.CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.setupStatus)
			})
			defer mockServer.Close()

			var stdout, stderr bytes.Buffer
			args := []string{
				"status",
			}

			if tt.messageID != "" {
				args = append(args, tt.messageID)
			}

			cmd := exec.Command("./smsgate", args...)
			cmd.Env = append([]string{}, os.Environ()...)
			cmd.Env = append(cmd.Env, fmt.Sprintf("ASG_ENDPOINT=%s", mockServer.URL))
			cmd.Env = append(cmd.Env, "ASG_USERNAME=testuser", "ASG_PASSWORD=testpass")
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			assert.Error(t, err)
			assert.Contains(t, stderr.String(), tt.expectedErr)
		})
	}
}

func TestMessageCommandHelp(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectHelp  bool
		expectUsage bool
	}{
		{
			name:        "send command help",
			args:        []string{"send", "--help"},
			expectHelp:  true,
			expectUsage: true,
		},
		{
			name:        "status command help",
			args:        []string{"status", "--help"},
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

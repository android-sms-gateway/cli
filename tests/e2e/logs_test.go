package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"testing"
	"time"

	"e2e/testutils"

	"github.com/stretchr/testify/assert"
)

func TestLogs(t *testing.T) {
	from := time.Date(2025, 1, 10, 11, 0, 0, 0, time.UTC)
	to := time.Date(2025, 1, 10, 12, 0, 0, 0, time.UTC)

	t.Run("get logs in json format", func(t *testing.T) {
		mockServer := testutils.CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/logs", r.URL.Path)

			q, err := url.ParseQuery(r.URL.RawQuery)
			assert.NoError(t, err)
			assert.Equal(t, from.Format(time.RFC3339), q.Get("from"))
			assert.Equal(t, to.Format(time.RFC3339), q.Get("to"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`[
				{
					"id": 101,
					"priority": "INFO",
					"module": "messages",
					"message": "Sent message",
					"context": {"messageId": "msg-1"},
					"createdAt": "2025-01-10T11:30:00Z"
				}
			]`))
		})
		defer mockServer.Close()

		var stdout, stderr bytes.Buffer
		cmd := exec.Command(
			"./smsgate",
			"--format", "json",
			"logs",
			"--from", from.Format(time.RFC3339),
			"--to", to.Format(time.RFC3339),
		)
		cmd.Env = append([]string{}, os.Environ()...)
		cmd.Env = append(cmd.Env, fmt.Sprintf("ASG_ENDPOINT=%s", mockServer.URL))
		cmd.Env = append(cmd.Env, "ASG_USERNAME=testuser", "ASG_PASSWORD=testpass")
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		assert.NoError(t, err, "stderr: %s", stderr.String())

		var out []map[string]any
		err = json.Unmarshal(stdout.Bytes(), &out)
		assert.NoError(t, err)
		assert.Len(t, out, 1)
		assert.Equal(t, float64(101), out[0]["id"])
	})

	t.Run("invalid range", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		cmd := exec.Command(
			"./smsgate",
			"logs",
			"--from", to.Format(time.RFC3339),
			"--to", from.Format(time.RFC3339),
		)
		cmd.Env = append([]string{}, os.Environ()...)
		cmd.Env = append(cmd.Env, "ASG_ENDPOINT=http://localhost:9999")
		cmd.Env = append(cmd.Env, "ASG_USERNAME=testuser", "ASG_PASSWORD=testpass")
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		assert.Error(t, err)
		assert.Contains(t, stderr.String(), "From date must be less than or equal to To date")
	})
}

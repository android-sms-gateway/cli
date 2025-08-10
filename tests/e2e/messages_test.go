package e2e

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendCommand(t *testing.T) {
	// Build the CLI binary
	cmd := exec.Command("go", "build", "-o", "tests/e2e/smsgate", "cmd/smsgate/smsgate.go")
	cmd.Dir = "../../"
	err := cmd.Run()
	assert.NoError(t, err)

	t.Cleanup(func() {
		cmd := exec.Command("rm", "smsgate")
		err := cmd.Run()
		assert.NoError(t, err)
	})

	// Run the CLI binary with the send command
	cmd = exec.Command("./smsgate", "send", "--help")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	assert.NoError(t, err)

	// Verify the output
	assert.Contains(t, out.String(), "Send message")
}

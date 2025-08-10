package e2e

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelpFlag(t *testing.T) {
	// Run the CLI binary with the --help flag
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("./smsgate", "--help")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	assert.NoError(t, err)

	// Verify the output
	assert.Contains(t, stdout.String(), "CLI interface for working with SMS Gateway for Android™")
}

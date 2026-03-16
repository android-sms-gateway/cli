package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	// Compute absolute path for the binary
	binPath, err := filepath.Abs("smsgate")
	if err != nil {
		panic(fmt.Sprintf("failed to get absolute path: %v", err))
	}

	// Build the CLI binary
	cmd := exec.Command("go", "build", "-o", binPath, "cmd/smsgate/smsgate.go")
	cmd.Dir = "../../"
	err = cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("failed to build binary: %v", err))
	}

	// Export the binary path for tests to use
	os.Setenv("SMSGATE_BIN", binPath)

	code := m.Run()

	// Cleanup: remove the built binary
	_ = os.Remove(binPath)

	os.Exit(code)
}

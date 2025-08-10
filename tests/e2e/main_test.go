package e2e

import (
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	// Build the CLI binary
	cmd := exec.Command("go", "build", "-o", "tests/e2e/smsgate", "cmd/smsgate/smsgate.go")
	cmd.Dir = "../../"
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	code := m.Run()

	_ = os.Remove("smsgate")

	os.Exit(code)
}

package batch_test

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/android-sms-gateway/cli/internal/commands/messages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestCommands_ContainsBatch(t *testing.T) {
	t.Parallel()

	_, ok := findBatchSendCommand(messages.Commands())
	assert.True(t, ok)
}

func newBatchSendFixture(t *testing.T) (string, *cli.Command, bool) {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "input.csv")
	require.NoError(t, os.WriteFile(path, []byte("Phone,Message\n+12025550123,Hello"), 0o600))

	cmd, ok := findBatchSendCommand(messages.Commands())
	return path, cmd, ok
}

func TestBatchSend_ValidateOnly(t *testing.T) {
	t.Parallel()

	path, cmd, ok := newBatchSendFixture(t)
	require.True(t, ok)

	ctx := newContext(t, cmd, []string{
		"--map", "phone=Phone,text=Message",
		"--validate-only",
		path,
	})

	require.NoError(t, cmd.Before(ctx))
	require.NoError(t, cmd.Action(ctx))
}

func TestBatchSend_InvalidMap(t *testing.T) {
	t.Parallel()

	path, cmd, ok := newBatchSendFixture(t)
	require.True(t, ok)

	ctx := newContext(t, cmd, []string{
		"--map", "text=Message",
		"--validate-only",
		path,
	})

	err := cmd.Before(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "map must include phone=<column>")
}

func findBatchSendCommand(cmds []*cli.Command) (*cli.Command, bool) {
	for _, cmd := range cmds {
		if cmd.Name != "batch" {
			continue
		}
		for _, sub := range cmd.Subcommands {
			if sub.Name == "send" {
				return sub, true
			}
		}
	}

	return nil, false
}

func newContext(t *testing.T, cmd *cli.Command, args []string) *cli.Context {
	t.Helper()

	set := flag.NewFlagSet(cmd.Name, flag.ContinueOnError)
	for _, cliFlag := range cmd.Flags {
		require.NoError(t, cliFlag.Apply(set))
	}
	require.NoError(t, set.Parse(args))

	return cli.NewContext(&cli.App{}, set, nil)
}

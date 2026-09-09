package messages_test

import (
	"flag"
	"testing"

	"github.com/android-sms-gateway/cli/internal/commands/messages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestSendCmd_MMSFlags(t *testing.T) {
	t.Parallel()

	cmd := findSendCommand(t)
	ctx := newSendContext(t, cmd, []string{
		"--mms",
		"--subject", "Hello",
		"--attachment", "a.png",
		"--attachment", "b.jpg",
		"--phones", "+12025550123",
	})

	assert.True(t, ctx.Bool("mms"))
	assert.Equal(t, "Hello", ctx.String("subject"))
	assert.Equal(t, []string{"a.png", "b.jpg"}, ctx.StringSlice("attachment"))
}

func TestSendBefore_MMS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name: "valid text only",
			args: []string{"--mms", "--phones", "+12025550123", "hello"},
		},
		{
			name: "valid attachment only",
			args: []string{"--mms", "--phones", "+12025550123", "--attachment", "a.png"},
		},
		{
			name:    "no text and no attachments",
			args:    []string{"--mms", "--phones", "+12025550123"},
			wantErr: "MMS message requires text or at least one attachment",
		},
		{
			name:    "mms and data are mutually exclusive",
			args:    []string{"--mms", "--data", "--phones", "+12025550123", "hello"},
			wantErr: "--data and --mms are mutually exclusive",
		},
		{
			name:    "subject without mms",
			args:    []string{"--subject", "Hello", "--phones", "+12025550123", "hello"},
			wantErr: "--subject and --attachment require --mms",
		},
		{
			name:    "attachment without mms",
			args:    []string{"--attachment", "a.png", "--phones", "+12025550123", "hello"},
			wantErr: "--subject and --attachment require --mms",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cmd := findSendCommand(t)
			ctx := newSendContext(t, cmd, tt.args)

			err := cmd.Before(ctx)
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func findSendCommand(t *testing.T) *cli.Command {
	t.Helper()

	for _, cmd := range messages.Commands() {
		if cmd.Name == "send" {
			return cmd
		}
	}

	t.Fatal("send command not found")
	return nil
}

func newSendContext(t *testing.T, cmd *cli.Command, args []string) *cli.Context {
	t.Helper()

	set := flag.NewFlagSet(cmd.Name, flag.ContinueOnError)
	for _, cliFlag := range cmd.Flags {
		require.NoError(t, cliFlag.Apply(set))
	}
	require.NoError(t, set.Parse(args))

	return cli.NewContext(&cli.App{}, set, nil)
}

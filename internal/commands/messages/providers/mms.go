package providers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/urfave/cli/v2"
)

const fallbackContentType = "application/octet-stream"

// MmsContentProvider fills the MmsMessage field of a Message.
// It parses text, subject, and attachment paths from CLI flags.
type MmsContentProvider struct{}

// PrepareMessage fills the MmsMessage field of msg.
// It reads text from the first CLI argument, subject from --subject flag,
// and attachment paths from --attachment flag.
// Blank text and subject are omitted (nil). MIME type is detected from
// file extension, falling back to application/octet-stream.
func (p MmsContentProvider) PrepareMessage(c *cli.Context, msg *smsgateway.Message) error {
	text := strings.TrimSpace(c.Args().Get(0))
	subject := c.String("subject")
	attachmentPaths := c.StringSlice("attachment")

	attachments, err := loadMmsAttachments(attachmentPaths)
	if err != nil {
		return err
	}

	msg.MmsMessage = newMmsMessage(text, subject, attachments)
	return nil
}

// newMmsMessage builds an MMS message. Blank text and subject are omitted.
func newMmsMessage(text string, subject string, attachments []smsgateway.MmsAttachment) *smsgateway.MmsMessage {
	mms := &smsgateway.MmsMessage{
		Subject:     nil,
		Text:        nil,
		Attachments: attachments,
	}

	if text != "" {
		mms.Text = &text
	}
	if subject != "" {
		mms.Subject = &subject
	}

	return mms
}

// loadMmsAttachments reads the given files and builds MMS attachments.
// The content type is detected from the file extension.
func loadMmsAttachments(paths []string) ([]smsgateway.MmsAttachment, error) {
	attachments := make([]smsgateway.MmsAttachment, 0, len(paths))
	for _, path := range paths {
		attachment, err := loadMmsAttachment(path)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, attachment)
	}

	return attachments, nil
}

// loadMmsAttachment reads a single file and builds an MMS attachment with
// base64-encoded content.
func loadMmsAttachment(path string) (smsgateway.MmsAttachment, error) {
	buffer := bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, &buffer)

	file, err := os.Open(path)
	if err != nil {
		return smsgateway.MmsAttachment{}, fmt.Errorf("failed to open attachment %q: %w", path, err)
	}
	defer file.Close()

	if _, cpErr := io.Copy(encoder, file); cpErr != nil {
		return smsgateway.MmsAttachment{}, fmt.Errorf("failed to read attachment %q: %w", path, cpErr)
	}
	if closeErr := encoder.Close(); closeErr != nil {
		return smsgateway.MmsAttachment{}, fmt.Errorf("failed to read attachment %q: %w", path, closeErr)
	}

	name := filepath.Base(path)

	return smsgateway.MmsAttachment{
		ContentType: detectContentType(path),
		Name:        &name,
		Data:        buffer.String(),
	}, nil
}

// detectContentType detects the MIME type from the file extension.
// It falls back to application/octet-stream for unknown extensions.
func detectContentType(name string) string {
	contentType := mime.TypeByExtension(filepath.Ext(name))
	if contentType == "" {
		return fallbackContentType
	}

	// Strip parameters (e.g. "text/plain; charset=utf-8" -> "text/plain").
	if params := strings.IndexByte(contentType, ';'); params >= 0 {
		contentType = contentType[:params]
	}

	return strings.TrimSpace(contentType)
}

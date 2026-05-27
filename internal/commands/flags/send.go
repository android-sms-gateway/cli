package flags

import (
	"fmt"
	"time"

	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

const (
	categoryOptions = "Options"
	categoryBody    = "Body"
)

func Send() []cli.Flag {
	return []cli.Flag{
		// Body fields
		&cli.StringFlag{
			Name:        "device-id",
			Aliases:     []string{"device", "deviceId"},
			Category:    categoryBody,
			Usage:       "Optional device ID for explicit selection",
			DefaultText: "auto",
		},
		&cli.UintFlag{
			Name:        "sim-number",
			Aliases:     []string{"simNumber", "sim"},
			Category:    categoryBody,
			Usage:       "SIM card index (one-based index, e.g. 1)",
			DefaultText: "see device settings",
		},
		&cli.BoolFlag{
			Name:     "delivery-report",
			Aliases:  []string{"deliveryReport"},
			Category: categoryBody,
			Usage:    "Enable delivery report",
			Value:    true,
		},
		&cli.IntFlag{
			Name:     "priority",
			Category: categoryBody,
			Usage:    "Priority, use >= 100 to bypass all limits and delays (-128 to 127)",
			Value:    0,
		},

		&cli.DurationFlag{
			Name:        "ttl",
			Category:    categoryOptions,
			Usage:       "Time to live (duration, e.g. 1h30m)",
			DefaultText: "unlimited",
		},
		&cli.TimestampFlag{
			Name:     "valid-until",
			Aliases:  []string{"validUntil"},
			Category: categoryOptions,
			Usage:    "Valid until (RFC3339 format, e.g. 2006-01-02T15:04:05Z07:00)",
			Layout:   time.RFC3339,
			Timezone: time.Local,
		},
		&cli.TimestampFlag{
			Name:     "schedule-at",
			Aliases:  []string{"scheduleAt"},
			Category: categoryOptions,
			Usage:    "Schedule message delivery at a specific time (RFC3339 format, e.g. 2006-01-02T15:04:05Z07:00)",
			Layout:   time.RFC3339,
			Timezone: time.Local,
		},
		&cli.BoolFlag{
			Name:     "skip-phone-validation",
			Aliases:  []string{"skipPhoneValidation"},
			Category: categoryOptions,
			Usage:    "Skip phone number validation",
			Value:    false,
		},
		&cli.UintFlag{
			Name:     "device-active-within",
			Aliases:  []string{"deviceActiveWithin"},
			Category: categoryOptions,
			Usage:    "Filter devices active within the specified number of hours",
			Value:    0,
		},
	}
}

type SendFlags struct {
	DeviceID            string
	SimNumber           *uint8
	DeliveryReport      bool
	Priority            smsgateway.MessagePriority
	TTL                 *uint64
	ValidUntil          *time.Time
	ScheduleAt          *time.Time
	SkipPhoneValidation bool
	DeviceActiveWithin  uint
}

func NewSendFlags(c *cli.Context) (*SendFlags, error) {
	const maxSIMNumber = 255

	fl := &SendFlags{
		DeviceID:            c.String("device-id"),
		DeliveryReport:      c.Bool("delivery-report"),
		ValidUntil:          c.Timestamp("valid-until"),
		SkipPhoneValidation: c.Bool("skip-phone-validation"),
		DeviceActiveWithin:  c.Uint("device-active-within"),
		ScheduleAt:          c.Timestamp("schedule-at"),

		SimNumber: nil,
		Priority:  smsgateway.PriorityDefault,
		TTL:       nil,
	}

	simNumber := c.Uint("sim-number")
	if simNumber > maxSIMNumber {
		return nil, fmt.Errorf(
			"%w: sim number must be between 0 and %d: %d",
			ErrValidationFailed,
			maxSIMNumber,
			simNumber,
		)
	}
	fl.SimNumber = lo.EmptyableToPtr(uint8(simNumber))

	ttl := c.Duration("ttl")
	if ttl < 0 {
		return nil, fmt.Errorf("%w: TTL must be positive: %s", ErrValidationFailed, ttl.String())
	}
	if ttl%time.Second != 0 {
		return nil, fmt.Errorf("%w: TTL must be a whole number of seconds: %s", ErrValidationFailed, ttl.String())
	}
	fl.TTL = lo.EmptyableToPtr(uint64(ttl / time.Second))

	validUntil := c.Timestamp("valid-until")
	if validUntil != nil && validUntil.Before(time.Now()) {
		return nil, fmt.Errorf(
			"%w: Valid Until must be in the future: %s",
			ErrValidationFailed,
			validUntil.Format(time.RFC3339),
		)
	}
	fl.ValidUntil = validUntil

	scheduleAt := c.Timestamp("schedule-at")
	if scheduleAt != nil && !scheduleAt.After(time.Now()) {
		return nil, fmt.Errorf(
			"%w: Schedule At must be in the future: %s",
			ErrValidationFailed,
			scheduleAt.Format(time.RFC3339),
		)
	}
	fl.ScheduleAt = scheduleAt

	priority := c.Int(
		"priority",
	)
	if priority < int(smsgateway.PriorityMinimum) ||
		priority > int(smsgateway.PriorityMaximum) {
		return nil, fmt.Errorf(
			"%w: Priority must be between %d and %d: %d",
			ErrValidationFailed,
			smsgateway.PriorityMinimum,
			smsgateway.PriorityMaximum,
			priority,
		)
	}
	fl.Priority = smsgateway.MessagePriority(priority)

	if fl.TTL != nil && fl.ValidUntil != nil {
		return nil, fmt.Errorf("%w: TTL and Valid Until flags are mutually exclusive", ErrValidationFailed)
	}

	return fl, nil
}

func (s SendFlags) Merge(src smsgateway.Message) smsgateway.Message {
	return smsgateway.Message{
		ID:           src.ID,
		Message:      src.Message,
		TextMessage:  src.TextMessage,
		DataMessage:  src.DataMessage,
		PhoneNumbers: src.PhoneNumbers,
		IsEncrypted:  src.IsEncrypted,

		DeviceID:           lo.CoalesceOrEmpty(src.DeviceID, s.DeviceID),
		SimNumber:          lo.CoalesceOrEmpty(src.SimNumber, s.SimNumber),
		WithDeliveryReport: lo.CoalesceOrEmpty(src.WithDeliveryReport, &s.DeliveryReport),
		Priority:           lo.CoalesceOrEmpty(src.Priority, s.Priority),
		TTL:                lo.CoalesceOrEmpty(src.TTL, s.TTL),
		ValidUntil:         lo.CoalesceOrEmpty(src.ValidUntil, s.ValidUntil),
		ScheduleAt:         lo.CoalesceOrEmpty(src.ScheduleAt, s.ScheduleAt),
	}
}

func (s SendFlags) Option() []smsgateway.SendOption {
	options := []smsgateway.SendOption{}

	if s.SkipPhoneValidation {
		options = append(options, smsgateway.WithSkipPhoneValidation(true))
	}
	if s.DeviceActiveWithin > 0 {
		options = append(options, smsgateway.WithDeviceActiveWithin(s.DeviceActiveWithin))
	}

	return options
}

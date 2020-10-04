package store

import (
	"strings"

	"github.com/pkg/errors"
)

const (
	FlagOnlyFailedBuilds = "only-failed"
)

// Contain the options for subscriptions
type SubscriptionFlags struct {
	OnlyFailedBuilds bool `json:"OnlyFailedBuilds"`
}

// Add a flag to the structure
func (s *SubscriptionFlags) AddFlag(flag string) error {
	switch flag { // nolint:gocritic // It's expected that more flags get added.
	case FlagOnlyFailedBuilds:
		s.OnlyFailedBuilds = true

	default:
		return errors.New("Unknown flag " + flag)
	}

	return nil
}

// Output the flags in a well-formatted string
func (s SubscriptionFlags) String() string {
	flags := []string{}

	if s.OnlyFailedBuilds {
		flag := "--" + FlagOnlyFailedBuilds
		flags = append(flags, flag)
	}

	if len(flags) == 0 {
		return "No flags set"
	}

	return strings.Join(flags, ",")
}

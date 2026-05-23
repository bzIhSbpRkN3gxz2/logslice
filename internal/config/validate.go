package config

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationError collects one or more field-level problems found during
// a call to Validate.
type ValidationError struct {
	Problems []string
}

func (e *ValidationError) Error() string {
	return "config validation failed: " + strings.Join(e.Problems, "; ")
}

// Validate performs semantic checks on cfg beyond what Parse enforces.
// It returns a *ValidationError listing every problem found, or nil when cfg
// is fully valid. This allows callers to surface all issues at once rather
// than fixing them one at a time.
func Validate(cfg *Config) error {
	if cfg == nil {
		return errors.New("config: nil Config")
	}

	var problems []string

	if cfg.Source == "" {
		problems = append(problems, "source must not be empty")
	}

	for i, out := range cfg.Outputs {
		if out == "" {
			problems = append(problems, fmt.Sprintf("output[%d] must not be empty", i))
		}
	}

	if cfg.TimeFormat != "" {
		// A minimal sanity check: the format must contain at least one digit
		// placeholder from the reference time (Mon Jan 2 15:04:05 MST 2006).
		refTokens := []string{"2006", "01", "02", "15", "04", "05", "MST", "Mon", "Jan"}
		found := false
		for _, tok := range refTokens {
			if strings.Contains(cfg.TimeFormat, tok) {
				found = true
				break
			}
		}
		if !found {
			problems = append(problems, "-fmt does not appear to contain a Go reference-time token")
		}
	}

	if len(problems) > 0 {
		return &ValidationError{Problems: problems}
	}
	return nil
}

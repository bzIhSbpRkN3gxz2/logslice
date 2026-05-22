package filter

import "errors"

// ErrInvalidRange is returned when Start is after End in a Range.
var ErrInvalidRange = errors.New("filter: start time must not be after end time")

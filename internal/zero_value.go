package internal

import (
	"strconv"
	"strings"
)

// IsZeroValue determines whether a frontmatter value counts as zero.
// A missing key (present == false) is treated as zero. ok is false when the
// value's type or content can't be safely interpreted as a number, in which
// case the caller should skip the file rather than guess.
func IsZeroValue(value any, present bool) (isZero bool, ok bool) {
	if !present {
		return true, true
	}
	switch v := value.(type) {
	case int:
		return v == 0, true
	case int64:
		return v == 0, true
	case float64:
		return v == 0, true
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return false, false
		}
		return f == 0, true
	default:
		return false, false
	}
}

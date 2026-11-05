package internal

import (
	"strconv"
	"strings"
	"time"
)

// ResolveDateArg resolves a date argument that is either an absolute date (2006-01-02) or an
// offset relative to today expressed as +N or -N days.
func ResolveDateArg(raw string) (string, error) {
	if relativeString, ok := strings.CutPrefix(raw, "-"); ok {
		offset, err := strconv.Atoi(relativeString)
		if err != nil {
			return "", err
		}
		return time.Now().AddDate(0, 0, offset*-1).Format(time.DateOnly), nil
	}
	if relativeString, ok := strings.CutPrefix(raw, "+"); ok {
		offset, err := strconv.Atoi(relativeString)
		if err != nil {
			return "", err
		}
		return time.Now().AddDate(0, 0, offset).Format(time.DateOnly), nil
	}
	return raw, nil
}

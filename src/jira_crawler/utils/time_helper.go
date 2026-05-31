package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ParsedTime struct {
	Time     time.Time
	TimeStr  string
	Epoch    int64
	EpochStr string
}

func ParseStringToTimeAndEpoch(value string) (ParsedTime, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return ParsedTime{}, nil
	}

	layouts := []string{
		"2006-01-02T15:04:05-0700",
		time.RFC3339,
		time.RFC3339Nano,
		"02/Jan/06 15:04",
		"02/Jan/2006 15:04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
	}

	var lastErr error

	for _, layout := range layouts {
		t, err := time.Parse(layout, value)
		if err == nil {
			return ParsedTime{
				Time:     t,
				TimeStr:  t.Format("2006-01-02T15:04:05-0700"),
				Epoch:    t.Unix(),
				EpochStr: strconv.FormatInt(t.Unix(), 10),
			}, nil
		}

		lastErr = err
	}

	return ParsedTime{}, fmt.Errorf("unsupported time format %q: %w", value, lastErr)
}

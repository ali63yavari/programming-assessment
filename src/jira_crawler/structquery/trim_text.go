package structquery

import (
	"fmt"
	"strings"
	"unicode"
)

type TrimType string

const (
	TrimNone    TrimType = "none"
	TrimAll     TrimType = "all"
	TrimSpace   TrimType = "space"
	TrimControl TrimType = "control"
)

func TrimText(trim TrimType, text string) (string, error) {
	switch trim {
	case TrimNone:
		return text, nil
	case TrimSpace:
		return strings.Map(trimSpace, text), nil
	case TrimControl:
		return strings.Map(trimControl, text), nil
	case TrimAll:
		return strings.Map(trimControl, strings.Map(trimSpace, text)), nil
	default:
		return text, fmt.Errorf("trim type '%s' not supported", trim)
	}
}

func trimSpace(r rune) rune {
	if unicode.IsSpace(r) {
		return -1
	}

	return r
}

func trimControl(r rune) rune {
	if unicode.IsControl(r) {
		return -1
	}

	return r
}

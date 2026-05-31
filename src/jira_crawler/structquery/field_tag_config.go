package structquery

import (
	"fmt"
	"reflect"
)

type ExtractMode string

const (
	ModeText      ExtractMode = "text"
	ModeAttr      ExtractMode = "attr"
	ModeHTML      ExtractMode = "html"
	ModeOuterHTML ExtractMode = "outer_html"
	ModeExists    ExtractMode = "exists"
	ModeCount     ExtractMode = "count"
)

func (em *ExtractMode) UnmarshalText(text []byte) error {
	v := ExtractMode(text)

	switch v {
	case ModeText, ModeAttr, ModeHTML, ModeOuterHTML, ModeExists, ModeCount:
		*em = v
		return nil
	default:
		return fmt.Errorf("invalid ExtractMode [%s]", text)
	}
}

type fieldTagConfig struct {
	Selector string
	Mode     ExtractMode
	Trim     TrimType
	Attr     string
	NonEmpty bool
	Required bool
	Each     bool
	Split    string
	Enums    []string
	Match    string
}

func newValidFieldTagConfig() *fieldTagConfig {
	return &fieldTagConfig{
		Selector: "",
		Mode:     ModeText,
		Trim:     TrimAll,
		Attr:     "",
		NonEmpty: false,
		Required: false,
	}
}

func (ftc fieldTagConfig) validate(fieldType reflect.Type) error {
	if ftc.Selector == "" {
		return fmt.Errorf("selector is required")
	}

	if !isValidMode(ftc.Mode) {
		return fmt.Errorf("invalid mode %q", ftc.Mode)
	}

	if ftc.Mode == ModeAttr && ftc.Attr == "" {
		return fmt.Errorf("attr is required when mode=attr")
	}

	if ftc.Mode != ModeAttr && ftc.Attr != "" {
		return fmt.Errorf("attr can only be used with mode=attr")
	}

	if ftc.Each && ftc.Split != "" {
		return fmt.Errorf("each and split cannot be used together")
	}

	if ftc.Each && fieldType.Kind() != reflect.Slice {
		return fmt.Errorf("each can only be used with slice fields")
	}

	if ftc.Split != "" && fieldType.Kind() != reflect.Slice {
		return fmt.Errorf("split can only be used with slice fields")
	}

	if ftc.Mode == ModeExists && fieldType.Kind() != reflect.Bool {
		return fmt.Errorf("mode=exists requires bool field")
	}

	if ftc.Mode == ModeCount && !isIntegerKind(fieldType.Kind()) {
		return fmt.Errorf("mode=count requires integer field")
	}

	if fieldType.Kind() == reflect.Slice {
		elemKind := fieldType.Elem().Kind()

		if elemKind == reflect.Struct {
			if ftc.Split != "" {
				return fmt.Errorf("split cannot be used with []struct")
			}

			if ftc.Mode == ModeCount || ftc.Mode == ModeExists {
				return fmt.Errorf("mode=%s cannot be used with []struct", ftc.Mode)
			}
		}
	}

	return nil
}

func isValidMode(mode ExtractMode) bool {
	switch mode {
	case ModeText, ModeAttr, ModeHTML, ModeOuterHTML, ModeExists, ModeCount:
		return true
	default:
		return false
	}
}

func isIntegerKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

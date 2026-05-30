package structquery

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const TagKey string = "sq"

func splitTagInlineConfig(tag string) []string {
	var parts []string
	var current strings.Builder

	var quote rune
	escaped := false

	for _, r := range tag {
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}

		if r == '\\' {
			escaped = true
			current.WriteRune(r)
			continue
		}

		if quote != 0 {
			current.WriteRune(r)

			if r == quote {
				quote = 0
			}

			continue
		}

		if r == '\'' || r == '"' {
			quote = r
			current.WriteRune(r)
			continue
		}

		if r == ',' || r == ';' {
			parts = append(parts, strings.TrimSpace(current.String()))
			current.Reset()
			continue
		}

		current.WriteRune(r)
	}

	if current.Len() > 0 {
		parts = append(parts, strings.TrimSpace(current.String()))
	}

	return parts
}

func tagContentIsValid(tagContent string) bool {
	if tagContent == "" || tagContent == "-" {
		return false
	}

	return true
}
func parseTagContent(tagContent string) (*FieldTagConfig, error) {
	if !tagContentIsValid(tagContent) {
		return nil, errors.New("invalid tag content")
	}
	rawConfigs := splitTagInlineConfig(tagContent)
	if len(rawConfigs) == 0 {
		return nil, errors.New("tag should include at least a valid rawConfig")
	}

	tc := FieldTagConfig{}
	for _, rawConfig := range rawConfigs {
		c := strings.Split(rawConfig, "=")
		if strings.TrimSpace(c[0]) == "selector" {
			tc.Selector = strings.TrimSpace(c[1])
			continue
		}

	}

	return &tc, nil
}

func parseEnum[TElmType any](s string) (TElmType, error) {
	var result TElmType

	quotedJSON := fmt.Sprintf("%q", s)

	err := json.Unmarshal([]byte(quotedJSON), &result)
	if err != nil {
		return result, err
	}
	return result, nil
}
func parseString[TElmType any](s string) (TElmType, error) {
	var result TElmType

	switch any(&result).(type) {
	case *int:
		val, err := strconv.Atoi(s)
		if err != nil {
			return result, err
		}
		*(any(&result).(*int)) = val

	case *float64:
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return result, err
		}
		*(any(&result).(*float64)) = val

	case *bool:
		val, err := strconv.ParseBool(s)
		if err != nil {
			return result, err
		}
		*(any(&result).(*bool)) = val

	case *string:
		*(any(&result).(*string)) = s

	default:
		//TODO: maybe in needs the parseEnum to e used when input is from Custom
		//ENUM type
		return result, fmt.Errorf("unsupported target type: %T", result)
	}

	return result, nil
}

func extractStringConfig[TElmType any](rawConfig, configName string) (
	TElmType,
	error,
) {
	var res TElmType
	if strings.TrimSpace(rawConfig) == "" {
		return res, fmt.Errorf("any field with name [%s] not found", configName)
	}

	c := strings.Split(rawConfig, "=")
	if len(c) != 2 {
		return res, fmt.Errorf(
			"[%s] config should be in form of [key=value]",
			configName,
		)
	}
	if len(strings.TrimSpace(c[1])) == 0 {
		//TODO: maybe it needs to investigated more
	}

	if strings.TrimSpace(c[0]) == strings.TrimSpace(configName) {
		v, err := parseString[TElmType](strings.TrimSpace(c[1]))
		if err != nil {
			return res, fmt.Errorf(
				"invalid value format:\n[%s]",
				err,
			)
		}

		return v, nil
	}

	return res, fmt.Errorf("any field with name [%s] not found", configName)
}

func extractBoolOrFlag(rawConfig, configName string) (bool, error) {
	c := strings.Split(rawConfig, "=")
	if strings.TrimSpace(c[0]) == strings.TrimSpace(configName) {
		if len(c) == 2 {
			if strings.TrimSpace(c[1]) == "false" {
				return false, nil
			}
		}

		return true, nil
	}

	return false, fmt.Errorf("any field with name [%s] not found", configName)
}

func ExtractInlineArray[TElmType any](
	rawConfig, configName,
	arraySeparator string,
) ([]TElmType, error) {
	var res []TElmType

	if len(strings.TrimSpace(arraySeparator)) == 0 {
		return res, fmt.Errorf(
			"invalid array separator character: [%s]",
			arraySeparator,
		)
	}

	c := strings.Split(rawConfig, "=")

	if len(c) != 2 {
		return res, fmt.Errorf(
			"[%s] config should be in form of [key=value]",
			configName,
		)
	}

	if len(strings.TrimSpace(c[1])) == 0 {
		return res, fmt.Errorf(
			"invalid value format: [%s]",
			c[1],
		)
	}

	sps := strings.Split(c[1], arraySeparator)
	for _, sp := range sps {
		v, err := parseString[TElmType](strings.TrimSpace(sp))
		if err != nil {
			continue
		}
		res = append(res, v)
	}

	return res, nil
}

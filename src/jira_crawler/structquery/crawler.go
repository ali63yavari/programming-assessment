package structquery

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Crawler interface {
	CrawlNow(output any) error
}
type crawler struct {
	url string
	doc *goquery.Document
}

func NewCrawlingPage(url string) (Crawler, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	crawlPage := &crawler{
		url: url,
		doc: doc,
	}

	return crawlPage, nil
}

func (c *crawler) CrawlNow(output any) error {
	return unmarshalStruct(c.doc.Selection, output)
}

func unmarshalStruct(root *goquery.Selection, output any) error {
	if root == nil {
		return fmt.Errorf("root node is null")
	}

	v := reflect.ValueOf(output)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return errors.New("output parameters is not a pointer or is null")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("output does not refer to a valid struct")
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		tagContent := fieldType.Tag.Get(TagKey)
		if !tagContentIsValid(tagContent) {
			continue
		}

		config, err := parseTagContent(tagContent)
		if err != nil {
			continue
		}
		if err := config.validate(fieldType.Type); err != nil {
			continue
		}

		selection := root.Find(config.Selector)

		if err := setField(fieldValue, config, selection); err != nil {
			return fmt.Errorf("field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

func setField(
	field reflect.Value,
	config *fieldTagConfig, selection *goquery.Selection,
) error {
	if !field.CanSet() {
		return fmt.Errorf("field cannot be set")
	}

	if config.Required && (selection == nil || selection.Length() == 0) {
		return fmt.Errorf(
			"html node `selection`: `%s` is null or empty",
			config.Selector,
		)
	}

	switch field.Kind() {
	case reflect.Slice:
		return setSliceField(field, config, selection)

	case reflect.Struct:
		if selection.Length() == 0 {
			if config.Required {
				return fmt.Errorf(
					"required selector %q did not match any element",
					config.Selector,
				)
			}
			return nil
		}

		return unmarshalStruct(selection.First(), field.Addr().Interface())

	default:
		return setSingleField(field, config, selection)
	}
}

func setSliceField(
	fieldValue reflect.Value,
	config *fieldTagConfig, root *goquery.Selection,
) error {
	elmType := fieldValue.Type().Elem()
	result := reflect.MakeSlice(fieldValue.Type(), 0, 5)

	sp := strings.TrimSpace(config.Split)
	if len(sp) != 0 {
		sps := strings.Split(root.Text(), sp)
		for _, item := range sps {
			raw := strings.TrimSpace(item)
			v, err := parseStringByType(raw, elmType)
			if err != nil {
				return err
			}
			result = reflect.Append(result, v)
		}
		fieldValue.Set(result)
		return nil
	}

	var firstErr error

	root.Each(
		func(i int, item *goquery.Selection) {
			if firstErr != nil {
				return
			}

			if elmType.Kind() == reflect.Struct {
				elem := reflect.New(elmType).Elem()

				err := unmarshalStruct(item, elem.Addr().Interface())
				if err != nil {
					firstErr = fmt.Errorf("item %d: %w", i, err)
					return
				}

				result = reflect.Append(result, elem)
				return
			}

			raw, err := extractFromSelection(item, config)
			if err != nil {
				firstErr = fmt.Errorf("item %d: %w", i, err)
				return
			}

			raw, err = TrimText(config.Trim, raw)
			if err != nil {
				firstErr = fmt.Errorf("item %d: %w", i, err)
				return
			}

			if err := validateExtractedValue(raw, config); err != nil {
				firstErr = fmt.Errorf("item %d: %w", i, err)
				return
			}

			converted, err := parseStringByType(raw, elmType)
			if err != nil {
				firstErr = fmt.Errorf("item %d: %w", i, err)
				return
			}

			result = reflect.Append(result, converted)
		},
	)

	if firstErr != nil {
		return firstErr
	}

	fieldValue.Set(result)
	return nil
}

func setSingleField(
	fieldValue reflect.Value,
	config *fieldTagConfig, root *goquery.Selection,
) error {
	if root.Length() == 0 {
		return nil
	}

	raw, err := extractFromSelection(root, config)
	if err != nil {
		return err
	}

	raw, err = TrimText(config.Trim, raw)
	if err != nil {
		return err
	}

	err = validateExtractedValue(raw, config)
	if err != nil {
		return err
	}

	finalV, err := parseStringByType(raw, fieldValue.Type())
	if err != nil {
		return err
	}

	fieldValue.Set(finalV)

	return nil
}

func outerHTML(selection *goquery.Selection) (string, error) {
	if selection == nil || selection.Length() == 0 {
		return "", nil
	}

	node := selection.Get(0)
	if node == nil {
		return "", nil
	}

	var buf bytes.Buffer

	if err := html.Render(&buf, node); err != nil {
		return "", fmt.Errorf("render outer html: %w", err)
	}

	return buf.String(), nil
}

func extractFromSelection(
	selection *goquery.Selection,
	config *fieldTagConfig,
) (string, error) {
	switch config.Mode {
	case ModeText:
		return selection.Text(), nil
	case ModeHTML:
		txt, err := selection.First().Html()
		if err != nil {
			return "", fmt.Errorf("html extraction failed: %w", err)
		}

		return txt, nil
	case ModeOuterHTML:
		return outerHTML(selection.First())
	case ModeAttr:
		value, exists := selection.Attr(config.Attr)
		if !exists {
			return "", nil
		}
		return value, nil

	default:
		return "", fmt.Errorf("unsupported extraction mode: %s", config.Mode)
	}
}

func validateExtractedValue(rawValue string, config *fieldTagConfig) error {
	s := strings.TrimSpace(rawValue)
	if config.NonEmpty {
		if len(s) == 0 {
			return fmt.Errorf(
				"config forced nonempty flag while extracted value" +
					" is empty",
			)
		}
	}

	if !slices.Contains(config.Enums, s) {
		return fmt.Errorf(
			"config forced enums [%q] but `%s` is not member of that",
			config.Enums,
			s,
		)
	}

	pattern := strings.TrimSpace(config.Match)
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return fmt.Errorf(
			"config forced regex match[%s] but evaluation failed: %w",
			pattern,
			err,
		)
	}

	if !matched {
		return fmt.Errorf(
			"config forced regex match but `%s` is not meet the pattern [%s]",
			s,
			pattern,
		)
	}

	return nil
}
func parseStringByType(s string, targetType reflect.Type) (reflect.Value, error) {
	value := reflect.New(targetType).Elem()

	switch targetType.Kind() {
	case reflect.Int:
		parsed, err := strconv.Atoi(s)
		if err != nil {
			return reflect.Value{}, err
		}
		value.SetInt(int64(parsed))
		return value, nil

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(s, 10, targetType.Bits())
		if err != nil {
			return reflect.Value{}, err
		}
		value.SetInt(parsed)
		return value, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsed, err := strconv.ParseUint(s, 10, targetType.Bits())
		if err != nil {
			return reflect.Value{}, err
		}
		value.SetUint(parsed)
		return value, nil

	case reflect.Float32, reflect.Float64:
		parsed, err := strconv.ParseFloat(s, targetType.Bits())
		if err != nil {
			return reflect.Value{}, err
		}
		value.SetFloat(parsed)
		return value, nil

	case reflect.Bool:
		parsed, err := strconv.ParseBool(s)
		if err != nil {
			return reflect.Value{}, err
		}
		value.SetBool(parsed)
		return value, nil

	case reflect.String:
		value.SetString(s)
		return value, nil

	default:
		return reflect.Value{}, fmt.Errorf("unsupported target type: %s", targetType)
	}
}

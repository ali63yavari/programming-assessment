package structquery

import (
	"errors"
	"reflect"

	"github.com/PuerkitoBio/goquery"
)

type Crawler interface {
	CrawlNow(output any) error
}
type crawler struct {
	url string
	doc *goquery.Document
}

func NewCrawlingPage(url string, doc *goquery.Document) Crawler {
	//it needs preparation and initialization of goquery for future works
	crawlPage := &crawler{
		url: url,
		doc: doc,
	}

	return crawlPage
}

func (c *crawler) CrawlNow(output any) error {
	v := reflect.ValueOf(output)

	if v.Kind() != reflect.Pointer || v.IsNil() {
		return errors.New("output parameters is not a pointer or is null")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("output does not refer to a valid struct")
	}

	t := v.Type()

	for i := 0; i <= v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		tagContent := fieldType.Tag.Get(TagKey)
		if !tagContentIsValid(tagContent) {
			continue
		}

		if !fieldValue.CanSet() {
			continue
		}

		config, err := parseTagContent(tagContent)
		if err != nil {
			continue
		}

		if fieldValue.Kind() == reflect.Slice {
			err = c.setSliceField(fieldValue, config, c.doc.Selection)
		} else {
			err = c.setSingleField(fieldValue, config, c.doc.Selection)
		}
	}

	return nil
}

func (c *crawler) setSliceField(
	fieldValue reflect.Value,
	config *FieldTagConfig, root *goquery.Selection,
) error {
	panic("not implemented")
}

func (c *crawler) setSingleField(
	fieldValue reflect.Value,
	config *FieldTagConfig, root *goquery.Selection,
) error {
	panic("not implemented")
}

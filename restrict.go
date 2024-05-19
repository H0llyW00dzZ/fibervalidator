// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package validator

import (
	"encoding/xml"
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// RestrictUnicode is a Restrictor implementation that restricts the use of Unicode characters
// in specified fields of the request body.
type RestrictUnicode struct {
	// Fields specifies the fields to check for Unicode characters.
	Fields []string
}

// Restrict implements the Restrictor interface for RestrictUnicode.
// It checks the specified fields in the request body for Unicode characters based on the content type.
func (r RestrictUnicode) Restrict(c *fiber.Ctx) error {
	contentType := c.Get(fiber.HeaderContentType)
	switch contentType {
	case fiber.MIMEApplicationJSON,
		fiber.MIMEApplicationJSONCharsetUTF8:
		return r.restrictJSON(c)
	case fiber.MIMEApplicationXML,
		fiber.MIMEApplicationXMLCharsetUTF8,
		fiber.MIMETextXML,
		fiber.MIMETextXMLCharsetUTF8:
		return r.restrictXML(c)
	default:
		return r.restrictOther(c)
	}
}

// restrictJSON checks the specified fields in the JSON request body for Unicode characters.
func (r RestrictUnicode) restrictJSON(c *fiber.Ctx) error {
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return NewError(fiber.StatusBadRequest, ErrInvalidJSONBody)
	}
	for _, field := range r.Fields {
		value, ok := body[field]
		if ok {
			if str, ok := value.(string); ok {
				if containsUnicode(str) {
					return NewError(fiber.StatusBadRequest, fmt.Sprintf(ErrUnicodeNotAllowedInField, field))
				}
			}
		}
	}
	return nil
}

// restrictXML checks the specified fields in the XML request body for Unicode characters.
func (r RestrictUnicode) restrictXML(c *fiber.Ctx) error {
	fields := make([]reflect.StructField, len(r.Fields))
	caser := cases.Title(language.English)
	for i, field := range r.Fields {
		fields[i] = reflect.StructField{
			Name: caser.String(field),
			Type: reflect.TypeOf(""),
			Tag:  reflect.StructTag(`xml:"` + field + `"`),
		}
	}
	bodyType := reflect.StructOf(fields)
	bodyValue := reflect.New(bodyType).Elem()

	if err := xml.Unmarshal(c.Body(), bodyValue.Addr().Interface()); err != nil {
		return NewError(fiber.StatusBadRequest, ErrInvalidXMLBody)
	}

	for _, field := range r.Fields {
		value := bodyValue.FieldByName(caser.String(field)).String()
		if containsUnicode(value) {
			return NewError(fiber.StatusBadRequest, fmt.Sprintf(ErrUnicodeNotAllowedInField, field))
		}
	}
	return nil
}

// restrictOther checks the specified fields in the request body of other content types for Unicode characters.
func (r RestrictUnicode) restrictOther(c *fiber.Ctx) error {
	body := string(c.Body())
	for _, field := range r.Fields {
		fieldValue := extractFieldValue(body, field, r)
		if containsUnicode(fieldValue) {
			return NewError(fiber.StatusBadRequest, fmt.Sprintf(ErrUnicodeNotAllowedInField, field))
		}
	}
	return nil
}

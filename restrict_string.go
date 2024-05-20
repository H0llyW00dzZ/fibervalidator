// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package validator

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// RestrictStringLength is a Restrictor implementation that restricts the length of string fields
// and allows setting an optional maximum length.
type RestrictStringLength struct {
	// Fields specifies the fields to check for string length validation.
	Fields []string

	// MaxLength specifies the maximum allowed length for the fields (optional).
	MaxLength *int
}

// Restrict implements the Restrictor interface for RestrictStringLength.
// It checks the specified fields in the request body for string length based on the content type.
func (r RestrictStringLength) Restrict(c *fiber.Ctx) error {
	return restrictByContentType(c, r.restrictJSON, r.restrictXML, r.restrictOther)
}

// restrictJSON checks the specified fields in the JSON request body for string length and maximum limit.
func (r RestrictStringLength) restrictJSON(c *fiber.Ctx) error {
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return NewError(fiber.StatusBadRequest, ErrInvalidJSONBody)
	}

	var invalidFields []string
	for _, field := range r.Fields {
		value, ok := body[field]
		if ok {
			if str, ok := value.(string); ok {
				if r.MaxLength != nil && len(str) > *r.MaxLength {
					return NewError(fiber.StatusBadRequest, fmt.Sprintf(ErrFieldExceedsMaximumLength, field, *r.MaxLength))
				}
			}
		}
	}

	if len(invalidFields) > 0 {
		return NewError(fiber.StatusBadRequest, fmt.Sprintf(ErrFieldsExceedMaximumLength, strings.Join(invalidFields, "', '")))
	}

	return nil
}

// restrictXML checks the specified fields in the XML request body for string length and maximum limit.
func (r RestrictStringLength) restrictXML(c *fiber.Ctx) error {
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

	var invalidFields []string
	for _, field := range r.Fields {
		value := bodyValue.FieldByName(caser.String(field)).String()
		if r.MaxLength != nil && len(value) > *r.MaxLength {
			return NewError(fiber.StatusBadRequest, fmt.Sprintf(ErrFieldExceedsMaximumLength, field, *r.MaxLength))
		}
	}

	if len(invalidFields) > 0 {
		return NewError(fiber.StatusBadRequest, fmt.Sprintf(ErrFieldsExceedMaximumLength, strings.Join(invalidFields, "', '")))
	}

	return nil
}

// restrictOther checks the specified fields in the request body of other content types for string length and maximum limit.
func (r RestrictStringLength) restrictOther(c *fiber.Ctx) error {
	body := string(c.Body())

	var invalidFields []string
	for _, field := range r.Fields {
		fieldValue := extractFieldValue(body, field, RestrictUnicode{Fields: r.Fields})
		if r.MaxLength != nil && len(fieldValue) > *r.MaxLength {
			return NewError(fiber.StatusBadRequest, fmt.Sprintf(ErrFieldExceedsMaximumLength, field, *r.MaxLength))
		}
	}

	if len(invalidFields) > 0 {
		return NewError(fiber.StatusBadRequest, fmt.Sprintf(ErrFieldsExceedMaximumLength, strings.Join(invalidFields, "', '")))
	}

	return nil
}

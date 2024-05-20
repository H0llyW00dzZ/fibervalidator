// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package validator

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// RestrictNumberOnly is a Restrictor implementation that restricts fields to contain only numbers
// and allows setting an optional maximum value and maximum number of digits.
type RestrictNumberOnly struct {
	// Fields specifies the fields to check for number-only validation.
	Fields []string

	// Max specifies the maximum allowed value for the fields (optional).
	Max *int

	// MaxDigits specifies the maximum number of digits allowed in the field value (optional).
	MaxDigits *int
}

// Restrict implements the Restrictor interface for RestrictNumberOnly.
// It checks the specified fields in the request body for numeric values and maximum limit based on the content type.
func (r RestrictNumberOnly) Restrict(c *fiber.Ctx) error {
	return restrictByContentType(c, r.restrictJSON, r.restrictXML, r.restrictOther)
}

// restrictJSON checks the specified fields in the JSON request body for numeric values and maximum limit.
func (r RestrictNumberOnly) restrictJSON(c *fiber.Ctx) error {
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return NewError(fiber.StatusBadRequest, ErrInvalidJSONBody)
	}

	var invalidFields []string
	for _, field := range r.Fields {
		value, ok := body[field]
		if ok {
			var num int
			var numStr string
			switch v := value.(type) {
			case string:
				if !isNumberOnly(v) {
					invalidFields = append(invalidFields, field)
					continue
				}
				numStr = v
				num, _ = strconv.Atoi(v)
			case float64:
				num = int(v)
				numStr = strconv.Itoa(num)
			default:
				invalidFields = append(invalidFields, field)
				continue
			}
			if r.MaxDigits != nil && len(numStr) > *r.MaxDigits {
				return NewError(fiber.StatusBadRequest, fmt.Sprintf(ErrFieldExceedsMaximumDigits, field, *r.MaxDigits))
			}
			if r.Max != nil && num > *r.Max {
				return NewError(fiber.StatusBadRequest, fmt.Sprintf(ErrFieldExceedsMaximumValue, field, *r.Max))
			}
		}
	}

	if len(invalidFields) > 0 {
		return NewError(fiber.StatusBadRequest, fmt.Sprintf(ErrFieldMustContainNumbersOnly, strings.Join(invalidFields, "', '")))
	}

	return nil
}

// restrictXML checks the specified fields in the XML request body for numeric values and maximum limit.
func (r RestrictNumberOnly) restrictXML(c *fiber.Ctx) error {
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
		if !isNumberOnly(value) {
			invalidFields = append(invalidFields, field)
		} else {
			num, _ := strconv.Atoi(value)
			if r.MaxDigits != nil && len(value) > *r.MaxDigits {
				return NewError(fiber.StatusBadRequest, fmt.Sprintf(ErrFieldExceedsMaximumDigits, field, *r.MaxDigits))
			}
			if r.Max != nil && num > *r.Max {
				return NewError(fiber.StatusBadRequest, fmt.Sprintf(ErrFieldExceedsMaximumValue, field, *r.Max))
			}
		}
	}

	if len(invalidFields) > 0 {
		return NewError(fiber.StatusBadRequest, fmt.Sprintf(ErrFieldMustContainNumbersOnly, strings.Join(invalidFields, "', '")))
	}

	return nil
}

// restrictOther checks the specified fields in the request body of other content types for numeric values and maximum limit.
func (r RestrictNumberOnly) restrictOther(c *fiber.Ctx) error {
	body := string(c.Body())

	var invalidFields []string
	for _, field := range r.Fields {
		fieldValue := extractFieldValueForNumberOnly(body, field)
		if !isNumberOnly(fieldValue) {
			invalidFields = append(invalidFields, field)
		} else {
			num, _ := strconv.Atoi(fieldValue)
			if r.MaxDigits != nil && len(fieldValue) > *r.MaxDigits {
				return NewError(fiber.StatusBadRequest, fmt.Sprintf(ErrFieldExceedsMaximumDigits, field, *r.MaxDigits))
			}
			if r.Max != nil && num > *r.Max {
				return NewError(fiber.StatusBadRequest, fmt.Sprintf(ErrFieldExceedsMaximumValue, field, *r.Max))
			}
		}
	}

	if len(invalidFields) > 0 {
		return NewError(fiber.StatusBadRequest, fmt.Sprintf(ErrFieldMustContainNumbersOnly, strings.Join(invalidFields, "', '")))
	}

	return nil
}

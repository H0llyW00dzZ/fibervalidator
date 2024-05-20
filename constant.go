// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package validator

const (
	// ErrInvalidJSONBody represents an error message for an invalid JSON request body.
	ErrInvalidJSONBody = "Invalid JSON request body"

	// ErrUnicodeNotAllowedInField represents an error message for Unicode characters not allowed in a specific field.
	ErrUnicodeNotAllowedInField = "Unicode characters are not allowed in the '%s' field"

	// ErrInvalidXMLBody represents an error message for an invalid XML request body.
	ErrInvalidXMLBody = "Invalid XML request body"
)

const (
	// ErrFieldMustContainNumbersOnly represents an error message for a field that must contain only numbers.
	ErrFieldMustContainNumbersOnly = "The '%s' field must contain only numbers"

	// ErrFieldExceedsMaximumValue represents an error message for a field that exceeds the maximum allowed value.
	ErrFieldExceedsMaximumValue = "The '%s' field must not exceed %d"
)

const (
	// Define the range of numeric characters
	numericStart = '0' + iota
	numericEnd   = '9'
)

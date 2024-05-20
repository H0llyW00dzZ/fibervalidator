// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package validator

import (
	"strings"
)

// containsUnicode checks if a string contains Unicode characters.
// It returns true if any character is outside the ASCII range.
//
// Note: This is a slick way to boost performance, leaving regex in the dust.
func containsUnicode(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] > 127 { // ASCII range is from 0 to 127
			return true
		}
	}
	return false
}

// extractFieldValue extracts the value of a specified field from the request body string.
func extractFieldValue(body, field string, r RestrictUnicode) string {
	// Perform a case-insensitive search for the field name
	fieldLower := strings.ToLower(field)
	bodyLower := strings.ToLower(body)

	// Find the position of the field name in the body
	startPos := strings.Index(bodyLower, fieldLower)
	if startPos == -1 {
		return ""
	}

	// Find the position of the value start (after the field name and any whitespace/separators)
	valueStartPos := startPos + len(fieldLower)
	for valueStartPos < len(body) && (body[valueStartPos] == ' ' || body[valueStartPos] == ':' || body[valueStartPos] == '=') {
		valueStartPos++
	}

	// Find the position of the value end (end of the body or next field)
	valueEndPos := len(body)
	for _, f := range r.Fields {
		if f != field {
			endPos := strings.Index(bodyLower[valueStartPos:], strings.ToLower(f))
			if endPos != -1 {
				valueEndPos = valueStartPos + endPos
				break
			}
		}
	}

	return strings.TrimSpace(body[valueStartPos:valueEndPos])
}

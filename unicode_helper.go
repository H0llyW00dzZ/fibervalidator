// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package validator

import (
	"regexp"
	"strings"
)

// containsUnicode checks if a string contains Unicode characters.
func containsUnicode(str string) bool {
	// Regular expression pattern to match Unicode characters
	pattern := "[^\\x00-\\x7F]"
	matched, _ := regexp.MatchString(pattern, str)
	return matched
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

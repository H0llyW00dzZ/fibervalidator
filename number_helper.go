// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package validator

import "strings"

// extractFieldValueForNumberOnly extracts the value of a specified field from the request body string.
func extractFieldValueForNumberOnly(body, field string) string {
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
	for i := valueStartPos; i < len(body); i++ {
		if body[i] == '&' || body[i] == '\n' || body[i] == '\r' {
			valueEndPos = i
			break
		}
	}

	return strings.TrimSpace(body[valueStartPos:valueEndPos])
}

// isNumberOnly checks if a string contains only numeric characters.
func isNumberOnly(str string) bool {
	for _, char := range str {
		if !isNumericChar(char) {
			return false
		}
	}
	return true
}

// isNumericChar checks if a rune is a numeric character.
func isNumericChar(char rune) bool {
	return char >= numericStart && char <= numericEnd
}

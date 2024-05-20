// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package validator

import "github.com/gofiber/fiber/v2"

// restrictByContentType is a helper function that determines the content type and calls the appropriate restrict function.
func restrictByContentType(c *fiber.Ctx, restrictJSON, restrictXML, restrictOther func(c *fiber.Ctx) error) error {
	contentType := c.Get(fiber.HeaderContentType)
	switch contentType {
	case fiber.MIMEApplicationJSON,
		fiber.MIMEApplicationJSONCharsetUTF8:
		return restrictJSON(c)
	case fiber.MIMEApplicationXML,
		fiber.MIMEApplicationXMLCharsetUTF8,
		fiber.MIMETextXML,
		fiber.MIMETextXMLCharsetUTF8:
		return restrictXML(c)
	default:
		return restrictOther(c)
	}
}

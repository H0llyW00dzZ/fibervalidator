// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package validator

import "github.com/gofiber/fiber/v2"

// Error represents a validation error.
type Error struct {
	Status  int
	Message string
}

// NewError creates a new Error instance.
func NewError(status int, message string) *Error {
	return &Error{
		Status:  status,
		Message: message,
	}
}

// Error implements the error interface for Error.
func (e *Error) Error() string {
	return e.Message
}

// DefaultErrorHandler is the default error handler function.
func DefaultErrorHandler(c *fiber.Ctx, err error) error {
	if e, ok := err.(*Error); ok {
		return restrictByContentType(c, jsonErrorHandler(e), xmlErrorHandler(e), defaultErrorHandler(e))
	}
	return err
}

// jsonErrorHandler formats the error as JSON.
func jsonErrorHandler(e *Error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Status(e.Status).JSON(fiber.Map{
			"error": e.Message,
		})
	}
}

// xmlErrorHandler formats the error as XML.
type xmlError struct {
	Error string `xml:"error"`
}

func xmlErrorHandler(e *Error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Status(e.Status).XML(xmlError{Error: e.Message})
	}
}

// defaultErrorHandler sends the error as plain text.
func defaultErrorHandler(e *Error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Status(e.Status).SendString(e.Message)
	}
}

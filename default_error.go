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
		return c.Status(e.Status).SendString(e.Message)
	}
	return err
}

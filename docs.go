// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

// Package validator provides a custom validator middleware for the Fiber web framework.
// It allows for flexible and extensible validation of incoming request bodies.
//
// # Installation
//
// To use this middleware in a Fiber project, Go must be installed and set up.
//
// 1. Install the package using Go modules:
//
//	go get github.com/H0llyW00dzZ/FiberValidator
//
// 2. Import the package in the Fiber application:
//
//	import "github.com/H0llyW00dzZ/FiberValidator"
//
// # Usage
//
// To use the validator middleware in a Fiber application, create a new instance of the middleware with the desired configuration and register it with the application.
//
//	package main
//
//	import (
//		"github.com/gofiber/fiber/v2"
//		"github.com/H0llyW00dzZ/FiberValidator"
//	)
//
//	func main() {
//		app := fiber.New()
//
//		app.Use(validator.New(validator.Config{
//			Rules: []validator.Restrictor{
//				validator.RestrictUnicode{
//					Fields: []string{"name", "email"},
//				},
//			},
//		}))
//
//		// Register routes and start the server
//		// ...
//
//		app.Listen(":3000")
//	}
//
// In the example above, a new instance of the validator middleware is created with a configuration that restricts the use of Unicode characters in the "name" and "email" fields of the request body.
//
// # Configuration
//
// The validator middleware accepts a [validator.Config] struct for configuration. The available options are:
//
//   - Rules: A slice of [validator.Restrictor] implementations that define the validation rules to be applied.
//   - Next: An optional function that determines whether to skip the validation middleware for a given request. If the function returns true, the middleware will be skipped.
//   - ErrorHandler: An optional custom error handler function that handles the error response. If not provided, the default error handler will be used.
//
// # Custom Validation Rules
//
// To define custom validation rules, implement the [validator.Restrictor] interface:
//
//	type Restrictor interface {
//		Restrict(c *fiber.Ctx) error
//	}
//
// The Restrict method takes a Fiber context and returns an error if the validation fails.
//
// Here's an example implementation that restricts the use of Unicode characters in specified fields:
//
//	type RestrictUnicode struct {
//		Fields []string
//	}
//
//	func (r RestrictUnicode) Restrict(c *fiber.Ctx) error {
//		// Parse the request body based on the Content-Type header
//		// ...
//
//		// Check each configured field for Unicode characters
//		// ...
//
//		return nil
//	}
//
// # Error Handling
//
// The validator middleware provides a default error handler that formats the error response based on the content type of the request. It supports JSON, XML, and plain text formats.
//
//   - For JSON requests, the error response is formatted as {"error": "Error message"}.
//   - For XML requests, the error response is formatted as <xmlError><error>Error message</error></xmlError>.
//   - For other content types, the error response is sent as plain text.
//
// You can customize the error handling behavior by providing a custom error handler function in the ErrorHandler field of the [validator.Config] struct. The custom error handler should have the following signature:
//
//	func(c *fiber.Ctx, err error) error
//
// Example custom error handler:
//
//	func customErrorHandler(c *fiber.Ctx, err error) error {
//		if e, ok := err.(*validator.Error); ok {
//			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
//				"custom_error": e.Message,
//			})
//		}
//		return err
//	}
//
// In this example, the custom error handler checks if the error is of type [*validator.Error] and returns a JSON response with a custom error format.
package validator

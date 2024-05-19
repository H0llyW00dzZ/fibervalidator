// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package validator

import "github.com/gofiber/fiber/v2"

// Restrictor is an interface for defining custom validation rules.
type Restrictor interface {
	Restrict(c *fiber.Ctx) error
}

// Config defines the configuration for the Validator middleware.
type Config struct {
	// Rules is a slice of Restrictor implementations to be used for validation.
	Rules []Restrictor

	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// ErrorHandler is a function that handles the error response.
	//
	// Optional. Default: DefaultErrorHandler
	ErrorHandler func(c *fiber.Ctx, err error) error
}

// ConfigDefault is the default configuration for the Validator middleware.
var ConfigDefault = Config{
	Rules:        nil,
	Next:         nil,
	ErrorHandler: DefaultErrorHandler,
}

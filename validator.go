// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package validator

import (
	"github.com/gofiber/fiber/v2"
)

// New creates a new Validator middleware with the provided configuration.
func New(config ...Config) fiber.Handler {
	cfg := ConfigDefault

	if len(config) > 0 {
		cfg = config[0]

		if cfg.ErrorHandler == nil {
			cfg.ErrorHandler = ConfigDefault.ErrorHandler
		}
	}

	return func(c *fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		for _, rule := range cfg.Rules {
			if err := rule.Restrict(c); err != nil {
				if cfg.ContextKey != "" {
					c.Locals(cfg.ContextKey, err)
				}
				return cfg.ErrorHandler(c, err)
			}
		}

		if cfg.ContextKey != "" {
			c.Locals(cfg.ContextKey, nil)
		}

		return c.Next()
	}
}

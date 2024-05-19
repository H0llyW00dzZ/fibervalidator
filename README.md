# Fiber Validator Middleware

This is a custom validator middleware for the Fiber web framework. It provides a flexible and extensible way to define and apply validation rules to incoming request bodies. The middleware allows for easy validation and sanitization of data, enforcement of specific field requirements, and ensures the integrity of the application's input.

The middleware currently supports the following features:
- Validation of request bodies in JSON and XML formats
- Restriction of Unicode characters in specified fields
- Customizable error handling
- Conditional validation skipping based on custom logic

More features and validation capabilities will be added in the future to enhance the middleware's functionality and cater to a wider range of validation scenarios.

## Installation

To use this middleware in a Fiber project, Go must be installed and set up.

1. Install the package using Go modules:

```shell
go get github.com/H0llyW00dzZ/FiberValidator
```

2. Import the package in the Fiber application:

```go
import "github.com/H0llyW00dzZ/FiberValidator"
```

## Usage

To use the validator middleware in a Fiber application, create a new instance of the middleware with the desired configuration and register it with the application.

```go
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/H0llyW00dzZ/FiberValidator"
)

func main() {
	app := fiber.New()

	app.Use(validator.New(validator.Config{
		Rules: []validator.Restrictor{
			validator.RestrictUnicode{
				Fields: []string{"name", "email"},
			},
		},
	}))

	// Register routes and start the server
	// ...

	app.Listen(":3000")
}
```

In the example above, a new instance of the validator middleware is created with a configuration that restricts the use of Unicode characters in the "name" and "email" fields of the request body.

### Configuration

The validator middleware accepts a `validator.Config` struct for configuration. The available options are:

- `Rules`: A slice of `validator.Restrictor` implementations that define the validation rules to be applied.
- `Next`: An optional function that determines whether to skip the validation middleware for a given request. If the function returns `true`, the middleware will be skipped.
- `ErrorHandler`: An optional custom error handler function that handles the error response. If not provided, the default error handler will be used.

### Custom Validation Rules

To define custom validation rules, implement the `validator.Restrictor` interface:

```go
type Restrictor interface {
	Restrict(c *fiber.Ctx) error
}
```

The `Restrict` method takes a Fiber context and returns an error if the validation fails.

Here's an example implementation that restricts the use of Unicode characters in specified fields:

```go
type RestrictUnicode struct {
	Fields []string
}

func (r RestrictUnicode) Restrict(c *fiber.Ctx) error {
	// Parse the request body based on the Content-Type header
	// ...

	// Check each configured field for Unicode characters
	// ...

	return nil
}
```

### Error Handling

By default, the validator middleware uses a default error handler that returns a JSON response with an "error" field containing the error message.

To provide a custom error handler function to handle the error response differently, the custom error handler should have the following signature:

```go
func(c *fiber.Ctx, err error) error
```

Example custom error handler:

```go
func customErrorHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
}
```

## Contributing

Contributions are welcome! If there are any issues or suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the BSD License. See the [LICENSE](LICENSE) file for details.

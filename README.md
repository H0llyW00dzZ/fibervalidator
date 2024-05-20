# Fiber Validator Middleware
[![Go Version](https://img.shields.io/badge/1.22.3-gray?style=flat&logo=go&logoWidth=15)](https://github.com/H0llyW00dzZ/FiberValidator/blob/master/go.mod#L3blob/master/go.mod#L3)
[![Go Reference](https://pkg.go.dev/badge/github.com/H0llyW00dzZ/FiberValidator.svg)](https://pkg.go.dev/github.com/H0llyW00dzZ/FiberValidator) [![Go Report Card](https://goreportcard.com/badge/github.com/H0llyW00dzZ/FiberValidator)](https://goreportcard.com/report/github.com/H0llyW00dzZ/FiberValidator)

This is a custom validator middleware for the Fiber web framework. It provides a flexible and extensible way to define and apply validation rules to incoming request bodies. The middleware allows for easy validation and sanitization of data, enforcement of specific field requirements, and ensures the integrity of the application's input.

The middleware currently supports the following features:
- Validation of request bodies in various formats, including JSON, XML, and other content types
- Restriction of Unicode characters in specified fields
- Customizable error handling based on content type
- Conditional validation skipping based on custom logic
- Restriction of fields to contain only numbers with an optional maximum value
- Storing validation results in the request context for advanced use cases
- Restriction of string length for specified fields with a configurable maximum length

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

The validator middleware provides a default error handler that formats the error response based on the content type of the request. It supports JSON, XML, and plain text formats.

- For JSON requests, the error response is formatted as `{"error": "Error message"}`.
- For XML requests, the error response is formatted as `<xmlError><error>Error message</error></xmlError>`.
- For other content types, the error response is sent as plain text.

You can customize the error handling behavior by providing a custom error handler function in the `ErrorHandler` field of the `validator.Config` struct. The custom error handler should have the following signature:

```go
func(c *fiber.Ctx, err error) error
```

Example custom error handler:

```go
func customErrorHandler(c *fiber.Ctx, err error) error {
	if e, ok := err.(*validator.Error); ok {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"custom_error": e.Message,
		})
	}
	return err
}
```

In this example, the custom error handler checks if the error is of type `*validator.Error` and returns a JSON response with a custom error format.

## Contributing

Contributions are welcome! If there are any issues or suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the BSD License. See the [LICENSE](LICENSE) file for details.


## Benchmark

```sh
goos: windows
goarch: amd64
pkg: github.com/H0llyW00dzZ/FiberValidator
cpu: AMD Ryzen 9 3900X 12-Core Processor            
BenchmarkValidatorWithSonicJSON/Valid_JSON_request-24         	   45967	     24768 ns/op	   16464 B/op	      86 allocs/op
BenchmarkValidatorWithStandardJSON/Valid_JSON_request-24      	   43248	     27835 ns/op	   16624 B/op	     112 allocs/op
BenchmarkValidatorWithDefaultXML/Valid_XML_request-24         	   28101	     42913 ns/op	   23223 B/op	     212 allocs/op
BenchmarkValidatorWithCustomXML/Valid_XML_request-24          	   28191	     43596 ns/op	   23248 B/op	     212 allocs/op
```

#### Updated:

```sh
goos: windows
goarch: amd64
pkg: github.com/H0llyW00dzZ/FiberValidator
cpu: AMD Ryzen 9 3900X 12-Core Processor            
BenchmarkValidatorWithSonicJSONSeafood/Valid_JSON_request-24         	   46785	     24696 ns/op	   16447 B/op	      86 allocs/op
BenchmarkValidatorWithStandardJSONSeafood/Valid_JSON_request-24      	   42541	     28542 ns/op	   16672 B/op	     112 allocs/op
BenchmarkValidatorWithDefaultXMLSeafood/Valid_XML_request-24         	   26637	     44806 ns/op	   23450 B/op	     213 allocs/op
BenchmarkValidatorWithCustomXMLSeafood/Valid_XML_request-24          	   26622	     45684 ns/op	   23458 B/op	     213 allocs/op
BenchmarkValidatorWithSonicJSON/Valid_JSON_request-24                	   50625	     24377 ns/op	   16410 B/op	      86 allocs/op
BenchmarkValidatorWithStandardJSON/Valid_JSON_request-24             	   42150	     27954 ns/op	   16626 B/op	     112 allocs/op
BenchmarkValidatorWithDefaultXML/Valid_XML_request-24                	   27764	     43721 ns/op	   23244 B/op	     212 allocs/op
BenchmarkValidatorWithCustomXML/Valid_XML_request-24                 	   27417	     43951 ns/op	   23256 B/op	     212 allocs/op
```

```sh
goos: windows
goarch: amd64
pkg: github.com/H0llyW00dzZ/FiberValidator
cpu: AMD Ryzen 9 3900X 12-Core Processor            
BenchmarkRestrictStringLengthLongDescriptionSonicJSON/Valid_JSON_long_description-24         	   55225	     23263 ns/op	   19344 B/op	      58 allocs/op
BenchmarkRestrictStringLengthLongDescriptionStandardJSON/Valid_JSON_long_description-24      	   48288	     24452 ns/op	   19234 B/op	      65 allocs/op
BenchmarkRestrictStringLengthLongDescriptionDefaultXML/Valid_XML_long_description-24         	   30114	     39411 ns/op	   25018 B/op	     111 allocs/op
BenchmarkRestrictStringLengthLongDescriptionCustomXML/Valid_XML_long_description-24          	   30322	     40180 ns/op	   25025 B/op	     111 allocs/op
BenchmarkValidatorWithSonicJSONSeafood/Valid_JSON_request-24                                 	   51158	     23802 ns/op	   16421 B/op	      86 allocs/op
BenchmarkValidatorWithStandardJSONSeafood/Valid_JSON_request-24                              	   43411	     27269 ns/op	   16652 B/op	     112 allocs/op
BenchmarkValidatorWithDefaultXMLSeafood/Valid_XML_request-24                                 	   26984	     45736 ns/op	   23451 B/op	     213 allocs/op
BenchmarkValidatorWithCustomXMLSeafood/Valid_XML_request-24                                  	   26899	     45503 ns/op	   23443 B/op	     213 allocs/op
BenchmarkValidatorWithSonicJSON/Valid_JSON_request-24                                        	   50936	     24697 ns/op	   16417 B/op	      86 allocs/op
BenchmarkValidatorWithStandardJSON/Valid_JSON_request-24                                     	   43172	     28678 ns/op	   16627 B/op	     112 allocs/op
BenchmarkValidatorWithDefaultXML/Valid_XML_request-24                                        	   27283	     43998 ns/op	   23262 B/op	     212 allocs/op
BenchmarkValidatorWithCustomXML/Valid_XML_request-24                                         	   27990	     43546 ns/op	   23264 B/op	     212 allocs/op
```

> [!NOTE]
> Based on the benchmark results, the following observations can be made:
>
> - The Sonic JSON encoder/decoder consistently outperforms the standard JSON encoder/decoder in terms of execution time, bytes allocated, and allocations per operation.
> - The custom XML encoder/decoder performs slightly slower than the default XML encoder/decoder in most cases, with similar memory usage.
> - The JSON benchmarks generally have better performance compared to the XML benchmarks, with lower execution time, bytes allocated, and allocations per operation.
>
> <p align="center">
>   <img src="https://i.imgur.com/PxjZ0Dz.png" alt="gopher run" />
> </p>
>
> Overall, the benchmarks using the Sonic JSON encoder/decoder demonstrate the best performance among the tested scenarios.

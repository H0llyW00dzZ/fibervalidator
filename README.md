# Fiber Validator Middleware
[![Go Version](https://img.shields.io/badge/1.22.3-gray?style=flat&logo=go&logoWidth=15)](https://github.com/H0llyW00dzZ/FiberValidator/blob/master/go.mod#L3blob/master/go.mod#L3)
[![Go Reference](https://pkg.go.dev/badge/github.com/H0llyW00dzZ/FiberValidator.svg)](https://pkg.go.dev/github.com/H0llyW00dzZ/FiberValidator) [![Go Report Card](https://goreportcard.com/badge/github.com/H0llyW00dzZ/FiberValidator)](https://goreportcard.com/report/github.com/H0llyW00dzZ/FiberValidator)

This is a custom validator middleware for the Fiber web framework. It provides a flexible and extensible way to define and apply validation rules to incoming request bodies. The middleware allows for easy validation and sanitization of data, enforcement of specific field requirements, and ensures the integrity of the application's input.

## Features

The middleware currently supports the following features:

### Request Body Validation
- Validation of request bodies in various formats, including JSON, XML, and other content types
- Customizable error handling based on content type

### Unicode Restriction
- Restriction of Unicode characters in specified fields

### Conditional Validation
- Conditional validation skipping based on custom logic

### Number Restriction
- Restriction of fields to contain only numbers with an optional maximum value

### String Length Restriction
- Restriction of string length for specified fields with a configurable maximum length

### Advanced Use Cases
- Storing validation results in the request context for advanced use cases

More features and validation capabilities will be added in the future to enhance the middleware's functionality and cater to a wider range of validation scenarios.

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

```sh
goos: windows
goarch: amd64
pkg: github.com/H0llyW00dzZ/FiberValidator
cpu: AMD Ryzen 9 3900X 12-Core Processor            
BenchmarkRestrictStringLengthLongDescriptionSonicJSON/Valid_JSON_long_description-24         	   51478	     21306 ns/op	   19331 B/op	      58 allocs/op
BenchmarkRestrictStringLengthLongDescriptionStandardJSON/Valid_JSON_long_description-24      	   47730	     25044 ns/op	   19240 B/op	      65 allocs/op
BenchmarkRestrictStringLengthLongDescriptionDefaultXML/Valid_XML_long_description-24         	   30450	     39450 ns/op	   25025 B/op	     111 allocs/op
BenchmarkRestrictStringLengthLongDescriptionCustomXML/Valid_XML_long_description-24          	   29910	     40017 ns/op	   25025 B/op	     111 allocs/op
BenchmarkValidatorWithSonicJSONSeafood/Valid_JSON_request-24                                 	   60594	     20457 ns/op	   14826 B/op	      68 allocs/op
BenchmarkValidatorWithStandardJSONSeafood/Valid_JSON_request-24                              	   50836	     23481 ns/op	   15081 B/op	      94 allocs/op
BenchmarkValidatorWithDefaultXMLSeafood/Valid_XML_request-24                                 	   30085	     40004 ns/op	   21815 B/op	     195 allocs/op
BenchmarkValidatorWithCustomXMLSeafood/Valid_XML_request-24                                  	   29608	     40108 ns/op	   21819 B/op	     195 allocs/op
BenchmarkValidatorWithSonicJSON/Valid_JSON_request-24                                        	   61047	     20826 ns/op	   14801 B/op	      68 allocs/op
BenchmarkValidatorWithStandardJSON/Valid_JSON_request-24                                     	   51811	     23220 ns/op	   15045 B/op	      94 allocs/op
BenchmarkValidatorWithDefaultXML/Valid_XML_request-24                                        	   31850	     38777 ns/op	   21632 B/op	     194 allocs/op
BenchmarkValidatorWithCustomXML/Valid_XML_request-24                                         	   30577	     39133 ns/op	   21633 B/op	     194 allocs/op
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

## Contributing

Contributions are welcome! If there are any issues or suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the BSD License. See the [LICENSE](LICENSE) file for details.

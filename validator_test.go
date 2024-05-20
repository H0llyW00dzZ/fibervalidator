// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package validator_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	validator "github.com/H0llyW00dzZ/FiberValidator"

	"github.com/gofiber/fiber/v2"
)

func TestValidatorWithDefaultErrorHandler(t *testing.T) {
	app := fiber.New()

	app.Use(validator.New(validator.Config{
		Rules: []validator.Restrictor{
			validator.RestrictUnicode{
				Fields: []string{"name", "email"},
			},
		},
	}))

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	testCases := []struct {
		name           string
		contentType    string
		requestBody    string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Valid JSON",
			contentType:    "application/json",
			requestBody:    `{"name":"gopher","email":"gopher@example.com"}`,
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "Invalid JSON - Unicode in name",
			contentType:    "application/json",
			requestBody:    `{"name":"Gøpher","email":"gopher@example.com"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Unicode characters are not allowed in the 'name' field",
		},
		{
			name:           "Invalid JSON - Unicode in email",
			contentType:    "application/json",
			requestBody:    `{"name":"Gopher","email":"gøpher@example.com"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Unicode characters are not allowed in the 'email' field",
		},
		{
			name:           "Valid XML",
			contentType:    "application/xml",
			requestBody:    `<data><name>Gopher</name><email>gopher@example.com</email></data>`,
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "Invalid XML - Unicode in name",
			contentType:    "application/xml",
			requestBody:    `<data><name>Gøpher</name><email>gopher@example.com</email></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Unicode characters are not allowed in the 'name' field",
		},
		{
			name:           "Invalid XML - Unicode in email",
			contentType:    "application/xml",
			requestBody:    `<data><name>Gopher</name><email>gøpher@example.com</email></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Unicode characters are not allowed in the 'email' field",
		},
		{
			name:           "Invalid Other Content-Type - Unicode in Name",
			contentType:    "text/plain",
			requestBody:    "name=Gøpher&email=gopher@example.com",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Unicode characters are not allowed in the 'name' field",
		},
		{
			name:           "Invalid Other Content-Type - Unicode in Email",
			contentType:    "text/plain",
			requestBody:    "name=Gopher&email=gøpher@example.com",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Unicode characters are not allowed in the 'email' field",
		},
		{
			name:           "Valid Other Content-Type",
			contentType:    "text/plain",
			requestBody:    `{"name":"gopher","email":"gopher@example.com"}`,
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", tc.contentType)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Unexpected error reading response body: %v", err)
			}

			if tc.expectedError != "" {
				if string(body) != tc.expectedError {
					t.Errorf("Expected error message '%s', got '%s'", tc.expectedError, string(body))
				}
			} else {
				if string(body) != "OK" {
					t.Errorf("Expected response body 'OK', got '%s'", string(body))
				}
			}
		})
	}
}

func TestValidatorWithCustomErrorHandler(t *testing.T) {
	app := fiber.New()

	app.Use(validator.New(validator.Config{
		Rules: []validator.Restrictor{
			validator.RestrictUnicode{
				Fields: []string{"custom"},
			},
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
		},
	}))

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	testCases := []struct {
		name           string
		contentType    string
		requestBody    string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Valid JSON",
			contentType:    "application/json",
			requestBody:    `{"custom":"value"}`,
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "Invalid JSON - Unicode in custom field",
			contentType:    "application/json",
			requestBody:    `{"custom":"vålue"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "Unicode characters are not allowed in the 'custom' field",
		},
		{
			name:           "Valid XML",
			contentType:    "application/xml",
			requestBody:    `<data><custom>value</custom></data>`,
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "Invalid XML - Unicode in custom field",
			contentType:    "application/xml",
			requestBody:    `<data><custom>vålue</custom></data>`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "Unicode characters are not allowed in the 'custom' field",
		},
		{
			name:           "Other Content-Type",
			contentType:    "text/plain",
			requestBody:    "custom=vålue",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "Unicode characters are not allowed in the 'custom' field",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", tc.contentType)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Unexpected error reading response body: %v", err)
			}

			if tc.expectedError != "" {
				if string(body) != tc.expectedError {
					t.Errorf("Expected error message '%s', got '%s'", tc.expectedError, string(body))
				}
			} else {
				if string(body) != "OK" {
					t.Errorf("Expected response body 'OK', got '%s'", string(body))
				}
			}
		})
	}
}

func TestValidatorWithNext(t *testing.T) {
	app := fiber.New()

	app.Use(validator.New(validator.Config{
		Rules: []validator.Restrictor{
			validator.RestrictUnicode{
				Fields: []string{"name", "email"},
			},
		},
		Next: func(c *fiber.Ctx) bool {
			return c.Method() == http.MethodGet
		},
	}))

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Skipped validation")
	})

	testCases := []struct {
		name           string
		method         string
		contentType    string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "POST request - validation applied",
			method:         http.MethodPost,
			contentType:    "application/json",
			requestBody:    `{"name":"Gøpher","email":"gopher@example.com"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Unicode characters are not allowed in the 'name' field",
		},
		{
			name:           "GET request - validation skipped",
			method:         http.MethodGet,
			contentType:    "application/json",
			requestBody:    `{"name":"Gøpher","email":"gopher@example.com"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "Skipped validation",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", tc.contentType)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Unexpected error reading response body: %v", err)
			}

			if strings.TrimSpace(string(body)) != tc.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tc.expectedBody, string(body))
			}
		})
	}
}

func TestValidatorWithInvalidJSONBody(t *testing.T) {
	app := fiber.New()

	app.Use(validator.New(validator.Config{
		Rules: []validator.Restrictor{
			validator.RestrictUnicode{
				Fields: []string{"name", "email"},
			},
			validator.RestrictNumberOnly{
				Fields: []string{"age", "score"},
			},
		},
	}))

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Invalid JSON - Incomplete body",
			requestBody:    `{"name":"Gopher","email":"gopher@example.com`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON request body",
		},
		{
			name:           "Invalid JSON - Unicode in name",
			requestBody:    `{"name":"Gøpher","email":"gopher@example.com"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Unicode characters are not allowed in the 'name' field",
		},
		{
			name:           "Invalid JSON - Non-numeric age",
			requestBody:    `{"name":"Gopher","email":"gopher@example.com","age":"abc"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "The 'age' fields must contain only numbers",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Unexpected error reading response body: %v", err)
			}

			if string(body) != tc.expectedError {
				t.Errorf("Expected error message '%s', got '%s'", tc.expectedError, string(body))
			}
		})
	}
}

func TestValidatorWithInvalidXMLBody(t *testing.T) {
	app := fiber.New()

	app.Use(validator.New(validator.Config{
		Rules: []validator.Restrictor{
			validator.RestrictUnicode{
				Fields: []string{"name", "email"},
			},
			validator.RestrictNumberOnly{
				Fields: []string{"age", "score"},
			},
		},
	}))

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Invalid XML - Incomplete body",
			requestBody:    `<data><name>Gopher</name><email>gopher@example.com</data>`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid XML request body",
		},
		{
			name:           "Invalid XML - Unicode in email",
			requestBody:    `<data><name>Gopher</name><email>gøpher@example.com</email></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Unicode characters are not allowed in the 'email' field",
		},
		{
			name:           "Invalid XML - Non-numeric score",
			requestBody:    `<data><name>Gopher</name><email>gopher@example.com</email><score>abc</score></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "The 'score' fields must contain only numbers",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/xml")
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Unexpected error reading response body: %v", err)
			}

			if string(body) != tc.expectedError {
				t.Errorf("Expected error message '%s', got '%s'", tc.expectedError, string(body))
			}
		})
	}
}

func ptr(i int) *int {
	return &i
}

func TestRestrictNumberOnly(t *testing.T) {
	app := fiber.New()

	app.Use(validator.New(validator.Config{
		Rules: []validator.Restrictor{
			validator.RestrictNumberOnly{
				Fields: []string{"age", "score", "seafood_price"},
				Max:    ptr(100),
			},
		},
	}))

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	testCases := []struct {
		name           string
		contentType    string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid JSON request",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"age":30,"score":80,"seafood_price":50}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "Invalid JSON request - non-numeric age",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"age":"abc","score":80,"seafood_price":50}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'age' fields must contain only numbers",
		},
		{
			name:           "Invalid JSON request - non-numeric score",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"age":30,"score":"def","seafood_price":50}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'score' fields must contain only numbers",
		},
		{
			name:           "Invalid JSON request - age exceeds maximum",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"age":120,"score":80,"seafood_price":50}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'age' field must not exceed 100",
		},
		{
			name:           "Valid XML request",
			contentType:    fiber.MIMEApplicationXML,
			requestBody:    `<data><age>30</age><score>80</score><seafood_price>50</seafood_price></data>`,
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "Invalid XML request - non-numeric age",
			contentType:    fiber.MIMEApplicationXML,
			requestBody:    `<data><age>abc</age><score>80</score><seafood_price>50</seafood_price></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'age' fields must contain only numbers",
		},
		{
			name:           "Invalid XML request - non-numeric score",
			contentType:    fiber.MIMEApplicationXML,
			requestBody:    `<data><age>30</age><score>def</score><seafood_price>50</seafood_price></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'score' fields must contain only numbers",
		},
		{
			name:           "Invalid XML request - score exceeds maximum",
			contentType:    fiber.MIMEApplicationXML,
			requestBody:    `<data><age>30</age><score>120</score><seafood_price>50</seafood_price></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'score' field must not exceed 100",
		},
		{
			name:           "Invalid Other Content-Type - age exceeds maximum",
			contentType:    "text/plain",
			requestBody:    "age=120&score=80&seafood_price=50",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'age' field must not exceed 100",
		},
		{
			name:           "Invalid Other Content-Type - age not numeric",
			contentType:    "text/plain",
			requestBody:    "age=gh0per&score=80&seafood_price=50",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'age' fields must contain only numbers",
		},
		{
			name:           "Invalid JSON request - non-numeric age and score",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"age":"abc","score":"xa","seafood_price":50}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'age', 'score' fields must contain only numbers",
		},
		{
			name:           "Invalid JSON request - non-numeric age, score, and seafood_price",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"age":"abc","score":"def","seafood_price":"ghi"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'age', 'score', 'seafood_price' fields must contain only numbers",
		},
		{
			name:           "Invalid JSON request - age and seafood_price exceed maximum",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"age":120,"score":80,"seafood_price":150}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'age' field must not exceed 100",
		},
		{
			name:           "Invalid XML request - non-numeric score and seafood_price",
			contentType:    fiber.MIMEApplicationXML,
			requestBody:    `<data><age>30</age><score>abc</score><seafood_price>def</seafood_price></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'score', 'seafood_price' fields must contain only numbers",
		},
		{
			name:           "Invalid XML request - age and score exceed maximum",
			contentType:    fiber.MIMEApplicationXML,
			requestBody:    `<data><age>120</age><score>150</score><seafood_price>80</seafood_price></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'age' field must not exceed 100",
		},
		{
			name:           "Invalid Other Content-Type - seafood_price exceeds maximum",
			contentType:    "text/plain",
			requestBody:    "age=30&score=80&seafood_price=120",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'seafood_price' field must not exceed 100",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", tc.contentType)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Unexpected error reading response body: %v", err)
			}

			if strings.TrimSpace(string(body)) != tc.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tc.expectedBody, string(body))
			}
		})
	}
}
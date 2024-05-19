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
		},
	}))

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Gopher","email":"gopher@example.com`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unexpected error reading response body: %v", err)
	}

	expectedError := "Invalid JSON request body"
	if string(body) != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, string(body))
	}
}

func TestValidatorWithInvalidXMLBody(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`<data><name>Gopher</name><email>gopher@example.com</data>`))
	req.Header.Set("Content-Type", "application/xml")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unexpected error reading response body: %v", err)
	}

	expectedError := "Invalid XML request body"
	if string(body) != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, string(body))
	}
}

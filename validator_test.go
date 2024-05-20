// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package validator_test

import (
	"encoding/json"
	"encoding/xml"
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
				if tc.contentType == fiber.MIMEApplicationJSON {
					var jsonResp map[string]string
					if err := json.Unmarshal(body, &jsonResp); err != nil {
						t.Fatalf("Unexpected error unmarshaling JSON response: %v", err)
					}
					if jsonResp["error"] != tc.expectedError {
						t.Errorf("Expected error message '%s', got '%s'", tc.expectedError, jsonResp["error"])
					}
				} else if tc.contentType == fiber.MIMEApplicationXML {
					var xmlResp struct {
						Error string `xml:"error"`
					}
					if err := xml.Unmarshal(body, &xmlResp); err != nil {
						t.Fatalf("Unexpected error unmarshaling XML response: %v", err)
					}
					if xmlResp.Error != tc.expectedError {
						t.Errorf("Expected error message '%s', got '%s'", tc.expectedError, xmlResp.Error)
					}
				} else {
					if string(body) != tc.expectedError {
						t.Errorf("Expected error message '%s', got '%s'", tc.expectedError, string(body))
					}
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
			expectedBody:   `{"error":"Unicode characters are not allowed in the 'name' field"}`,
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
			expectedError:  "The 'age' field must contain only numbers",
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

			var jsonResp map[string]string
			if err := json.Unmarshal(body, &jsonResp); err != nil {
				t.Fatalf("Unexpected error unmarshaling JSON response: %v", err)
			}
			if jsonResp["error"] != tc.expectedError {
				t.Errorf("Expected error message '%s', got '%s'", tc.expectedError, jsonResp["error"])
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
			expectedError:  `<xmlError><error>Invalid XML request body</error></xmlError>`,
		},
		{
			name:           "Invalid XML - Unicode in email",
			requestBody:    `<data><name>Gopher</name><email>gøpher@example.com</email></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  `<xmlError><error>Unicode characters are not allowed in the &#39;email&#39; field</error></xmlError>`,
		},
		{
			name:           "Invalid XML - Non-numeric score",
			requestBody:    `<data><name>Gopher</name><email>gopher@example.com</email><score>abc</score></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  `<xmlError><error>The &#39;score&#39; field must contain only numbers</error></xmlError>`,
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
			expectedBody:   `{"error":"The 'age' field must contain only numbers"}`,
		},
		{
			name:           "Invalid JSON request - non-numeric score",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"age":30,"score":"def","seafood_price":50}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"The 'score' field must contain only numbers"}`,
		},
		{
			name:           "Invalid JSON request - age exceeds maximum",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"age":120,"score":80,"seafood_price":50}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"The 'age' field must not exceed 100"}`,
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
			expectedBody:   `<xmlError><error>The &#39;age&#39; field must contain only numbers</error></xmlError>`,
		},
		{
			name:           "Invalid XML request - non-numeric score",
			contentType:    fiber.MIMEApplicationXML,
			requestBody:    `<data><age>30</age><score>def</score><seafood_price>50</seafood_price></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `<xmlError><error>The &#39;score&#39; field must contain only numbers</error></xmlError>`,
		},
		{
			name:           "Invalid XML request - score exceeds maximum",
			contentType:    fiber.MIMEApplicationXML,
			requestBody:    `<data><age>30</age><score>120</score><seafood_price>50</seafood_price></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `<xmlError><error>The &#39;score&#39; field must not exceed 100</error></xmlError>`,
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
			expectedBody:   "The 'age' field must contain only numbers",
		},
		{
			name:           "Invalid JSON request - non-numeric age and score",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"age":"abc","score":"xa","seafood_price":50}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"The 'age', 'score' field must contain only numbers"}`,
		},
		{
			name:           "Invalid JSON request - non-numeric age, score, and seafood_price",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"age":"abc","score":"def","seafood_price":"ghi"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"The 'age', 'score', 'seafood_price' field must contain only numbers"}`,
		},
		{
			name:           "Invalid JSON request - age and seafood_price exceed maximum",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"age":120,"score":80,"seafood_price":150}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"The 'age' field must not exceed 100"}`,
		},
		{
			name:           "Invalid XML request - non-numeric score and seafood_price",
			contentType:    fiber.MIMEApplicationXML,
			requestBody:    `<data><age>30</age><score>abc</score><seafood_price>def</seafood_price></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "<xmlError><error>The &#39;score&#39;, &#39;seafood_price&#39; field must contain only numbers</error></xmlError>",
		},
		{
			name:           "Invalid XML request - age and score exceed maximum",
			contentType:    fiber.MIMEApplicationXML,
			requestBody:    `<data><age>120</age><score>150</score><seafood_price>80</seafood_price></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "<xmlError><error>The &#39;age&#39; field must not exceed 100</error></xmlError>",
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

func TestRestrictNumberWithMaxDigitsOnly(t *testing.T) {
	app := fiber.New()

	app.Use(validator.New(validator.Config{
		Rules: []validator.Restrictor{
			validator.RestrictNumberOnly{
				Fields:    []string{"grilled_fish", "lobster_roll", "seafood_platter", "bali_juice"},
				Max:       ptr(100),
				MaxDigits: ptr(3),
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
			requestBody:    `{"grilled_fish":30,"lobster_roll":80,"seafood_platter":50,"bali_juice":10}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "Invalid JSON request - grilled_fish exceeds maximum digits",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"grilled_fish":1234,"lobster_roll":80,"seafood_platter":50,"bali_juice":10}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"The 'grilled_fish' field must not exceed 3 digits"}`,
		},
		{
			name:           "Invalid XML request - lobster_roll exceeds maximum digits",
			contentType:    fiber.MIMEApplicationXML,
			requestBody:    `<data><grilled_fish>30</grilled_fish><lobster_roll>1234</lobster_roll><seafood_platter>50</seafood_platter><bali_juice>10</bali_juice></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `<xmlError><error>The &#39;lobster_roll&#39; field must not exceed 3 digits</error></xmlError>`,
		},
		{
			name:           "Invalid Other Content-Type - seafood_platter exceeds maximum digits",
			contentType:    "text/plain",
			requestBody:    "grilled_fish=30&lobster_roll=80&seafood_platter=1234&bali_juice=10",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'seafood_platter' field must not exceed 3 digits",
		},
		{
			name:           "Invalid JSON request - grilled_fish and lobster_roll exceed maximum digits",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"grilled_fish":1234,"lobster_roll":5678,"seafood_platter":50,"bali_juice":10}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"The 'grilled_fish' field must not exceed 3 digits"}`,
		},
		{
			name:           "Invalid XML request - grilled_fish, lobster_roll, and seafood_platter exceed maximum digits",
			contentType:    fiber.MIMEApplicationXML,
			requestBody:    `<data><grilled_fish>1234</grilled_fish><lobster_roll>5678</lobster_roll><seafood_platter>9012</seafood_platter><bali_juice>10</bali_juice></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `<xmlError><error>The &#39;grilled_fish&#39; field must not exceed 3 digits</error></xmlError>`,
		},
		{
			name:           "Invalid Other Content-Type - grilled_fish and seafood_platter exceed maximum digits",
			contentType:    "text/plain",
			requestBody:    "grilled_fish=1234&lobster_roll=80&seafood_platter=5678&bali_juice=10",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'grilled_fish' field must not exceed 3 digits",
		},
		{
			name:           "Invalid JSON request - bali_juice exceeds maximum digits",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"grilled_fish":30,"lobster_roll":80,"seafood_platter":50,"bali_juice":1234}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"The 'bali_juice' field must not exceed 3 digits"}`,
		},
		{
			name:           "Invalid XML request - bali_juice exceeds maximum digits",
			contentType:    fiber.MIMEApplicationXML,
			requestBody:    `<data><grilled_fish>30</grilled_fish><lobster_roll>80</lobster_roll><seafood_platter>50</seafood_platter><bali_juice>1234</bali_juice></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `<xmlError><error>The &#39;bali_juice&#39; field must not exceed 3 digits</error></xmlError>`,
		},
		{
			name:           "Invalid Other Content-Type - bali_juice exceeds maximum digits",
			contentType:    "text/plain",
			requestBody:    "grilled_fish=30&lobster_roll=80&seafood_platter=50&bali_juice=1234",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'bali_juice' field must not exceed 3 digits",
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

func TestValidatorWithContextKey(t *testing.T) {
	app := fiber.New()

	app.Use(validator.New(validator.Config{
		Rules: []validator.Restrictor{
			validator.RestrictUnicode{
				Fields: []string{"name", "email"},
			},
		},
		ContextKey: "validationResult",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		},
	}))

	app.Post("/", func(c *fiber.Ctx) error {
		validationResult := c.Locals("validationResult")
		if validationResult != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Validation failed")
		}
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
			name:           "Valid request",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"name":"Gopher","email":"gopher@example.com"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "Invalid request - Unicode in name",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"name":"Gøpher","email":"gopher@example.com"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Unicode characters are not allowed in the 'name' field",
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

func TestRestrictStringLength(t *testing.T) {
	app := fiber.New()

	app.Use(validator.New(validator.Config{
		Rules: []validator.Restrictor{
			validator.RestrictStringLength{
				Fields:    []string{"name", "description"},
				MaxLength: ptr(50),
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
			requestBody:    `{"name":"Gopher","description":"A short description"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "Invalid JSON request - name exceeds maximum length",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"name":"Gopher with a very long name that exceeds the maximum length","description":"A short description"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"The 'name' field must not exceed 50 characters"}`,
		},
		{
			name:           "Invalid JSON request - description exceeds maximum length",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"name":"Gopher","description":"A very long description that exceeds the maximum length limit"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"The 'description' field must not exceed 50 characters"}`,
		},
		{
			name:           "Invalid JSON request - invalid JSON body",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"name":"Gopher","description":"A short description"`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid JSON request body"}`,
		},
		{
			name:           "Invalid JSON request - multiple fields exceed maximum length",
			contentType:    fiber.MIMEApplicationJSON,
			requestBody:    `{"name":"Gopher with a very long name that exceeds the maximum length","description":"A very long description that exceeds the maximum length limit"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"The 'name' field must not exceed 50 characters"}`,
		},
		{
			name:           "Valid XML request",
			contentType:    fiber.MIMEApplicationXML,
			requestBody:    `<data><name>Gopher</name><description>A short description</description></data>`,
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "Invalid XML request - name exceeds maximum length",
			contentType:    fiber.MIMEApplicationXML,
			requestBody:    `<data><name>Gopher with a very long name that exceeds the maximum length</name><description>A short description</description></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `<xmlError><error>The &#39;name&#39; field must not exceed 50 characters</error></xmlError>`,
		},
		{
			name:           "Invalid XML request - description exceeds maximum length",
			contentType:    fiber.MIMEApplicationXML,
			requestBody:    `<data><name>Gopher</name><description>A very long description that exceeds the maximum length limit</description></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `<xmlError><error>The &#39;description&#39; field must not exceed 50 characters</error></xmlError>`,
		},
		{
			name:           "Invalid XML request - invalid XML body",
			contentType:    fiber.MIMEApplicationXML,
			requestBody:    `<data><name>Gopher</name><description>A short description</description>`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `<xmlError><error>Invalid XML request body</error></xmlError>`,
		},
		{
			name:           "Invalid XML request - multiple fields exceed maximum length",
			contentType:    fiber.MIMEApplicationXML,
			requestBody:    `<data><name>Gopher with a very long name that exceeds the maximum length</name><description>A very long description that exceeds the maximum length limit</description></data>`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `<xmlError><error>The &#39;name&#39; field must not exceed 50 characters</error></xmlError>`,
		},
		{
			name:           "Valid Other Content-Type request",
			contentType:    "text/plain",
			requestBody:    "name=Gopher&description=A short description",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "Invalid Other Content-Type - name exceeds maximum length",
			contentType:    "text/plain",
			requestBody:    "name=Gopher with a very long name that exceeds the maximum length&description=A short description",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'name' field must not exceed 50 characters",
		},
		{
			name:           "Invalid Other Content-Type - description exceeds maximum length",
			contentType:    "text/plain",
			requestBody:    "name=Gopher&description=A very long description that exceeds the maximum length limit",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'description' field must not exceed 50 characters",
		},
		{
			name:           "Invalid Other Content-Type - multiple fields exceed maximum length",
			contentType:    "text/plain",
			requestBody:    "name=Gopher with a very long name that exceeds the maximum length&description=A very long description that exceeds the maximum length limit",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The 'name' field must not exceed 50 characters",
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

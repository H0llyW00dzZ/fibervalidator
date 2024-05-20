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
	"github.com/bytedance/sonic"
	"github.com/clbanning/mxj"
	"github.com/gofiber/fiber/v2"
)

func BenchmarkValidatorWithSonicJSONSeafood(b *testing.B) {
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	app.Use(validator.New(validator.Config{
		Rules: []validator.Restrictor{
			validator.RestrictUnicode{
				Fields: []string{"name", "description"},
			},
			validator.RestrictNumberOnly{
				Fields:    []string{"price", "quantity"},
				Max:       ptr(99999999),
				MaxDigits: ptr(99999999),
			},
		},
	}))

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	testCase := struct {
		name        string
		contentType string
		requestBody string
	}{
		name:        "Valid JSON request",
		contentType: fiber.MIMEApplicationJSON,
		requestBody: `{"name":"Lobster","description":"Fresh Maine lobster","price":50,"quantity":10}`,
	}

	b.Run(testCase.name, func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testCase.requestBody))
			req.Header.Set("Content-Type", testCase.contentType)
			resp, err := app.Test(req)
			if err != nil {
				b.Fatalf("Unexpected error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			}

			_, err = io.ReadAll(resp.Body)
			if err != nil {
				b.Fatalf("Unexpected error reading response body: %v", err)
			}
		}
	})
}

func BenchmarkValidatorWithStandardJSONSeafood(b *testing.B) {
	app := fiber.New()

	app.Use(validator.New(validator.Config{
		Rules: []validator.Restrictor{
			validator.RestrictUnicode{
				Fields: []string{"name", "description"},
			},
			validator.RestrictNumberOnly{
				Fields:    []string{"price", "quantity"},
				Max:       ptr(99999999),
				MaxDigits: ptr(99999999),
			},
		},
	}))

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	testCase := struct {
		name        string
		contentType string
		requestBody string
	}{
		name:        "Valid JSON request",
		contentType: fiber.MIMEApplicationJSON,
		requestBody: `{"name":"Lobster","description":"Fresh Maine lobster","price":50,"quantity":10}`,
	}

	b.Run(testCase.name, func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testCase.requestBody))
			req.Header.Set("Content-Type", testCase.contentType)
			resp, err := app.Test(req)
			if err != nil {
				b.Fatalf("Unexpected error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			}

			_, err = io.ReadAll(resp.Body)
			if err != nil {
				b.Fatalf("Unexpected error reading response body: %v", err)
			}
		}
	})
}

func BenchmarkValidatorWithDefaultXMLSeafood(b *testing.B) {
	app := fiber.New()

	app.Use(validator.New(validator.Config{
		Rules: []validator.Restrictor{
			validator.RestrictUnicode{
				Fields: []string{"name", "description"},
			},
			validator.RestrictNumberOnly{
				Fields:    []string{"price", "quantity"},
				Max:       ptr(99999999),
				MaxDigits: ptr(99999999),
			},
		},
	}))

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	testCase := struct {
		name        string
		contentType string
		requestBody string
	}{
		name:        "Valid XML request",
		contentType: fiber.MIMEApplicationXML,
		requestBody: `<data><name>Lobster</name><description>Fresh Maine lobster</description><price>50</price><quantity>10</quantity></data>`,
	}

	b.Run(testCase.name, func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testCase.requestBody))
			req.Header.Set("Content-Type", testCase.contentType)
			resp, err := app.Test(req)
			if err != nil {
				b.Fatalf("Unexpected error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			}

			_, err = io.ReadAll(resp.Body)
			if err != nil {
				b.Fatalf("Unexpected error reading response body: %v", err)
			}
		}
	})
}

func BenchmarkValidatorWithCustomXMLSeafood(b *testing.B) {
	app := fiber.New(fiber.Config{
		XMLEncoder: customXMLMarshal,
	})

	app.Use(validator.New(validator.Config{
		Rules: []validator.Restrictor{
			validator.RestrictUnicode{
				Fields: []string{"name", "description"},
			},
			validator.RestrictNumberOnly{
				Fields:    []string{"price", "quantity"},
				Max:       ptr(99999999),
				MaxDigits: ptr(99999999),
			},
		},
	}))

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	testCase := struct {
		name        string
		contentType string
		requestBody string
	}{
		name:        "Valid XML request",
		contentType: fiber.MIMEApplicationXML,
		requestBody: `<data><name>Lobster</name><description>Fresh Maine lobster</description><price>50</price><quantity>10</quantity></data>`,
	}

	b.Run(testCase.name, func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testCase.requestBody))
			req.Header.Set("Content-Type", testCase.contentType)
			resp, err := app.Test(req)
			if err != nil {
				b.Fatalf("Unexpected error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			}

			_, err = io.ReadAll(resp.Body)
			if err != nil {
				b.Fatalf("Unexpected error reading response body: %v", err)
			}
		}
	})
}
func BenchmarkValidatorWithSonicJSON(b *testing.B) {
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	app.Use(validator.New(validator.Config{
		Rules: []validator.Restrictor{
			validator.RestrictUnicode{
				Fields: []string{"name", "email"},
			},
			validator.RestrictNumberOnly{
				Fields: []string{"age", "score"},
				Max:    ptr(100),
			},
		},
	}))

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	testCase := struct {
		name        string
		contentType string
		requestBody string
	}{
		name:        "Valid JSON request",
		contentType: fiber.MIMEApplicationJSON,
		requestBody: `{"name":"Gopher","email":"gopher@example.com","age":30,"score":80}`,
	}

	b.Run(testCase.name, func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testCase.requestBody))
			req.Header.Set("Content-Type", testCase.contentType)
			resp, err := app.Test(req)
			if err != nil {
				b.Fatalf("Unexpected error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			}

			_, err = io.ReadAll(resp.Body)
			if err != nil {
				b.Fatalf("Unexpected error reading response body: %v", err)
			}
		}
	})
}

func BenchmarkValidatorWithStandardJSON(b *testing.B) {
	app := fiber.New()

	app.Use(validator.New(validator.Config{
		Rules: []validator.Restrictor{
			validator.RestrictUnicode{
				Fields: []string{"name", "email"},
			},
			validator.RestrictNumberOnly{
				Fields: []string{"age", "score"},
				Max:    ptr(100),
			},
		},
	}))

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	testCase := struct {
		name        string
		contentType string
		requestBody string
	}{
		name:        "Valid JSON request",
		contentType: fiber.MIMEApplicationJSON,
		requestBody: `{"name":"Gopher","email":"gopher@example.com","age":30,"score":80}`,
	}

	b.Run(testCase.name, func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testCase.requestBody))
			req.Header.Set("Content-Type", testCase.contentType)
			resp, err := app.Test(req)
			if err != nil {
				b.Fatalf("Unexpected error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			}

			_, err = io.ReadAll(resp.Body)
			if err != nil {
				b.Fatalf("Unexpected error reading response body: %v", err)
			}
		}
	})
}

func BenchmarkValidatorWithDefaultXML(b *testing.B) {
	app := fiber.New()

	app.Use(validator.New(validator.Config{
		Rules: []validator.Restrictor{
			validator.RestrictUnicode{
				Fields: []string{"name", "email"},
			},
			validator.RestrictNumberOnly{
				Fields: []string{"age", "score"},
				Max:    ptr(100),
			},
		},
	}))

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	testCase := struct {
		name        string
		contentType string
		requestBody string
	}{
		name:        "Valid XML request",
		contentType: fiber.MIMEApplicationXML,
		requestBody: `<data><name>Gopher</name><email>gopher@example.com</email><age>30</age><score>80</score></data>`,
	}

	b.Run(testCase.name, func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testCase.requestBody))
			req.Header.Set("Content-Type", testCase.contentType)
			resp, err := app.Test(req)
			if err != nil {
				b.Fatalf("Unexpected error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			}

			_, err = io.ReadAll(resp.Body)
			if err != nil {
				b.Fatalf("Unexpected error reading response body: %v", err)
			}
		}
	})
}

func BenchmarkValidatorWithCustomXML(b *testing.B) {
	app := fiber.New(fiber.Config{
		XMLEncoder: customXMLMarshal,
	})

	app.Use(validator.New(validator.Config{
		Rules: []validator.Restrictor{
			validator.RestrictUnicode{
				Fields: []string{"name", "email"},
			},
			validator.RestrictNumberOnly{
				Fields: []string{"age", "score"},
				Max:    ptr(100),
			},
		},
	}))

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	testCase := struct {
		name        string
		contentType string
		requestBody string
	}{
		name:        "Valid XML request",
		contentType: fiber.MIMEApplicationXML,
		requestBody: `<data><name>Gopher</name><email>gopher@example.com</email><age>30</age><score>80</score></data>`,
	}

	b.Run(testCase.name, func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testCase.requestBody))
			req.Header.Set("Content-Type", testCase.contentType)
			resp, err := app.Test(req)
			if err != nil {
				b.Fatalf("Unexpected error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			}

			_, err = io.ReadAll(resp.Body)
			if err != nil {
				b.Fatalf("Unexpected error reading response body: %v", err)
			}
		}
	})
}

func customXMLMarshal(v interface{}) ([]byte, error) {
	return mxj.AnyXmlIndent(v, "", "  ")
}

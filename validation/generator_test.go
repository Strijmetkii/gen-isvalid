package validation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerator(t *testing.T) {
	// Create a temporary test file
	dir := t.TempDir()
	testFile := filepath.Join(dir, "test.go")

	// Create test content
	content := `package test

// TestService is a test service
//go:generate go run ../cmd/gen/main.go
type TestService struct {
	Client *Client
	Config *Config
	Name string
}

// Client is a test client
type Client struct {}

// Config is a test config
type Config struct {}
`

	// Write test content to file
	err := os.WriteFile(testFile, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create generator
	generator := NewGenerator(testFile)

	// Generate code
	err = generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate code: %v", err)
	}

	// Check that output file exists
	outputFile := generator.OutputFile
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file not created: %v", err)
	}

	// Read generated code
	generatedCode, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read generated code: %v", err)
	}

	// Check generated code contents
	codeStr := string(generatedCode)

	// Check that the package is correct
	if !strings.Contains(codeStr, "package test") {
		t.Errorf("Generated code doesn't have correct package")
	}

	// Check that the parameter struct is generated
	if !strings.Contains(codeStr, "type TestServiceParams struct") {
		t.Errorf("Generated code doesn't have parameter struct")
	}

	// Check that the constructor is generated
	if !strings.Contains(codeStr, "func NewTestService(params TestServiceParams)") {
		t.Errorf("Generated code doesn't have constructor")
	}

	// Check that the validation function is generated
	if !strings.Contains(codeStr, "func isValidTestServiceParams(params TestServiceParams)") {
		t.Errorf("Generated code doesn't have validation function")
	}

	// Check that pointer fields are validated
	if !strings.Contains(codeStr, `if params.Client == nil {`) {
		t.Errorf("Generated code doesn't validate client pointer")
	}

	if !strings.Contains(codeStr, `if params.Config == nil {`) {
		t.Errorf("Generated code doesn't validate config pointer")
	}
}

func TestNoStructs(t *testing.T) {
	// Create a temporary test file
	dir := t.TempDir()
	testFile := filepath.Join(dir, "test.go")

	// Create test content with no structs that have go:generate directives
	content := `package test

// TestService is a test service
type TestService struct {
	client *Client
	config *Config
}

// Client is a test client
type Client struct {}

// Config is a test config
type Config struct {}
`

	// Write test content to file
	err := os.WriteFile(testFile, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create generator
	generator := NewGenerator(testFile)

	// Generate code should fail
	err = generator.Generate()
	if err == nil {
		t.Fatalf("Expected error when no structs with go:generate directive found")
	}

	if !strings.Contains(err.Error(), "no structs with go:generate directive found") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestGenerics(t *testing.T) {
	// Create a temporary test file
	dir := t.TempDir()
	testFile := filepath.Join(dir, "test_generics.go")

	// Create test content with generics
	content := `package test

// GenericService is a generic service
//go:generate go run ../cmd/gen/main.go
type GenericService[T any] struct {
	Repository *Repository[T]
	Logger *Logger
}

// Repository is a generic repository
type Repository[T any] interface {
	Get(id string) (T, error)
	Save(item T) error
}

// Logger is a test logger
type Logger struct {}

// MultiParamService has multiple type parameters
//go:generate go run ../cmd/gen/main.go
type MultiParamService[K comparable, V any] struct {
	Store *Store[K, V]
	Config *Config
}

// Store is a key-value store
type Store[K comparable, V any] interface {
	Get(key K) (V, error)
	Set(key K, value V) error
}

// Config is a test config
type Config struct {}
`

	// Write test content to file
	err := os.WriteFile(testFile, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create generator
	generator := NewGenerator(testFile)

	// Generate code
	err = generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate code: %v", err)
	}

	// Check that output file exists
	outputFile := generator.OutputFile
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file not created: %v", err)
	}

	// Read generated code
	generatedCode, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read generated code: %v", err)
	}

	// Check generated code contents
	codeStr := string(generatedCode)

	// Check that the generics are correctly handled
	if !strings.Contains(codeStr, "type GenericServiceParams[T any]") {
		t.Errorf("Generated code doesn't have generic parameter struct")
	}

	if !strings.Contains(codeStr, "func NewGenericService[T any](params GenericServiceParams[T])") {
		t.Errorf("Generated code doesn't have correct generic constructor")
	}

	if !strings.Contains(codeStr, "func isValidGenericServiceParams[T any](params GenericServiceParams[T])") {
		t.Errorf("Generated code doesn't have correct generic validation function")
	}

	// Check multiple type parameters
	if !strings.Contains(codeStr, "type MultiParamServiceParams[K comparable, V any]") {
		t.Errorf("Generated code doesn't handle multiple type parameters correctly")
	}

	// Check validation of pointer fields in generics
	if !strings.Contains(codeStr, "if params.Repository == nil {") {
		t.Errorf("Generated code doesn't validate pointer fields in generics")
	}
}

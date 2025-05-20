# Validation Code Generator for Go

This tool automatically generates parameter validation code for Go struct types. It examines Go source files for struct definitions with the `//go:generate` directive and generates corresponding validation infrastructure.

## Features

- Automatically generates parameter structs from your service structs
- Creates constructor functions with validation checks
- Validates pointer fields for nil values
- Organizes parameters in a clean, maintainable way
- Works with Go's built-in `go generate` tool
- Supports generic types and interfaces
- Skips generation if output file already exists (can be overridden)

## Installation

```bash
go install github.com/strijmetkii/gen-isvalid/cmd/gen@latest
```

## Usage

1. Add a `//go:generate` directive before your struct definitions:

```go
// ExampleService is a service for interacting with an API
//go:generate go run github.com/strijmetkii/gen-isvalid/cmd/gen
type ExampleService struct {
    Client *Client
    Config *Config
    Timeout int
}
```

2. Run `go generate` in your project:

```bash
go generate ./...
```

This will create a `*_gen.go` file with your validation code:

```go
// ExampleServiceParams is the parameter struct for creating a ExampleService
type ExampleServiceParams struct {
    Client *Client
    Config *Config
    Timeout int
}

// NewExampleService creates a new ExampleService
func NewExampleService(params ExampleServiceParams) (*ExampleService, error) {
    if err := isValidExampleServiceParams(params); err != nil {
        return nil, err
    }

    return &ExampleService{
        Client: params.Client,
        Config: params.Config,
        Timeout: params.Timeout,
    }, nil
}

// isValidExampleServiceParams validates the ExampleServiceParams
func isValidExampleServiceParams(params ExampleServiceParams) error {
    var errs []error
    if params.Client == nil {
        errs = append(errs, errors.New("Client is required"))
    }
    if params.Config == nil {
        errs = append(errs, errors.New("Config is required"))
    }
    return errors.Join(errs...)
}
```

## Command Line Options

You can also run the generator directly with these options:

```
Usage:
  gen [options]

Options:
  -input string
        Path to the input Go file (default is the file that triggered go:generate)
  -output string
        Path to the output Go file (default is <input>_gen.go)
  -force
        Force regeneration even if output file exists
```

## Generic Types Support

The generator fully supports Go's generic types:

```go
// GenericService is a service with a generic type parameter
//go:generate go run github.com/strijmetkii/gen-isvalid/cmd/gen
type GenericService[T any] struct {
    Repository *Repository[T]
    Logger Logger
}
```

This will generate:

```go
// GenericServiceParams is the parameter struct for creating a GenericService
type GenericServiceParams[T any] struct {
    Repository *Repository[T]
    Logger Logger
}

// NewGenericService creates a new GenericService
func NewGenericService[T any](params GenericServiceParams[T]) (*GenericService[T], error) {
    if err := isValidGenericServiceParams[T](params); err != nil {
        return nil, err
    }

    return &GenericService[T]{
        Repository: params.Repository,
        Logger: params.Logger,
    }, nil
}

// isValidGenericServiceParams validates the GenericServiceParams
func isValidGenericServiceParams[T any](params GenericServiceParams[T]) error {
    var errs []error
    if params.Repository == nil {
        errs = append(errs, errors.New("Repository is required"))
    }
    return errors.Join(errs...)
}
```

## Architecture

The generator is structured into several key components:

### 1. Command Line Interface (`cmd/gen/main.go`)

- Handles command-line arguments and environment variables
- Creates and configures the generator
- Reports success or errors to the user

### 2. Generator (`validation/generator.go`)

- Core logic for parsing Go source files and generating validation code
- Finds struct definitions with go:generate directives
- Extracts field information (name, type, pointer status)
- Uses templates to generate validation code
- Handles generic type parameters

### 3. Templates

- Uses Go's text/template package to generate code
- Produces parameter structs, constructor functions, and validation logic

## Implementation Details

### AST Parsing

The generator uses Go's abstract syntax tree (AST) packages to parse source files:

```go
fset := token.NewFileSet()
node, err := parser.ParseFile(fset, inputFile, nil, parser.ParseComments)
```

It then walks the AST to find struct declarations with go:generate directives:

```go
for _, decl := range node.Decls {
    genDecl, ok := decl.(*ast.GenDecl)
    if !ok || genDecl.Tok != token.TYPE {
        continue
    }
    
    // Process type declarations...
}
```

### Field Analysis

For each struct field, the generator:

1. Checks if the field is exported
2. Determines if it's a pointer type
3. Extracts the underlying type name
4. Creates appropriate validation for pointer fields

### Code Generation

The generator uses Go templates to produce consistent code:

```go
tmpl, err := template.New("validation").Parse(codeTemplate)
var buf bytes.Buffer
err = tmpl.Execute(&buf, data)
```

The generated code follows these patterns:

1. Parameter structs mirror the original struct fields
2. Constructor functions validate parameters and create the struct
3. Validation functions check for nil pointers and other requirements

## Development

To build and test the generator locally:

```bash
# Build the generator
make build

# Run the example (skips if files exist)
make example

# Force regeneration of example code
make force-example

# Run tests
make test

# Clean generated files
make clean
```

## License

MIT 
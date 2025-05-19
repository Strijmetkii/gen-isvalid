package validation

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Generator manages the validation code generation process
type Generator struct {
	// InputFile is the path to the input Go file
	InputFile string
	// OutputFile is the path to the output generated code file
	OutputFile string
	// PackageName is the name of the package for the generated code
	PackageName string
}

// StructInfo contains information about a struct for which validation code will be generated
type StructInfo struct {
	// Name is the name of the struct
	Name string
	// Fields are the fields of the struct
	Fields []FieldInfo
	// PackageName is the package of the struct
	PackageName string
	// TypeParams are the type parameters of the struct if it's generic
	TypeParams string
	// IsGeneric indicates if the struct is a generic type
	IsGeneric bool
}

// FieldInfo contains information about a struct field
type FieldInfo struct {
	// Name is the name of the field
	Name string
	// Type is the type of the field
	Type string
	// IsPointer indicates if the field is a pointer type
	IsPointer bool
}

// NewGenerator creates a new generator for the given input file
func NewGenerator(inputFile string) *Generator {
	dir, filename := filepath.Split(inputFile)
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	outputFile := filepath.Join(dir, base+"_gen.go")

	return &Generator{
		InputFile:  inputFile,
		OutputFile: outputFile,
	}
}

// Generate parses the input file and generates the validation code
func (g *Generator) Generate() error {
	// Parse the input file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, g.InputFile, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parsing file: %w", err)
	}

	g.PackageName = node.Name.Name

	// Find structs with the go:generate comment
	var structs []StructInfo
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// Check if the struct has our go:generate directive
			if !hasGenerateDirective(genDecl.Doc) {
				continue
			}

			// Check for type parameters (generics)
			typeParams := ""
			isGeneric := false
			if typeSpec.TypeParams != nil && len(typeSpec.TypeParams.List) > 0 {
				isGeneric = true
				typeParams = extractTypeParams(typeSpec.TypeParams)
			}

			// Extract struct info
			structInfo := StructInfo{
				Name:        typeSpec.Name.Name,
				PackageName: g.PackageName,
				Fields:      make([]FieldInfo, 0, len(structType.Fields.List)),
				TypeParams:  typeParams,
				IsGeneric:   isGeneric,
			}

			// Extract field info
			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					// Skip embedded fields
					continue
				}

				fieldName := field.Names[0].Name

				// Skip unexported fields
				if !ast.IsExported(fieldName) {
					continue
				}

				fieldType := ""
				isPointer := false

				// Check if the field is a pointer
				switch t := field.Type.(type) {
				case *ast.StarExpr:
					isPointer = true
					// Get the underlying type
					fieldType = extractType(t.X)
				default:
					fieldType = extractType(field.Type)
				}

				structInfo.Fields = append(structInfo.Fields, FieldInfo{
					Name:      fieldName,
					Type:      fieldType,
					IsPointer: isPointer,
				})
			}

			structs = append(structs, structInfo)
		}
	}

	if len(structs) == 0 {
		return fmt.Errorf("no structs with go:generate directive found")
	}

	// Generate the code
	code, err := g.generateCode(structs)
	if err != nil {
		return fmt.Errorf("generating code: %w", err)
	}

	// Format the code
	formattedCode, err := format.Source([]byte(code))
	if err != nil {
		return fmt.Errorf("formatting generated code: %w", err)
	}

	// Write the code to the output file
	err = os.WriteFile(g.OutputFile, formattedCode, 0o644)
	if err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	return nil
}

// extractTypeParams extracts the type parameters from a type parameter list
func extractTypeParams(typeParams *ast.FieldList) string {
	var params []string
	for _, param := range typeParams.List {
		for _, name := range param.Names {
			paramType := extractType(param.Type)
			params = append(params, name.Name+" "+paramType)
		}
	}
	return "[" + strings.Join(params, ", ") + "]"
}

// extractType extracts the type string from an AST expression
func extractType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		pkg := t.X.(*ast.Ident).Name
		sel := t.Sel.Name
		return pkg + "." + sel
	case *ast.StarExpr:
		return "*" + extractType(t.X)
	case *ast.ArrayType:
		if t.Len == nil {
			return "[]" + extractType(t.Elt)
		}
		return "[" + extractArrayLen(t.Len) + "]" + extractType(t.Elt)
	case *ast.MapType:
		return "map[" + extractType(t.Key) + "]" + extractType(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.IndexExpr:
		return extractType(t.X) + "[" + extractType(t.Index) + "]"
	case *ast.IndexListExpr:
		var indices []string
		for _, index := range t.Indices {
			indices = append(indices, extractType(index))
		}
		return extractType(t.X) + "[" + strings.Join(indices, ", ") + "]"
	default:
		return fmt.Sprintf("%T", expr)
	}
}

// extractArrayLen extracts the length of an array from an AST expression
func extractArrayLen(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.BasicLit:
		return t.Value
	default:
		return ""
	}
}

// hasGenerateDirective checks if the comment group contains our go:generate directive
func hasGenerateDirective(commentGroup *ast.CommentGroup) bool {
	if commentGroup == nil {
		return false
	}

	for _, comment := range commentGroup.List {
		if strings.Contains(comment.Text, "//go:generate") {
			return true
		}
	}

	return false
}

// generateCode generates the validation code for the given structs
func (g *Generator) generateCode(structs []StructInfo) (string, error) {
	funcMap := template.FuncMap{
		"split":      strings.Split,
		"splitN":     strings.SplitN,
		"trimSuffix": strings.TrimSuffix,
		"subtract": func(a, b int) int {
			return a - b
		},
		"extractTypeParamNames": extractTypeParamNames,
	}

	tmpl, err := template.New("validation").Funcs(funcMap).Parse(codeTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"PackageName": g.PackageName,
		"Structs":     structs,
	})
	if err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

// extractTypeParamNames extracts just the type parameter names from a full type parameter string
func extractTypeParamNames(typeParams string) string {
	// Remove the outer brackets
	inner := strings.TrimPrefix(typeParams, "[")
	inner = strings.TrimSuffix(inner, "]")

	// Split by comma to get each parameter
	params := strings.Split(inner, ", ")

	var paramNames []string
	for _, param := range params {
		// Split by space and take the first part (the name)
		parts := strings.Split(param, " ")
		paramNames = append(paramNames, parts[0])
	}

	return "[" + strings.Join(paramNames, ", ") + "]"
}

// Code template for the generated validation code
const codeTemplate = `// Code generated by validation-gen; DO NOT EDIT.

package {{.PackageName}}

import (
	"errors"
)

{{range .Structs}}
// {{.Name}}Params is the parameter struct for creating a {{.Name}}
type {{.Name}}Params{{if .IsGeneric}}{{.TypeParams}}{{end}} struct {
{{- range .Fields}}
	{{.Name}} {{if .IsPointer}}*{{end}}{{.Type}}
{{- end}}
}

// New{{.Name}} creates a new {{.Name}}
func New{{.Name}}{{if .IsGeneric}}{{.TypeParams}}{{end}}(params {{.Name}}Params{{if .IsGeneric}}{{extractTypeParamNames .TypeParams}}{{end}}) (*{{.Name}}{{if .IsGeneric}}{{extractTypeParamNames .TypeParams}}{{end}}, error) {
	if err := isValid{{.Name}}Params{{if .IsGeneric}}{{extractTypeParamNames .TypeParams}}{{end}}(params); err != nil {
		return nil, err
	}

	return &{{.Name}}{{if .IsGeneric}}{{extractTypeParamNames .TypeParams}}{{end}}{
{{- range .Fields}}
		{{.Name}}: params.{{.Name}},
{{- end}}
	}, nil
}

// isValid{{.Name}}Params validates the {{.Name}}Params
func isValid{{.Name}}Params{{if .IsGeneric}}{{.TypeParams}}{{end}}(params {{.Name}}Params{{if .IsGeneric}}{{extractTypeParamNames .TypeParams}}{{end}}) error {
	var errs []error
{{- range .Fields}}
{{- if .IsPointer}}
	if params.{{.Name}} == nil {
		errs = append(errs, errors.New("{{.Name}} is required"))
	}
{{- end}}
{{- end}}
	return errors.Join(errs...)
}
{{end}}

{{define "split"}}{{$s := index . 0}}{{$sep := index . 1}}{{$limit := index . 2}}{{if eq $limit "0"}}{{$s | split $sep}}{{else}}{{$s | splitN $sep $limit}}{{end}}{{end}}

{{define "trimSuffix"}}{{$s := index . 0}}{{$suffix := index . 1}}{{$s | trimSuffix $suffix}}{{end}}

{{define "subtract"}}{{$a := index . 0}}{{$b := index . 1}}{{$a | subtract $b}}{{end}}
`

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/stijmetkii/validation-gen/validation"
)

func main() {
	// Default input is the file that triggered go:generate
	defaultInput := os.Getenv("GOFILE")
	if defaultInput == "" {
		defaultInput = "./example.go"
	}

	// Parse flags
	inputFile := flag.String("input", defaultInput, "Path to the input Go file")
	outputFile := flag.String("output", "", "Path to the output Go file (default is <input>_gen.go)")
	flag.Parse()

	// If the output file is not specified, derive it from the input file
	if *outputFile == "" {
		dir, filename := filepath.Split(*inputFile)
		base := filepath.Base(filename)
		ext := filepath.Ext(base)
		name := base[:len(base)-len(ext)]
		*outputFile = filepath.Join(dir, name+"_gen.go")
	}

	// Create and run the generator
	generator := validation.NewGenerator(*inputFile)
	if *outputFile != "" {
		generator.OutputFile = *outputFile
	}

	if err := generator.Generate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated %s from %s\n", generator.OutputFile, generator.InputFile)
}

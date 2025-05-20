.PHONY: build test clean example force-example

# Build the code generator
build:
	go build -o bin/gen ./cmd/gen

# Install the code generator locally
install:
	go install ./cmd/gen

# Run tests
test:
	go test -v ./...

# Generate example code
example:
	go run ./cmd/gen/main.go -input ./example/service.go
	go run ./cmd/gen/main.go -input ./example/generic_service.go

# Force regeneration of example code
force-example:
	go run ./cmd/gen/main.go -input ./example/service.go -force
	go run ./cmd/gen/main.go -input ./example/generic_service.go -force

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f example/*_gen.go

# Display help information
help:
	@echo "Available targets:"
	@echo "  build        - Build the code generator"
	@echo "  install      - Install the code generator locally"
	@echo "  test         - Run tests"
	@echo "  example      - Generate example code (skips if files exist)"
	@echo "  force-example - Generate example code (overwrite existing files)"
	@echo "  clean        - Clean build artifacts" 
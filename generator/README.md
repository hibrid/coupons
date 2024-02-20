# Generator

The `generator` package provides functionality for generating and validating codes based on specified patterns, character sets, prefixes, and suffixes. This README document serves as a comprehensive guide to understanding the purpose, usage, and examples of the `generator` package.

## Purpose

The `generator` package is designed to facilitate the creation of codes according to user-defined patterns and configurations. It offers the following key features:

- Generation of codes based on specified patterns, character sets, prefixes, and suffixes.
- Validation of generated codes to ensure they conform to the specified patterns and configurations.

## Usage

To use the `generator` package, follow these steps:

1. Import the package:

```go
import "your/package/generator"
```

2. Create a new generator instance with default or custom options using `NewWithOptions` or `Default` functions:

```go
// Create a generator with default options
g, err := generator.Default()
if err != nil {
    // Handle error
}

// Create a generator with custom options
g, err := generator.NewWithOptions(
    generator.SetMinimumLength(6),
    generator.SetGenerateCount(10),
    generator.SetPattern("####-####"),
)
if err != nil {
    // Handle error
}
```

3. Generate codes using the `Run` method:

```go
codes, err := g.Run()
if err != nil {
    // Handle error
}
```

4. Validate codes using the `Validate` method:

```go
validatedCode, err := g.Validate("ABCD-EFGH")
if err != nil {
    // Handle error
}
```

## Examples

### Generating Codes

```go
// Create a generator with custom options
g, _ := generator.NewWithOptions(
    generator.SetMinimumLength(6),
    generator.SetGenerateCount(5),
    generator.SetPattern("####-####"),
)

// Generate codes
codes, _ := g.Run()
```

### Validating Codes

```go
// Create a generator with custom options
g, _ := generator.NewWithOptions(
    generator.SetPattern("####-####"),
)

// Validate a code
validatedCode, _ := g.Validate("ABCD-EFGH")
```

## Use Cases

1. **Basic Code Generation**: Generate codes without any specified pattern, using default options.

2. **Custom Pattern Generation**: Generate codes with custom patterns and lengths.

3. **Validation**: Validate generated codes to ensure they meet specified criteria.

4. **Prefix and Suffix Handling**: Generate codes with prefixes and suffixes appended.

5. **Combination Counting**: Determine the number of unique combinations possible for a given pattern.

## Running Tests

To run tests for the `generator` package, execute the following command:

```bash
go test -v
```
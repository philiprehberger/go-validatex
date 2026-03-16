# go-validatex

[![CI](https://github.com/philiprehberger/go-validatex/actions/workflows/ci.yml/badge.svg)](https://github.com/philiprehberger/go-validatex/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/philiprehberger/go-validatex.svg)](https://pkg.go.dev/github.com/philiprehberger/go-validatex)
[![License](https://img.shields.io/github/license/philiprehberger/go-validatex)](LICENSE)

Struct validation library for Go using struct tags. Zero external dependencies.

## Installation

```bash
go get github.com/philiprehberger/go-validatex
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/philiprehberger/go-validatex"
)

type User struct {
	Name  string `validate:"required,min=3,max=50"`
	Email string `validate:"required,email"`
	Role  string `validate:"oneof=admin|user|guest"`
	Age   int    `validate:"min=18,max=120"`
}

func main() {
	u := User{Name: "Al", Email: "bad", Role: "superadmin", Age: 10}
	err := validatex.Validate(u)
	if err != nil {
		for _, ve := range validatex.Errors(err) {
			fmt.Printf("Field %s (%s): %s\n", ve.Field, ve.Rule, ve.Message)
		}
	}
}
```

### Custom Rules

```go
validatex.Register("even", func(value any, param string) error {
	v, ok := value.(int)
	if !ok {
		return fmt.Errorf("expected int")
	}
	if v%2 != 0 {
		return fmt.Errorf("must be even")
	}
	return nil
})

type Config struct {
	Workers int `validate:"even,min=2"`
}
```

## Built-in Rules

| Rule | Applies To | Description |
|------|-----------|-------------|
| `required` | string, int, float | Non-empty string or non-zero number |
| `min=N` | string, int, float | Min length (string) or min value (number) |
| `max=N` | string, int, float | Max length (string) or max value (number) |
| `len=N` | string | Exact length |
| `email` | string | Must contain `@` and `.` |
| `url` | string | Must start with `http://` or `https://` |
| `oneof=a\|b\|c` | string | Must be one of the listed values |
| `pattern=regex` | string | Must match the regular expression |

## API

| Function / Type | Description |
|-----------------|-------------|
| `Validate(v any) error` | Validate struct fields using `validate` tags |
| `Register(name, fn)` | Register a custom validation rule |
| `Errors(err) []ValidationError` | Extract individual errors from a validation error |
| `ValidationError` | Single field validation failure (Field, Rule, Message) |
| `ValidationErrors` | Slice of ValidationError, implements `error` |

## Development

```bash
go test ./...
go vet ./...
```

## License

MIT

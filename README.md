# go-validatex

[![CI](https://github.com/philiprehberger/go-validatex/actions/workflows/ci.yml/badge.svg)](https://github.com/philiprehberger/go-validatex/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/philiprehberger/go-validatex.svg)](https://pkg.go.dev/github.com/philiprehberger/go-validatex)
[![Last updated](https://img.shields.io/github/last-commit/philiprehberger/go-validatex)](https://github.com/philiprehberger/go-validatex/commits/main)

Struct validation library for Go using struct tags. Zero external dependencies

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

### Custom Messages

Use the `msg` tag option to provide a custom error message for the preceding rule:

```go
type User struct {
	Name  string `validate:"required,msg=Name is required"`
	Email string `validate:"email,msg=Please enter a valid email"`
}
```

### Nested Structs

Nested structs with `validate` tags are validated automatically. Error field paths use dot notation:

```go
type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
}

type User struct {
	Name    string  `validate:"required"`
	Address Address
}

err := validatex.Validate(User{Name: "Alice", Address: Address{}})
// errors: Address.Street: is required, Address.City: is required
```

### ValidateField

Validate a single value against a rules string without defining a struct:

```go
err := validatex.ValidateField("bad", "email")
err = validatex.ValidateField("", "required,min=3")
err = validatex.ValidateField(42, "min=10,max=100")
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

### Built-in Rules

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
| `uuid` | string | Must be a valid UUID (8-4-4-4-12 hex format) |
| `ip` | string | Must be a valid IPv4 or IPv6 address |
| `ipv4` | string | Must be a valid IPv4 address |
| `ipv6` | string | Must be a valid IPv6 address |
| `alpha` | string | Must contain only ASCII letters |
| `numeric` | string | Must contain only ASCII digits |
| `alphanum` | string | Must contain only ASCII letters and digits |
| `contains=X` | string | Must contain substring X |
| `excludes=X` | string | Must not contain substring X |

## API

| Function / Type | Description |
|-----------------|-------------|
| `Validate(v any) error` | Validate struct fields using `validate` tags (supports nested structs) |
| `ValidateField(value any, rules string) error` | Validate a single value against a rules string |
| `Register(name, fn)` | Register a custom validation rule |
| `Errors(err) []ValidationError` | Extract individual errors from a validation error |
| `ValidationError` | Single field validation failure (Field, Rule, Message) |
| `ValidationErrors` | Slice of ValidationError, implements `error` |

## Development

```bash
go test ./...
go vet ./...
```

## Support

If you find this project useful:

⭐ [Star the repo](https://github.com/philiprehberger/go-validatex)

🐛 [Report issues](https://github.com/philiprehberger/go-validatex/issues?q=is%3Aissue+is%3Aopen+label%3Abug)

💡 [Suggest features](https://github.com/philiprehberger/go-validatex/issues?q=is%3Aissue+is%3Aopen+label%3Aenhancement)

❤️ [Sponsor development](https://github.com/sponsors/philiprehberger)

🌐 [All Open Source Projects](https://philiprehberger.com/open-source-packages)

💻 [GitHub Profile](https://github.com/philiprehberger)

🔗 [LinkedIn Profile](https://www.linkedin.com/in/philiprehberger)

## License

[MIT](LICENSE)

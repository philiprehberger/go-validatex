package validatex

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ruleFunc is the internal signature for a validation rule. It receives the
// reflected field value, the field name (for error reporting), and the
// parameter string from the tag (empty string when the rule takes no param).
// It returns a *ValidationError on failure or nil on success.
type ruleFunc func(fieldValue reflect.Value, fieldName string, param string) *ValidationError

// registry holds all registered validation rules keyed by name.
var registry = map[string]ruleFunc{}

func init() {
	registerBuiltins()
}

// Register adds a custom validation rule that can be referenced by name in
// struct tags. The provided function receives the field value as any and the
// parameter string from the tag. Return a non-nil error to indicate a
// validation failure.
func Register(name string, fn func(value any, param string) error) {
	registry[name] = func(fv reflect.Value, fieldName string, param string) *ValidationError {
		err := fn(fv.Interface(), param)
		if err != nil {
			return &ValidationError{
				Field:   fieldName,
				Rule:    name,
				Message: err.Error(),
			}
		}
		return nil
	}
}

func registerBuiltins() {
	registry["required"] = ruleRequired
	registry["min"] = ruleMin
	registry["max"] = ruleMax
	registry["email"] = ruleEmail
	registry["url"] = ruleURL
	registry["oneof"] = ruleOneof
	registry["len"] = ruleLen
	registry["pattern"] = rulePattern
}

func ruleRequired(fv reflect.Value, fieldName string, _ string) *ValidationError {
	failed := false
	switch fv.Kind() {
	case reflect.String:
		if fv.String() == "" {
			failed = true
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fv.Int() == 0 {
			failed = true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if fv.Uint() == 0 {
			failed = true
		}
	case reflect.Float32, reflect.Float64:
		if fv.Float() == 0 {
			failed = true
		}
	default:
		if fv.IsZero() {
			failed = true
		}
	}
	if failed {
		return &ValidationError{Field: fieldName, Rule: "required", Message: "is required"}
	}
	return nil
}

func ruleMin(fv reflect.Value, fieldName string, param string) *ValidationError {
	n, err := strconv.Atoi(param)
	if err != nil {
		return &ValidationError{Field: fieldName, Rule: "min", Message: fmt.Sprintf("invalid min param: %s", param)}
	}
	switch fv.Kind() {
	case reflect.String:
		if len(fv.String()) < n {
			return &ValidationError{Field: fieldName, Rule: "min", Message: fmt.Sprintf("must be at least %d characters", n)}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fv.Int() < int64(n) {
			return &ValidationError{Field: fieldName, Rule: "min", Message: fmt.Sprintf("must be at least %d", n)}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if fv.Uint() < uint64(n) {
			return &ValidationError{Field: fieldName, Rule: "min", Message: fmt.Sprintf("must be at least %d", n)}
		}
	case reflect.Float32, reflect.Float64:
		if fv.Float() < float64(n) {
			return &ValidationError{Field: fieldName, Rule: "min", Message: fmt.Sprintf("must be at least %d", n)}
		}
	}
	return nil
}

func ruleMax(fv reflect.Value, fieldName string, param string) *ValidationError {
	n, err := strconv.Atoi(param)
	if err != nil {
		return &ValidationError{Field: fieldName, Rule: "max", Message: fmt.Sprintf("invalid max param: %s", param)}
	}
	switch fv.Kind() {
	case reflect.String:
		if len(fv.String()) > n {
			return &ValidationError{Field: fieldName, Rule: "max", Message: fmt.Sprintf("must be at most %d characters", n)}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fv.Int() > int64(n) {
			return &ValidationError{Field: fieldName, Rule: "max", Message: fmt.Sprintf("must be at most %d", n)}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if fv.Uint() > uint64(n) {
			return &ValidationError{Field: fieldName, Rule: "max", Message: fmt.Sprintf("must be at most %d", n)}
		}
	case reflect.Float32, reflect.Float64:
		if fv.Float() > float64(n) {
			return &ValidationError{Field: fieldName, Rule: "max", Message: fmt.Sprintf("must be at most %d", n)}
		}
	}
	return nil
}

func ruleEmail(fv reflect.Value, fieldName string, _ string) *ValidationError {
	if fv.Kind() != reflect.String {
		return nil
	}
	s := fv.String()
	if !strings.Contains(s, "@") || !strings.Contains(s, ".") {
		return &ValidationError{Field: fieldName, Rule: "email", Message: "must be a valid email address"}
	}
	return nil
}

func ruleURL(fv reflect.Value, fieldName string, _ string) *ValidationError {
	if fv.Kind() != reflect.String {
		return nil
	}
	s := fv.String()
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		return &ValidationError{Field: fieldName, Rule: "url", Message: "must be a valid URL"}
	}
	return nil
}

func ruleOneof(fv reflect.Value, fieldName string, param string) *ValidationError {
	if fv.Kind() != reflect.String {
		return nil
	}
	options := strings.Split(param, "|")
	s := fv.String()
	for _, opt := range options {
		if s == opt {
			return nil
		}
	}
	return &ValidationError{Field: fieldName, Rule: "oneof", Message: fmt.Sprintf("must be one of: %s", strings.Join(options, ", "))}
}

func ruleLen(fv reflect.Value, fieldName string, param string) *ValidationError {
	n, err := strconv.Atoi(param)
	if err != nil {
		return &ValidationError{Field: fieldName, Rule: "len", Message: fmt.Sprintf("invalid len param: %s", param)}
	}
	if fv.Kind() != reflect.String {
		return nil
	}
	if len(fv.String()) != n {
		return &ValidationError{Field: fieldName, Rule: "len", Message: fmt.Sprintf("must be exactly %d characters", n)}
	}
	return nil
}

func rulePattern(fv reflect.Value, fieldName string, param string) *ValidationError {
	if fv.Kind() != reflect.String {
		return nil
	}
	re, err := regexp.Compile(param)
	if err != nil {
		return &ValidationError{Field: fieldName, Rule: "pattern", Message: fmt.Sprintf("invalid pattern: %s", param)}
	}
	if !re.MatchString(fv.String()) {
		return &ValidationError{Field: fieldName, Rule: "pattern", Message: fmt.Sprintf("must match pattern %s", param)}
	}
	return nil
}

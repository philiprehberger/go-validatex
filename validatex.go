// Package validatex provides struct validation via tags for Go.
package validatex

import (
	"fmt"
	"reflect"
	"strings"
)

// ValidationError represents a single validation failure for a struct field.
type ValidationError struct {
	Field   string
	Rule    string
	Message string
}

// Error returns a human-readable representation of the validation error in the
// format "Field: Message".
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors is a slice of ValidationError values returned when one or
// more fields fail validation.
type ValidationErrors []ValidationError

// Error returns a multi-line string listing every validation error.
func (e ValidationErrors) Error() string {
	lines := make([]string, len(e))
	for i, ve := range e {
		lines[i] = ve.Error()
	}
	return strings.Join(lines, "\n")
}

// Validate checks the exported fields of a struct against rules declared in
// their `validate` tags. It returns a ValidationErrors value containing all
// failures, or nil when every field passes. Non-struct values (and non-pointer-
// to-struct values) produce an error.
func Validate(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return fmt.Errorf("validatex: cannot validate nil pointer")
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("validatex: expected struct, got %s", rv.Kind())
	}

	rt := rv.Type()
	var errs ValidationErrors

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		fv := rv.Field(i)
		rules := parseTag(tag)

		for _, rule := range rules {
			fn, ok := registry[rule.name]
			if !ok {
				errs = append(errs, ValidationError{
					Field:   field.Name,
					Rule:    rule.name,
					Message: fmt.Sprintf("unknown rule: %s", rule.name),
				})
				continue
			}
			if ve := fn(fv, field.Name, rule.param); ve != nil {
				errs = append(errs, *ve)
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// Errors extracts individual ValidationError values from an error. If the
// error is a ValidationErrors, each element is returned. Otherwise nil is
// returned.
func Errors(err error) []ValidationError {
	if err == nil {
		return nil
	}
	if ve, ok := err.(ValidationErrors); ok {
		return []ValidationError(ve)
	}
	return nil
}

// rule is a parsed tag directive such as "min=3" (name="min", param="3") or
// "required" (name="required", param="").
type rule struct {
	name  string
	param string
}

// parseTag splits a comma-separated tag value into individual rules, each
// optionally carrying a parameter after the first "=".
func parseTag(tag string) []rule {
	parts := strings.Split(tag, ",")
	rules := make([]rule, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		name, param, _ := strings.Cut(p, "=")
		rules = append(rules, rule{name: name, param: param})
	}
	return rules
}

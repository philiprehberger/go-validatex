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
// to-struct values) produce an error. Nested structs with validate tags are
// validated recursively, with error field paths using dot notation.
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

	errs := validateStruct(rv, "")
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// validateStruct performs recursive struct validation. The prefix is prepended
// to field names using dot notation for nested structs.
func validateStruct(rv reflect.Value, prefix string) ValidationErrors {
	rt := rv.Type()
	var errs ValidationErrors

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fv := rv.Field(i)

		fieldName := field.Name
		if prefix != "" {
			fieldName = prefix + "." + field.Name
		}

		// Recurse into nested structs that have validate tags on their fields.
		fk := fv.Kind()
		if fk == reflect.Struct && field.Type != reflect.TypeOf(struct{}{}) {
			if hasValidateTags(field.Type) {
				nested := validateStruct(fv, fieldName)
				errs = append(errs, nested...)
			}
		}

		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		rules := parseTag(tag)

		for _, rule := range rules {
			fn, ok := registry[rule.name]
			if !ok {
				errs = append(errs, ValidationError{
					Field:   fieldName,
					Rule:    rule.name,
					Message: fmt.Sprintf("unknown rule: %s", rule.name),
				})
				continue
			}
			if ve := fn(fv, fieldName, rule.param); ve != nil {
				if rule.msg != "" {
					ve.Message = rule.msg
				}
				errs = append(errs, *ve)
			}
		}
	}

	return errs
}

// hasValidateTags reports whether a struct type has any fields with validate tags.
func hasValidateTags(t reflect.Type) bool {
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get("validate") != "" {
			return true
		}
		// Check nested structs recursively.
		ft := t.Field(i).Type
		if ft.Kind() == reflect.Struct && hasValidateTags(ft) {
			return true
		}
	}
	return false
}

// ValidateField validates a single value against a rules string (e.g.
// "required,min=3,max=50"). It returns a ValidationErrors value containing
// all failures, or nil when the value passes all rules.
func ValidateField(value any, rules string) error {
	fv := reflect.ValueOf(value)
	parsed := parseTag(rules)
	var errs ValidationErrors

	for _, r := range parsed {
		fn, ok := registry[r.name]
		if !ok {
			errs = append(errs, ValidationError{
				Field:   "value",
				Rule:    r.name,
				Message: fmt.Sprintf("unknown rule: %s", r.name),
			})
			continue
		}
		if ve := fn(fv, "value", r.param); ve != nil {
			if r.msg != "" {
				ve.Message = r.msg
			}
			errs = append(errs, *ve)
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
	msg   string // custom error message from msg= tag option
}

// parseTag splits a comma-separated tag value into individual rules, each
// optionally carrying a parameter after the first "=". A "msg=..." option
// applies to the immediately preceding rule.
func parseTag(tag string) []rule {
	parts := strings.Split(tag, ",")
	rules := make([]rule, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		name, param, _ := strings.Cut(p, "=")
		if name == "msg" && len(rules) > 0 {
			rules[len(rules)-1].msg = param
			continue
		}
		rules = append(rules, rule{name: name, param: param})
	}
	return rules
}

package validatex

import (
	"fmt"
	"net"
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
	registry["uuid"] = ruleUUID
	registry["ip"] = ruleIP
	registry["ipv4"] = ruleIPv4
	registry["ipv6"] = ruleIPv6
	registry["alpha"] = ruleAlpha
	registry["numeric"] = ruleNumeric
	registry["alphanum"] = ruleAlphanum
	registry["contains"] = ruleContains
	registry["excludes"] = ruleExcludes
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

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func ruleUUID(fv reflect.Value, fieldName string, _ string) *ValidationError {
	if fv.Kind() != reflect.String {
		return nil
	}
	if !uuidRegex.MatchString(fv.String()) {
		return &ValidationError{Field: fieldName, Rule: "uuid", Message: "must be a valid UUID"}
	}
	return nil
}

func ruleIP(fv reflect.Value, fieldName string, _ string) *ValidationError {
	if fv.Kind() != reflect.String {
		return nil
	}
	if net.ParseIP(fv.String()) == nil {
		return &ValidationError{Field: fieldName, Rule: "ip", Message: "must be a valid IP address"}
	}
	return nil
}

func ruleIPv4(fv reflect.Value, fieldName string, _ string) *ValidationError {
	if fv.Kind() != reflect.String {
		return nil
	}
	ip := net.ParseIP(fv.String())
	if ip == nil || ip.To4() == nil {
		return &ValidationError{Field: fieldName, Rule: "ipv4", Message: "must be a valid IPv4 address"}
	}
	return nil
}

func ruleIPv6(fv reflect.Value, fieldName string, _ string) *ValidationError {
	if fv.Kind() != reflect.String {
		return nil
	}
	ip := net.ParseIP(fv.String())
	if ip == nil || ip.To4() != nil {
		return &ValidationError{Field: fieldName, Rule: "ipv6", Message: "must be a valid IPv6 address"}
	}
	return nil
}

func ruleAlpha(fv reflect.Value, fieldName string, _ string) *ValidationError {
	if fv.Kind() != reflect.String {
		return nil
	}
	for _, r := range fv.String() {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
			return &ValidationError{Field: fieldName, Rule: "alpha", Message: "must contain only letters"}
		}
	}
	return nil
}

func ruleNumeric(fv reflect.Value, fieldName string, _ string) *ValidationError {
	if fv.Kind() != reflect.String {
		return nil
	}
	for _, r := range fv.String() {
		if r < '0' || r > '9' {
			return &ValidationError{Field: fieldName, Rule: "numeric", Message: "must contain only digits"}
		}
	}
	return nil
}

func ruleAlphanum(fv reflect.Value, fieldName string, _ string) *ValidationError {
	if fv.Kind() != reflect.String {
		return nil
	}
	for _, r := range fv.String() {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return &ValidationError{Field: fieldName, Rule: "alphanum", Message: "must contain only letters and digits"}
		}
	}
	return nil
}

func ruleContains(fv reflect.Value, fieldName string, param string) *ValidationError {
	if fv.Kind() != reflect.String {
		return nil
	}
	if !strings.Contains(fv.String(), param) {
		return &ValidationError{Field: fieldName, Rule: "contains", Message: fmt.Sprintf("must contain %q", param)}
	}
	return nil
}

func ruleExcludes(fv reflect.Value, fieldName string, param string) *ValidationError {
	if fv.Kind() != reflect.String {
		return nil
	}
	if strings.Contains(fv.String(), param) {
		return &ValidationError{Field: fieldName, Rule: "excludes", Message: fmt.Sprintf("must not contain %q", param)}
	}
	return nil
}

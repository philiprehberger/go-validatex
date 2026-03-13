package validatex

import (
	"fmt"
	"strings"
	"testing"
)

func TestRequiredStringEmpty(t *testing.T) {
	type S struct {
		Name string `validate:"required"`
	}
	err := Validate(S{Name: ""})
	if err == nil {
		t.Fatal("expected error for empty required string")
	}
	errs := Errors(err)
	if len(errs) != 1 || errs[0].Rule != "required" {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestRequiredStringNonEmpty(t *testing.T) {
	type S struct {
		Name string `validate:"required"`
	}
	if err := Validate(S{Name: "alice"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequiredIntZero(t *testing.T) {
	type S struct {
		Age int `validate:"required"`
	}
	err := Validate(S{Age: 0})
	if err == nil {
		t.Fatal("expected error for zero required int")
	}
}

func TestRequiredIntNonZero(t *testing.T) {
	type S struct {
		Age int `validate:"required"`
	}
	if err := Validate(S{Age: 25}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMinStringLength(t *testing.T) {
	type S struct {
		Name string `validate:"min=3"`
	}
	if err := Validate(S{Name: "ab"}); err == nil {
		t.Fatal("expected error for short string")
	}
	if err := Validate(S{Name: "abc"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMaxStringLength(t *testing.T) {
	type S struct {
		Name string `validate:"max=5"`
	}
	if err := Validate(S{Name: "abcdef"}); err == nil {
		t.Fatal("expected error for long string")
	}
	if err := Validate(S{Name: "abcde"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMinNumericValue(t *testing.T) {
	type S struct {
		Age int `validate:"min=18"`
	}
	if err := Validate(S{Age: 10}); err == nil {
		t.Fatal("expected error for value below min")
	}
	if err := Validate(S{Age: 18}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMaxNumericValue(t *testing.T) {
	type S struct {
		Age int `validate:"max=100"`
	}
	if err := Validate(S{Age: 101}); err == nil {
		t.Fatal("expected error for value above max")
	}
	if err := Validate(S{Age: 100}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEmailValid(t *testing.T) {
	type S struct {
		Email string `validate:"email"`
	}
	if err := Validate(S{Email: "user@example.com"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEmailInvalid(t *testing.T) {
	type S struct {
		Email string `validate:"email"`
	}
	if err := Validate(S{Email: "not-an-email"}); err == nil {
		t.Fatal("expected error for invalid email")
	}
}

func TestURLValid(t *testing.T) {
	type S struct {
		URL string `validate:"url"`
	}
	if err := Validate(S{URL: "https://example.com"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := Validate(S{URL: "http://example.com"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestURLInvalid(t *testing.T) {
	type S struct {
		URL string `validate:"url"`
	}
	if err := Validate(S{URL: "ftp://example.com"}); err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestOneofValid(t *testing.T) {
	type S struct {
		Role string `validate:"oneof=admin|user|guest"`
	}
	if err := Validate(S{Role: "admin"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOneofInvalid(t *testing.T) {
	type S struct {
		Role string `validate:"oneof=admin|user|guest"`
	}
	if err := Validate(S{Role: "superadmin"}); err == nil {
		t.Fatal("expected error for value not in oneof")
	}
}

func TestLenExact(t *testing.T) {
	type S struct {
		Code string `validate:"len=5"`
	}
	if err := Validate(S{Code: "12345"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := Validate(S{Code: "1234"}); err == nil {
		t.Fatal("expected error for wrong length")
	}
	if err := Validate(S{Code: "123456"}); err == nil {
		t.Fatal("expected error for wrong length")
	}
}

func TestPatternValid(t *testing.T) {
	type S struct {
		Zip string `validate:"pattern=^[0-9]{5}$"`
	}
	if err := Validate(S{Zip: "12345"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPatternInvalid(t *testing.T) {
	type S struct {
		Zip string `validate:"pattern=^[0-9]{5}$"`
	}
	if err := Validate(S{Zip: "abcde"}); err == nil {
		t.Fatal("expected error for non-matching pattern")
	}
}

func TestMultipleRulesOnOneField(t *testing.T) {
	type S struct {
		Name string `validate:"required,min=3,max=10"`
	}
	err := Validate(S{Name: ""})
	if err == nil {
		t.Fatal("expected errors")
	}
	errs := Errors(err)
	// empty string fails required and min=3
	if len(errs) < 2 {
		t.Fatalf("expected at least 2 errors, got %d: %v", len(errs), errs)
	}
}

func TestMultipleFieldsWithErrors(t *testing.T) {
	type S struct {
		Name  string `validate:"required"`
		Email string `validate:"email"`
	}
	err := Validate(S{Name: "", Email: "bad"})
	if err == nil {
		t.Fatal("expected errors")
	}
	errs := Errors(err)
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(errs), errs)
	}
	fields := map[string]bool{}
	for _, e := range errs {
		fields[e.Field] = true
	}
	if !fields["Name"] || !fields["Email"] {
		t.Fatalf("expected errors for Name and Email, got %v", errs)
	}
}

func TestNoValidateTagSkipsField(t *testing.T) {
	type S struct {
		Name    string `validate:"required"`
		Skipped string
	}
	if err := Validate(S{Name: "ok", Skipped: ""}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNonStructReturnsError(t *testing.T) {
	err := Validate("not a struct")
	if err == nil {
		t.Fatal("expected error for non-struct input")
	}
	if Errors(err) != nil {
		t.Fatal("non-struct error should not be ValidationErrors")
	}
}

func TestPointerToStructWorks(t *testing.T) {
	type S struct {
		Name string `validate:"required"`
	}
	s := &S{Name: "alice"}
	if err := Validate(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s2 := &S{Name: ""}
	if err := Validate(s2); err == nil {
		t.Fatal("expected error for empty required string via pointer")
	}
}

func TestCustomRuleViaRegister(t *testing.T) {
	Register("even", func(value any, param string) error {
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("expected int")
		}
		if v%2 != 0 {
			return fmt.Errorf("must be even")
		}
		return nil
	})

	type S struct {
		Num int `validate:"even"`
	}
	if err := Validate(S{Num: 4}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := Validate(S{Num: 3}); err == nil {
		t.Fatal("expected error for odd number")
	}
}

func TestValidationErrorsErrorFormat(t *testing.T) {
	errs := ValidationErrors{
		{Field: "Name", Rule: "required", Message: "is required"},
		{Field: "Email", Rule: "email", Message: "must be a valid email address"},
	}
	s := errs.Error()
	if !strings.Contains(s, "Name: is required") {
		t.Fatalf("expected Name error in output, got: %s", s)
	}
	if !strings.Contains(s, "Email: must be a valid email address") {
		t.Fatalf("expected Email error in output, got: %s", s)
	}
	lines := strings.Split(s, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %s", len(lines), s)
	}
}

func TestErrorsHelperExtractsIndividualErrors(t *testing.T) {
	errs := ValidationErrors{
		{Field: "A", Rule: "required", Message: "is required"},
		{Field: "B", Rule: "min", Message: "must be at least 3 characters"},
	}
	extracted := Errors(errs)
	if len(extracted) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(extracted))
	}
	if extracted[0].Field != "A" || extracted[1].Field != "B" {
		t.Fatalf("unexpected fields: %v", extracted)
	}
}

func TestErrorsHelperNilReturnsNil(t *testing.T) {
	if Errors(nil) != nil {
		t.Fatal("expected nil for nil error")
	}
}

func TestErrorsHelperNonValidationErrorsReturnsNil(t *testing.T) {
	if Errors(fmt.Errorf("some error")) != nil {
		t.Fatal("expected nil for non-ValidationErrors")
	}
}

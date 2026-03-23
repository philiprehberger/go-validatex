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

// --- UUID rule tests ---

func TestUUIDValid(t *testing.T) {
	type S struct {
		ID string `validate:"uuid"`
	}
	if err := Validate(S{ID: "550e8400-e29b-41d4-a716-446655440000"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUUIDInvalid(t *testing.T) {
	type S struct {
		ID string `validate:"uuid"`
	}
	cases := []string{"not-a-uuid", "550e8400e29b41d4a716446655440000", "550e8400-e29b-41d4-a716"}
	for _, c := range cases {
		if err := Validate(S{ID: c}); err == nil {
			t.Fatalf("expected error for invalid UUID: %s", c)
		}
	}
}

// --- IP rule tests ---

func TestIPValid(t *testing.T) {
	type S struct {
		Addr string `validate:"ip"`
	}
	for _, addr := range []string{"192.168.1.1", "::1", "2001:db8::1"} {
		if err := Validate(S{Addr: addr}); err != nil {
			t.Fatalf("unexpected error for %s: %v", addr, err)
		}
	}
}

func TestIPInvalid(t *testing.T) {
	type S struct {
		Addr string `validate:"ip"`
	}
	if err := Validate(S{Addr: "not-an-ip"}); err == nil {
		t.Fatal("expected error for invalid IP")
	}
}

// --- IPv4 rule tests ---

func TestIPv4Valid(t *testing.T) {
	type S struct {
		Addr string `validate:"ipv4"`
	}
	if err := Validate(S{Addr: "192.168.1.1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIPv4RejectsIPv6(t *testing.T) {
	type S struct {
		Addr string `validate:"ipv4"`
	}
	if err := Validate(S{Addr: "::1"}); err == nil {
		t.Fatal("expected error for IPv6 address with ipv4 rule")
	}
}

func TestIPv4Invalid(t *testing.T) {
	type S struct {
		Addr string `validate:"ipv4"`
	}
	if err := Validate(S{Addr: "not-an-ip"}); err == nil {
		t.Fatal("expected error for invalid IP with ipv4 rule")
	}
}

// --- IPv6 rule tests ---

func TestIPv6Valid(t *testing.T) {
	type S struct {
		Addr string `validate:"ipv6"`
	}
	if err := Validate(S{Addr: "::1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := Validate(S{Addr: "2001:db8::1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIPv6RejectsIPv4(t *testing.T) {
	type S struct {
		Addr string `validate:"ipv6"`
	}
	if err := Validate(S{Addr: "192.168.1.1"}); err == nil {
		t.Fatal("expected error for IPv4 address with ipv6 rule")
	}
}

// --- Alpha rule tests ---

func TestAlphaValid(t *testing.T) {
	type S struct {
		Name string `validate:"alpha"`
	}
	if err := Validate(S{Name: "HelloWorld"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAlphaInvalid(t *testing.T) {
	type S struct {
		Name string `validate:"alpha"`
	}
	for _, v := range []string{"hello123", "hello world", "hello!"} {
		if err := Validate(S{Name: v}); err == nil {
			t.Fatalf("expected error for non-alpha string: %s", v)
		}
	}
}

func TestAlphaEmpty(t *testing.T) {
	type S struct {
		Name string `validate:"alpha"`
	}
	// Empty string has no non-alpha chars, should pass
	if err := Validate(S{Name: ""}); err != nil {
		t.Fatalf("unexpected error for empty string: %v", err)
	}
}

// --- Numeric rule tests ---

func TestNumericValid(t *testing.T) {
	type S struct {
		Code string `validate:"numeric"`
	}
	if err := Validate(S{Code: "12345"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNumericInvalid(t *testing.T) {
	type S struct {
		Code string `validate:"numeric"`
	}
	if err := Validate(S{Code: "123abc"}); err == nil {
		t.Fatal("expected error for non-numeric string")
	}
}

// --- Alphanum rule tests ---

func TestAlphanumValid(t *testing.T) {
	type S struct {
		Code string `validate:"alphanum"`
	}
	if err := Validate(S{Code: "abc123"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAlphanumInvalid(t *testing.T) {
	type S struct {
		Code string `validate:"alphanum"`
	}
	if err := Validate(S{Code: "abc-123"}); err == nil {
		t.Fatal("expected error for non-alphanum string")
	}
}

// --- Contains rule tests ---

func TestContainsValid(t *testing.T) {
	type S struct {
		Bio string `validate:"contains=hello"`
	}
	if err := Validate(S{Bio: "say hello world"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestContainsInvalid(t *testing.T) {
	type S struct {
		Bio string `validate:"contains=hello"`
	}
	if err := Validate(S{Bio: "say hi world"}); err == nil {
		t.Fatal("expected error when string does not contain substring")
	}
}

// --- Excludes rule tests ---

func TestExcludesValid(t *testing.T) {
	type S struct {
		Bio string `validate:"excludes=badword"`
	}
	if err := Validate(S{Bio: "this is clean"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExcludesInvalid(t *testing.T) {
	type S struct {
		Bio string `validate:"excludes=badword"`
	}
	if err := Validate(S{Bio: "contains badword here"}); err == nil {
		t.Fatal("expected error when string contains excluded substring")
	}
}

// --- Custom error message (msg) tests ---

func TestCustomMessageOnRequired(t *testing.T) {
	type S struct {
		Name string `validate:"required,msg=Name is required"`
	}
	err := Validate(S{Name: ""})
	if err == nil {
		t.Fatal("expected error")
	}
	errs := Errors(err)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Message != "Name is required" {
		t.Fatalf("expected custom message, got: %s", errs[0].Message)
	}
}

func TestCustomMessageOnMin(t *testing.T) {
	type S struct {
		Name string `validate:"min=3,msg=Too short"`
	}
	err := Validate(S{Name: "ab"})
	if err == nil {
		t.Fatal("expected error")
	}
	errs := Errors(err)
	if errs[0].Message != "Too short" {
		t.Fatalf("expected custom message, got: %s", errs[0].Message)
	}
}

func TestCustomMessageDoesNotAffectOtherRules(t *testing.T) {
	type S struct {
		Name string `validate:"required,msg=Name is required,min=3"`
	}
	// Empty string fails both required (custom msg) and min (default msg)
	err := Validate(S{Name: ""})
	if err == nil {
		t.Fatal("expected error")
	}
	errs := Errors(err)
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(errs), errs)
	}
	if errs[0].Message != "Name is required" {
		t.Fatalf("expected custom message on first error, got: %s", errs[0].Message)
	}
	if errs[1].Message == "Name is required" {
		t.Fatal("second error should not have the custom message from the first rule")
	}
}

// --- Nested struct validation tests ---

func TestNestedStructValidation(t *testing.T) {
	type Address struct {
		Street string `validate:"required"`
		City   string `validate:"required"`
	}
	type User struct {
		Name    string  `validate:"required"`
		Address Address
	}
	err := Validate(User{Name: "Alice", Address: Address{Street: "", City: ""}})
	if err == nil {
		t.Fatal("expected errors for nested struct")
	}
	errs := Errors(err)
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(errs), errs)
	}
	fields := map[string]bool{}
	for _, e := range errs {
		fields[e.Field] = true
	}
	if !fields["Address.Street"] || !fields["Address.City"] {
		t.Fatalf("expected dot-notation fields, got: %v", errs)
	}
}

func TestNestedStructValid(t *testing.T) {
	type Address struct {
		Street string `validate:"required"`
	}
	type User struct {
		Name    string  `validate:"required"`
		Address Address
	}
	if err := Validate(User{Name: "Alice", Address: Address{Street: "123 Main"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeeplyNestedStruct(t *testing.T) {
	type Zip struct {
		Code string `validate:"required,len=5"`
	}
	type Address struct {
		Zip Zip
	}
	type User struct {
		Address Address
	}
	err := Validate(User{Address: Address{Zip: Zip{Code: "123"}}})
	if err == nil {
		t.Fatal("expected error for deeply nested struct")
	}
	errs := Errors(err)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if errs[0].Field != "Address.Zip.Code" {
		t.Fatalf("expected Address.Zip.Code, got: %s", errs[0].Field)
	}
}

func TestNestedStructWithoutValidateTagsIsSkipped(t *testing.T) {
	type Meta struct {
		Internal string
	}
	type User struct {
		Name string `validate:"required"`
		Meta Meta
	}
	// Meta has no validate tags, should not cause issues
	if err := Validate(User{Name: "Alice"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- ValidateField tests ---

func TestValidateFieldRequired(t *testing.T) {
	if err := ValidateField("hello", "required"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := ValidateField("", "required"); err == nil {
		t.Fatal("expected error for empty required field")
	}
}

func TestValidateFieldMultipleRules(t *testing.T) {
	if err := ValidateField("ab", "required,min=3"); err == nil {
		t.Fatal("expected error")
	}
	errs := Errors(ValidateField("ab", "required,min=3"))
	if len(errs) != 1 {
		t.Fatalf("expected 1 error (min), got %d: %v", len(errs), errs)
	}
	if errs[0].Rule != "min" {
		t.Fatalf("expected min rule failure, got: %s", errs[0].Rule)
	}
}

func TestValidateFieldEmail(t *testing.T) {
	if err := ValidateField("user@example.com", "email"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := ValidateField("bad", "email"); err == nil {
		t.Fatal("expected error for invalid email")
	}
}

func TestValidateFieldUUID(t *testing.T) {
	if err := ValidateField("550e8400-e29b-41d4-a716-446655440000", "uuid"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := ValidateField("not-uuid", "uuid"); err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestValidateFieldWithCustomMessage(t *testing.T) {
	err := ValidateField("", "required,msg=field cannot be blank")
	if err == nil {
		t.Fatal("expected error")
	}
	errs := Errors(err)
	if errs[0].Message != "field cannot be blank" {
		t.Fatalf("expected custom message, got: %s", errs[0].Message)
	}
}

func TestValidateFieldNumeric(t *testing.T) {
	if err := ValidateField(42, "min=10,max=100"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := ValidateField(5, "min=10"); err == nil {
		t.Fatal("expected error for value below min")
	}
}

func TestValidateFieldPasses(t *testing.T) {
	if err := ValidateField("hello", "required,min=3,max=10"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

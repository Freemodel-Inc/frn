package frn

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func RegisterValidation(validate *validator.Validate) {
	fn := func(fl validator.FieldLevel) bool {
		raw := fl.Field().Interface()
		if raw == nil {
			return true
		}

		var id ID
		switch v := raw.(type) {
		case ID:
			id = v
		case *ID:
			id = *v
		default:
			return true
		}
		if id == "" {
			return true
		}

		return isValidID(id, fl.Param())
	}

	err := validate.RegisterValidation("frn", fn, false)
	if err != nil {
		panic(err)
	}
}

func Validate(id ID, patterns ...string) error {
	if id == "" {
		return fmt.Errorf("ID not set")
	}

	for _, pattern := range patterns {
		if isValidID(id, pattern) {
			return nil
		}
	}

	return fmt.Errorf("ID invalid: expected on of %v", strings.Join(patterns, ", "))

}

func isValidID(id ID, pattern string) bool {
	parts := strings.Split(pattern, "/")
	switch len(parts) {
	case 1:
		return !id.HasChild() && isMatch(parts[0], id.Type())
	case 2:
		return id.HasChild() && isMatch(parts[0], id.Type()) && isMatch(parts[1], id.Child().Type())
	default:
		return false
	}
}

func isMatch(want string, got Type) bool {
	if want == "" {
		return true
	}
	return want == got.String()
}

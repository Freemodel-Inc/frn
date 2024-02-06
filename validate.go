package frn

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var re = regexp.MustCompile(`^([^/#]+)?(/)?([^/#]+)?(#)?([^#]+)?$`)

func RegisterValidation(validate *validator.Validate) {
	fn := func(fl validator.FieldLevel) bool {
		var ids []ID
		switch v := fl.Field().Interface().(type) {
		case ID:
			ids = append(ids, v)
		case []ID:
			ids = append(ids, v...)
		case IDSet:
			ids = append(ids, v...)
		case *ID:
			if v == nil {
				return true
			}
			ids = append(ids, *v)
		default:
			return true
		}
		if len(ids) == 0 {
			return true
		}

		param := fl.Param()
		for _, id := range ids {
			if !isValidID(id, param) {
				return false
			}
		}

		return true
	}

	err := validate.RegisterValidation("frn", fn, true)
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
	if id == "" {
		return true
	}
	if pattern == "" {
		return id.IsValid()
	}

	patternMatch := re.FindStringSubmatch(pattern)
	if len(patternMatch) == 0 {
		return false
	}

	// e.g. entity/card_tx#receipt
	for index, value := range patternMatch {
		switch index {
		case 1: // entity
			if value != "" {
				if id.Type().String() != value {
					return false // parent type mismatch
				}
			}
		case 2: // /
			switch {
			case value == "" && id.HasChild():
				return false
			case value != "" && !id.HasChild():
				return false
			}
		case 3: // card_tx
			if value != "" {
				if id.Child().Type().String() != value {
					return false // child type mismatch
				}
			}
		case 4: // #
			switch {
			case value == "" && id.HasPath():
				return false
			case value != "" && !id.HasPath():
				return false
			}
		case 5: // receipt
			if value != "" {
				want, _, _ := id.Path()
				if value != want {
					return false
				}
			}
		}
	}

	return true
}

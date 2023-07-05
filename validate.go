package frn

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var re = regexp.MustCompile(`^([^/#]+)?(/)?([^/#]+)?(#)?([^#]+)?$`)

func RegisterValidation(validate *validator.Validate) {
	fn := func(fl validator.FieldLevel) bool {
		var id ID
		switch v := fl.Field().Interface().(type) {
		case ID:
			id = v
		case *ID:
			if v == nil {
				return true
			}
			id = *v
		default:
			return true
		}
		if id == "" {
			return true
		}

		return isValidID(id, fl.Param())
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
					fmt.Println("1", id.Type().String(), value)
					encoder := json.NewEncoder(os.Stdout)
					encoder.SetIndent("", "  ")
					_ = encoder.Encode(patternMatch)
					return false // parent type mismatch
				}
			}
		case 2: // /
			switch {
			case value == "" && id.HasChild():
				fmt.Println("2a")
				encoder := json.NewEncoder(os.Stdout)
				encoder.SetIndent("", "  ")
				_ = encoder.Encode(patternMatch)
				return false
			case value != "" && !id.HasChild():
				fmt.Println("2b")
				return false
			}
		case 3: // card_tx
			if value != "" {
				if id.Child().Type().String() != value {
					fmt.Println("3")
					return false // child type mismatch
				}
			}
		case 4: // #
			switch {
			case value == "" && id.HasPath():
				fmt.Println("4a")
				return false
			case value != "" && !id.HasPath():
				fmt.Println("4b")
				return false
			}
		case 5: // receipt
			if value != "" {
				want, _, _ := id.Path()
				if value != want {
					fmt.Println("5")
					return false
				}
			}
		}
	}

	return true
}

func isPathMatch(want string, id ID) bool {
	got, _, _ := id.Path()
	return want == got
}

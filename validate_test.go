package frn

import (
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/tj/assert"
)

func TestValidator(t *testing.T) {
	validate := validator.New()
	RegisterValidation(validate)

	t.Run("required", func(t *testing.T) {
		type Example struct {
			Value ID `validate:"required"`
		}

		err := validate.Struct(Example{Value: ""})
		assert.NotNil(t, err)

		err = validate.Struct(Example{Value: "blah"})
		assert.Nil(t, err)
	})

	t.Run("is frn", func(t *testing.T) {
		type Example struct {
			Value ID `validate:"frn"`
		}

		err := validate.Struct(Example{Value: ""})
		assert.Nil(t, err)

		err = validate.Struct(Example{Value: "blah"})
		assert.NotNil(t, err)

		err = validate.Struct(Example{Value: "fm:dev:blah:123"})
		assert.Nil(t, err)
	})

	t.Run("parent", func(t *testing.T) {
		type Example struct {
			Value ID `validate:"frn=blah"`
		}

		testCases := map[string]struct {
			Value   ID
			WantErr bool
		}{
			"empty string": {
				Value:   "",
				WantErr: false,
			},
			"ok": {
				Value:   "fm:dev:blah:123",
				WantErr: false,
			},
			"fails": {
				Value:   "fm:dev:boom:123",
				WantErr: true,
			},
		}

		for label, tc := range testCases {
			t.Run(label, func(t *testing.T) {
				err := validate.Struct(Example{Value: tc.Value})
				if tc.WantErr {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err)
				}
			})
		}
	})

	t.Run("parent with path", func(t *testing.T) {
		type Example struct {
			Value ID `validate:"frn=blah#key"`
		}

		testCases := map[string]struct {
			Value   ID
			WantErr bool
		}{
			"empty string": {
				Value:   "",
				WantErr: false,
			},
			"ok": {
				Value:   "fm:dev:blah:123/key/value",
				WantErr: false,
			},
			"bad path": {
				Value:   "fm:dev:blah:123/nope/value",
				WantErr: true,
			},
			"child": {
				Value:   "fm:dev:blah:123:child:456/key/value",
				WantErr: true,
			},
			"fails": {
				Value:   "fm:dev:boom:123",
				WantErr: true,
			},
		}

		for label, tc := range testCases {
			t.Run(label, func(t *testing.T) {
				err := validate.Struct(Example{Value: tc.Value})
				if tc.WantErr {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err)
				}
			})
		}
	})

	t.Run("live parent with path", func(t *testing.T) {
		type Example struct {
			Value ID `validate:"required,frn=project#account"`
		}

		testCases := map[string]struct {
			Value   ID
			WantErr bool
		}{
			"ok": {
				Value:   "dev:crm:project:4/account/owner-ap",
				WantErr: false,
			},
		}

		for label, tc := range testCases {
			t.Run(label, func(t *testing.T) {
				err := validate.Struct(Example{Value: tc.Value})
				if tc.WantErr {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err)
				}
			})
		}
	})

	t.Run("child", func(t *testing.T) {
		type Example struct {
			Value ID `validate:"frn=/blah"`
		}

		testCases := map[string]struct {
			Value   ID
			WantErr bool
		}{
			"empty string": {
				Value:   "",
				WantErr: false,
			},
			"ok": {
				Value:   "fm:dev:do-not-matter:123:blah:456",
				WantErr: false,
			},
			"fails": {
				Value:   "fm:dev:blah:123",
				WantErr: true,
			},
		}

		for label, tc := range testCases {
			t.Run(label, func(t *testing.T) {
				err := validate.Struct(Example{Value: tc.Value})
				if tc.WantErr {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err)
				}
			})
		}
	})

	t.Run("parent and child", func(t *testing.T) {
		type Example struct {
			Value ID `validate:"frn=parent/child"`
		}

		testCases := map[string]struct {
			Value   ID
			WantErr bool
		}{
			"empty string": {
				Value:   "",
				WantErr: false,
			},
			"ok": {
				Value:   "fm:dev:parent:123:child:456",
				WantErr: false,
			},
			"ok - tertiary": {
				Value:   "fm:dev:parent:123:child:456/key/value",
				WantErr: true,
			},
			"bad child": {
				Value:   "fm:dev:parent:123:other:456",
				WantErr: true,
			},
			"bad parent": {
				Value:   "fm:dev:other:123:child:456",
				WantErr: true,
			},
		}

		for label, tc := range testCases {
			t.Run(label, func(t *testing.T) {
				err := validate.Struct(Example{Value: tc.Value})
				if tc.WantErr {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err)
				}
			})
		}
	})

	t.Run("parent and child and path", func(t *testing.T) {
		type Example struct {
			Value ID `validate:"frn=parent/child#key"`
		}

		testCases := map[string]struct {
			Value   ID
			WantErr bool
		}{
			"empty string": {
				Value:   "",
				WantErr: false,
			},
			"ok": {
				Value:   "fm:dev:parent:123:child:456/key/value",
				WantErr: false,
			},
			"ok - tertiary": {
				Value:   "fm:dev:parent:123:child:456/key/value",
				WantErr: false,
			},
			"bad child": {
				Value:   "fm:dev:parent:123:other:456",
				WantErr: true,
			},
			"bad parent": {
				Value:   "fm:dev:other:123:child:456",
				WantErr: true,
			},
		}

		for label, tc := range testCases {
			t.Run(label, func(t *testing.T) {
				err := validate.Struct(Example{Value: tc.Value})
				if tc.WantErr {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err)
				}
			})
		}
	})

	t.Run("no child", func(t *testing.T) {
		type Example struct {
			Value ID `validate:"frn=parent"`
		}

		testCases := map[string]struct {
			Value   ID
			WantErr bool
		}{
			"ok": {
				Value:   "fm:dev:parent:123",
				WantErr: false,
			},
			"ok - tertiary": {
				Value:   "fm:dev:parent:123/key/value",
				WantErr: true,
			},
			"err": {
				Value:   "fm:dev:parent:123:child:456",
				WantErr: true,
			},
		}

		for label, tc := range testCases {
			t.Run(label, func(t *testing.T) {
				err := validate.Struct(Example{Value: tc.Value})
				if tc.WantErr {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err)
				}
			})
		}
	})

	t.Run("ptr", func(t *testing.T) {
		type Example struct {
			Value *ID `validate:"frn=parent"`
		}

		testCases := map[string]struct {
			Value   *ID
			WantErr bool
		}{
			"nil": {},
		}

		for label, tc := range testCases {
			t.Run(label, func(t *testing.T) {
				err := validate.Struct(Example{Value: tc.Value})
				if tc.WantErr {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err)
				}
			})
		}
	})
}

func TestValidate(t *testing.T) {
	validate := validator.New()
	RegisterValidation(validate)

	testCases := map[string]struct {
		Value   ID
		Pattern string
		WantErr bool
	}{
		"empty string": {
			Value:   "",
			Pattern: "blah",
			WantErr: true, // Validate requires id to be present
		},
		"ok": {
			Value:   "fm:dev:blah:123",
			Pattern: "blah",
			WantErr: false,
		},
		"ok - tertiary": {
			Value:   "fm:dev:blah:123/key/value",
			Pattern: "blah",
			WantErr: true,
		},
		"fails": {
			Value:   "fm:dev:boom:123",
			Pattern: "blah",
			WantErr: true,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {

			err := Validate(tc.Value, strings.Split(tc.Pattern, ",")...)
			if tc.WantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

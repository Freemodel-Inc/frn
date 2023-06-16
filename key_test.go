package frn

import (
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/tj/assert"
)

const (
	TypeEntity  Type = "entity"
	TypeEvent   Type = "event"
	TypeProject Type = "project"
)

func TestPartString(t *testing.T) {
	var (
		id    = "abc"
		idSub = "def"
		ns    = NewNamespace("", ServiceCRM)
	)
	testCases := map[string]struct {
		Key    ID
		Index  int
		Want   string
		WantOk bool
	}{
		"0": {
			Key:    ns.New(TypeEntity, id),
			Index:  0,
			Want:   "fm",
			WantOk: true,
		},
		"1": {
			Key:    ns.New(TypeEntity, id),
			Index:  1,
			Want:   ServiceCRM.String(),
			WantOk: true,
		},
		"2": {
			Key:    ns.New(TypeEntity, id),
			Index:  2,
			Want:   TypeEntity.String(),
			WantOk: true,
		},
		"3": {
			Key:    ns.New(TypeEntity, id),
			Index:  3,
			Want:   id,
			WantOk: true,
		},
		"4": {
			Key:    ns.New(TypeEntity, id),
			Index:  4,
			WantOk: false,
		},
		"4 sub": {
			Key:    ns.NewWithChild(TypeProject, id, TypeEvent, idSub),
			Index:  4,
			Want:   TypeEvent.String(),
			WantOk: true,
		},
		"5 sub": {
			Key:    ns.NewWithChild(TypeProject, id, TypeEvent, idSub),
			Index:  5,
			Want:   idSub,
			WantOk: true,
		},
		"6 sub": {
			Key:    ns.NewWithChild(TypeProject, id, TypeEvent, idSub),
			Index:  6,
			WantOk: false,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			got, ok := tc.Key.partString(tc.Index)
			assert.Equal(t, tc.WantOk, ok)
			assert.Equal(t, tc.Want, got)
		})
	}
}

func TestID_Sub(t *testing.T) {
	var (
		id    = ksuid.New().String()
		idSub = ksuid.New().String()
		ns    = NewNamespace("", ServiceCRM)
		got   = ns.New(TypeEntity, id).Sub(TypeEvent, idSub)
		want  = ns.NewWithChild(TypeEntity, id, TypeEvent, idSub)
	)

	assert.Equal(t, want, got)
}

func TestParse(t *testing.T) {
	testCases := map[string]struct {
		ID            ID
		Parent        ID
		Child         ID
		Value         string
		TertiaryKey   string
		TertiaryValue string
	}{
		"empty": {
			ID: "",
		},
		"parent": {
			ID:     "fm:crm:project:1",
			Parent: "fm:crm:project:1",
			Value:  "1",
		},
		"parent tertiary": {
			ID:            "fm:crm:project:1/key/value",
			Parent:        "fm:crm:project:1",
			Value:         "1",
			TertiaryKey:   "key",
			TertiaryValue: "value",
		},
		"child": {
			ID:            "fm:crm:project:1:contract:2",
			Parent:        "fm:crm:project:1",
			Child:         "fm:crm:contract:2",
			Value:         "1",
			TertiaryKey:   "",
			TertiaryValue: "",
		},
		"child tertiary": {
			ID:            "fm:crm:project:1:contract:2/key/value",
			Parent:        "fm:crm:project:1",
			Child:         "fm:crm:contract:2",
			Value:         "1",
			TertiaryKey:   "key",
			TertiaryValue: "value",
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			assert.Equal(t, tc.Parent, tc.ID.Parent())
			assert.Equal(t, tc.Child, tc.ID.Child())
			assert.Equal(t, tc.Value, tc.ID.Value())

			key, value, ok := tc.ID.Path()
			assert.Equal(t, tc.TertiaryKey, key)
			assert.Equal(t, tc.TertiaryValue, value)
			assert.Equal(t, key != "" || value != "", ok)
		})
	}
}

func TestID(t *testing.T) {
	var (
		id     = ksuid.New().String()
		idSub  = ksuid.New().String()
		ns     = NewNamespace("", ServiceCRM)
		v      = ns.New(TypeEntity, id).Sub(TypeEvent, idSub)
		child  = v.Child()
		parent = v.Parent()
	)

	assert.Equal(t, ns, v.Namespace())
	assert.Equal(t, ns, child.Namespace())
	assert.Equal(t, ns, parent.Namespace())

	assert.True(t, v.HasChild())
	assert.Equal(t, id, v.Value())
	assert.Equal(t, ServiceCRM, v.Service())
	assert.Equal(t, TypeEntity, v.Type())
	assert.Equal(t, parent.String()+sep, v.ChildPrefix())

	assert.False(t, parent.HasChild())
	assert.Equal(t, id, parent.Value())
	assert.Equal(t, ServiceCRM, parent.Service())
	assert.Equal(t, TypeEntity, parent.Type())
	assert.Equal(t, parent.String()+sep, parent.ChildPrefix())

	assert.False(t, child.HasChild())
	assert.Equal(t, idSub, child.Value())
	assert.Equal(t, ServiceCRM, child.Service())
	assert.Equal(t, TypeEvent, child.Type())
	assert.Equal(t, child.String()+sep, child.ChildPrefix())
}

func TestID_IsValid(t *testing.T) {
	testCases := map[string]struct {
		ID            ID
		IsValid       bool
		TertiaryKey   string
		TertiaryValue string
	}{
		"parent": {
			ID:      "namespace:service:type:value",
			IsValid: true,
		},
		"parent - tertiary": {
			ID:            "namespace:service:type:value/tertiary-key/tertiary-value",
			IsValid:       true,
			TertiaryKey:   "tertiary-key",
			TertiaryValue: "tertiary-value",
		},
		"parent no namespace": {
			ID:      ":service:type:value",
			IsValid: false,
		},
		"parent no service": {
			ID:      "namespace::type:value",
			IsValid: false,
		},
		"parent small service": {
			ID:      "namespace:s:type:value",
			IsValid: true,
		},
		"parent no type": {
			ID:      "namespace:service::value",
			IsValid: false,
		},
		"parent no value": {
			ID:      "namespace:service:type:",
			IsValid: false,
		},
		"parent small value": {
			ID:      "namespace:service:type:v",
			IsValid: true,
		},
		"child": {
			ID:      "namespace:service:type:value:sub-type:sub-value",
			IsValid: true,
		},
		"child no namespace": {
			ID:      ":service:type:value:sub-type:sub-value",
			IsValid: false,
		},
		"child no sub-type": {
			ID:      "namespace:service:type:value::sub-value",
			IsValid: false,
		},
		"child no sub-value": {
			ID:      "namespace:service:type:value:sub-type:",
			IsValid: true,
		},
		"child small sub-value": {
			ID:      "namespace:service:type:value:sub-type:s",
			IsValid: true,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			assert.Equal(t, tc.IsValid, tc.ID.IsValid())
			key, value, ok := tc.ID.Path()
			assert.Equal(t, tc.TertiaryKey, key)
			assert.Equal(t, tc.TertiaryValue, value)
			assert.Equal(t, key != "" || value != "", ok)
		})
	}
}

func TestID_Base(t *testing.T) {
	testCases := map[string]struct {
		ID   ID
		Want ID
	}{
		"parent": {
			ID:   "fm:crm:project:1",
			Want: "fm:crm:project:1",
		},
		"parent with path": {
			ID:   "fm:crm:project:1/key/value",
			Want: "fm:crm:project:1",
		},
		"child": {
			ID:   "fm:crm:project:1:invoice:2",
			Want: "fm:crm:project:1:invoice:2",
		},
		"child with path": {
			ID:   "fm:crm:project:1:invoice:2/key/value",
			Want: "fm:crm:project:1:invoice:2",
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			assert.Equal(t, tc.Want, tc.ID.Base())
		})
	}
}

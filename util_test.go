package frn

import (
	"testing"

	"github.com/tj/assert"
)

func TestReShape(t *testing.T) {
	testCases := map[string]struct {
		Shape string
		Want  []string
	}{
		"empty": {
			Shape: "",
			Want:  nil,
		},
		"unary": {
			Shape: "project",
			Want:  []string{"project", "project", "", "", "", ""},
		},
		"binary": {
			Shape: "project/contract",
			Want:  []string{"project/contract", "project", "/contract", "contract", "", ""},
		},
		"tertiary": {
			Shape: "project/contract#change",
			Want:  []string{"project/contract#change", "project", "/contract", "contract", "#change", "change"},
		},
		"tertiary - alt": {
			Shape: "project#change",
			Want:  []string{"project#change", "project", "", "", "#change", "change"},
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			got := reShape.FindStringSubmatch(tc.Shape)
			assert.Equal(t, tc.Want, got)
		})
	}
}

func TestParentShape(t *testing.T) {
	testCases := map[string]struct {
		Shape string
		Want  []string
	}{
		"unary": {
			Shape: "project",
			Want:  []string{"", "", ""},
		},
		"binary": {
			Shape: "project/contract",
			Want:  []string{"project", "", ""},
		},
		"tertiary": {
			Shape: "project/contract#change",
			Want:  []string{"project", "contract", ""},
		},
		"tertiary - alt": {
			Shape: "project#change",
			Want:  []string{"project", "", ""},
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			got := ParentShape(ShapeSlice(tc.Shape))
			assert.Equal(t, tc.Want, got)
		})
	}
}

func TestSampleViaShape(t *testing.T) {
	testCases := map[string]struct {
		ID     ID
		Shape  string
		Want   ID
		WantOk bool
	}{
		"binary": {
			ID:     "dev:crm:project:1",
			Shape:  "project/contract",
			Want:   "dev:crm:project:1:contract:_",
			WantOk: true,
		},
		"tertiary": {
			ID:     "dev:crm:project:1:contract:2",
			Shape:  "project/contract#work_item",
			Want:   "dev:crm:project:1:contract:2/work_item/_",
			WantOk: true,
		},
		"tertiary - alt": {
			ID:     "dev:crm:project:1",
			Shape:  "project#work_item",
			Want:   "dev:crm:project:1/work_item/_",
			WantOk: true,
		},
		"bad, tertiary": {
			ID:     "dev:crm:project:1:entity:2",
			Shape:  "project/contract#work_item",
			WantOk: false,
		},
		"bad, tertiary - alt": {
			ID:     "dev:crm:user:1",
			Shape:  "project#work_item",
			WantOk: false,
		},
	}

	ns := NewNamespace("dev", ServiceCRM)
	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			got, ok := SampleViaShape(ns, tc.ID, tc.Shape)
			if !tc.WantOk {
				assert.False(t, ok)
				return
			}
			assert.True(t, ok)
			assert.Equal(t, tc.Want, got)
		})
	}
}

func TestShapeSliceValue(t *testing.T) {
	testCases := map[string]struct {
		Slice []string
		Want  string
	}{
		"all": {
			Slice: []string{"project", "contract", "work_item"},
			Want:  "project/contract#work_item",
		},
		"parent child": {
			Slice: []string{"project", "contract", ""},
			Want:  "project/contract",
		},
		"parent tertiary": {
			Slice: []string{"project", "", "work_item"},
			Want:  "project#work_item",
		},
		"parent": {
			Slice: []string{"project", "", ""},
			Want:  "project",
		},
		"too long": {
			Slice: []string{"1", "2", "3", "4"},
			Want:  "",
		},
		"too short": {
			Slice: []string{"1", "2"},
			Want:  "",
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			got := ShapeSliceValue(tc.Slice)
			assert.Equal(t, tc.Want, got)
		})
	}
}

package frn

import (
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/tj/assert"
)

func TestPartString(t *testing.T) {
	id := "abc"
	idSub := "def"
	testCases := map[string]struct {
		Key    ID
		Index  int
		Want   string
		WantOk bool
	}{
		"0": {
			Key:    New(ServiceCRM, TypeEntity, id),
			Index:  0,
			Want:   namespace,
			WantOk: true,
		},
		"1": {
			Key:    New(ServiceCRM, TypeEntity, id),
			Index:  1,
			Want:   ServiceCRM.String(),
			WantOk: true,
		},
		"2": {
			Key:    New(ServiceCRM, TypeEntity, id),
			Index:  2,
			Want:   TypeEntity.String(),
			WantOk: true,
		},
		"3": {
			Key:    New(ServiceCRM, TypeEntity, id),
			Index:  3,
			Want:   id,
			WantOk: true,
		},
		"4": {
			Key:    New(ServiceCRM, TypeEntity, id),
			Index:  4,
			WantOk: false,
		},
		"4 sub": {
			Key:    NewSubType(ServiceCRM, TypeProject, id, TypeEvent, idSub),
			Index:  4,
			Want:   TypeEvent.String(),
			WantOk: true,
		},
		"5 sub": {
			Key:    NewSubType(ServiceCRM, TypeProject, id, TypeEvent, idSub),
			Index:  5,
			Want:   idSub,
			WantOk: true,
		},
		"6 sub": {
			Key:    NewSubType(ServiceCRM, TypeProject, id, TypeEvent, idSub),
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
		got   = New(ServiceCRM, TypeEntity, id).Sub(TypeEvent, idSub)
		want  = NewSubType(ServiceCRM, TypeEntity, id, TypeEvent, idSub)
	)

	assert.Equal(t, want, got)
}

func TestID(t *testing.T) {
	var (
		id     = ksuid.New().String()
		idSub  = ksuid.New().String()
		v      = New(ServiceCRM, TypeEntity, id).Sub(TypeEvent, idSub)
		child  = v.Child()
		parent = v.Parent()
	)

	assert.True(t, v.HasChild())
	assert.Equal(t, id, v.ID())
	assert.Equal(t, ServiceCRM, v.Service())
	assert.Equal(t, TypeEntity, v.Type())
	assert.Equal(t, parent.String()+sep, v.ChildPrefix())

	assert.False(t, parent.HasChild())
	assert.Equal(t, id, parent.ID())
	assert.Equal(t, ServiceCRM, parent.Service())
	assert.Equal(t, TypeEntity, parent.Type())
	assert.Equal(t, parent.String()+sep, parent.ChildPrefix())

	assert.False(t, child.HasChild())
	assert.Equal(t, idSub, child.ID())
	assert.Equal(t, ServiceCRM, child.Service())
	assert.Equal(t, TypeEvent, child.Type())
	assert.Equal(t, child.String()+sep, child.ChildPrefix())
}

package frn

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/tj/assert"
)

func TestKey_MarshalDynamoDBAttributeValue(t *testing.T) {
	want := NewNamespace("", ServiceCRM).New(TypeEntity, "abc")
	item, err := dynamodbattribute.Marshal(want)
	assert.Nil(t, err)

	var got ID
	err = dynamodbattribute.Unmarshal(item, &got)
	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestKeySet_MarshalDynamoDBAttributeValue(t *testing.T) {
	ns := NewNamespace("", ServiceCRM)
	want := KeySet{
		ns.New(TypeEntity, "a"),
		ns.New(TypeEntity, "b"),
		ns.New(TypeEntity, "c"),
	}

	item, err := dynamodbattribute.Marshal(want)
	assert.Nil(t, err)

	var got KeySet
	err = dynamodbattribute.Unmarshal(item, &got)
	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestNewKeySet(t *testing.T) {
	tests := map[string]struct {
		Input []string
		Want  KeySet
	}{
		"empty strings": {
			Input: []string{"", "", ""},
			Want:  nil,
		},
		"strip empty strings": {
			Input: []string{"a", "b", ""},
			Want:  KeySet{"a", "b"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := NewKeySet(tc.Input...)
			assert.Equal(t, tc.Want, got)
		})
	}
}

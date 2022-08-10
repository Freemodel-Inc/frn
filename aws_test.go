package frn

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/tj/assert"
)

func TestKey_MarshalDynamoDBAttributeValue(t *testing.T) {
	want := New(ServiceCRM, TypeEntity, "abc")
	item, err := dynamodbattribute.Marshal(want)
	assert.Nil(t, err)

	var got ID
	err = dynamodbattribute.Unmarshal(item, &got)
	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestKeySet_MarshalDynamoDBAttributeValue(t *testing.T) {
	want := KeySet{
		New(ServiceCRM, TypeEntity, "a"),
		New(ServiceCRM, TypeEntity, "b"),
		New(ServiceCRM, TypeEntity, "c"),
	}

	item, err := dynamodbattribute.Marshal(want)
	assert.Nil(t, err)

	var got KeySet
	err = dynamodbattribute.Unmarshal(item, &got)
	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

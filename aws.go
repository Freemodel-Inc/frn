package frn

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func (id ID) MarshalDynamoDBAttributeValue(item *dynamodb.AttributeValue) error {
	if item == nil {
		return fmt.Errorf("unable to marshal key: nil AttributeValue")
	}

	item.S = aws.String(id.String())

	return nil
}

func (id *ID) UnmarshalDynamoDBAttributeValue(item *dynamodb.AttributeValue) error {
	if item.S == nil {
		return fmt.Errorf("unable to marshal key: nil string")
	}

	*id = ID(aws.StringValue(item.S))

	return nil
}

type KeySet []ID

func NewKeySet(ss ...string) KeySet {
	var (
		seen = map[string]struct{}{}
		kk   KeySet
	)
	for _, k := range ss {
		if _, ok := seen[k]; ok {
			continue
		}
		kk = append(kk, ID(k))
		seen[k] = struct{}{}
	}

	return kk
}

func (kk KeySet) Contains(key ID) bool {
	for _, k := range kk {
		if k == key {
			return true
		}
	}
	return false
}

func (kk KeySet) MarshalDynamoDBAttributeValue(item *dynamodb.AttributeValue) error {
	var (
		seen = map[ID]struct{}{}
		ss   []string
	)
	for _, k := range kk {
		if _, ok := seen[k]; ok {
			continue
		}
		ss = append(ss, k.String())
		seen[k] = struct{}{}
	}

	item.SS = aws.StringSlice(ss)

	return nil
}

// Strings exports keys as string slice
func (kk KeySet) Strings() []string {
	var ss []string
	for _, k := range kk {
		ss = append(ss, k.String())
	}
	return ss
}

func (kk *KeySet) UnmarshalDynamoDBAttributeValue(item *dynamodb.AttributeValue) error {
	if item.SS == nil {
		return nil
	}

	*kk = NewKeySet(aws.StringValueSlice(item.SS)...)

	return nil
}

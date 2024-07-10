package frn

import (
	"regexp"
	"slices"

	"github.com/segmentio/ksuid"
)

var reShape = regexp.MustCompile(`^([a-zA-Z0-9_]+)(/([a-zA-Z0-9_]+))?(#([a-zA-Z0-9_]+))?`)

// NewValue generates a new value for an id
func NewValue() string {
	return ksuid.New().String()
}

// SampleViaShape generates a sample id in the shape requested using the potential parent id as a base (if necessary)
func SampleViaShape(ns Namespace, potentialParentID ID, s string) (ID, bool) {
	shape := ShapeSlice(s)
	return SampleViaShapeSlice(ns, potentialParentID, shape)
}

func SampleViaShapeSlice(ns Namespace, potentialParentID ID, shape []string) (ID, bool) {
	have := potentialParentID.ShapeSlice()
	want := ParentShape(shape)
	if !slices.Equal(have, want) {
		return "", false // id cannot be a parent of shape
	}

	if len(shape) == 3 {
		for i := 2; i >= 0; i-- {
			switch {
			case i == 2 && shape[i] != "":
				return potentialParentID.WithPath(shape[i], "_"), true
			case i == 1 && shape[i] != "":
				return potentialParentID.Sub(Type(shape[i]), "_"), true
			case i == 0:
				return ns.New(Type(shape[i]), "_"), true
			}
		}
	}

	return "", false
}

// ParentShape returns the logical parent shape for the given id
func ParentShape(shape []string) []string {
	parent := make([]string, 0, len(shape))
	parent = append(parent, shape...)

	if len(shape) == 3 {
		for i := 2; i >= 0; i-- {
			if parent[i] != "" {
				parent[i] = ""
				break
			}
		}
	}

	return parent
}

// ShapeSlice takes a shape and returns a slice of 3 elements, one for each part (primary, secondary, and tertiary)
func ShapeSlice(shape string) []string {
	ss := [3]string{}
	if match := reShape.FindStringSubmatch(shape); len(match) == 6 {
		ss[0] = match[1]
		ss[1] = match[3]
		ss[2] = match[5]
	}
	return ss[:]
}

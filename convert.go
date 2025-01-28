package frn

import "strings"

// FromStringSlice converts a string slice into an ID slice.
func FromStringSlice(ss ...string) (ids []ID) {
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		ids = append(ids, ID(s))
	}
	return ids
}

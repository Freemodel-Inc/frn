package frn

import (
	"strconv"
	"sync/atomic"
)

type Sequence struct {
	value int64
}

func (s *Sequence) Next() string {
	return strconv.FormatInt(atomic.AddInt64(&s.value, 1), 10)
}

// NewSequence returns a sequence that is useful for generator deterministic ids
func NewSequence(initialValue int64) *Sequence {
	return &Sequence{
		value: initialValue,
	}
}

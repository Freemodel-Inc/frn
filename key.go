package frn

import (
	"github.com/segmentio/ksuid"
	"strings"
)

const sep = ":"

type IDFactoryFunc func(v string) ID

func (fn IDFactoryFunc) NewID() ID {
	return fn(ksuid.New().String())
}

func (fn IDFactoryFunc) WithValue(v string) ID {
	return fn(v)
}

type IDFactory interface {
	NewID() ID
	WithValue(v string) ID
}

// Namespace consists of prefix plus service e.g. frm:crm
type Namespace string

func NewNamespace(env string, s Service) Namespace {
	if env == "" || env == "prd" {
		env = "fm"
	}
	return Namespace(env + sep + s.String())
}

func (n Namespace) IDFactory(t Type) IDFactoryFunc {
	return func(v string) ID {
		return ID(n.String() + sep + t.String() + sep + v)
	}
}

func (n Namespace) New(t Type, id string) ID {
	return ID(n.String() + sep + t.String() + sep + id)
}

func (n Namespace) NewWithChild(t Type, id string, st Type, idSub string) ID {
	return ID(n.String() + sep + t.String() + sep + id + sep + st.String() + sep + idSub)
}

func (n Namespace) String() string {
	return string(n)
}

type ID string

func (id ID) part(i int) (begin, end int, ok bool) {
	for index, s := range id {
		if s != ':' {
			continue
		}

		switch i {
		case 0:
			return begin, index, true
		default:
			begin = index + 1
		}
		i--
	}

	if i == 0 {
		return begin, len(id), true
	}

	return 0, 0, false
}

func (id ID) partString(i int) (string, bool) {
	begin, end, ok := id.part(i)
	if !ok {
		return "", false
	}

	return string(id[begin:end]), true
}

func (id ID) Child() ID {
	st, ok := id.partString(4)
	if !ok {
		return ""
	}

	si, ok := id.partString(5)
	if !ok {
		return ""
	}

	return id.Namespace().New(Type(st), si)
}

// ChildPrefix returns prefix all children of parent must begin with
func (id ID) ChildPrefix() string {
	if id.HasChild() {
		return id.Parent().ChildPrefix()
	}
	return id.String() + sep
}

func (id ID) HasChild() bool {
	_, ok := id.partString(4)
	return ok
}

func (id ID) Value() string {
	s, _ := id.partString(3)
	return s
}

func (id ID) IsEmpty() bool {
	return id == ""
}

func (id ID) IsPresent() bool {
	return !id.IsEmpty()
}

func (id ID) Namespace() Namespace {
	s := id.String()
	a := strings.Index(s, sep)
	if a == -1 {
		return ""
	}

	b := strings.Index(s[a+1:], sep)
	if b == -1 {
		return ""
	}

	return Namespace(id[:a+b+1])
}

func (id ID) Parent() ID {
	if id.HasChild() {
		return id.Namespace().New(id.Type(), id.Value())
	}
	return id
}

func (id ID) Service() Service {
	partString, _ := id.partString(1)
	return Service(partString)
}

func (id ID) String() string {
	return string(id)
}

func (id ID) Sub(st Type, idSub string) ID {
	return ID(id.String() + sep + st.String() + sep + idSub)
}

func (id ID) Type() Type {
	s, _ := id.partString(2)
	return Type(s)
}

func (id ID) WithChild(child ID) ID {
	return id.Sub(child.Type(), child.Value())
}

type IDMap map[ID]struct{}

func (vv IDMap) Slice() IDSet {
	var idSet IDSet
	for v := range vv {
		idSet = append(idSet, v)
	}
	return idSet
}

type IDSet []ID

// Contains returns true if id provided part of set
func (vv IDSet) Contains(want ID) bool {
	for _, v := range vv {
		if v == want {
			return true
		}
	}
	return false
}

func (vv IDSet) Where(fn func(ID) bool) IDSet {
	var results IDSet
	for _, v := range vv {
		if fn(v) {
			results = append(results, v)
		}
	}
	return results
}

type Service string

func (s Service) String() string {
	return string(s)
}

const (
	ServiceCRM        Service = "crm"
	ServiceOnboarding Service = "onboarding"
	ServiceSystem     Service = "system"
)

type Type string

func (t Type) String() string {
	return string(t)
}

func First(ids ...ID) ID {
	for _, id := range ids {
		if id != "" {
			return id
		}
	}
	return ""
}

func Ptr(id ID) *ID {
	if id == "" {
		return nil
	}
	return &id
}

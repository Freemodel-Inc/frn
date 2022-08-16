package frn

import "github.com/segmentio/ksuid"

const sep = ":"

var defaultNamespace = NewNamespace("")

type IDFactoryFunc func() ID

func (fn IDFactoryFunc) NewID() ID {
	return fn()
}

type IDFactory interface {
	NewID() ID
}

type Namespace string

func NewNamespace(env string) Namespace {
	if env == "" || env == "prd" {
		return "fm"
	}
	return Namespace(env)
}

func (n Namespace) IDFactory(s Service, t Type) IDFactoryFunc {
	return func() ID {
		return ID(n.String() + sep + s.String() + sep + t.String() + sep + ksuid.New().String())
	}
}

func (n Namespace) String() string {
	return string(n)
}

func (n Namespace) New(s Service, t Type, id string) ID {
	return ID(n.String() + sep + s.String() + sep + t.String() + sep + id)
}

func (n Namespace) NewWithChild(s Service, t Type, id string, st Type, idSub string) ID {
	return ID(n.String() + sep + s.String() + sep + t.String() + sep + id + sep + st.String() + sep + idSub)
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

	return New(id.Service(), Type(st), si)
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

func (id ID) ID() string {
	s, _ := id.partString(3)
	return s
}

func (id ID) IsEmpty() bool {
	return id == ""
}

func (id ID) Parent() ID {
	return New(id.Service(), id.Type(), id.ID())
}

func (id ID) Service() Service {
	partString, _ := id.partString(1)
	return Service(partString)
}

func (id ID) String() string {
	return string(id)
}

func (id ID) Sub(st Type, idSub string) ID {
	return NewSubType(id.Service(), id.Type(), id.ID(), st, idSub)
}

func (id ID) Type() Type {
	s, _ := id.partString(2)
	return Type(s)
}

func (id ID) WithChild(child ID) ID {
	return id.Sub(child.Type(), child.ID())
}

type Service string

func (s Service) String() string {
	return string(s)
}

const (
	ServiceCRM        Service = "crm"
	ServiceOnboarding Service = "onboarding"
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

func New(s Service, t Type, id string) ID {
	return defaultNamespace.New(s, t, id)
}

func NewSubType(s Service, t Type, id string, st Type, idSub string) ID {
	return defaultNamespace.NewWithChild(s, t, id, st, idSub)
}

func Ptr(id ID) *ID {
	if id == "" {
		return nil
	}
	return &id
}

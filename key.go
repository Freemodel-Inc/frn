package frn

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"regexp"
	"strings"

	"github.com/segmentio/ksuid"
)

const (
	sep     = ":" // sep separates parent and child parts of the FRN
	pathSep = "/" // separates parent and child ids from pathSep data
)

var (
	// see test case, TestID_IsValid, for examples
	// e.g. fm:crm:contact:1234
	reValid = regexp.MustCompile(`^([a-zA-Z0-9\-_]+:){3}[a-zA-Z0-9\-_]+(:[a-zA-Z0-9\-_]+:[a-zA-Z0-9\-_]*)?(/[a-z0-9\-_/]+)?$`)
)

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

// Env returns the environment e.g. fm:crm => fm
func (n Namespace) Env() string {
	if index := strings.Index(string(n), sep); index > 0 {
		return string(n)[:index]
	}
	return ""
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

// Base returns the id sans any path elements e.g. fm:crm:contact:1234/key/value => fm:crm:contact:1234
func (id ID) Base() ID {
	if index := strings.Index(id.String(), pathSep); index != -1 {
		return id[:index]
	}
	return id
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

	// strip trailing path from child
	if index := strings.Index(si, pathSep); index != -1 {
		si = si[:index]
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

func (id ID) HasPath() bool {
	s := id.String()
	index := strings.Index(s, pathSep)
	if index == -1 {
		return false
	}
	return true
}

// In returns true if the id is explicitly within the provided set of ids
func (id ID) In(wants ...ID) bool {
	for _, want := range wants {
		if id == want {
			return true
		}
	}
	return false
}

func (id ID) Value() string {
	s, _ := id.partString(3)
	if index := strings.Index(s, pathSep); index != -1 {
		return s[:index]
	}
	return s
}

func (id ID) IsEmpty() bool {
	return id == ""
}

// IsParentType returns true if the parent id is of the provided type
func (id ID) IsParentType(want Type) bool {
	return id.Type() == want
}

// IsChildType returns true if the child id is of the provided type
func (id ID) IsChildType(want Type) bool {
	return id.HasChild() && id.Child().Type() == want
}

func (id ID) IsPresent() bool {
	return !id.IsEmpty()
}

func (id ID) IsValid() bool {
	return reValid.MatchString(id.String())
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

	// strip path if present from parent
	if index := strings.Index(id.String(), pathSep); index != -1 {
		return id[0:index]
	}

	return id
}

// Path extracts the tertiary values from the id
func (id ID) Path() (head, tail string, ok bool) {
	s := id.String()
	index := strings.Index(s, pathSep)
	if index == -1 {
		return "", "", false
	}

	parts := strings.SplitN(s[index+1:], pathSep, 2)
	if len(parts) < 2 {
		if parts[0] == "" {
			return "", "", false
		}
		return parts[0], "", true
	}

	return parts[0], parts[1], true
}

func (id ID) Service() Service {
	partString, _ := id.partString(1)
	return Service(partString)
}

// Shape returns the shape of the id e.g. frm:crm:project:1234 => project, frm:crm:entity:1:card_tx:2/fund_request/3 => entity/card_tx#fund_request
func (id ID) Shape() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(id.Type().String())
	if id.HasChild() {
		buf.WriteString("/")
		buf.WriteString(id.Child().Type().String())
	}
	if head, _, ok := id.Path(); ok {
		buf.WriteString("#")
		buf.WriteString(head)
	}
	return buf.String()
}

// ShapeSlice returns the shape as a slice with 3 elements e.g. project/contract -> ["project", "contract", ""]
func (id ID) ShapeSlice() []string {
	shape := [3]string{}
	shape[0] = id.Type().String()
	if id.HasChild() {
		shape[1] = id.Child().Type().String()
	}
	if head, _, ok := id.Path(); ok {
		shape[2] = head
	}
	return shape[:]
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

// WithPath returns the tertiary form of the id e.g. frm:crm:contact:123:address:456/a/b/c
// currently only supports two levels of nesting e.g. a or a/b, but not a/b/c
// WithPath replaces any existing path values
func (id ID) WithPath(head string, tail ...string) ID {
	s := string(id)
	if index := strings.Index(s, pathSep); index != -1 {
		s = s[:index]
	}

	buf := bytes.NewBuffer(make([]byte, 0, 64))
	buf.WriteString(s)
	buf.WriteString(pathSep)
	buf.WriteString(head)
	for _, t := range tail {
		if t == "" {
			continue
		}
		buf.WriteString(pathSep)
		buf.WriteString(t)
	}
	return ID(buf.String())
}

func (id ID) WithPathStringer(head fmt.Stringer, tail ...fmt.Stringer) ID {
	var ss []string
	for _, t := range tail {
		ss = append(ss, t.String())
	}
	return id.WithPath(head.String(), ss...)
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
//
//goland:noinspection GoMixedReceiverTy
//goland:noinspection GoMixedReceiverTypes
func (vv IDSet) Contains(want ID) bool {
	for _, v := range vv {
		if v == want {
			return true
		}
	}
	return false
}

//goland:noinspection GoMixedReceiverTypes
func (vv IDSet) MarshalDynamoDBAttributeValue(item *dynamodb.AttributeValue) error {
	if len(vv) == 0 {
		item.NULL = aws.Bool(true)
		return nil
	}

	var ss []*string
	for _, v := range vv {
		if v == "" {
			continue
		}
		ss = append(ss, aws.String(v.String()))
	}

	item.SS = ss

	return nil
}

// Trim returns a new IDSet with the blank ids removed
//
//goland:noinspection GoMixedReceive
//goland:noinspection GoMixedReceiverTypes
func (vv IDSet) Trim() IDSet {
	var idSet IDSet
	for _, v := range vv {
		if v == "" {
			continue
		}
		idSet = append(idSet, v)
	}
	return idSet
}

//goland:noinspection GoMixedReceiverTypes
func (vv *IDSet) UnmarshalDynamoDBAttributeValue(item *dynamodb.AttributeValue) error {
	if aws.BoolValue(item.NULL) || item.SS == nil {
		return nil
	}

	var set IDSet
	for _, s := range item.SS {
		set = append(set, ID(aws.StringValue(s)))
	}

	*vv = set

	return nil
}

//goland:noinspection GoMixedReceiverTypes
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
	ServiceAgentDB    Service = "adb"
	ServiceCRM        Service = "crm"
	ServiceFinance    Service = "fin"
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

// WithPath provides a syntactic sugar wrapper around WithPath that's more akin to how it's used in practice
func WithPath[T string | int64 | int32 | int](id ID, head Type, tail ...T) ID {
	var ss []string
	for _, t := range tail {
		ss = append(ss, fmt.Sprintf("%v", t))
	}
	return id.WithPath(head.String(), ss...)
}

package object

import (
	"bytes"
	"example.com/writing-an-interpreter/ast"
	"fmt"
	"hash/fnv"
	"strings"
)

type ObjectType string

const (
	INTEGER      = "INTEGER"
	BOOLEAN      = "BOOLEAN"
	STRING       = "STRING"
	NULL         = "NULL"
	ARRAY        = "ARRAY"
	MAP          = "MAP"
	FUNCTION     = "FUNCTION"
	RETURN_VALUE = "RETURN_VALUE"
	ERROR        = "ERROR"
	BUILTIN      = "BUILTIN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL
}

func (n *Null) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

type Function struct {
	Parameters  []*ast.Identifier
	Body        *ast.BlockStatement
	Environment *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION
}

func (f *Function) Inspect() string {
	var out bytes.Buffer
	var params []string

	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING
}

func (s *String) Inspect() string {
	return s.Value
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN
}

func (b *Builtin) Inspect() string {
	return "built-in function"
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType {
	return ARRAY
}

func (a *Array) Inspect() string {
	var out bytes.Buffer
	var elements []string

	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteRune('[')
	out.WriteString(strings.Join(elements, ", "))
	out.WriteRune(']')

	return out.String()
}

type MapPair struct {
	Key   Object
	Value Object
}

type Map struct {
	Pairs map[HashKey]MapPair
}

func (m *Map) Type() ObjectType {
	return MAP
}

func (m *Map) Inspect() string {
	var out bytes.Buffer
	var pairs []string

	for _, pair := range m.Pairs {
		pairs = append(pairs, pair.Key.Inspect()+": "+pair.Value.Inspect())
	}

	out.WriteRune('{')
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteRune('}')

	return out.String()
}

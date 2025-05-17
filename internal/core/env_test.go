package core

import (
	"ocala/internal/tt"
	"testing"
)

func TestNil(t *testing.T) {
	tt.Eq(t, "<NIL>", NIL.Inspect())
	tt.Eq(t, Value(NIL), NIL.Dup())
}

func TestUndefined(t *testing.T) {
	tt.Eq(t, "<UNDEFINED>", UNDEFINED.Inspect())
	tt.Eq(t, Value(UNDEFINED), UNDEFINED.Dup())
}

func TestInt(t *testing.T) {
	tt.Eq(t, "1234", Int(1234).Inspect())
	tt.Eq(t, Value(Int(1234)), Int(1234).Dup())
}

func TestBinary(t *testing.T) {
	tt.Eq(t, "&[0 1 2 3]", (&Binary{0, 1, 2, 3}).Inspect())

	a := &Binary{0, 1, 2, 3}
	b := a.Dup().(*Binary)
	tt.EqSlice(t, []byte(*a), []byte(*b))
	tt.True(t, a != b)

	(*b)[1] = 100
	tt.Eq(t, 1, (*a)[1])
}

func TestKeyword(t *testing.T) {
	tt.Eq(t, "test", Intern("test").String())
	tt.Eq(t, ":test", Intern("test").Inspect())
	tt.Eq(t, Value(Intern("test")), Intern("test").Dup())

	tt.True(t, Intern("test") == Intern("test")) //lint:ignore SA4000 test
	tt.True(t, Intern("test") != NewKeyword("test"))
	tt.True(t, Intern("test") != Gensym("test"))
	tt.True(t, NewKeyword("test") != NewKeyword("test")) //lint:ignore SA4000 test
	tt.True(t, NewKeyword("test") != Gensym("test"))
	tt.True(t, Gensym("test") != Gensym("test")) //lint:ignore SA4000 test

	tt.True(t, Intern("test").MatchId(InternalId(Intern("test"))) != nil)
	tt.True(t, Intern("test").MatchId(InternalId(Intern("tete"))) == nil)

	tt.True(t, Intern("test").MatchExpr(&Vec{InternalId(Intern("test"))}) != nil)
	tt.True(t, Intern("test").MatchId(&Vec{}) == nil)
}

func TestStr(t *testing.T) {
	tt.Eq(t, "test", NewStr("test").String())
	tt.Eq(t, `"test"`, NewStr("test").Inspect())
	tt.Eq(t, `"test\0\x01\x19\x7f\xff\a\b\e\f\n\r\t\v\\'\""`,
		NewStr("test\x00\x01\x19\x7f\xff\a\b\x1b\f\n\r\t\v\\'\"").Inspect())

	a := NewStr("test")
	b := a.Dup().(*Str)
	tt.True(t, a != b)
	tt.True(t, *a == *b)
	tt.True(t, a.Intern() == Intern("test"))
}

func TestVec(t *testing.T) {
	tt.Eq(t, "[0 1 2]", (&Vec{Int(0), Int(1), Int(2)}).Inspect())
	tt.Eq(t, "[0 <nil>]", (&Vec{Int(0), nil}).Inspect())

	a := &Vec{Int(0), Intern("test"), NewStr("test")}
	tt.True(t, a.At(0) == Int(0))
	tt.True(t, a.At(1) == Intern("test"))
	tt.True(t, a.At(2).Inspect() == `"test"`)
	tt.True(t, a.AtOrUndef(0) == Int(0))
	tt.True(t, a.AtOrUndef(3) == UNDEFINED)

	b := a.Dup().(*Vec)
	tt.True(t, a != b)
	tt.True(t, a.At(0) == b.At(0))
	tt.True(t, a.At(1) == b.At(1))
	tt.True(t, a.At(2) != b.At(2))
	tt.True(t, a.At(2).Inspect() == b.At(2).Inspect())

	tt.True(t, a.Size() == 3)
	tt.True(t, a.Push(Int(2)) == a)
	tt.True(t, a.Push(Int(3), Int(4)) == a)
	tt.True(t, a.Size() == 6)
	tt.True(t, a.At(5) == Int(4))

	c := &Vec{Int(0), &Vec{}, &Vec{Int(1), Int(2)}, Int(3), &Vec{Int(4), &Vec{Int(5), Int(6)}}}
	tt.True(t, c.Flatten().Size() == 7)
	tt.EqSlice(t, []Value(*c.Flatten()),
		[]Value{Int(0), Int(1), Int(2), Int(3), Int(4), Int(5), Int(6)})

	tt.True(t, (&Vec{InternalId(Intern("test"))}).ExprTag().Name == Intern("test"))
	tt.True(t, (&Vec{Int(1)}).ExprTag() == nil)
	tt.True(t, (&Vec{}).ExprTag() == nil)

	nsid := &Identifier{Name: Intern("test"), Namespace: Intern("ns")}
	tt.True(t, (&Vec{nsid}).ExprTagName() == nil)
	tt.True(t, (&Vec{InternalId(Intern("test"))}).ExprTagName() == Intern("test"))
}

func TestIdentifier(t *testing.T) {
	ns := Intern("ns")
	k := Intern("test")

	a := InternalId(k)
	nsid := InternalId(k)
	nsid.Namespace = ns
	phid := InternalId(k)
	phid.PlaceHolder = "%="

	tt.Eq(t, "test", a.Inspect())
	tt.Eq(t, "ns:test", nsid.Inspect())
	tt.Eq(t, "%=test", phid.Inspect())

	b := a.Dup().(*Identifier)
	tt.True(t, a != b)
	tt.True(t, a.Token != b.Token)
	tt.True(t, a.Name == b.Name)

	c := a.Expand(Intern("expand"))
	tt.True(t, a != c)
	tt.True(t, a.Token != c.Token)
	tt.True(t, a == c.ExpandedBy)
	tt.True(t, c.Name == Intern("expand"))
}

func TestSig(t *testing.T) {
	x, y := Intern("X"), Intern("Y")
	a := &Sig{
		Required:    []*Keyword{x, y},
		Results:     []*Keyword{y},
		Invalidated: []*Keyword{x},
	}
	tt.Eq(t, "<Sig:[X Y] => [Y] ! [X]>", a.Inspect())

	a.IsInline = true
	tt.Eq(t, "<Sig:-* [X Y] => [Y] ! [X]>", a.Inspect())

	b := a.Dup().(*Sig)
	tt.True(t, a != b)
	tt.True(t, a.Equals(b))
	tt.EqSlice(t, a.Required, b.Required)
	tt.EqSlice(t, a.Results, b.Results)
	tt.EqSlice(t, a.Invalidated, b.Invalidated)

	b.Invalidated = []*Keyword{}
	tt.True(t, !a.Equals(b))
}

func TestMacro(t *testing.T) {
	a := &Macro{}
	tt.Prefix(t, "<Macro:", a.Inspect())
	tt.Eq(t, Value(a), a.Dup())
}

func TestConstexpr(t *testing.T) {
	a := InternalConstexpr(Int(0))
	b := InternalConstexpr(Int(0))
	tt.Prefix(t, "(0)", a.Inspect())
	tt.True(t, a != b)
	tt.True(t, Value(a) != a.Dup())
}

func TestConstFn(t *testing.T) {
	a := &ConstFn{Body: Int(0), Args: []*Keyword{Intern("a")}}
	tt.Prefix(t, "[a](0)", a.Inspect())
	tt.Eq(t, Value(a), a.Dup())
}

func TestSyntaxFn(t *testing.T) {
	a := SyntaxFn(func(cc *Compiler, env *Env, e *Vec) Value { return NIL })
	tt.Prefix(t, "<Syntax:", a.Inspect())
	tt.True(t, nil != a.Dup())
}

func TestBCode(t *testing.T) {
	a := BCode{}
	tt.Eq(t, "00:00:00:00", a.Inspect())
	tt.Eq(t, Value(a), a.Dup())
}

func TestOperand(t *testing.T) {
	a := &Operand{Kind: Intern("A")}
	tt.Eq(t, "{A}", a.Inspect())
	tt.True(t, Value(a) != a.Dup())

	a = &Operand{Kind: Intern("%W"), A0: InternalConstexpr(Int(0))}
	tt.Eq(t, "{%W:(0)}", a.Inspect())
}

func TestInst(t *testing.T) {
	a := NewInst(&Vec{InternalId(Intern("expr"))}, InstLabel, Int(1))
	tt.Eq(t, "0 [1]", a.Inspect())
	tt.True(t, Value(a) != a.Dup())
	tt.Eq(t, "expr", a.ExprTag().String())
}

func TestModule(t *testing.T) {
	a := NewModule(Intern("ModA"), NewEnv(nil))
	tt.Prefix(t, "<Module:", a.Inspect())
	tt.Eq(t, Value(a), a.Dup())
}

func TestNamed(t *testing.T) {
	a := &Named{Name: Intern("named"), Kind: NmLabel, Value: NIL}
	tt.Prefix(t, "<Named:", a.Inspect())
	tt.True(t, Value(a) != a.Dup())
}

func TestLabel(t *testing.T) {
	a := &Label{}
	tt.Prefix(t, "<Label:", a.Inspect())
	tt.True(t, Value(a) != a.Dup())
	tt.True(t, !a.LinkedToData())

	a.Link = NewInst(nil, InstData)
	tt.True(t, a.LinkedToData())

	a.Link = NewInst(nil, InstDS)
	tt.True(t, a.LinkedToData())
}

func TestInline(t *testing.T) {
	a := &Inline{}
	tt.Prefix(t, "<Inline:", a.Inspect())
	tt.True(t, Value(a) != a.Dup())
}

func TestEnv(t *testing.T) {
	a := NewEnv(nil)
	tt.Prefix(t, "<Env:", a.Inspect())
	tt.Eq(t, Value(a), a.Dup())
}

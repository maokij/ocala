package core

import (
	"ocala/internal/tt"
	"testing"
)

func TestScanner(t *testing.T) {
	t.Run("LineBytes", func(t *testing.T) {
		src := ""
		s := &Scanner{Path: "-", Text: []byte(src)}
		s.Init()
		tt.Eq(t, "", string(s.LineBytes(0, true)))
		tt.Eq(t, "", string(s.LineBytes(0, false)))

		src = "012"
		s = &Scanner{Path: "-", Text: []byte(src)}
		s.Init()
		tt.Eq(t, "012", string(s.LineBytes(0, true)))
		tt.Eq(t, "012", string(s.LineBytes(0, false)))

		src = "012\n"
		s = &Scanner{Path: "-", Text: []byte(src)}
		s.Init()
		tt.Eq(t, "012", string(s.LineBytes(0, true)))
		tt.Eq(t, "012\n", string(s.LineBytes(0, false)))
		tt.Eq(t, "", string(s.LineBytes(1, true)))
		tt.Eq(t, "", string(s.LineBytes(1, false)))

		src = "012\n456"
		s = &Scanner{Path: "-", Text: []byte(src)}
		s.Init()
		tt.Eq(t, "012", string(s.LineBytes(0, true)))
		tt.Eq(t, "012\n", string(s.LineBytes(0, false)))
		tt.Eq(t, "456", string(s.LineBytes(1, true)))
		tt.Eq(t, "456", string(s.LineBytes(1, false)))
	})
	t.Run("OnError", func(t *testing.T) {
		var handled *InternalError
		s := &Scanner{}
		s.OnError = func(err *InternalError) { handled = err }
		s.scanError("handled-error")
		tt.Eq(t, "handled-error", handled.message)
	})
	t.Run("SubStringFrom", func(t *testing.T) {
		src := "01234567"
		s := &Scanner{Path: "-", Text: []byte(src)}
		s.Pos = 6
		tt.Eq(t, "012345", s.SubStringFrom(0))
		tt.Eq(t, "345", s.SubStringFrom(3))
		tt.Eq(t, "", s.SubStringFrom(6))
	})
}

func TestToken(t *testing.T) {
	s := &Scanner{Path: "-", Text: []byte("a\nb\nc\n")}
	s.Init()
	a := &Token{From: s, Pt: Pt{Pos: 4, Line: 2}}
	tt.Prefix(t, "<Token:", a.Inspect())
	tt.True(t, Value(a) != a.Dup())

	iid := InternalId(Intern("internal"))
	tt.Eq(t, "", iid.Token.PtPrefix())

	id := Intern("original").ToId(a)
	ex1 := id.Expand(Intern("expanded-1"))
	ex2 := ex1.Expand(Intern("expanded-2"))
	tt.Eq(t, "-:3:0: ", id.Token.PtPrefix())
	tt.Eq(t, "-:3:0: ", ex1.Token.PtPrefix())
	tt.Eq(t, "-:3:0: ", ex2.Token.PtPrefix())

	tt.Eq(t, "  at -:3:0\n   |c\n   |^-- ??\n",
		id.Token.FormatAsErrorLine("at"))
	tt.Eq(t, "  from internal identifier expanded-1\n",
		ex1.Token.FormatAsErrorLine("from"))
	tt.Eq(t, "  from internal identifier expanded-2\n",
		ex2.Token.FormatAsErrorLine("from"))
}

func TestFindToken(t *testing.T) {
	token := &Token{From: InternalParser}
	e := &Vec{&Identifier{Token: token}}
	tt.Eq(t, token, FindToken(token))
	tt.Eq(t, token, FindToken(e.At(0)))
	tt.Eq(t, token, FindToken(e))
	tt.Eq(t, token, FindToken(NewInst(e, InstDS)))
	tt.Eq(t, token, FindToken(&Operand{From: e}))
	tt.Eq(t, token, FindToken(&Constexpr{Token: token}))
	tt.Eq(t, token, FindToken(&Named{Token: token}))

	tt.Eq(t, nil, FindToken(&Identifier{}))
	tt.Eq(t, nil, FindToken(&Vec{}))
	tt.Eq(t, nil, FindToken(&Constexpr{}))
	tt.Eq(t, nil, FindToken(Int(0)))
	tt.Eq(t, nil, FindToken(NewStr("")))
}

func TestInternalError(t *testing.T) {
	cc := &Compiler{}
	s := &Scanner{Path: "-", Text: []byte("a\nb\nc\n")}
	s.Init()
	a := &Token{From: s, Pt: Pt{Pos: 4, Line: 2}}
	b := &Token{From: s, Pt: Pt{Pos: 1, Line: 0}}
	err := cc.ErrorAt(a, nil, b)
	err.SetMessage("error")
	tt.EqSlice(t, []*Token{a, b}, err.Tokens())

	tt.EqText(t, tt.Unindent(`
		-:3:0: compile error: error
		[error #0]
		  at -:3:0
		   |c
		   |^-- ??
		[error #1]
		  at -:1:1
		   |a
		   | ^-- ??
    `), string(err.FullMessage()))
}

package file

import (
	"ocala/core"
	tt "ocala/testutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	Init("../../bin/a.out")
	m.Run()
}

func TestInit(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		err := Init("../../bin/a.out")
		tt.Eq(t, nil, err)
	})

	t.Run("error", func(t *testing.T) {
		err := Init("")
		tt.Eq(t, "invalid installation", err.Error())
	})
}

func TestCheckCode(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		errors := CheckCode("-", []byte("arch z80"), NewCompileOptions())
		tt.Eq(t, 0, len(errors))
	})

	t.Run("ok: test1.oc", func(t *testing.T) {
		cc := BuildCompiler("z80")
		src := tt.Must(os.ReadFile("testdata/test1.oc"))
		errors := checkSyntax(cc, "-", src, NewCompileOptions())
		tt.Eq(t, 0, len(errors))
	})

	t.Run("ok: include", func(t *testing.T) {
		errors := CheckCode("testdata/test.oc", tt.UnindentBytes(`
			arch z80
			include "./test3.oc"
		`), NewCompileOptions())
		tt.Eq(t, 0, len(errors))
	})

	t.Run("error: no arch", func(t *testing.T) {
		errors := CheckCode("-", nil, NewCompileOptions())
		tt.Eq(t, 1, len(errors))

		ei := ErrorInfoOf(errors[0])
		tt.Eq(t, invalidToken, ei[0].Token)
		tt.EqText(t, "arch required", ei[0].Message)
	})

	t.Run("error: unknown arch", func(t *testing.T) {
		errors := CheckCode("-", []byte("arch unknown"), NewCompileOptions())
		tt.Eq(t, 1, len(errors))

		ei := ErrorInfoOf(errors[0])
		tt.Eq(t, invalidToken, ei[0].Token)
		tt.EqText(t, "unknown arch: unknown", ei[0].Message)
	})

	t.Run("error: syntax", func(t *testing.T) {
		errors := CheckCode("-", tt.UnindentBytes(`
			arch z80
			const const
			proc () { ); ) }
			A <- ~ 1
		`), NewCompileOptions())
		tt.Eq(t, 4, len(errors))

		ei := ErrorInfoOf(errors[0])
		tt.Eq(t, 1, ei[0].Token.Line)
		tt.EqText(t, tt.Unindent(`
			parse error: unexpected CONST, expected identifier
			[error #0]
			  at -:2:6
			   |const const
			   |      ^-- ??
		`), ei[0].Message)
	})

	t.Run("error: test2.oc", func(t *testing.T) {
		cc := BuildCompiler("z80")
		src := tt.Must(os.ReadFile("testdata/test2.oc"))
		errors := checkSyntax(cc, "-", src, NewCompileOptions())
		tt.True(t, len(errors) > 0)
	})

	t.Run("error: semantic", func(t *testing.T) {
		errors := CheckCode("-", tt.UnindentBytes(`
			arch z80; flat!
			macro m() { A <- HL }
			m
		`), NewCompileOptions())
		tt.Eq(t, 1, len(errors))

		ei := ErrorInfoOf(errors[0])
		tt.Eq(t, 1, ei[0].Token.Line)
		tt.EqText(t, tt.Unindent(`
			compile error: cannot use HL as operand#2 for LD
			[error #0]
			  at -:2:17
			   |macro m() { A <- HL }
			   |                 ^-- ??
			  from -:3:0
			   |m
			   |^-- ??
		`), ei[0].Message)
	})
}

type testhandler struct {
	files map[string]*File
}

func newTesthandler() *testhandler {
	return &testhandler{files: map[string]*File{}}
}

func (h *testhandler) addFile(path string) *File {
	if f := h.files[path]; f != nil {
		return f
	}
	text := tt.Must(os.ReadFile(path))
	f := NewFile(path, 0, NewCompileOptions())
	f.Update(text)
	h.files[path] = f
	f.Analyze(h.addFile)
	return f
}

func TestFile(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		f := NewFile("-", 0, NewCompileOptions())
		tt.Eq(t, "-", f.Path)

		f.Update(tt.UnindentBytes(`
			arch z80; flat!
			macro mac() { NOP }
			proc fn() { data d001 = byte; RET }
			const a = 1
			const f(x) = x + 1
			data d = byte [a]
			module mod { data d = byte }
			L0:
			struct s { a byte; b byte}
			if 1 { NOP }
		`))

		f.Analyze(newTesthandler().addFile)
		tt.EqText(t, tt.Unindent(`
			<*file.BlockNode:
			    <*file.MacroNode:
			    >
			    <*file.ProcNode:
			        <*file.DataNode:
			        >>
			    <*file.ConstNode:
			    >
			    <*file.ConstFnNode:
			    >
			    <*file.DataNode:
			    >
			    <*file.ModuleNode:
			        <*file.DataNode:
			        >>
			    <*file.LabelNode:
			    >
			    <*file.StructNode:
			        <*file.StructFieldNode:
			        >
			        <*file.StructFieldNode:
			        >>
			    <*file.BlockNode:
			    >>
		`), f.Node.Inspect())

		x := f.TokenIndexAt(68)
		token := f.Tokens[x]
		tt.Eq(t, "RET", TokenString(token))

		env := EnvAt(token)
		tt.True(t, env != nil)
		tt.Eq(t, core.Intern("d001"), env.Lookup(core.Intern("d001")).Name)
	})

	t.Run("ok: arch fallback", func(t *testing.T) {
		f := NewFile("-", 0, NewCompileOptions())
		f.Update([]byte("A <- 1"))
		f.Parse()
		tt.True(t, f.Node != nil)
	})

	t.Run("ok: include", func(t *testing.T) {
		f := NewFile("testdata/test.oc", 0, NewCompileOptions())
		f.Update(tt.UnindentBytes(`
			arch z80
			include "./test3.oc"
		`))
		f.Analyze(newTesthandler().addFile)
		tt.True(t, f.Node != nil)
	})

	t.Run("ok: test1.oc", func(t *testing.T) {
		src := tt.Must(os.ReadFile("testdata/test1.oc"))
		f := NewFile("-", 0, NewCompileOptions())
		f.Update(src)
		f.Analyze(newTesthandler().addFile)
		tt.True(t, f.Node != nil)
	})

	t.Run("ok: test2.oc", func(t *testing.T) {
		src := tt.Must(os.ReadFile("testdata/test2.oc"))
		f := NewFile("-", 0, NewCompileOptions())
		f.Update(src)
		f.Analyze(newTesthandler().addFile)
		tt.True(t, f.Node != nil)
	})
}

func TestSafeGetText(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		f := NewFile("-", 0, NewCompileOptions())
		f.Update([]byte("arch z80"))
		tt.Eq(t, "arch z80", string(f.SafeGetText()))
	})
}

func TestTokenIndexAt(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		f := NewFile("-", 0, NewCompileOptions())
		tt.Eq(t, "-", f.Path)
		f.Update([]byte("arch z80"))

		tt.Eq(t, -1, f.TokenIndexAt(0))
		tt.Eq(t, -1, f.TokenIndexAt(1000))
	})
}

func TestFunctions(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		tt.Eq(t, true, IsDotToken(&core.Token{Kind: tkDTMI}))
		tt.Eq(t, false, IsDotToken(&core.Token{Kind: tkSC}))
		tt.Eq(t, true, IsDotLikeToken(&core.Token{Kind: tkDTMI}))
		tt.Eq(t, true, IsDotLikeToken(&core.Token{Kind: tkDOT_OPERATOR}))
		tt.Eq(t, false, IsDotLikeToken(&core.Token{Kind: tkSC}))
	})
}

func TestSigDetail(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		kw := func(s ...string) []*core.Keyword {
			r := []*core.Keyword{}
			for _, i := range s {
				r = append(r, core.Intern(i))
			}
			return r
		}
		es := []struct {
			want string
			v    *core.Sig
		}{
			{"", &core.Sig{}},
			{"-*", &core.Sig{IsInline: true}},
			{"A", &core.Sig{Required: kw("A")}},
			{"A B", &core.Sig{Required: kw("A", "B")}},
			{"=> A", &core.Sig{Results: kw("A")}},
			{"=> A B", &core.Sig{Results: kw("A", "B")}},
			{"!", &core.Sig{Invalidated: []*core.Keyword{core.KwUNDER}}},
			{"! A", &core.Sig{Invalidated: kw("A")}},
			{"! A B", &core.Sig{Invalidated: kw("A", "B")}},
			{"-* A B => A B ! A B", &core.Sig{
				IsInline:    true,
				Required:    kw("A", "B"),
				Results:     kw("A", "B"),
				Invalidated: kw("A", "B"),
			}},
		}
		for x, i := range es {
			got := SigDetail(i.v)
			tt.Eq(t, i.want, got, x, i.v)
		}
	})
}

func TestAST(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		node := &NodeBase{}
		sub := &NodeBase{}
		node.AddChildNode(sub)
		sub.parent = node
		tt.Eq(t, "<Node:<Node:>>", node.Inspect())
		tt.Eq(t, Node(node), sub.Parent())
		tt.True(t, core.Value(node) != node.Dup())

		f := NewFile("-", 0, NewCompileOptions())
		f.Update(tt.UnindentBytes(`
			macro m001(a b) {}
			proc p001(A => B) {}
			const c001(a b) = 0
		`))
		f.Analyze(newTesthandler().addFile)
		nodes := f.Node.Children()
		tt.Eq(t, "m001(a, b)", nodes[0].(*MacroNode).Detail())
		tt.Eq(t, "p001(A => B)", nodes[1].(*ProcNode).Detail())
		tt.Eq(t, "c001(a, b)", nodes[2].(*ConstFnNode).Detail())
	})
}

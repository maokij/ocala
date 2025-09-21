package ttarch

import (
	"bytes"
	"maps"
	. "ocala/internal/core" //lint:ignore ST1001 core
	_ "ocala/internal/mos6502"
	_ "ocala/internal/z80"
	"os"
	"regexp"
	"strings"
	"testing"
)

func init() {
	RegisterArchs(ArchMap{
		"ttarch":     BuildCompiler,
		"ttarch+ext": BuildCompilerExt,
	})
}

func NewGenerator(src string) *Generator {
	return &Generator{
		InReader:  bytes.NewBuffer([]byte(src)),
		OutWriter: &bytes.Buffer{},
		ErrWriter: &bytes.Buffer{},
		OutPath:   "-",
		ListText:  &[]byte{},
	}
}

var IdTestid = InternalId(Intern("testid"))
var reArch = regexp.MustCompile(`\A\s*arch +([^;\r\n]+)`)

func BuildGenerator(base string, src string) *Generator {
	m := reArch.FindStringSubmatch(src)

	var cc *Compiler
	if len(m) == 0 {
		cc = NewCompiler(base)
	} else {
		a := strings.Split(m[1], " ")
		if len(a) == 2 {
			cc = NewCompiler(a[0] + a[1])
		} else {
			cc = NewCompiler(m[1])
		}
	}

	g := NewGenerator(src)
	g.SetCompiler(cc)
	return g
}

func ExpectCompileOk(t *testing.T, base string, src string) []byte {
	g := BuildGenerator(base, src)
	binary, _, mes := DoCompile(g, "-")
	if mes != "" {
		t.Helper()
		t.Fatalf("%s\n%s", mes, src)
	}
	return binary
}

func ExpectCompileError(t *testing.T, base string, src string) string {
	g := BuildGenerator(base, src)
	binary, _, mes := DoCompile(g, "-")
	if len(binary) != 0 || mes == "" {
		t.Helper()
		t.Fatalf("%s\n%s", mes, src)
	}
	return mes
}

func ExpectGenListOk(t *testing.T, base string, src string) string {
	g := BuildGenerator(base, src)
	g.GenList = true
	_, list, mes := DoCompile(g, "-")
	if len(list) == 0 || mes != "" {
		t.Helper()
		t.Fatalf("%s\n%s", mes, src)
	}
	return string(list)
}

func DoCompile(g *Generator, path string) ([]byte, []byte, string) {
	binary := func() []byte {
		defer g.HandlePanic()
		return g.GenerateBin(g.Compile(path, g.InReader.(*bytes.Buffer).Bytes()))
	}()
	return binary, *g.ListText, g.ErrorMessage()
}

func WarningMessages(g *Generator) []string {
	warns := []string{}
	for _, i := range g.Warnings {
		warns = append(warns, i.Error())
	}
	return warns
}

func CompileTestFile(base string, path string) ([]byte, []byte, string) {
	s, err := os.ReadFile(path + ".oc")
	if err != nil {
		return []byte{}, []byte{1}, err.Error()
	}

	b, err := os.ReadFile(path + ".dat")
	if err != nil {
		return []byte{}, []byte{1}, err.Error()
	}

	g := BuildGenerator(base, string(s))
	a, _, mes := DoCompile(g, "-")
	return a, b, mes
}

func BuildCompiler() *Compiler {
	return &Compiler{
		Arch:            "ttarch",
		InstMap:         instMap,
		InstAliases:     instAliases,
		SyntaxMap:       syntaxMap,
		CtxOpMap:        ctxOpMap,
		ExprToOperand:   exprToOperand,
		AsmOperands:     asmOperands,
		CollectRegs:     collectRegs,
		KwRegA:          kwRegA,
		TokenWords:      tokenWords,
		AdjustOperand:   adjustOperand,
		OperandToNamed:  operandToNamed,
		BMaps:           bmaps,
		TokenAliases:    tokenAliases,
		IsValidProcTail: isValidProcTail,
		AdjustInline:    adjustInline,
		OptimizeBCode:   optimizeBCode,
	}
}

func BuildCompilerExt() *Compiler {
	cc := BuildCompiler()
	cc.Variant = "+ext"
	cc.AsmOperands = maps.Clone(cc.AsmOperands)
	cc.TokenWords = MergeTokenWords(cc.TokenWords, tokenWordsExt)
	cc.InstMap = MergeInstMap(cc.InstMap, instMapExt)
	cc.CtxOpMap = MergeCtxOpMap(cc.CtxOpMap, ctxOpMapExt)
	return cc
}

var syntaxMap = map[*Keyword]SyntaxFn{
	Intern("db"): SyntaxFn(sDb),
	Intern("dw"): SyntaxFn(sDw),
	Intern("ds"): SyntaxFn(sDs),

	Intern("link-as-tests"): SyntaxFn(sLinkAsTests),
	Intern("expect"):        SyntaxFn(sExpect),
	KwOptimize:              SyntaxFn(sOptimize),
}

func exprToOperand(cc *Compiler, e Value) *Operand {
	switch e := e.(type) {
	case *Operand:
		return e
	case *Constexpr:
		return &Operand{From: e, Kind: kwImmNN, A0: e}
	case *Identifier:
		if cc.IsReg(e.Name) || cc.IsCond(e.Name) {
			return &Operand{From: e, Kind: e.Name}
		}
	case *Vec:
		switch e.ExprTagName() {
		case KwTpl:
			return &Operand{From: e, Kind: kwRegPQ, A0: e.At(1), A1: e.At(2)}
		case KwMem:
			if e.Size() > 2 {
				break
			}
			a := e.OperandAt(1)
			if a.Kind == kwImmNN {
				return &Operand{From: e, Kind: kwMemNN, A0: a.A0}
			} else if a.Kind == kwRegX {
				return &Operand{From: e, Kind: kwMemX, A0: a.A0}
			} else if a.Kind == kwRegY {
				return &Operand{From: e, Kind: kwMemY, A0: a.A0}
			}
		}
	}
	return InvalidOperand
}

func adjustOperand(cc *Compiler, e *Operand, n int, etag *Identifier) {
	switch e.Kind {
	case kwImmN:
		if n < -128 || n > 255 {
			e.Kind = kwImmNN
		}
	case kwImmNN:
		if n >= -128 && n <= 255 {
			e.Kind = kwImmN
		}
	}
}

func operandToNamed(cc *Compiler, v Value) *Named {
	return OperandA0ToNamed(cc, v, kwImmNN)
}

func isValidProcTail(cc *Compiler, inst *Inst) bool {
	switch inst.Args[0] {
	case kwJMP, KwJump:
		return len(inst.Args) == 2
	case kwRET, KwReturn:
		return len(inst.Args) == 1
	}
	return false
}

func adjustInline(cc *Compiler, env *Env, insts []*Inst) {
	ci := insts[0]
	for _, i := range insts {
		switch i.Kind {
		case InstLabel, InstCode:
			ci = i
			switch i.Args[0] {
			case kwRET, KwReturn:
				a := i.ExprTag().Expand(KwEndInline).ToConstexpr(env)
				i.Args = append(
					[]Value{KwJump, &Operand{From: i.From, Kind: kwImmNN, A0: a}},
					i.Args[1:]...)
			}
		}
	}

	if !ci.MatchCode(kwJMP, KwJump) || len(ci.Args) != 2 {
		cc.ErrorAt(ci).With("invalid inline proc tail")
	}
	if a := ci.Args[1].(*Operand); a.Kind == kwImmNN &&
		KwEndInline.MatchId(GetConstBody(a.A0)) != nil {
		ci.CommentOut()
	}
}

func optimizeBCode(cc *Compiler, inst *Inst, bcodes []BCode, commit bool) []BCode {
	if commit {
		switch kw := inst.Args[0].(*Keyword); kw {
		case KwJump:
			inst.Args[0] = kwJMP
		case KwCall:
			inst.Args[0] = kwJSR
		case KwReturn:
			inst.Args[0] = kwRET
		}
	}
	return bcodes
}

func collectRegs(regs []*Keyword, reg *Keyword) []*Keyword {
	return append(regs, reg)
}

// SYNTAX: (db a)
func sDb(cc *Compiler, env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, -1, CtModule|CtProc, cc)
	v := &Vec{etag.Expand(KwDataList)}
	for _, i := range (*e)[1:] {
		v.Push(cc.CompileExpr(env, i))
	}
	inst := NewInst(e, InstData, ByteType, Int(1), NIL, Int(DataSizeAuto), v)
	return cc.EmitCode(inst)
}

// SYNTAX: (dw a)
func sDw(cc *Compiler, env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, -1, CtModule|CtProc, cc)
	v := &Vec{etag.Expand(KwDataList)}
	for _, i := range (*e)[1:] {
		v.Push(cc.CompileExpr(env, i))
	}
	inst := NewInst(e, InstData, WordType, Int(1), NIL, Int(DataSizeAuto), v)
	return cc.EmitCode(inst)
}

// SYNTAX: (ds a)
func sDs(cc *Compiler, env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtModule|CtProc, cc)
	e1 := CheckValue(e.At(1), ConstexprT, "count", etag, cc)
	n := EvalConstAs(e1, env, IntT, "count", etag, cc)
	if n < 1 {
		cc.ErrorAt(etag).With("invalid repeat count %d", n)
	}
	return cc.EmitCode(NewInst(e, InstDS, ByteType, n))
}

var linkAsTests = []byte(`link { org 0 0 1; merge text { data byte ["ok"] }; ` +
	`org 0 0 0; merge text _; merge bss _ }`)

// SYNTAX: (link-as-tests )
func sLinkAsTests(cc *Compiler, env *Env, e *Vec) Value {
	CheckExpr(e, 1, 1, CtModule|CtProc, cc)
	cc.CompileExpr(env, cc.Parse("@", linkAsTests))
	return NIL
}

// SYNTAX: (expect )
func sExpect(cc *Compiler, env *Env, e *Vec) Value {
	etag, n := CheckExpr(e, 2, 3, CtModule|CtProc, cc)

	a := e.At(1).(*Constexpr).Body
	b := a
	if n == 2 {
		a = Value(Int(1))
	} else {
		b = e.At(2).(*Constexpr).Body
	}
	cond := &Constexpr{Token: etag.Token, Body: &Vec{etag.Expand(KwEqlOp), a, b}}
	s := NewStr("assertion failed: " + cond.Inspect())
	inst := NewInst(e, InstAssert, cc.CompileExpr(env, cond), s)
	return cc.EmitCode(inst)
}

// SYNTAX: (optimize ...)
func sOptimize(cc *Compiler, env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, -1, CtModule|CtProc, cc)
	k := CheckConstPlainId(e.At(1), "kind", etag, cc)

	CheckToplevelEnv(env, etag, cc)
	cc.ProcessDefaultOptimizeOption(env, e, k, etag)
	return NIL
}

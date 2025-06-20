package ttarch

import (
	"bytes"
	. "ocala/internal/core" //lint:ignore ST1001 core
	"os"
	"strings"
)

func BuildGenerator(cc *Compiler, src string) *Generator {
	g := &Generator{
		InReader:  bytes.NewBuffer([]byte(src)),
		OutWriter: &bytes.Buffer{},
		ErrWriter: &bytes.Buffer{},
		OutPath:   "-",
		ListText:  &[]byte{},
		Archs:     map[string]func() *Compiler{"ttarch": BuildCompiler},
	}
	if cc != nil {
		g.SetCompiler(cc)
	}
	return g
}

func Compile(cc *Compiler, src string) ([]byte, string) {
	g := BuildGenerator(cc, src)
	binary, _, mes := DoCompile(g, "-")
	return binary, mes
}

func GenList(cc *Compiler, src string) (string, string) {
	g := BuildGenerator(cc, src)
	g.GenList = true
	_, list, mes := DoCompile(g, "-")
	return string(list), mes
}

func DoCompile(g *Generator, path string) ([]byte, []byte, string) {
	binary := func() []byte {
		defer g.HandlePanic()
		return g.GenerateBin(g.Compile(path, g.InReader.(*bytes.Buffer).Bytes()))
	}()
	return binary, *g.ListText, g.ErrorMessage()
}

func CompileTestFile(cc *Compiler, path string) ([]byte, []byte, string) {
	s, err := os.ReadFile(path + ".oc")
	if err != nil {
		return []byte{}, []byte{1}, err.Error()
	}

	b, err := os.ReadFile(path + ".dat")
	if err != nil {
		return []byte{}, []byte{1}, err.Error()
	}

	a, mes := Compile(cc, string(s))
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
		OperandToAsm:    operandToAsm,
		CollectRegs:     collectRegs,
		KwRegA:          kwRegA,
		TokenWords:      tokenWords,
		AdjustOperand:   adjustOperand,
		BMaps:           bmaps,
		TokenAliases:    tokenAliases,
		IsValidProcTail: isValidProcTail,
		AdjustInline:    adjustInline,
	}
}

var kwJMP = Intern("JMP")
var kwJSR = Intern("JSR")
var kwRET = Intern("RET")

var syntaxMap = map[*Keyword]SyntaxFn{
	KwJump:       SyntaxFn(sJump), // for loop syntax
	KwCall:       SyntaxFn(sCall),
	Intern("db"): SyntaxFn(sDb),
	Intern("dw"): SyntaxFn(sDw),
	Intern("ds"): SyntaxFn(sDs),
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

func adjustOperand(cc *Compiler, e *Operand, etag *Identifier) {
	c, ok := e.A0.(*Constexpr)
	if !ok {
		return
	}

	n := EvalConstAs(c, c.Env, IntT, "operand", etag, cc)
	cc.Constvals[c] = n

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

func operandToAsm(g *Generator, e *Operand) string {
	a := operandToAsmMap[e.Kind]
	if a.t {
		return strings.Replace(a.s, "%", g.ValueToAsm(nil, e.A0), 1)
	}
	return a.s
}

func isValidProcTail(cc *Compiler, inst *Inst) bool {
	return inst.MatchCode(kwJMP, kwRET)
}

func adjustInline(cc *Compiler, insts []*Inst) {
	ci := insts[0]
	for _, i := range insts {
		switch i.Kind {
		case InstLabel, InstCode:
			ci = i
			switch i.Args[0] {
			case kwRET:
				a := i.ExprTag().Expand(KwEndInline).ToConstexpr(nil)
				i.Args = []Value{kwJMP, &Operand{From: i.From, Kind: kwImmNN, A0: a}}
			}
		}
	}

	if !ci.MatchCode(kwJMP) {
		cc.ErrorAt(ci).With("invalid inline proc tail")
	}
	if a := ci.Args[1].(*Operand); a.Kind == kwImmNN &&
		KwEndInline.MatchId(GetConstBody(a.A0)) != nil {
		*ci = *NewInst(ci.From, InstMisc, KwUNDER)
	}
}

func collectRegs(regs []*Keyword, reg *Keyword) []*Keyword {
	return append(regs, reg)
}

// SYNTAX: (#.jump addr cond)
func sJump(cc *Compiler, env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtProc, cc)
	jump := &Vec{etag.ExpandedBy.Expand(kwJMP), e.At(1)}
	if e.At(2) != NIL {
		cc.ErrorAt(etag).With("conditional jump is not supported")
	}
	return cc.CompileExpr(env, jump)
}

// SYNTAX: (#.call addr cond)
func sCall(cc *Compiler, env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtProc, cc)
	call := &Vec{etag.ExpandedBy.Expand(kwJSR), e.At(1)}
	if e.At(2) != NIL {
		cc.ErrorAt(etag).With("conditional call is not supported")
	}
	return cc.CompileExpr(env, call)
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

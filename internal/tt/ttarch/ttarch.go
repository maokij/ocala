package ttarch

import (
	"bytes"
	"maps"
	. "ocala/internal/core" //lint:ignore ST1001 core
	"os"
	"slices"
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
		AsmOperands:     asmOperands,
		CollectRegs:     collectRegs,
		KwRegA:          kwRegA,
		TokenWords:      tokenWords,
		AdjustOperand:   adjustOperand,
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
	maps.Copy(cc.AsmOperands, asmOperandsExt)

	cc.TokenWords = slices.Clone(cc.TokenWords)
	cc.TokenWords[0] = append(cc.TokenWords[0], tokenWordsExt[0]...)
	cc.TokenWords[1] = append(cc.TokenWords[1], tokenWordsExt[1]...)
	cc.TokenWords[2] = append(cc.TokenWords[2], tokenWordsExt[2]...)
	cc.TokenWords[3] = append(cc.TokenWords[3], tokenWordsExt[3]...)

	cc.InstMap = MergeInstMap(cc.InstMap, instMapExt)
	cc.CtxOpMap = MergeCtxOpMap(cc.CtxOpMap, ctxOpMapExt)
	return cc
}

var syntaxMap = map[*Keyword]SyntaxFn{
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

func isValidProcTail(cc *Compiler, inst *Inst) bool {
	return inst.MatchCode(kwJMP, kwRET, KwJump)
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

func optimizeBCode(cc *Compiler, inst *Inst, bcodes []BCode, commit bool) []BCode {
	if inst.MatchCode(KwJump, KwCall) && commit {
		n := len(inst.Args)
		switch kw := inst.Args[0].(*Keyword); kw {
		case KwJump:
			if n == 2 {
				inst.Args[0] = kwJMP
			} else {
				inst.Args[0] = kwBCO
			}
		case KwCall:
			inst.Args[0] = kwJSR
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

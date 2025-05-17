package mos6502

import (
	. "ocala/internal/core" //lint:ignore ST1001 core
	"strings"
)

func BuildCompiler() *Compiler {
	return &Compiler{
		Arch:            "mos6502",
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
		OppositeConds:   oppositeConds,
		IsValidProcTail: isValidProcTail,
		AdjustInline:    adjustInline,
	}
}

var kwJMP = Intern("JMP")
var kwJSR = Intern("JSR")
var kwRTS = Intern("RTS")
var kwRTI = Intern("RTI")
var kwOutOfRangeMem = Intern("out-of-range")

var syntaxMap = map[*Keyword]SyntaxFn{
	KwJump: SyntaxFn(sJump), // for loop syntax
	KwCall: SyntaxFn(sCall),
}

//

func exprToOperand(cc *Compiler, e Value) *Operand {
	switch e := e.(type) {
	case *Operand:
		return e
	case *Constexpr:
		return &Operand{Kind: kwImmNN, A0: e}
	case *Identifier: // REG
		if cc.IsReg(e.Name) || cc.IsCond(e.Name) {
			return &Operand{Kind: e.Name}
		}
	case *Vec:
		if e.ExprTagName() == KwMem && e.Size() <= 3 {
			a := e.OperandAt(1)
			b := e.OperandAt(2)
			if a.Kind == kwImmNN {
				if b == NoOperand {
					return &Operand{Kind: kwMemAN, A0: a.A0}
				} else if b.Kind == kwRegX {
					return &Operand{Kind: kwMemAX, A0: a.A0}
				} else if b.Kind == kwRegY {
					return &Operand{Kind: kwMemAY, A0: a.A0}
				}
			} else if a.Kind == kwMemAN {
				if b == NoOperand {
					return &Operand{Kind: kwMemIN, A0: a.A0}
				} else if b.Kind == kwRegY {
					return &Operand{Kind: kwMemIY, A0: a.A0}
				}
			} else if a.Kind == kwMemAX {
				if b == NoOperand {
					return &Operand{Kind: kwMemIX, A0: a.A0}
				}
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
	case kwMemZN:
		if n < -128 || n > 255 {
			e.Kind = kwMemAN
		}
	case kwMemZX:
		if n < -128 || n > 255 {
			e.Kind = kwMemAX
		}
	case kwMemZY:
		if n < -128 || n > 255 {
			e.Kind = kwMemAY
		}
	case kwMemAN:
		if n >= -128 && n <= 255 {
			e.Kind = kwMemZN
		}
	case kwMemAX:
		if n >= -128 && n <= 255 {
			e.Kind = kwMemZX
		}
	case kwMemAY:
		if n >= -128 && n <= 255 {
			e.Kind = kwMemZY
		}
	case kwMemIX:
		if n < -128 || n > 255 {
			e.Kind = kwOutOfRangeMem
		}
	case kwMemIY:
		if n < -128 || n > 255 {
			e.Kind = kwOutOfRangeMem
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
	return inst.MatchCode(kwJMP, kwRTS, kwRTI)
}

func adjustInline(cc *Compiler, insts []*Inst) {
	ci := insts[0]
	for _, i := range insts {
		switch i.Kind {
		case InstLabel, InstCode:
			ci = i
			switch i.Args[0] {
			case kwRTS:
				a := i.ExprTag().Expand(KwEndInline).ToConstexpr(nil)
				i.Args = []Value{kwJMP, &Operand{Kind: kwMemAN, A0: a}}
			case kwRTI:
				cc.RaiseCompileError(i.ExprTag(), "unsupported instruction in inline code")
			}
		}
	}

	if !ci.MatchCode(kwJMP) {
		cc.RaiseCompileError(ci.ExprTag(), "invalid inline proc tail")
	}
	if a := ci.Args[1].(*Operand); a.Kind == kwMemAN &&
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
	jump := &Vec{etag.ExpandedBy.Expand(kwJMP), &Vec{etag.ExpandedBy.Expand(KwMem), e.At(1)}}
	if e.At(2) != NIL {
		cc.RaiseCompileError(etag, "conditional jump is not supported")
	}
	return cc.CompileExpr(env, jump)
}

// SYNTAX: (#.call addr cond)
func sCall(cc *Compiler, env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtProc, cc)
	call := &Vec{etag.ExpandedBy.Expand(kwJSR), &Vec{etag.ExpandedBy.Expand(KwMem), e.At(1)}}
	if e.At(2) != NIL {
		cc.RaiseCompileError(etag, "conditional call is not supported")
	}
	return cc.CompileExpr(env, call)
}

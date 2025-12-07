package mos6502

import (
	. "ocala/core" //lint:ignore ST1001 core
)

func init() {
	RegisterArchs(ArchMap{
		"6502":    BuildCompiler,
		"mos6502": BuildCompiler,
	})
}

func BuildCompiler() *Compiler {
	return &Compiler{
		Arch:            "mos6502",
		InstMap:         instMap,
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

var kwOutOfRangeMem = Intern("out-of-range")
var kwNearJump = Intern("near-jump")

var syntaxMap = map[*Keyword]SyntaxFn{
	KwOptimize: SyntaxFn(sOptimize),
}

//

func exprToOperand(cc *Compiler, e Value) *Operand {
	switch e := e.(type) {
	case *Operand:
		return e
	case *Constexpr:
		return &Operand{From: e, Kind: kwImmNN, A0: e}
	case *Identifier: // REG
		if cc.IsReg(e.Name) || cc.IsCond(e.Name) {
			return &Operand{From: e, Kind: e.Name}
		}
	case *Vec:
		if e.ExprTagName() == KwMem && e.Size() <= 3 {
			a := e.OperandAt(1)
			b := e.OperandAt(2)
			if a.Kind == kwImmNN {
				if b == NoOperand {
					return &Operand{From: e, Kind: kwMemAN, A0: a.A0}
				} else if b.Kind == kwRegX {
					return &Operand{From: e, Kind: kwMemAX, A0: a.A0}
				} else if b.Kind == kwRegY {
					return &Operand{From: e, Kind: kwMemAY, A0: a.A0}
				}
			} else if a.Kind == kwMemAN {
				if b == NoOperand {
					return &Operand{From: e, Kind: kwMemIN, A0: a.A0}
				} else if b.Kind == kwRegY {
					return &Operand{From: e, Kind: kwMemIY, A0: a.A0}
				}
			} else if a.Kind == kwMemAX {
				if b == NoOperand {
					return &Operand{From: e, Kind: kwMemIX, A0: a.A0}
				}
			}
		}
	}
	return &Operand{From: e, Kind: KwInvalidOperand}
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

func operandToNamed(cc *Compiler, v Value) *Named {
	return OperandA0ToNamed(cc, v, kwImmNN)
}

func isValidProcTail(cc *Compiler, inst *Inst) bool {
	switch inst.Args[0] {
	case kwJMP, KwJump:
		return len(inst.Args) == 2
	case kwRTS, kwRTI, KwReturn:
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
			case kwRTS, KwReturn:
				a := i.ExprTag().Expand(KwEndInline).ToConstexpr(env)
				i.Args = append(
					[]Value{KwJump, &Operand{From: i.From, Kind: kwImmNN, A0: a}},
					i.Args[1:]...)
			case kwRTI:
				cc.ErrorAt(i).With("unsupported instruction in inline code")
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

func updateInst(inst *Inst, kw *Keyword) {
	inst.Args[1].(*Operand).Kind = kwMemAN
	inst.Args = []Value{kw, inst.Args[1]}
}

var optJumpNames = map[byte]*Keyword{
	0x10: kwBPL,
	0x30: kwBMI,
	0x50: kwBVC,
	0x70: kwBVS,
	0x90: kwBCC,
	0xB0: kwBCS,
	0xD0: kwBNE,
	0xF0: kwBEQ,
	0x4C: kwJMP,
}

var optJumpCodes = map[byte]byte{
	0x10: 0x30,
	0x30: 0x10,
	0x50: 0x70,
	0x70: 0x50,
	0x90: 0xB0,
	0xB0: 0x90,
	0xD0: 0xF0,
	0xF0: 0xD0,
}

func optimizeBCode(cc *Compiler, inst *Inst, bcodes []BCode, commit bool) []BCode {
	n := len(inst.Args)
	switch inst.Args[0].(*Keyword) {
	case KwJump:
		code := bcodes[0].A0
		if n == 3 {
			c := inst.Args[1].(*Operand).A0.(*Constexpr)
			v := cc.Constvals[c].(Int)
			d := int(v) - cc.Pc
			if d > 0 { // forward
				d -= inst.Size
			} else { // backward
				d -= 2
			}
			if d >= -128 && d <= 127 {
				code = optJumpCodes[code]
				bcodes = []BCode{{Kind: BcByte, A0: code}, {Kind: BcByte, A0: byte(d)}}
			}
		}
		if commit && len(bcodes) < 5 {
			updateInst(inst, optJumpNames[code])
		}
	case KwCall:
		if commit {
			updateInst(inst, kwJSR)
		}
	case KwReturn:
		if commit && len(bcodes) == 1 {
			inst.Args[0] = kwRTS
		}
	}
	return bcodes
}

func noOptimizeBCode(cc *Compiler, inst *Inst, bcodes []BCode, commit bool) []BCode {
	if commit {
		switch inst.Args[0].(*Keyword) {
		case KwJump:
			if len(inst.Args) == 2 {
				updateInst(inst, kwJMP)
			}
		case KwCall:
			updateInst(inst, kwJSR)
		case KwReturn:
			if len(inst.Args) == 1 {
				inst.Args[0] = kwRTS
			}
		}
	}
	return bcodes
}

func collectRegs(regs []*Keyword, reg *Keyword) []*Keyword {
	return append(regs, reg)
}

// SYNTAX: (optimize ...)
func sOptimize(cc *Compiler, env *Env, e *Vec) Value {
	etag, n := CheckExpr(e, 2, -1, CtModule|CtProc, cc)
	k := CheckConstPlainId(e.At(1), "kind", etag, cc)

	CheckToplevelEnv(env, etag, cc)
	switch k.Name {
	case kwNearJump:
		v := Int(1)
		if n == 3 {
			v = CheckAndEvalConstAs(e.At(2), env, IntT, "option", etag, cc)
		}

		cc.OptimizeBCode = noOptimizeBCode
		if v != Int(0) {
			cc.OptimizeBCode = optimizeBCode
		}
	default:
		cc.ProcessDefaultOptimizeOption(env, e, k, etag)
	}
	return NIL
}

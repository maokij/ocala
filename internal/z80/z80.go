package z80

import (
	"maps"
	. "ocala/internal/core" //lint:ignore ST1001 core
	"slices"
)

func BuildCompiler() *Compiler {
	return &Compiler{
		Arch:            "z80",
		InstMap:         instMap,
		SyntaxMap:       syntaxMap,
		CtxOpMap:        ctxOpMap,
		ExprToOperand:   exprToOperand,
		AsmOperands:     asmOperands,
		TrimAsmOperand:  true,
		CollectRegs:     collectRegs,
		KwRegA:          kwRegA,
		TokenWords:      tokenWords,
		AdjustOperand:   adjustOperand,
		BMaps:           bmaps,
		TokenAliases:    tokenAliases,
		IsValidProcTail: isValidProcTail,
		AdjustInline:    adjustInline,
		OptimizeBCode:   optimizeBCodeForward,
	}
}

func BuildCompilerUndocumented() *Compiler {
	cc := BuildCompiler()
	cc.Variant = "+undocumented"
	cc.AsmOperands = maps.Clone(cc.AsmOperands)
	maps.Copy(cc.AsmOperands, asmOperandsUndocumented)

	cc.TokenWords = slices.Clone(cc.TokenWords)
	cc.TokenWords[0] = append(cc.TokenWords[0], tokenWordsUndocumented[0]...)
	cc.TokenWords[1] = append(cc.TokenWords[1], tokenWordsUndocumented[1]...)
	cc.TokenWords[2] = append(cc.TokenWords[2], tokenWordsUndocumented[2]...)
	cc.TokenWords[3] = append(cc.TokenWords[3], tokenWordsUndocumented[3]...)

	cc.InstMap = MergeInstMap(cc.InstMap, instMapUndocumented)
	cc.CtxOpMap = MergeCtxOpMap(cc.CtxOpMap, ctxOpMapUndocumented)
	return cc
}

var kwNearJump = Intern("near-jump")

var syntaxMap = map[*Keyword]SyntaxFn{
	KwOptimize:      SyntaxFn(sOptimize),
	Intern("#.LDP"): SyntaxFn(sLdp),
}

//

func isReg16(a *Keyword) bool {
	switch a {
	case kwRegAF, kwAltAF, kwRegBC, kwRegDE, kwRegHL, kwRegSP, kwRegIX, kwRegIY:
		return true
	}
	return false
}

func toMemReg(a *Keyword) *Keyword {
	switch a {
	case kwRegC:
		return kwMemC
	case kwRegBC:
		return kwMemBC
	case kwRegDE:
		return kwMemDE
	case kwRegHL:
		return kwMemHL
	case kwRegIX:
		return kwMemIX
	case kwRegIY:
		return kwMemIY
	case kwRegSP:
		return kwMemSP
	}
	return nil
}

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
		switch e.ExprTagName() {
		case KwTpl: // (tpl ...)
			return &Operand{From: e, Kind: kwRegPQ, A0: e.At(1), A1: e.At(2)}
		case KwMem:
			if e.Size() > 3 {
				break
			}
			a := e.OperandAt(1)
			b := e.OperandAt(2)
			if a.Kind == kwImmNN {
				if b == NoOperand {
					return &Operand{From: e, Kind: kwMemNN, A0: a.A0}
				}
			} else if k := toMemReg(a.Kind); k != nil {
				if k == kwMemIX || k == kwMemIY {
					if b == NoOperand {
						return &Operand{From: e, Kind: k, A0: InternalConstexpr(Int(0))}
					} else if b.Kind == kwImmNN {
						return &Operand{From: e, Kind: k, A0: b.A0}
					}
				} else if b == NoOperand {
					return &Operand{From: e, Kind: k}
				}
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
	case kwMemN:
		if n < -128 || n > 255 {
			e.Kind = kwMemNN
		}
	case kwMemNN:
		if n >= -128 && n <= 255 {
			e.Kind = kwMemN
		}
	}
}

func isValidProcTail(cc *Compiler, inst *Inst) bool {
	switch inst.Args[0] {
	case kwJP, kwJR, KwJump:
		return len(inst.Args) == 2
	case kwRET, kwRETI, kwRETN:
		return len(inst.Args) == 1
	}
	return false
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
				i.Args = append(i.Args, &Operand{From: i.From, Kind: kwImmNN, A0: a})
				i.Args[0] = kwJP
			case kwRETI, kwRETN:
				cc.ErrorAt(i).With("unsupported instruction in inline code")
			}
		}
	}

	if !ci.MatchCode(kwJP, kwJR) {
		cc.ErrorAt(ci).With("invalid inline proc tail")
	}
	if a := ci.Args[1].(*Operand); a.Kind == kwImmNN &&
		KwEndInline.MatchId(GetConstBody(a.A0)) != nil {
		*ci = *NewInst(ci.From, InstMisc, KwUNDER)
	}
}

func updateInst(inst *Inst, kw *Keyword) {
	if len(inst.Args) == 2 {
		inst.Args = []Value{kw, inst.Args[1]}
	} else {
		inst.Args = []Value{kw, inst.Args[2], inst.Args[1]}
	}
}

var optJumpCodes = map[byte]byte{
	0xc3: 0x18, // JP
	0xc2: 0x20, // JP NZ
	0xca: 0x28, // JP Z
	0xd2: 0x30, // JP NC
	0xda: 0x38, // JP C
}

func optimizeBCode(cc *Compiler, inst *Inst, bcodes []BCode, commit bool, full bool) []BCode {
	switch inst.Args[0].(*Keyword) {
	case KwJump:
		if code := optJumpCodes[bcodes[0].A0]; code != 0 {
			c := inst.Args[1].(*Operand).A0.(*Constexpr)
			v := cc.Constvals[c].(Int)
			d := int(v) - cc.Pc
			if d > 0 { // forward
				d -= inst.Size
			} else { // backward
				d -= 2
			}
			if d >= -128 && d <= 127 && (full || d >= 0) {
				bcodes = []BCode{{Kind: BcByte, A0: code}, {Kind: BcByte, A0: byte(d)}}
			}
		}

		if commit {
			kw := kwJP
			if len(bcodes) == 2 {
				kw = kwJR
			}
			updateInst(inst, kw)
		}
	case KwCall:
		if commit {
			updateInst(inst, kwCALL)
		}
	}
	return bcodes
}

func noOptimizeBCode(cc *Compiler, inst *Inst, bcodes []BCode, commit bool) []BCode {
	if commit {
		switch inst.Args[0].(*Keyword) {
		case KwJump:
			updateInst(inst, kwJP)
		case KwCall:
			updateInst(inst, kwCALL)
		}
	}
	return bcodes
}

func optimizeBCodeForward(cc *Compiler, inst *Inst, bcodes []BCode, commit bool) []BCode {
	return optimizeBCode(cc, inst, bcodes, commit, false)
}

func optimizeBCodeFull(cc *Compiler, inst *Inst, bcodes []BCode, commit bool) []BCode {
	return optimizeBCode(cc, inst, bcodes, commit, true)
}

func collectRegs(regs []*Keyword, reg *Keyword) []*Keyword {
	if isReg16(reg) && reg != kwRegIX && reg != kwRegIY {
		h, l, _ := splitReg(reg)
		if h != nil {
			return append(regs, h, l)
		}
	}
	return append(regs, reg)
}

func splitReg(a *Keyword) (*Keyword, *Keyword, bool) {
	switch a {
	case kwRegBC:
		return kwRegB, kwRegC, true
	case kwRegDE:
		return kwRegD, kwRegE, true
	case kwRegHL:
		return kwRegH, kwRegL, true
	default:
		return nil, nil, false
	}
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
		if v == Int(1) {
			cc.OptimizeBCode = optimizeBCodeForward
		} else if v > Int(1) {
			cc.OptimizeBCode = optimizeBCodeFull
		}
	default:
		cc.ErrorAt(etag).With("unknown optimizer: %s", k)
	}
	return NIL
}

// SYNTAX: (#.LDP ...)
func sLdp(cc *Compiler, env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtProc, cc)
	a := cc.ExprToOperand(cc, e.At(1))
	b := cc.ExprToOperand(cc, e.At(2))

	if a.Kind == kwRegPQ {
		h, l, ok := splitReg(b.Kind)
		if ok {
			cc.CompileExpr(env, &Vec{etag.Expand(kwLD), a.A0, etag.Expand(h)})
			cc.CompileExpr(env, &Vec{etag.Expand(kwLD), a.A1, etag.Expand(l)})
			return NIL
		}
	} else if b.Kind == kwRegPQ {
		h, l, ok := splitReg(a.Kind)
		if ok {
			cc.CompileExpr(env, &Vec{etag.Expand(kwLD), etag.Expand(h), b.A0})
			cc.CompileExpr(env, &Vec{etag.Expand(kwLD), etag.Expand(l), b.A1})
			return NIL
		}
	}
	cc.ErrorAt(etag).With("invalid operands for LDP(#1 = %s, #2 = %s)", a.Kind, b.Kind)
	return NIL
}

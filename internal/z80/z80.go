package z80

import (
	. "ocala/internal/core" //lint:ignore ST1001 core
	"strings"
)

func BuildCompiler() *Compiler {
	return &Compiler{
		Arch:            "z80",
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

var kwLD = Intern("LD")
var kwJP = Intern("JP")
var kwJR = Intern("JR")
var kwRET = Intern("RET")
var kwRETI = Intern("RETI")
var kwRETN = Intern("RETN")
var kwCALL = Intern("CALL")
var kwNearJump = Intern("near-jump")

var syntaxMap = map[*Keyword]SyntaxFn{
	KwJump:          SyntaxFn(sJump), // for loop syntax
	KwCall:          SyntaxFn(sCall),
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
		return &Operand{Kind: kwImmNN, A0: e}
	case *Identifier: // REG
		if cc.IsReg(e.Name) || cc.IsCond(e.Name) {
			return &Operand{Kind: e.Name}
		}
	case *Vec:
		switch e.ExprTagName() {
		case KwTpl: // (tpl ...)
			return &Operand{Kind: kwRegPQ, A0: e.At(1), A1: e.At(2)}
		case KwMem:
			if e.Size() > 3 {
				break
			}
			a := e.OperandAt(1)
			b := e.OperandAt(2)
			if a.Kind == kwImmNN {
				if b == NoOperand {
					return &Operand{Kind: kwMemNN, A0: a.A0}
				}
			} else if k := toMemReg(a.Kind); k != nil {
				if k == kwMemIX || k == kwMemIY {
					if b == NoOperand {
						return &Operand{Kind: k, A0: InternalConstexpr(Int(0))}
					} else if b.Kind == kwImmNN {
						return &Operand{Kind: k, A0: b.A0}
					}
				} else if b == NoOperand {
					return &Operand{Kind: k}
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

func operandToAsm(g *Generator, e *Operand) string {
	a := operandToAsmMap[e.Kind]
	s := a.s
	if a.t {
		s = strings.Replace(a.s, "%", g.ValueToAsm(nil, e.A0), 1)
	}
	if s[0] == '0' && s[1] == '+' && s[2] == ' ' && s[3] != '(' {
		s = s[3:]
	}
	return s
}

func isValidProcTail(cc *Compiler, inst *Inst) bool {
	switch inst.Args[0] {
	case kwJP, kwJR:
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
				i.Args = append(i.Args, &Operand{Kind: kwImmNN, A0: a})
				i.Args[0] = kwJP
			case kwRETI, kwRETN:
				cc.RaiseCompileError(i.ExprTag(), "unsupported instruction in inline code")
			}
		}
	}

	if !ci.MatchCode(kwJP, kwJR) {
		cc.RaiseCompileError(ci.ExprTag(), "invalid inline proc tail")
	}
	if a := ci.Args[1].(*Operand); a.Kind == kwImmNN &&
		KwEndInline.MatchId(GetConstBody(a.A0)) != nil {
		*ci = *NewInst(ci.From, InstMisc, KwUNDER)
	}
}

var optJumps = map[byte]byte{
	0xc3: 0x18, // JP
	0xc2: 0x20, // JP NZ
	0xca: 0x28, // JP Z
	0xd2: 0x30, // JP NC
	0xda: 0x38, // JP C
}

func optimizeBCode(cc *Compiler, inst *Inst, bcodes []BCode, commit bool) []BCode {
	if op := inst.Args[0].(*Keyword); op != kwJP {
		return bcodes
	}

	n, ok := optJumps[bcodes[0].A0]
	if !ok || inst.From.ExprTagName() == KwVolatile {
		return bcodes
	}

	x := 1
	if n != 0x18 {
		x = 2
	}

	c := inst.Args[x].(*Operand).A0.(*Constexpr)
	v := cc.Constvals[c].(Int)
	d := int(v) - cc.Pc - 2

	if d < -128 || d > 127 {
		return bcodes
	}

	if commit {
		inst.Args[0] = kwJR
	}

	return []BCode{{Kind: BcByte, A0: n}, {Kind: BcByte, A0: byte(d)}}
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

// SYNTAX: (#.jump addr cond)
func sJump(cc *Compiler, env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtProc, cc)
	jump := &Vec{etag.ExpandedBy.Expand(kwJP)}
	if e.At(2) != NIL {
		jump.Push(e.At(2))
	}
	jump.Push(e.At(1))
	return cc.CompileExpr(env, jump)
}

// SYNTAX: (#.call addr cond)
func sCall(cc *Compiler, env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtProc, cc)
	call := &Vec{etag.ExpandedBy.Expand(kwCALL)}
	if e.At(2) != NIL {
		call.Push(e.At(2))
	}
	call.Push(e.At(1))
	return cc.CompileExpr(env, call)
}

// SYNTAX: (optimize ...)
func sOptimize(cc *Compiler, env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, -1, CtModule|CtProc, cc)
	k := CheckConstPlainId(e.At(1), "kind", etag, cc)

	CheckToplevelEnv(env, etag, cc)
	optimizer := cc.GetOptimizer()
	switch k.Name {
	case kwNearJump:
		optimizer.OptimizeBCode = optimizeBCode
	default:
		cc.RaiseCompileError(etag, "unknown optimizer: %s", k.String())
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
	cc.RaiseCompileError(etag, "invalid operands for LDP(#1 = %s, #2 = %s)", a.Kind, b.Kind)
	return NIL
}

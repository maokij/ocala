package core

import (
	"fmt"
)

// SPECIAL: (__PC__)
func (cc *Compiler) sCurloc(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 1, 1, CtConstexpr, cc)
	CheckPhase(PhLink, etag, cc)
	return Int(cc.Pc)
}

// SPECIAL: (__ORG__)
func (cc *Compiler) sCurorg(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 1, 1, CtConstexpr, cc)
	CheckPhase(PhLink, etag, cc)
	return Int(cc.Org)
}

// SPECIAL: (loaded-as-main?)
func (cc *Compiler) sLoadedAsMain(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 1, 1, CtConstexpr, cc)
	CheckPhase(PhCompile, etag, cc)
	return BoolInt(cc.MainPath == cc.InPath)
}

// SPECIAL: (<reserved>)
func (cc *Compiler) sReserved(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 1, 1, CtConstexpr, cc)
	cc.RaiseCompileError(etag, "reserved but undefined")
	return NIL
}

// SPECIAL: (__FILE__)
// func (cc *Compiler) sFilename(env *Env, e *Vec) Value

// SYNTAX: (include path)
// func (cc *Compiler) sInclude(env *Env, e *Vec) Value

// SYNTAX: (load-file path)
// func (cc *Compiler) sLoadFile(env *Env, e *Vec) Value

// SYNTAX: (arch s)
func (cc *Compiler) sArch(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtModule|CtProc, cc)
	arch := CheckConst(e.At(1), IdentifierT, "arch name", etag, cc)

	CheckToplevelEnv(env, etag, cc)
	if cc.Arch != arch.String() {
		cc.RaiseCompileError(etag, "current target arch is %s", cc.Arch)
	}
	return NIL
}

// SYNTAX: (align a)
func (cc *Compiler) sAlign(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtModule|CtProc, cc)
	a := CheckAndEvalConstAs(e.At(1), env, IntT, "argument", etag, cc)
	return cc.EmitCode(NewInst(e, InstAlign, a))
}

// SYNTAX: (#.label l)
func (cc *Compiler) sLabel(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtModule|CtProc, cc)
	a := CheckConst(e.At(1), IdentifierT, "label name", etag, cc)
	nm := cc.InstallNamed(env, a, NmLabel, &Label{})
	return cc.EmitCode(NewInst(e, InstLabel, nm))
}

// SYNTAX: (#.tpl ...)
func (cc *Compiler) sTpl(env *Env, e *Vec) Value {
	return e
}

// SYNTAX: (#.prog ...)
func (cc *Compiler) sProg(env *Env, e *Vec) Value {
	return cc.sBlock(env.Enter(), e)
}

// SYNTAX: (#.block ...)
func (cc *Compiler) sBlock(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 1, -1, 0, cc)
	for _, i := range (*e)[1:] {
		if _, ok := i.(*Vec); !ok {
			cc.RaiseCompileError(etag, "invalid do form(%T)", i)
		}
		cc.CompileExpr(env, i)
	}
	return NIL
}

// SYNTAX: (do a)
func (cc *Compiler) sDo(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, 0, cc)
	v := CheckBlockForm(e.At(1), "body", etag, cc)
	return cc.CompileExpr(env, v)
}

// SYNTAX: (<volatile> a)
func (cc *Compiler) sVolatile(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, -1, CtProc, cc)
	op := CheckConst(e.At(1), IdentifierT, "op", etag, cc)
	nm := cc.LookupNamed(env, op)
	if nm == nil || nm.Kind != NmInst {
		cc.RaiseCompileError(etag, "unknown inst name %s", op)
	}

	cc.CompileExpr(env, NewVec(append([]Value{op}, (*e)[2:]...)))
	insts := cc.CodeStack[len(cc.CodeStack)-1]
	insts[len(insts)-1].From = e
	return NIL
}

// SYNTAX: (apply fn ...)
func (cc *Compiler) sApply(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, -1, CtModule|CtProc, cc)
	f := CheckConst(e.At(1), IdentifierT, "callable", etag, cc)
	return cc.CompileExpr(env, NewVec(append([]Value{f}, (*e)[2:]...)))
}

// SYNTAX: (loop body)
func (cc *Compiler) sLoop(env *Env, e *Vec) Value {
	etag, n := CheckExpr(e, 2, -1, CtProc, cc)

	env = env.Enter()
	begNm := cc.InstallNamed(env, etag.Expand(KwUBEG), NmLabel, &Label{})
	endNm := cc.InstallNamed(env, etag.Expand(KwUEND), NmLabel, &Label{})
	condNm := cc.InstallNamed(env, etag.Expand(KwUCOND), NmLabel, &Label{})

	cc.EmitCode(NewInst(e, InstLabel, begNm))
	cc.CompileExpr(env, CheckBlockForm(e.At(1), "body", etag, cc))
	cc.EmitCode(NewInst(e, InstLabel, condNm))
	if n == 2 {
		a := etag.Expand(KwUBEG).ToConstexpr(nil)
		cc.CompileExpr(env, &Vec{etag.Expand(KwJump), a, NIL})
	} else {
		a := CheckConst(e.At(2), IdentifierT, "loop condition", etag, cc)
		cc.CompileExpr(env, NewVec(append([]Value{a}, (*e)[3:]...)))
	}
	return cc.EmitCode(NewInst(e, InstLabel, endNm))
}

func (cc *Compiler) sIfCond(env *Env, e *Vec) Value {
	etag, n := CheckExpr(e, 3, 5, CtProc, cc)
	cond := CheckValue(e.At(1), IdentifierT, "cond", etag, cc)
	thenBody := CheckBlockForm(e.At(2), "body", etag, cc)

	if !cc.IsCond(cond.Name) {
		cc.RaiseCompileError(cond, "invalid condition")
	}

	endThenId := Gensym("end-then").ToId(etag.Token)
	endThenNm := cc.InstallNamed(env, endThenId, NmLabel, &Label{})
	cc.CompileExpr(env, &Vec{
		etag.Expand(KwWith),
		endThenId.ToConstexpr(env),
		&Vec{etag.Expand(KwJumpUnless), cond, NIL},
	})
	cc.CompileExpr(env, thenBody)
	if n == 3 {
		return cc.EmitCode(NewInst(e, InstLabel, endThenNm))
	}

	elseTag := CheckConst(e.At(3), IdentifierT, "else", etag, cc)
	if KwElse.MatchId(elseTag) == nil {
		cc.RaiseCompileError(elseTag, "invalid if form, expected `else`")
	} else if n != 5 {
		cc.RaiseCompileError(elseTag, "invalid if form, else body required")
	}
	elseBody := CheckBlockForm(e.At(4), "body", etag, cc)

	endElseId := Gensym("end-else").ToId(etag.Token)
	endElseNm := cc.InstallNamed(env, endElseId, NmLabel, &Label{})
	cc.CompileExpr(env, &Vec{etag.Expand(KwJump), endElseId.ToConstexpr(env), NIL})
	cc.EmitCode(NewInst(e, InstLabel, endThenNm))
	cc.CompileExpr(env, elseBody)
	return cc.EmitCode(NewInst(e, InstLabel, endElseNm))
}

// SYNTAX: (if cond then `else` else)
func (cc *Compiler) sIf(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, -1, CtModule|CtProc, cc)

	if _, ok := e.At(1).(*Constexpr); !ok {
		return cc.sIfCond(env, e)
	}

	var matched Value
	body := CheckBlockForm(e.At(2), "then-body", etag, cc)
	if EvalConstAs(e.At(1), env, IntT, "cond", etag, cc) != Int(0) {
		matched = body
	}

	for v := (*e)[3:]; len(v) > 0; { // else {} ... | else if {} ...
		etag := CheckConst(v[0], IdentifierT, "else", etag, cc)
		if KwElse.MatchId(etag) == nil {
			cc.RaiseCompileError(etag, "invalid if form, expected `else`")
		} else if len(v) < 2 {
			cc.RaiseCompileError(etag, "invalid if form, else body required")
		}

		if elif, ok := v[1].(*Constexpr); ok && KwIf.MatchId(elif.Body) != nil {
			if len(v) < 4 {
				cc.RaiseCompileError(etag, "invalid else-if form")
			}

			cond := CheckValue(v[2], ConstexprT, "cond", etag, cc)
			body := CheckBlockForm(v[3], "then-body", etag, cc)
			if matched == nil && EvalConstAs(cond, env, IntT, "cond", etag, cc) != Int(0) {
				matched = body
			}
			v = v[4:]
			continue
		} else if len(v) == 2 {
			body := CheckBlockForm(v[1], "else-body", etag, cc)
			if matched == nil {
				matched = body
			}
			break
		}
		cc.RaiseCompileError(etag, "invalid if-else form")
	}

	if matched != nil { // ok
		return cc.CompileExpr(env, matched)
	}
	return NIL
}

// SYNTAX: (case value rest ...)
func (cc *Compiler) sCase(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 4, -1, CtModule|CtProc, cc)
	a := CheckAndEvalConst(e.At(1), env, "value", etag, cc)

	var matched Value
	for v := (*e)[2:]; len(v) > 0; {
		etag := CheckConst(v[0], IdentifierT, "when/else", etag, cc)
		if KwWhen.MatchId(etag) != nil {
			if len(v) < 3 {
				cc.RaiseCompileError(etag, "invalid case-when form")
			}

			b := CheckValue(v[1], ConstexprT, "pattern", etag, cc)
			body := CheckBlockForm(v[2], "when-body", etag, cc)
			if matched == nil && equalValue(a, EvalConst(b, env, etag, cc)) {
				matched = body
			}
			v = v[3:]
			continue
		} else if KwElse.MatchId(etag) != nil {
			if len(v) != 2 {
				cc.RaiseCompileError(etag, "invalid case-else form")
			}

			body := CheckBlockForm(v[1], "else-body", etag, cc)
			if matched == nil {
				matched = body
			}
			break
		}
		cc.RaiseCompileError(etag, "invalid case form. when/else required")
	}

	if matched != nil {
		return cc.CompileExpr(env, matched)
	}
	return NIL
}

// SYNTAX: (when cond then body else elbody)
func (cc *Compiler) sWhen(env *Env, e *Vec) Value {
	etag, n := CheckExpr(e, 4, 6, CtProc, cc)
	condBody := CheckBlockForm(e.At(1), "cond", etag, cc)

	thenTag := CheckConst(e.At(2), IdentifierT, "then", etag, cc)
	if KwThen.MatchId(thenTag) == nil {
		cc.RaiseCompileError(thenTag, "invalid when form, expected `then`")
	}
	thenBody := CheckBlockForm(e.At(3), "body", etag, cc)

	thenId := Gensym("then").ToId(etag.Token)
	thenNm := cc.InstallNamed(env, thenId, NmLabel, &Label{})
	endThenId := Gensym("end-then").ToId(etag.Token)
	endThenNm := cc.InstallNamed(env, endThenId, NmLabel, &Label{})

	{ // cond part
		env := env.Enter()
		and := SyntaxFn(func(cc *Compiler, env *Env, e *Vec) Value {
			etag, _ := CheckExpr(e, 2, 2, CtProc, cc)
			cond := CheckValue(e.At(1), IdentifierT, "cond", etag, cc)
			return cc.CompileExpr(env, &Vec{
				etag.Expand(KwWith),
				endThenId.ToConstexpr(env),
				&Vec{etag.Expand(KwJumpUnless), cond, NIL},
			})
		})
		cc.InstallNamed(env, etag.Expand(KwWhenAnd), NmSyntax, and)

		or := SyntaxFn(func(cc *Compiler, env *Env, e *Vec) Value {
			etag, _ := CheckExpr(e, 2, 2, CtProc, cc)
			cond := CheckValue(e.At(1), IdentifierT, "cond", etag, cc)
			return cc.CompileExpr(env, &Vec{
				etag.Expand(KwWith),
				thenId.ToConstexpr(env),
				&Vec{etag.Expand(KwJumpIf), cond, NIL},
			})
		})
		cc.InstallNamed(env, etag.Expand(KwWhenOr), NmSyntax, or)

		n := condBody.Size()
		if n == 1 {
			cc.RaiseCompileError(etag, "at least one condition required")
		} else if KwWhenAnd.MatchExpr(condBody.At(n-1)) == nil {
			cc.RaiseCompileError(etag, "the last condition must be an AND expression")
		}
		cc.CompileExpr(env, condBody)
	}

	cc.EmitCode(NewInst(e, InstLabel, thenNm))
	cc.CompileExpr(env, thenBody)
	if n == 4 {
		return cc.EmitCode(NewInst(e, InstLabel, endThenNm))
	}

	elseTag := CheckConst(e.At(4), IdentifierT, "else", etag, cc)
	if KwElse.MatchId(elseTag) == nil {
		cc.RaiseCompileError(elseTag, "invalid when form, expected `else`")
	} else if n != 6 {
		cc.RaiseCompileError(elseTag, "invalid when form, else body required")
	}
	elseBody := CheckBlockForm(e.At(5), "body", etag, cc)

	endElseId := Gensym("end-else").ToId(etag.Token)
	endElseNm := cc.InstallNamed(env, endElseId, NmLabel, &Label{})
	cc.CompileExpr(env, &Vec{etag.Expand(KwJump), endElseId.ToConstexpr(env), NIL})
	cc.EmitCode(NewInst(e, InstLabel, endThenNm))
	cc.CompileExpr(env, elseBody)
	return cc.EmitCode(NewInst(e, InstLabel, endElseNm))
}

// SYNTAX: (alias new old)
func (cc *Compiler) sAlias(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtModule|CtProc, cc)
	new := CheckConst(e.At(1), IdentifierT, "new name", etag, cc)
	old := CheckConst(e.At(2), IdentifierT, "old name", etag, cc)

	nm := cc.LookupNamed(env, old)
	if nm == nil {
		cc.RaiseCompileError(etag, "unknown name %s", old.String())
	} else if nm.Kind != NmMacro && nm.Kind != NmInline && nm.Kind != NmLabel {
		cc.RaiseCompileError(etag, "aliases are not allowed for this type(%s)", NamedKindLabels[nm.Kind])
	}
	cc.InstallNamed(env, new, nm.Kind, nm.Value)

	return NIL
}

// SYNTAX: (#.module name body)
func (cc *Compiler) sModule(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtModule|CtProc, cc)
	name := CheckConst(e.At(1), IdentifierT, "module name", etag, cc)
	body := CheckBlockForm(e.At(2), "module body", etag, cc)

	CheckToplevelEnv(env, etag, cc)
	module, section := cc.Module, cc.Section
	cc.Module = NewModule(name.Name, nil)
	cc.Section = cc.Module.Sections[KwTEXT]
	cc.InstallNamed(env, name, NmModule, cc.Module)

	cc.EnterContext(CtModule)
	cc.EnterCodeBlock()
	cc.Module.Env = env.Enter()
	cc.Module.Env.Module = cc.Module
	cc.sBlock(cc.Module.Env, body)
	cc.EmitCodeToSection(cc.Section, cc.LeaveCodeBlock()...)
	cc.LeaveContext()
	cc.Module, cc.Section = module, section
	return NIL
}

// SYNTAX: (section name)
func (cc *Compiler) sSection(env *Env, e *Vec) Value {
	etag, n := CheckExpr(e, 2, 3, CtModule|CtProc, cc)
	id := CheckConstPlainId(e.At(1), "section name", etag, cc)

	CheckToplevelEnvIfCtProc(env, etag, cc)
	if cc.Section.Name == id.Name {
		return NIL
	}

	section := cc.Section
	cc.EmitCodeToSection(section, cc.LeaveCodeBlock()...)
	cc.EnterCodeBlock()
	cc.Section = cc.Module.FindOrNewSection(id.Name)
	if n == 3 {
		v := CheckBlockForm(e.At(2), "body", etag, cc)
		cc.CompileExpr(env, v)
		cc.CompileExpr(env, &Vec{etag, etag.Expand(section.Name).ToConstexpr(env)})
	}
	return NIL
}

// SYNTAX: (link ...)
func (cc *Compiler) sLink(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtModule|CtProc, cc)
	body := CheckBlockForm(e.At(1), "link body", etag, cc)

	CheckToplevelEnv(env, etag, cc)
	if cc.link != nil {
		cc.RaiseCompileError(etag, "link is already registered")
	}
	cc.link = body
	return NIL
}

// SYNTAX: (flat! ...)
func (cc *Compiler) sFlatMode(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 1, 1, CtModule, cc)

	CheckToplevelEnv(env, etag, cc)
	cc.Section = cc.Module.Sections[KwTEXT]
	cc.EnterContext(CtProc)
	cc.EnterCodeBlock()

	hook := cc.hooks.beforeLink
	cc.hooks.beforeLink = func(cc *Compiler) {
		cc.EmitCodeToSection(cc.Section, cc.LeaveCodeBlock()...)
		cc.LeaveContext()
		hook(cc)
	}
	return NIL
}

// SYNTAX: (pragma ...)
func (cc *Compiler) sPragma(env *Env, e *Vec) Value {
	etag, n := CheckExpr(e, 2, -1, 0, cc)
	k := CheckConstPlainId(e.At(1), "pragma", etag, cc)

	switch k.Name {
	case KwListConstants:
		CheckExpr(e, 3, 4, CtModule, cc)

		if n == 4 {
			v := CheckAndEvalConstAs(e.At(3), env, StrT, "message", etag, cc)
			cc.EmitCodeTo(0, NewInst(e, InstMisc, KwComment, v))
		}
		v := CheckAndEvalConstAs(e.At(2), env, IntT, "value", etag, cc)
		cc.EmitCodeTo(0, NewInst(e, InstMisc, KwListConstants, v))
	case KwComment:
		CheckExpr(e, 3, -1, CtModule|CtProc, cc)
		v := CheckAndEvalConstAs(e.At(2), env, StrT, "comment", etag, cc)
		r := []Value{KwComment, v}
		for _, i := range (*e)[3:] {
			r = append(r, cc.CompileExpr(env, i))
		}
		cc.EmitCode(NewInst(e, InstMisc, r...))
	default:
		cc.RaiseCompileError(etag, "unknown pragma: %s", k.String())
	}
	return NIL
}

// SYNTAX: (#.macro name args vars rest body)
func (cc *Compiler) sMacro(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 6, 6, CtModule|CtProc, cc)
	name := CheckConst(e.At(1), IdentifierT, "macro name", etag, cc)
	as := CheckValue(e.At(2), VecT, "macro parameters", etag, cc)
	vs := CheckValue(e.At(3), VecT, "macro variables", etag, cc)
	rest := e.At(4) != NIL
	body := CheckBlockForm(e.At(5), "macro body", etag, cc)

	if as.Size() == 0 && rest {
		cc.RaiseCompileError(etag, "rest parameter name required")
	}

	opts := []Value{}
	args := []*Keyword{}
	for x, a := range *as {
		a := a.(*Vec)
		k := CheckConstPlainId(a.At(0), "macro parameter", etag, cc)
		v := a.At(1)
		if rest && x == len(*as)-1 {
			if v != NIL {
				cc.RaiseCompileError(k, "the rest parameter cannot have default value")
			}
			opts = append(opts, v)
		} else if len(opts) > 0 && v == NIL {
			cc.RaiseCompileError(k, "default value required")
		}
		if v != NIL {
			opts = append(opts, v)
		}
		args = append(args, k.Name)
	}

	vars := []Value{}
	for x, n := 0, vs.Size(); x < n; x += 2 {
		k := CheckConstPlainId((*vs)[x], "macro variable", etag, cc)
		v := (*vs)[x+1]
		if v != NIL {
			CheckValue(v, ConstexprT, "macro variable body", etag, cc)
		}
		vars = append(vars, k.Name, v)
	}

	value := &Macro{Args: args, Vars: vars, Opts: opts, Rest: rest, Body: body}
	cc.InstallNamed(env, name, NmMacro, value)
	return NIL
}

// SYNTAX: (#.proc name sig body)
func (cc *Compiler) sProc(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 4, 4, CtModule|CtProc, cc)
	name := CheckConst(e.At(1), IdentifierT, "proc name", etag, cc)
	sigv := CheckValue(e.At(2), VecT, "signature", etag, cc)
	sig := cc.newSig(sigv)

	if addr, ok := e.At(3).(*Constexpr); ok {
		if sig.IsInline {
			cc.RaiseCompileError(etag, "cannot bind inline proc")
		}
		addr = cc.CompileExpr(env, addr).(*Constexpr)
		nm := cc.InstallNamed(env, name, NmLabel, &Label{Sig: sig, At: addr})
		return cc.EmitCodeByEnv(env, NewInst(e, InstBind, nm))
	}

	body := CheckBlockForm(e.At(3), "proc body", etag, cc)
	if sig.IsInline {
		cc.InstallNamed(env, name, NmInline, &Inline{Body: body, Sig: sig})
		return NIL
	}

	cc.EnterContext(CtProc)
	cc.EnterCodeBlock()

	nm := cc.InstallNamed(env, name, NmLabel, &Label{Sig: sig})
	cc.EmitCode(NewInst(e, InstLabel, nm, NIL))

	env = env.Enter()
	procNm := cc.InstallNamed(env, etag.Expand(KwPROCNAME), NmConst, name.ToConstexpr(env))
	cc.EmitCode(NewInst(e, InstConst, procNm))

	cc.sBlock(env, body)
	cc.EmitCode(NewInst(e, InstMisc, KwEndProc))
	cc.EmitCode(cc.LeaveCodeBlock()...)
	cc.LeaveContext()
	return NIL
}

// SYNTAX: (#.callproc name cond sig op)
func (cc *Compiler) sCallproc(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 5, 5, CtProc, cc)
	id := CheckConst(e.At(1), IdentifierT, "callee name", etag, cc)
	sigv := CheckValue(e.At(3), VecT, "signature", etag, cc)
	op := CheckValue(e.At(4), KeywordT, "opcode", etag, cc) // for tco

	name := cc.CompileExpr(env, e.At(1)) // bind env
	cond := e.At(2)
	sig := cc.newSig(sigv)

	if sig.IsInline {
		inst := NewInst(e, InstMisc, KwInline, env, id, NIL)
		cc.InlineInsts = append(cc.InlineInsts, inst)
		cc.EmitCode(inst)
	} else {
		cc.CompileExpr(env, &Vec{etag.Expand(op), name, cond})
	}
	return cc.EmitCode(NewInst(e, InstMisc, KwCallproc, name, sig))
}

// SYNTAX: (fallthrough ...)
func (cc *Compiler) sFallthrough(env *Env, e *Vec) Value {
	CheckExpr(e, 1, 1, CtProc, cc)
	return cc.EmitCode(NewInst(e, InstMisc, KwFallthrough))
}

// SYNTAX: (tco a)
func (cc *Compiler) sTco(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtProc, cc)
	body := CheckBlockForm(e.At(1), "body", etag, cc)

	if body.Size() != 2 {
		cc.RaiseCompileError(etag, "only one statement is allowed within the tco form")
	}

	call := KwCallproc.MatchExpr(body.At(1))
	if call == nil {
		cc.RaiseCompileError(etag, "proc call required")
	}
	if call.At(2) != NIL {
		cc.RaiseCompileError(etag, "the conditional call is not a tail call")
	}
	call.SetAt(4, KwJump) // original is KwCall
	return cc.CompileExpr(env, call)
}

// SYNTAX: (#.const name args value)
func (cc *Compiler) sConst(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 4, 4, CtModule|CtProc, cc)
	name := CheckConst(e.At(1), IdentifierT, "constant name", etag, cc)
	value := CheckValue(cc.CompileExpr(env, e.At(3)), ConstexprT, "constant value", etag, cc)

	as := CheckValue(e.At(2), VecT, "function parameters", etag, cc)
	if as.Size() > 0 {
		opts := []Value{}
		args := []*Keyword{}
		for _, a := range *as {
			a := a.(*Vec)
			k := CheckConstPlainId(a.At(0), "function parameter", etag, cc)

			v := a.At(1)
			if len(opts) > 0 && v == NIL {
				cc.RaiseCompileError(k, "default value required")
			}
			if v != NIL {
				opts = append(opts, cc.CompileExpr(env, v))
			}
			args = append(args, k.Name)
		}
		cc.InstallNamed(env, name, NmFunc, NewConstFn(args, opts, value))
		return NIL
	}

	nm := cc.InstallNamed(env, name, NmConst, value)
	return cc.EmitCodeByEnv(env, NewInst(e, InstConst, nm))
}

// SYNTAX: (#.data name type (values op count section addr))
func (cc *Compiler) sData(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 4, 4, CtModule|CtProc, cc)

	var nm *Named
	info := &Label{}

	cc.EnterCodeBlock()
	if name := e.At(1); name != NIL {
		name := CheckConst(name, IdentifierT, "data name", etag, cc)
		nm = cc.InstallNamed(env, name, NmLabel, info)
		cc.EmitCode(NewInst(e, InstLabel, nm))
	}

	t := CheckConstPlainId(e.At(2), "data type", etag, cc).Name
	if t != KwByte && t != KwWord {
		cc.RaiseCompileError(etag, "invalid data type %s", t.String())
	}

	body := CheckValue(e.At(3), VecT, "data body", etag, cc)
	// op := CheckConstAs(body.At(1), env, IdentifierT, "op", etag, cc)
	n := CheckAndEvalConstAs(body.At(2), env, IntT, "count", etag, cc)
	if n < 1 {
		cc.RaiseCompileError(etag, "invalid repeat count %d", n)
	}

	switch values := body.At(0).(type) {
	case *Constexpr:
		info.Link = NewInst(e, InstBlob, EvalConstAs(values, env, BlobT, "blob", etag, cc))
	case *Vec:
		info.Link = NewInst(e, InstData, t, n, cc.sVec(env, values))
	default:
		info.Link = NewInst(e, InstDS, t, n)
	}
	cc.EmitCode(info.Link)

	if body.At(3) != NIL {
		name := CheckConstPlainId(body.At(3), "section", etag, cc)
		section := cc.Module.FindOrNewSection(name.Name)
		cc.EmitCodeToSection(section, cc.LeaveCodeBlock()...)
	} else if body.At(4) != NIL {
		if nm == nil {
			cc.RaiseCompileError(etag, "data name required")
		}
		if info.Link.Kind != InstDS {
			cc.RaiseCompileError(etag, "the addressed data cannot contain any elements")
		}

		addr := CheckValue(body.At(4), ConstexprT, "data address", etag, cc)
		info.At = cc.CompileExpr(env, addr).(*Constexpr)
		cc.LeaveCodeBlock()
		cc.EmitCodeByEnv(env, NewInst(e, InstBind, nm))
	} else {
		cc.EmitCode(cc.LeaveCodeBlock()...)
	}

	return e
}

// SYNTAX: (#.vec ...)
func (cc *Compiler) sVec(env *Env, e *Vec) Value {
	CheckExpr(e, 1, -1, CtModule|CtProc, cc)

	r := &Vec{}
	for _, i := range (*e)[1:] {
		i := cc.CompileExpr(env, i)
		if v, ok := i.(*Vec); ok {
			r.Push([]Value(*v)...)
		} else {
			r.Push(i)
		}
	}
	return r
}

// SYNTAX: (#.BYTE n)
func (cc *Compiler) sByte(env *Env, e *Vec) Value {
	CheckExpr(e, 2, 2, CtProc, cc)
	a := cc.CompileExpr(env, e.At(1))
	if v, ok := a.(*Operand); ok {
		a = v.A0
	}
	return cc.EmitCode(NewInst(e, InstData, KwByte, Int(1), &Vec{a}))
}

// SYNTAX: (#.REP n body)
func (cc *Compiler) sRep(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtProc, cc)
	a := cc.ExprToOperand(cc, e.At(1))
	n := CheckAndEvalConstAs(a.A0, env, IntT, "loop counter", etag, cc)

	for range int(n) {
		cc.CompileExpr(env, e.At(2).Dup())
	}
	return NIL
}

// SYNTAX: (#.INVALID ...)
func (cc *Compiler) sInvalid(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtProc, cc)
	cc.RaiseCompileError(etag, "invalid operands")
	return NIL
}

// SYNTAX: (#.mem a)
func (cc *Compiler) sMem(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, -1, CtProc, cc)

	r := &Vec{etag}
	for _, v := range (*e)[1:] {
		r.Push(cc.ExprToOperand(cc, cc.CompileExpr(env, v)))
	}
	return cc.ExprToOperand(cc, r)
}

// SYNTAX: (#.valueof a)
func (cc *Compiler) sValueOf(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtModule|CtProc, cc)
	a := CheckValue(e.At(1), ConstexprT, "value", etag, cc)
	v := EvalConst(a, env, etag, cc)

	switch v.(type) {
	case Int, *Str, *Blob:
		return &Constexpr{Token: a.Token, Body: v, Env: env}
	}
	return v
}

// SYNTAX: (#.field a b)
func (cc *Compiler) sField(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 1, -1, CtConstexpr, cc)
	cc.RaiseCompileError(etag, "<field> form not implemented yet")
	return NIL
}

func expandOperatorTemplate(op *Identifier, arg []*Operand, e *Vec) Value {
	r := &Vec{}
	for _, i := range *e {
		switch i := i.(type) {
		case Int:
			r.Push(&Constexpr{Token: op.Token, Body: i})
		case *Keyword:
			r.Push(op.Expand(i))
		case *Operand:
			r.Push(i.Dup())
		case *Vec:
			switch v := i.At(0).(type) {
			case Int:
				a, b := arg[v], i.At(1)
				if b != nil {
					a.Kind = b.(*Keyword)
				}
				r.Push(a)
			case *Vec:
				s := &Vec{op.Expand(KwBlock)}
				for _, i := range *i {
					s.Push(expandOperatorTemplate(op, arg, i.(*Vec)))
				}
				r.Push(s)
			}
		}
	}
	return r
}

// SYNTAX: (#.with context &rest ops)
func (cc *Compiler) sWith(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, -1, CtProc, cc)
	t := cc.ExprToOperand(cc, cc.CompileExpr(env, e.At(1)))

	for _, i := range (*e)[2:] {
		i := i.(*Vec)
		op := CheckValue(i.At(0), IdentifierT, "operator", etag, cc)

		u := &Operand{Kind: nil}
		if op.Name == KwDot {
			cc.CompileExpr(env, i.At(1))
			continue
		} else if i.At(1) != NIL {
			u = cc.ExprToOperand(cc, cc.CompileExpr(env, i.At(1)))
		}

		m1 := cc.CtxOpMap[op.Name]
		if m1 == nil {
			cc.RaiseCompileError(op, "unknown operator: %s", op)
		}

		m2 := m1[t.Kind]
		if m2 == nil {
			m2 = m1[KwAny]
		}
		if m2 == nil {
			cc.RaiseCompileError(op, "cannot use %s as first operand for '%s'", t.Kind, op)
		}

		body := m2[u.Kind]
		if body == nil {
			body = m2[KwAny]
		}
		if body == nil {
			if u == nil {
				cc.RaiseCompileError(op, "[BUG] %s require second operand", op)
			} else {
				cc.RaiseCompileError(op, "cannot use %s as second operand for '%s'", u.Kind, op)
			}
		}

		for _, i := range body {
			i := expandOperatorTemplate(op, []*Operand{t, u}, NewVec(i))
			cc.CompileExpr(env, i)
		}
	}
	return t
}

// SYNTAX: (compile-error s)
func (cc *Compiler) sCompileError(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtModule|CtProc, cc)
	s := CheckConst(e.At(1), StrT, "error message", etag, cc)
	if etag.ExpandedBy != nil {
		etag = etag.ExpandedBy
	}
	cc.RaiseCompileError(etag, s.String())
	return NIL
}

// SYNTAX: (assert expr message)
func (cc *Compiler) sAssert(env *Env, e *Vec) Value {
	etag, n := CheckExpr(e, 2, 3, CtModule|CtProc, cc)
	cond := CheckValue(e.At(1), ConstexprT, "cond", etag, cc)

	s := Value(NewStr(fmt.Sprintf("assertion %s failed", cond.Inspect())))
	if n > 2 {
		s = CheckConst(e.At(2), StrT, "message", etag, cc)
	}

	cc.EmitCode(NewInst(e, InstAssert, cc.CompileExpr(env, cond), s))
	return NIL
}

// SYNTAX: (import s)
func (cc *Compiler) sImport(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtModule|CtProc, cc)
	id := CheckConst(e.At(1), IdentifierT, "module name", etag, cc)

	mod := cc.FindModule(env, id.Name, id)
	if mod == nil {
		cc.RaiseCompileError(id, "unknown namespace %s", id.String())
	}
	env.InsertEnv(mod.Env)
	return NIL
}

// SYNTAX: (expand-loop n a)
func (cc *Compiler) sExpandLoop(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtProc, cc)
	n := CheckAndEvalConstAs(e.At(1), env, IntT, "loop counter", etag, cc)
	body := CheckBlockForm(e.At(2), "loop body", etag, cc)

	for range int(n) {
		cc.CompileExpr(env, body.Dup())
	}
	return NIL
}

// SYNTAX: (*patch* name delta)
func (cc *Compiler) sPatch(env *Env, e *Vec) Value {
	etag, n := CheckExpr(e, 2, 3, CtProc, cc)
	name := CheckConst(e.At(1), IdentifierT, "name", etag, cc)

	nm := cc.LookupNamed(env, name)
	CheckReserved(nm, name, cc)

	delta := Value(Int(0))
	if n == 3 {
		delta = CheckValue(e.At(2), ConstexprT, "delta", etag, cc).Body
		if KwByte.MatchId(delta) != nil {
			delta = Int(-1)
		} else if KwWord.MatchId(delta) != nil {
			delta = Int(-2)
		}
	}

	id := Gensym("patch." + name.String()).ToId(etag.Token)
	label := cc.InstallNamed(env, id, NmLabel, &Label{})

	nm.Value.(*Label).At = &Constexpr{Env: env, Body: &Vec{etag.Expand(KwPlusOp), id, delta}}
	return cc.EmitCode(NewInst(e, InstLabel, label))
}

// SYNTAX: (make-counter name value)
func (cc *Compiler) sMakeCounter(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtModule|CtProc, cc)
	name := CheckConst(e.At(1), IdentifierT, "name", etag, cc)
	value := CheckAndEvalConstAs(e.At(2), env, IntT, "start", etag, cc)

	getter := SyntaxFn(func(cc *Compiler, env *Env, e *Vec) Value {
		v := value
		value += 1
		return v
	})
	cc.InstallNamed(env, name, NmSpecial, getter)

	fname := Gensym("make-counter").ToId(etag.Token)
	resetter := SyntaxFn(func(cc *Compiler, env *Env, e *Vec) Value {
		value = e.At(2).(Int)
		return value
	})
	cc.InstallNamed(env, fname, NmFunc, resetter)

	lhs := Gensym("counter").ToId(etag.Token)
	rhs := &Constexpr{Token: etag.Token, Body: &Vec{fname, name, value}, Env: env}
	nm := cc.InstallNamed(env, lhs, NmConst, rhs)
	return cc.EmitCodeByEnv(env, NewInst(e, InstConst, nm))
}

// SYNTAX: (debug-inspect a)
func (cc *Compiler) sDebugInspect(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtModule|CtProc, cc)
	a := e.At(1)
	if _, ok := a.(*Constexpr); ok {
		a = EvalConst(a, env, etag, cc)
	}
	fmt.Fprintln(cc.g.ErrWriter, "[DEBUG]", a.Inspect())
	return NIL
}

////////////////////////////////////////////////////////////

// SYNTAX: (quote a)
func (cc *Compiler) sQuote(env *Env, e *Vec) Value {
	CheckExpr(e, 2, 2, CtConstexpr, cc)
	return e.At(1)
}

// SYNTAX: (#.exprdata a)
func (cc *Compiler) sExprdata(env *Env, e *Vec) Value {
	CheckExpr(e, 2, 2, CtConstexpr, cc)
	return e
}

// SYNTAX: (#.invalid-expansion ...)
func (cc *Compiler) sInvalidExpansion(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 1, 1, CtConstexpr, cc)
	cc.RaiseCompileError(etag, "invalid expansion")
	return NIL
}

// SYNTAX: (&& a b)
func (cc *Compiler) sAnd(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtConstexpr, cc)
	a := EvalConst(e.At(1), env, etag, cc)

	if a == Int(0) {
		return a
	}
	return EvalConst(e.At(2), env, etag, cc)
}

// SYNTAX: (|| a b)
func (cc *Compiler) sOr(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtConstexpr, cc)
	a := EvalConst(e.At(1), env, etag, cc)

	if a != Int(0) {
		return a
	}
	return EvalConst(e.At(2), env, etag, cc)
}

// SYNTAX: (sizeof a)
func (cc *Compiler) sSizeof(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtConstexpr, cc)
	name := CheckValue(e.At(1), IdentifierT, "label", etag, cc)

	nm := cc.LookupNamed(env, name)
	if nm == nil || nm.Kind != NmLabel {
		cc.RaiseCompileError(etag, "unknown label name %s", name.String())
	}

	v := nm.Value.(*Label)
	if !v.LinkedToData() {
		cc.RaiseCompileError(etag, "%s is not a data", name.String())
	}
	return Int(v.Link.size)
}

// SYNTAX: (nameof a)
func (cc *Compiler) sNameof(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtConstexpr, cc)
	name := CheckValue(e.At(1), IdentifierT, "name", etag, cc)
	return NewStr(name.String())
}

// SYNTAX: (nametypeof a)
func (cc *Compiler) sNametypeof(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtConstexpr, cc)
	name := CheckValue(e.At(1), IdentifierT, "name", etag, cc)

	nm := cc.LookupNamed(env, name)
	if nm == nil {
		cc.RaiseCompileError(etag, "unknown name %s", name.String())
	}
	return NewStr(NamedKindLabels[nm.Kind])
}

// SYNTAX: (exprtypeof a)
func (cc *Compiler) sExprtypeOf(env *Env, e *Vec) Value {
	CheckExpr(e, 2, 2, CtConstexpr, cc)

	switch e.At(1).(type) {
	case Int:
		return NewStr("int")
	case *Identifier:
		return NewStr("identifier")
	case *Str:
		return NewStr("str")
	}
	return NewStr("unknown")
}

// SYNTAX: (defined? a)
func (cc *Compiler) sDefinedp(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtConstexpr, cc)
	name := CheckValue(e.At(1), IdentifierT, "name", etag, cc)
	return BoolInt(cc.LookupNamed(env, name) != nil)
}

////////////////////////////////////////////////////////////

// FUN: (#.make-id ...)
func (cc *Compiler) fMakeId(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 1, -1, CtConstexpr, cc)
	s := ""
	for _, i := range (*e)[1:] {
		switch i := i.(type) {
		case *Identifier:
			s += string(*i.Name)
		case *Str:
			s += string(*i)
		case *Constexpr:
			switch i := i.Body.(type) {
			case *Identifier:
				s += string(*i.Name)
			case *Str:
				s += string(*i)
			default:
				cc.RaiseCompileError(etag, "invalid fragment")
			}
		default:
			cc.RaiseCompileError(etag, "[BUG] invalid fragment")
		}
	}
	return Intern(s).ToId(etag.Token)
}

// FUN: (* a b)
func (cc *Compiler) fMul(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "left value", etag, cc)
	b := CheckValue(e.At(2), IntT, "right value", etag, cc)

	return a * b
}

// FUN: (/ a b)
func (cc *Compiler) fDiv(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "left value", etag, cc)
	b := CheckValue(e.At(2), IntT, "right value", etag, cc)

	return a / b
}

// FUN: (% a b)
func (cc *Compiler) fMod(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "left value", etag, cc)
	b := CheckValue(e.At(2), IntT, "right value", etag, cc)

	return a % b
}

// FUN: (+ a b)
func (cc *Compiler) fAdd(env *Env, e *Vec) Value {
	etag, n := CheckExpr(e, 2, 3, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "left value", etag, cc)

	if n == 2 {
		return a
	}

	b := CheckValue(e.At(2), IntT, "right value", etag, cc)
	return a + b
}

// FUN: (- a b)
func (cc *Compiler) fSub(env *Env, e *Vec) Value {
	etag, n := CheckExpr(e, 2, 3, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "left value", etag, cc)

	if n == 2 {
		return -a
	}

	b := CheckValue(e.At(2), IntT, "right value", etag, cc)
	return a - b
}

// FUN: (<< a b)
func (cc *Compiler) fLsl(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "left value", etag, cc)
	b := CheckValue(e.At(2), IntT, "right value", etag, cc)

	return a << b
}

// FUN: (>> a b)
func (cc *Compiler) fAsr(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "left value", etag, cc)
	b := CheckValue(e.At(2), IntT, "right value", etag, cc)

	return a >> b
}

// FUN: (>>> a b)
func (cc *Compiler) fLsr(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "left value", etag, cc)
	b := CheckValue(e.At(2), IntT, "right value", etag, cc)

	return Int(uint64(a) >> uint64(b))
}

// FUN: (< a b)
func (cc *Compiler) fLt(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "left value", etag, cc)
	b := CheckValue(e.At(2), IntT, "right value", etag, cc)

	return BoolInt(a < b)
}

// FUN: (<= a b)
func (cc *Compiler) fLe(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "left value", etag, cc)
	b := CheckValue(e.At(2), IntT, "right value", etag, cc)

	return BoolInt(a <= b)
}

// FUN: (> a b)
func (cc *Compiler) fGt(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "left value", etag, cc)
	b := CheckValue(e.At(2), IntT, "right value", etag, cc)

	return BoolInt(a > b)
}

// FUN: (>= a b)
func (cc *Compiler) fGe(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "left value", etag, cc)
	b := CheckValue(e.At(2), IntT, "right value", etag, cc)

	return BoolInt(a >= b)
}

func equalValue(a, b Value) bool {
	if s, ok := a.(*Str); ok {
		if t, ok := b.(*Str); ok {
			return *s == *t
		}
	}
	return a == b
}

// FUN: (== a b)
func (cc *Compiler) fEql(env *Env, e *Vec) Value {
	CheckExpr(e, 3, 3, CtConstexpr, cc)
	a := e.At(1)
	b := e.At(2)

	return BoolInt(equalValue(a, b))
}

// FUN: (!= a b)
func (cc *Compiler) fNotEql(env *Env, e *Vec) Value {
	return BoolInt(cc.fEql(env, e) == Int(0))
}

// FUN: (& a b)
func (cc *Compiler) fAnd(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "left value", etag, cc)
	b := CheckValue(e.At(2), IntT, "right value", etag, cc)

	return a & b
}

// FUN: (| a b)
func (cc *Compiler) fOr(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "left value", etag, cc)
	b := CheckValue(e.At(2), IntT, "right value", etag, cc)

	return a | b
}

// FUN: (^ a b)
func (cc *Compiler) fXor(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "left value", etag, cc)
	b := CheckValue(e.At(2), IntT, "right value", etag, cc)

	return a ^ b
}

// FUN: (~ a)
func (cc *Compiler) fNot(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "left value", etag, cc)

	return ^a
}

// FUN: (! a)
func (cc *Compiler) fLogicalNot(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "left value", etag, cc)

	return BoolInt(a == Int(0))
}

// FUN: (byte a)
func (cc *Compiler) fByte(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "value", etag, cc)
	return Int(byte(a))
}

// FUN: (word a)
func (cc *Compiler) fWord(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "value", etag, cc)
	return Int(uint16(a))
}

// FUN: (lobyte a)
func (cc *Compiler) fLobyte(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "value", etag, cc)
	return Int(a & 0xff)
}

// FUN: (hibyte a)
func (cc *Compiler) fHibyte(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "value", etag, cc)
	return Int((a >> 8) & 0xff)
}

// FUN: (asword h l)
func (cc *Compiler) fAsWord(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtConstexpr, cc)
	h := CheckValue(e.At(1), IntT, "high byte", etag, cc)
	l := CheckValue(e.At(2), IntT, "low byte", etag, cc)
	return Int(((h & 0xff) << 8) | (l & 0xff))
}

// FUN: (unuse? a)
func (cc *Compiler) fUnusep(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtConstexpr, cc)
	form := cc.unwrapExprdata(e.At(1), etag)

	c, ok := form.(*Constexpr)
	return BoolInt(ok && KwUNDER.MatchId(c.Body) != nil)
}

// FUN: (use? a)
func (cc *Compiler) fUsep(env *Env, e *Vec) Value {
	return BoolInt(cc.fUnusep(env, e) == Int(0))
}

// FUN: (formtypeof a)
func (cc *Compiler) fFormtypeof(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtConstexpr, cc)
	form := cc.unwrapExprdata(e.At(1), etag)

	switch e := form.(type) {
	case *Constexpr:
		return NewStr("constexpr")
	case *Identifier:
		kw := e.Name
		if cc.IsReg(kw) {
			return NewStr("reg")
		} else if cc.IsCond(kw) {
			return NewStr("cond")
		}
	case *Vec:
		switch e.ExprTagName() {
		case KwBlock, KwProg:
			return NewStr("block-form")
		case KwMem:
			return NewStr("mem-form")
		}
	}
	return NewStr("unknown")
}

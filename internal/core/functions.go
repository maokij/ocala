package core

import (
	"fmt"
	"slices"
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
	cc.ErrorAt(etag).With("reserved but undefined")
	return NIL
}

func sCannotEval(cc *Compiler, env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 1, 1, CtConstexpr, cc)
	cc.ErrorAt(etag).With("[BUG] invalid evaluation")
	return NIL
}

// SPECIAL: (__FILE__)
// func (cc *Compiler) sFilename(env *Env, e *Vec) Value

// SYNTAX: (include path)
// func (cc *Compiler) sInclude(env *Env, e *Vec) Value

// SYNTAX: (load-file path)
// func (cc *Compiler) sLoadFile(env *Env, e *Vec) Value

// SYNTAX: (arch name variant)
func (cc *Compiler) sArch(env *Env, e *Vec) Value {
	etag, n := CheckExpr(e, 2, 3, CtModule|CtProc, cc)
	arch := CheckConst(e.At(1), IdentifierT, "arch name", etag, cc)

	CheckToplevelEnv(env, etag, cc)
	if cc.Arch == arch.String() {
		if n == 2 {
			return NIL // ok
		}

		variant := CheckConst(e.At(2), IdentifierT, "variant name", etag, cc)
		if cc.Variant == variant.String() {
			return NIL // ok
		}
	}

	cc.ErrorAt(arch, etag).With("the current target arch is %s", cc.FullArchName())
	return NIL
}

// SYNTAX: (align a)
func (cc *Compiler) sAlign(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtModule|CtProc, cc)
	n := CheckAndEvalConstAs(e.At(1), env, IntT, "argument", etag, cc)
	if n < 1 || (n&(n-1)) != 0 {
		cc.ErrorAt(etag).With("the alignment size must be power of 2")
	}
	return cc.EmitCode(NewInst(e, InstAlign, n))
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
			cc.ErrorAt(i, etag).With("invalid do form(%T)", i)
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
		cc.CompileExpr(env, &Vec{etag.Expand(KwJump), a})
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
		cc.ErrorAt(cond, etag).With("invalid condition")
	}

	endThenId := etag.Expand(Gensym("end-then"))
	endThenNm := cc.InstallNamed(env, endThenId, NmLabel, &Label{})
	cc.CompileExpr(env, &Vec{
		etag.Expand(KwWith),
		endThenId.ToConstexpr(env),
		&Vec{etag.Expand(KwJumpUnlessOp), cond},
	})
	cc.CompileExpr(env, thenBody)
	if n == 3 {
		return cc.EmitCode(NewInst(e, InstLabel, endThenNm))
	}

	elseTag := CheckConst(e.At(3), IdentifierT, "else", etag, cc)
	if KwElse.MatchId(elseTag) == nil {
		cc.ErrorAt(elseTag, etag).With("invalid if form, expected `else`")
	} else if n != 5 {
		cc.ErrorAt(elseTag, etag).With("invalid if form, else body required")
	}
	elseBody := CheckBlockForm(e.At(4), "body", etag, cc)

	endElseId := etag.Expand(Gensym("end-else"))
	endElseNm := cc.InstallNamed(env, endElseId, NmLabel, &Label{})
	cc.CompileExpr(env, &Vec{etag.Expand(KwJump), endElseId.ToConstexpr(env)})
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
		elseTag := CheckConst(v[0], IdentifierT, "else", etag, cc)
		if KwElse.MatchId(elseTag) == nil {
			cc.ErrorAt(elseTag, etag).With("invalid if form, expected `else`")
		} else if len(v) < 2 {
			cc.ErrorAt(elseTag, etag).With("invalid if form, else body required")
		} else if elseIfTag := KwIf.MatchConstId(v[1]); elseIfTag != nil {
			if len(v) < 4 {
				cc.ErrorAt(elseIfTag, etag).With("invalid else-if form")
			}
			cond := CheckValue(v[2], ConstexprT, "cond", elseIfTag, cc)
			body := CheckBlockForm(v[3], "then-body", elseIfTag, cc)
			if matched == nil && EvalConstAs(cond, env, IntT, "cond", elseIfTag, cc) != Int(0) {
				matched = body
			}
			v = v[4:]
		} else if len(v) == 2 {
			body := CheckBlockForm(v[1], "else-body", etag, cc)
			if matched == nil {
				matched = body
			}
			v = v[:0]
		} else {
			cc.ErrorAt(elseTag, etag).With("invalid if-else form")
		}
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
		clauseTag := CheckConstPlainId(v[0], "clause tag", etag, cc)
		switch clauseTag.Name {
		case KwWhen:
			if len(v) < 3 {
				cc.ErrorAt(clauseTag, etag).With("invalid case-when form")
			}

			b := CheckValue(v[1], ConstexprT, "pattern", clauseTag, cc)
			body := CheckBlockForm(v[2], "when-body", clauseTag, cc)
			if matched == nil && equalValue(a, EvalConst(b, env, clauseTag, cc)) {
				matched = body
			}
			v = v[3:]
		case KwElse:
			if len(v) != 2 {
				cc.ErrorAt(clauseTag, etag).With("invalid case-else form")
			}

			body := CheckBlockForm(v[1], "else-body", etag, cc)
			if matched == nil {
				matched = body
			}
			v = v[:0]
		default:
			cc.ErrorAt(clauseTag, etag).With("invalid case form. when/else required")
		}
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
		cc.ErrorAt(thenTag, etag).With("invalid when form, expected `then`")
	}
	thenBody := CheckBlockForm(e.At(3), "then-body", etag, cc)

	thenId := etag.Expand(Gensym("then"))
	thenNm := cc.InstallNamed(env, thenId, NmLabel, &Label{})
	endThenId := etag.Expand(Gensym("end-then"))
	endThenNm := cc.InstallNamed(env, endThenId, NmLabel, &Label{})

	{ // cond part
		env := env.Enter()
		and := SyntaxFn(func(cc *Compiler, env *Env, e *Vec) Value {
			etag, _ := CheckExpr(e, 2, 2, CtProc, cc)
			cond := CheckValue(e.At(1), IdentifierT, "cond", etag, cc)
			return cc.CompileExpr(env, &Vec{
				etag.Expand(KwWith),
				endThenId.ToConstexpr(env),
				&Vec{etag.Expand(KwJumpUnlessOp), cond},
			})
		})
		cc.InstallNamed(env, etag.Expand(KwWhenAnd), NmSyntax, and)

		or := SyntaxFn(func(cc *Compiler, env *Env, e *Vec) Value {
			etag, _ := CheckExpr(e, 2, 2, CtProc, cc)
			cond := CheckValue(e.At(1), IdentifierT, "cond", etag, cc)
			return cc.CompileExpr(env, &Vec{
				etag.Expand(KwWith),
				thenId.ToConstexpr(env),
				&Vec{etag.Expand(KwJumpIfOp), cond},
			})
		})
		cc.InstallNamed(env, etag.Expand(KwWhenOr), NmSyntax, or)

		n := condBody.Size()
		if n == 1 {
			cc.ErrorAt(condBody, etag).With("at least one condition required")
		} else if tail := condBody.At(n - 1); KwWhenAnd.MatchExpr(tail) == nil {
			cc.ErrorAt(tail, etag).With("the last condition must be an AND expression")
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
		cc.ErrorAt(elseTag, etag).With("invalid when form, expected `else`")
	} else if n != 6 {
		cc.ErrorAt(elseTag, etag).With("invalid when form, else body required")
	}
	elseBody := CheckBlockForm(e.At(5), "else-body", etag, cc)

	endElseId := etag.Expand(Gensym("end-else"))
	endElseNm := cc.InstallNamed(env, endElseId, NmLabel, &Label{})
	cc.CompileExpr(env, &Vec{etag.Expand(KwJump), endElseId.ToConstexpr(env)})
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
		cc.ErrorAt(old, etag).With("unknown name %s", old)
	} else if nm.Kind != NmMacro && nm.Kind != NmInline && nm.Kind != NmLabel {
		cc.ErrorAt(old, etag).With("aliases are not allowed for this type(%s)", NamedKindLabels[nm.Kind])
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
		cc.ErrorAt(etag).With("link is already registered")
	}
	cc.link = body
	return NIL
}

// SYNTAX: (link/keep ...)
func (cc *Compiler) sLinkKeep(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, -1, CtModule|CtProc, cc)
	r := []Value{KwComment, NIL}
	for _, i := range (*e)[1:] {
		CheckConst(i, IdentifierT, "name", etag, cc)
		r = append(r, cc.CompileExpr(env, i))
	}
	return cc.EmitCode(NewInst(e, InstMisc, r...))
}

// SYNTAX: (link/with ...)
func (cc *Compiler) sLinkWith(env *Env, e *Vec) Value {
	etag, n := CheckExpr(e, 3, -1, CtModule|CtProc, cc)

	r := &Vec{}
	for _, i := range (*e)[1 : n-1] {
		id := CheckConst(i, IdentifierT, "dependency", etag, cc)
		r.Push(id)
	}

	body := CheckBlockForm(e.At(n-1), "body", etag, cc)
	if body.ExprTag().Name != KwProg {
		cc.ErrorAt(body, etag).With("body must be scoped block")
	}

	nm := cc.InstallNamed(env, etag.Expand(Gensym("link/with")), NmLabel, &Label{})
	cc.EmitCode(NewInst(e, InstMisc, KwBeginDep, nm, env, r, NIL))
	for _, i := range (*body)[1:] {
		id, _ := AsTaggedVec(i)
		if id == nil || !id.IsPlain() {
			cc.ErrorAt(i, etag).With("invalid form")
		}
		if kw := id.Name; kw != KwAssert && kw != KwAlign && kw != KwLinkKeep {
			cc.ErrorAt(id, etag).With("unsupported form '%s'", id)
		}
		cc.CompileExpr(env, i)
	}
	return cc.EmitCode(NewInst(e, InstMisc, KwEndDep))
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
			i := CheckValue(i, ConstexprT, "comment element", etag, cc)
			r = append(r, cc.CompileExpr(env, i))
		}
		cc.EmitCode(NewInst(e, InstMisc, r...))
	default:
		cc.ErrorAt(k, etag).With("unknown pragma: %s", k)
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
		cc.ErrorAt(e.At(4), etag).With("rest parameter name required")
	}

	opts := []Value{}
	args := []*Keyword{}
	for x, a := range *as {
		a := a.(*Vec)
		k := CheckConstPlainId(a.At(0), "macro parameter", etag, cc)
		v := a.At(1)
		if rest && x == len(*as)-1 {
			if v != NIL {
				cc.ErrorAt(k, etag).With("the rest parameter cannot have default value")
			}
			opts = append(opts, v)
		} else if len(opts) > 0 && v == NIL {
			cc.ErrorAt(k, etag).With("default value required")
		}
		if v != NIL {
			opts = append(opts, v)
		}
		if slices.Index(args, k.Name) > -1 {
			cc.ErrorAt(k, etag).With("parameter %s is already defined", k.Name)
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
		if slices.Index(args, k.Name) > -1 || slices.Index(vars, Value(k.Name)) > -1 {
			cc.ErrorAt(k, etag).With("variable %s is already defined", k.Name)
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
			cc.ErrorAt(addr, etag).With("cannot bind inline proc")
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
	cc.EmitCode(NewInst(e, InstMisc, KwBeginProc, NIL)) // #1(NIL): last inst link
	cc.EmitCode(NewInst(e, InstLabel, nm))

	env = env.Enter()
	procNm := cc.InstallNamed(env, etag.Expand(KwPROCNAME), NmLabel, &Label{})
	cc.EmitCode(NewInst(e, InstLabel, procNm))

	cc.sBlock(env, body)
	cc.EmitCode(NewInst(e, InstMisc, KwEndProc, nm))
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
	} else if cond != NIL {
		cc.CompileExpr(env, &Vec{etag.Expand(op), name, cond})
	} else {
		cc.CompileExpr(env, &Vec{etag.Expand(op), name})
	}
	return cc.EmitCode(NewInst(e, InstMisc, KwCallproc, name, sig))
}

// SYNTAX: (fallthrough ...)
func (cc *Compiler) sFallthrough(env *Env, e *Vec) Value {
	CheckExpr(e, 1, 1, CtProc, cc)
	return cc.EmitCode(NewInst(e, InstMisc, KwFallthrough, NIL)) // NIL(#1): next proc nm
}

// SYNTAX: (tco a)
func (cc *Compiler) sTco(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtProc, cc)

	body := CheckBlockForm(e.At(1), "body", etag, cc)
	if body.Size() != 2 {
		cc.ErrorAt(body, etag).With("only one statement is allowed within the tco form")
	}

	call := KwCallproc.MatchExpr(body.At(1))
	if call == nil {
		cc.ErrorAt(body.At(1), etag).With("proc call required")
	}
	if cond := call.At(2); cond != NIL {
		cc.ErrorAt(cond, etag).With("the conditional call is not a tail call")
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
				cc.ErrorAt(k, name).With("default value required")
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

	body := CheckValue(e.At(3), VecT, "data body", etag, cc)
	if KwMulOp.MatchId(body.At(1)) == nil {
		cc.ErrorAt(body.At(1), etag).With("repeat operator must be `*`")
	}
	repeat := CheckAndEvalConstAs(body.At(2), env, IntT, "repeat count", etag, cc)
	if repeat < 1 {
		cc.ErrorAt(body.At(2), etag).With("invalid repeat count %d", repeat)
	}

	values := body.At(0)
	switch values := values.(type) {
	case *Constexpr:
		t, size := cc.evaluateType(env, e.At(2), DataSizeAuto, etag)
		if t != ByteType || size != DataSizeAuto {
			cc.ErrorAt(e.At(2), etag).With("invalid blob type")
		}
		v := EvalConstAs(values, env, BlobT, "blob", etag, cc)
		info.Link = NewInst(e, InstBlob, ByteType, NIL, v)
	case *Vec:
		defaultSize := DataSizeSingle
		if KwDataList.MatchId(values.At(0)) != nil {
			defaultSize = DataSizeAuto
		}
		t, size := cc.evaluateType(env, e.At(2), defaultSize, etag)
		v := cc.CompileExpr(env, values)
		info.Link = NewInst(e, InstData, t, repeat, NIL, Int(size), v)
	default: // NIL
		t, size := cc.evaluateType(env, e.At(2), DataSizeSingle, etag)
		if size > 1 {
			repeat *= Int(size)
		}
		info.Link = NewInst(e, InstDS, t, repeat, NIL)
	}
	cc.EmitCode(info.Link)

	if body.At(3) != NIL {
		name := CheckConstPlainId(body.At(3), "section", etag, cc)
		section := cc.Module.FindOrNewSection(name.Name)
		cc.EmitCodeToSection(section, cc.LeaveCodeBlock()...)
	} else if body.At(4) != NIL {
		if nm == nil {
			cc.ErrorAt(body.At(4), etag).With("the addressed data must be named")
		}
		if info.Link.Kind != InstDS {
			cc.ErrorAt(values, etag).With("the addressed data cannot contain any elements")
		}

		addr := CheckValue(body.At(4), ConstexprT, "data address", etag, cc)
		info.At = cc.CompileExpr(env, addr).(*Constexpr)
		cc.LeaveCodeBlock()
		cc.EmitCodeByEnv(env, NewInst(e, InstBind, nm))
	} else {
		cc.EmitCode(cc.LeaveCodeBlock()...)
	}

	return info
}

// SYNTAX: (#.datalist ...)
func (cc *Compiler) sDataList(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 1, -1, CtModule|CtProc, cc)
	r := &Vec{etag}
	for _, i := range (*e)[1:] {
		r.Push(cc.CompileExpr(env, i))
	}
	return r
}

// SYNTAX: (#.structdata ...)
func (cc *Compiler) sStructData(env *Env, e *Vec) Value {
	return cc.sDataList(env, e)
}

// SYNTAX: (#.struct isinner name field ...)
func (cc *Compiler) sStruct(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, -1, CtModule|CtProc, cc)

	t := NewDatatype(nil)
	if e.At(1) == NIL {
		t.Name = CheckConstPlainId(e.At(2), "struct name", etag, cc)
		cc.InstallNamed(env, t.Name, NmDatatype, t)
	} else {
		if e.At(2) != NIL {
			cc.ErrorAt(e.At(2), etag).With("named inner struct is not allowed")
		}
		t.Name = etag.Expand(Gensym("struct"))
	}

	for _, i := range (*e)[3:] {
		i := i.(*Vec)
		id := CheckConstPlainId(i.At(0), "field name", etag, cc)
		if t.Map[id.Name] != nil {
			cc.ErrorAt(id, etag).With("the field %s is already defined", id)
		}

		u, size := cc.evaluateType(env, i.At(1), DataSizeSingle, id)
		t.AddField(id.Name, u, size)
	}
	return t
}

// SYNTAX: (#.array base size)
func (cc *Compiler) sArray(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, 3, CtModule|CtProc, cc)
	id := CheckConstPlainId(e.At(1), "element type", etag, cc)

	u := GetBuiltinTypeById(id)
	if u == nil {
		cc.ErrorAt(id, etag).With("only builtin types can be used as array element types")
	}

	n := CheckAndEvalConstAs(e.At(2), env, IntT, "size", etag, cc)
	if n <= 0 {
		cc.ErrorAt(e.At(2), etag).With("invalid data length")
	}

	t := NewDatatype(etag.Expand(Gensym("array")))
	t.AddField(KwUNDER, u, int(n))
	return t
}

// SYNTAX: (#.BYTE n)
func (cc *Compiler) sByte(env *Env, e *Vec) Value {
	CheckExpr(e, 2, 2, CtProc, cc)
	a := cc.CompileExpr(env, e.At(1))
	if v, ok := a.(*Operand); ok {
		a = v.A0
	}
	return cc.EmitCode(NewInst(e, InstData, ByteType, Int(1), NIL, Int(DataSizeSingle), a))
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
	cc.ErrorAt(etag).With("invalid operands")
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

// SYNTAX: (#.field a ...)
func (cc *Compiler) sFieldOffset(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, -1, CtConstexpr, cc)
	id := CheckValue(e.At(1), IdentifierT, "base", etag, cc)

	offset := 0
	nm := cc.LookupNamed(env, id)
	var t *Datatype
	switch nm.Kind {
	case NmDatatype:
		t = nm.Value.(*Datatype)
	case NmLabel:
		if v := nm.Value.(*Label); v.LinkedToData() {
			t = v.Link.Args[0].(*Datatype)
			offset = int(EvalConstAs(e.At(1), env, IntT, "data", etag, cc))
		}
	}

	for _, v := range (*e)[2:] {
		if t == nil || len(t.Map) == 0 {
			cc.ErrorAt(id).With("%s is not a struct type", id)
		}

		id = CheckValue(v, IdentifierT, "field name", etag, cc)
		field := t.GetField(id.Name)
		if field == nil {
			cc.ErrorAt(id).With("unknown field %s", id)
		}

		offset += field.Offset
		t = field.Datatype
	}
	return Int(offset)
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

		u := NoOperand
		if op.Name == KwDot {
			cc.CompileExpr(env, i.At(1))
			continue
		} else if i.At(1) != NIL {
			u = cc.ExprToOperand(cc, cc.CompileExpr(env, i.At(1)))
		}

		m1 := cc.CtxOpMap[op.Name]
		if m1 == nil {
			cc.ErrorAt(op).With("unknown operator: %s", op)
		}

		m2 := m1[t.Kind]
		if m2 == nil {
			m2 = m1[KwAny]
		}
		if m2 == nil {
			cc.ErrorAt(t, op).With("cannot use %s as first operand for '%s'", t.Kind, op)
		}

		body := m2[u.Kind]
		if body == nil {
			body = m2[KwAny]
		}
		if body == nil {
			if u == NoOperand {
				cc.ErrorAt(u, op).With("[BUG] %s require second operand", op)
			} else {
				cc.ErrorAt(u, op).With("cannot use %s as second operand for '%s'", u.Kind, op)
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
	s := CheckConst(e.At(1), StrT, "message", etag, cc)
	if etag.ExpandedBy != nil {
		etag = etag.ExpandedBy
	}
	cc.ErrorAt(etag).With("%s", s)
	return NIL
}

// SYNTAX: (warn s)
func (cc *Compiler) sWarn(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtModule|CtProc, cc)
	s := CheckConst(e.At(1), StrT, "message", etag, cc)
	if etag.ExpandedBy != nil {
		etag = etag.ExpandedBy
	}
	cc.WarnAt(etag).With("%s", s)
	return NIL
}

// SYNTAX: (assert expr message)
func (cc *Compiler) sAssert(env *Env, e *Vec) Value {
	etag, n := CheckExpr(e, 2, 3, CtModule|CtProc, cc)
	cond := CheckValue(e.At(1), ConstexprT, "cond", etag, cc)

	s := NewStr("assertion failed")
	if n == 3 {
		s = CheckAndEvalConstAs(e.At(2), env, StrT, "message", etag, cc)
	}
	return cc.EmitCode(NewInst(e, InstAssert, cc.CompileExpr(env, cond), s))
}

// SYNTAX: (import s)
func (cc *Compiler) sImport(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtModule|CtProc, cc)
	id := CheckConstPlainId(e.At(1), "module name", etag, cc)

	mod := cc.FindModule(env, id.Name, id)
	if mod == nil {
		cc.ErrorAt(id, etag).With("unknown namespace %s", id)
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
		if id, ok := delta.(*Identifier); ok {
			if t := GetBuiltinTypeById(id); t != nil {
				delta = Int(-t.Size)
			}
		}
		CheckValue(delta, IntT, "delta", etag, cc)
	}

	id := etag.Expand(Gensym("patch." + name.String()))
	labelNm := cc.InstallNamed(env, id, NmLabel, &Label{})

	nm.Value.(*Label).At = &Constexpr{Env: env, Body: &Vec{etag.Expand(KwPlusOp), id, delta}}
	cc.EmitCode(NewInst(e, InstMisc, KwPatchAnchor, nm))
	return cc.EmitCode(NewInst(e, InstLabel, labelNm))
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

	fname := etag.Expand(Gensym("make-counter"))
	resetter := SyntaxFn(func(cc *Compiler, env *Env, e *Vec) Value {
		value = e.At(2).(Int)
		return value
	})
	cc.InstallNamed(env, fname, NmFunc, resetter)

	lhs := etag.Expand(Gensym("counter"))
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
	cc.ErrorAt(etag).With("invalid expansion")
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

func (cc *Compiler) sFieldSize(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 3, -1, CtConstexpr, cc)
	id := CheckValue(e.At(1), IdentifierT, "base", etag, cc)
	nm := cc.LookupNamed(env, id)

	var t *Datatype
	switch nm.Kind {
	case NmDatatype:
		t = nm.Value.(*Datatype)
	case NmLabel:
		if v := nm.Value.(*Label); v.LinkedToData() {
			t = v.Link.Args[0].(*Datatype)
		}
	}

	n := 1
	for _, v := range (*e)[2:] {
		if t == nil || len(t.Map) == 0 {
			cc.ErrorAt(id).With("%s is not a struct type", id)
		}

		id = CheckValue(v, IdentifierT, "field name", etag, cc)
		field := t.GetField(id.Name)
		if field == nil {
			cc.ErrorAt(id).With("unknown field %s", id)
		}

		t = field.Datatype
		n = field.Size
	}
	if n < 0 {
		n = 1
	}
	return Int(t.Size * n)
}

// SYNTAX: (sizeof a)
func (cc *Compiler) sSizeof(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtConstexpr, cc)

	if e := KwField.MatchExpr(e.At(1)); e != nil {
		return cc.sFieldSize(env, e)
	}

	name := CheckValue(e.At(1), IdentifierT, "label", etag, cc)
	nm := cc.LookupNamed(env, name)
	if nm == nil {
		cc.ErrorAt(name, etag).With("unknown label name %s", name)
	}

	switch nm.Kind {
	case NmLabel:
		v := nm.Value.(*Label)
		if !v.LinkedToData() {
			cc.ErrorAt(name, etag).With("%s is not a data label", name)
		}
		return Int(v.Link.Size)
	case NmDatatype:
		v := nm.Value.(*Datatype)
		return Int(v.Size)
	}

	cc.ErrorAt(name, etag).With("sizeof requires a data label or struct")
	return NIL
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
		cc.ErrorAt(name, etag).With("unknown name %s", name)
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
				cc.ErrorAt(etag).With("invalid fragment")
			}
		default:
			cc.ErrorAt(etag).With("[BUG] invalid fragment")
		}
	}
	return etag.Expand(Intern(s))
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
	a := CheckValue(e.At(1), IntT, "value", etag, cc)

	return ^a
}

// FUN: (! a)
func (cc *Compiler) fLogicalNot(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtConstexpr, cc)
	a := CheckValue(e.At(1), IntT, "value", etag, cc)

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

// FUN: (opcode asm size)
func (cc *Compiler) fOpcode(_ *Env, e *Vec) Value {
	etag, n := CheckExpr(e, 2, 3, CtConstexpr, cc)
	asm := CheckValue(e.At(1), StrT, "mnemonic", etag, cc)
	size := 1
	if n == 3 {
		size = int(CheckValue(e.At(2), IntT, "size", etag, cc))
	}

	cc.nested = append([]Value{etag}, cc.nested...)
	r := cc.Parse("<opcode>", []byte(*asm)).AtOrUndef(1)
	v := CheckValue(r, VecT, "mnemonic", etag, cc)

	op := v.ExprTag()
	if !op.IsPlain() {
		cc.ErrorAt(etag).With("unknown mnemonic")
	}

	found, ok := cc.InstMap[op.Name]
	if !ok {
		cc.ErrorAt(etag).With("unknown mnemonic")
	}

	env := cc.Builtins.Enter()
	ab := []*Operand{}
	tab := found.(InstPat)
	for _, i := range (*v)[1:] {
		i := cc.ExprToOperand(cc, cc.CompileExpr(env, i))
		if c, ok := i.A0.(*Constexpr); ok {
			n := CheckConst(c, IntT, "operand", etag, cc)
			cc.Constvals[c] = n
			cc.AdjustOperand(cc, i, int(n), etag)
		}

		ab = append(ab, i)
		tab, ok = tab[i.Kind].(InstPat)
		if !ok {
			break
		}
	}

	body, ok := tab[nil].(InstDat)
	if !ok || body[0].Kind == BcTemp {
		cc.ErrorAt(etag).With("invalid operands for '%s'", op)
	}

	inst := NewInst(e, InstMisc, KwUNDER)
	code := uint64(0)
	for x := range min(size, len(body), 8) {
		code = (code << 8) | uint64(cc.g.expandBCode(inst, body[x], ab, 0))
	}
	cc.nested = cc.nested[1:]

	return Int(code)
}

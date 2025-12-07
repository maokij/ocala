package core

import (
	"fmt"
)

var nmReturn = &Named{Name: KwReturn}

func OperandToNamedLabel(cc *Compiler, a Value) (*Named, *Label) {
	if nm := cc.OperandToNamed(cc, a); nm != nil && nm.Kind == NmLabel {
		return nm, nm.Value.(*Label)
	}
	return nil, nil
}

func (cc *Compiler) removeUnreachableCode(insts []*Inst) {
	unreachable := false
	for _, i := range insts {
		n := len(i.Args)
		switch i.Kind {
		case InstMisc:
			if unreachable && i.Args[0] == KwPatchAnchor {
				cc.ErrorAt(i).With("the patch address `%s` is unreachable", i.Args[1].(*Named).Name)
			}
		case InstConst, InstBind, InstAssert,
			InstOrg, InstData, InstDS, InstBlob, InstAlign:
			// NOP
		case InstCode:
			if unreachable {
				i.CommentOut()
			} else if i.MatchCode(KwReturn) && n == 1 {
				unreachable = true
			} else if i.MatchCode(KwJump) && n == 2 {
				_, label := OperandToNamedLabel(cc, i.Args[1])
				if label != nil && !label.IsComputed() {
					unreachable = true
				}
			}
		case InstLabel:
			unreachable = false
		default:
			cc.ErrorAt(i).With("[BUG] invalid inst kind %v", i)
		}
	}
}

func (cc *Compiler) simplifyRedundantJump(insts []*Inst) {
	type handle struct {
		nm *Named
	}
	labels := map[*Named]*handle{}

	next := func(nm *Named) *Named {
		if found := labels[nm]; found != nil {
			return found.nm
		}
		return nil
	}

	detectCycle := func(nm *Named) bool {
		for a, b := nm, nm; b != nil; {
			a = next(a)
			b = next(next(b))

			if a != nil && a == b {
				return true
			}
		}
		return false
	}

	var m *handle
	jumps := []*Inst{}
	for _, i := range insts {
		n := len(i.Args)
		switch i.Kind {
		case InstMisc, InstConst, InstBind, InstAssert:
			// NOP
		case InstOrg, InstData, InstDS, InstBlob, InstAlign:
			m = nil
		case InstCode:
			if i.MatchCode(KwReturn) {
				if n == 1 && m != nil {
					m.nm = nmReturn
				}
			} else if i.MatchCode(KwJump) {
				nm, label := OperandToNamedLabel(cc, i.Args[1])
				if label != nil && !label.IsComputed() {
					jumps = append(jumps, i)
					if n == 2 && m != nil {
						m.nm = nm
					}
				}
			}
			m = nil
		case InstLabel:
			if m == nil {
				m = &handle{}
			}
			labels[i.Args[0].(*Named)] = m
		default:
			cc.ErrorAt(i).With("[BUG] invalid inst kind %v", i)
		}
	}

	resolved := map[*Named]*Named{}
	for k, v := range labels {
		if detectCycle(k) {
			to := labels[k].nm
			cc.WarnAt(k).With("jump cycle detected: %s -> %s", k.Name, to.Name)
		} else {
			nm := v.nm
			for a := nm; a != nil; a = next(a) {
				nm = a
			}
			resolved[k] = nm
		}
	}

	for _, i := range jumps {
		nm := cc.OperandToNamed(cc, i.Args[1])
		a := resolved[nm]
		if a == nil {
			continue
		}

		if a == nmReturn {
			i.Args = append([]Value{KwReturn}, i.Args[2:]...)
			i.Comment = fmt.Sprintf("// %s", nm.AsmName)
		} else if a != nm {
			c := a.Name.ToId(a.Token).ToConstexpr(a.Env)
			i.Args[1] = cc.ExprToOperand(cc, c)
			i.Comment = fmt.Sprintf("// %s", nm.AsmName)
		}
	}
}

func (cc *Compiler) removeRedundantCode(insts []*Inst) {
	type marker struct {
		nm   *Named
		cond *Keyword
		inst *Inst
		id   int
	}
	nextLabels := map[*Named]int{}
	nextSize := 0

	eqlabel := func(labels map[*Named]int, a, b *Named) bool {
		if a == b {
			return true
		}
		if la := labels[a]; la != 0 {
			return la == labels[b]
		}
		return false
	}

	for changed := true; changed; {
		changed = false
		size := nextSize
		labels := nextLabels
		nextSize = 0
		nextLabels = map[*Named]int{}

		m := marker{}
		for _, i := range insts {
			n := len(i.Args)
			switch i.Kind {
			case InstMisc, InstConst, InstBind, InstAssert:
			// NOP
			case InstOrg, InstData, InstDS, InstBlob, InstAlign:
				m = marker{}
			case InstCode:
				if i.MatchCode(KwReturn) {
					var cond *Keyword
					if n == 1 {
						// L0: return;    L1: return => L0: L1: return
						// L0: return Z?; L1: return => L0: L1: return
						if m.nm == nmReturn {
							m.inst.CommentOut()
							changed = true
						}
					} else if n == 2 {
						// L0: return Z?; L1: return Z? -> L0: L1: return Z?
						cond = i.Args[1].(*Operand).Kind
						if m.nm == nmReturn && m.cond == cond {
							m.inst.CommentOut()
							changed = true
						}
					}
					m = marker{nm: nmReturn, cond: cond, inst: i}
					continue
				} else if i.MatchCode(KwJump) {
					nm, label := OperandToNamedLabel(cc, i.Args[1])
					if label != nil && !label.IsComputed() {
						var cond *Keyword
						if n == 2 {
							// L0: jump L;    L1: jump L => L0: L1: jump L
							// L0: jump L Z?; L1: jump L => L0: L1: jump L
							if eqlabel(labels, nm, m.nm) {
								m.inst.CommentOut()
								changed = true
							}
						} else if n == 3 {
							// L0: jump L Z?; L1: jump L Z? => L0: L1: jump L Z?
							cond = i.Args[2].(*Operand).Kind
							if eqlabel(labels, nm, m.nm) && m.cond == cond {
								m.inst.CommentOut()
								changed = true
							}
						}
						m = marker{nm: nm, cond: cond, inst: i}
						continue
					}
				}
				m = marker{}
			case InstLabel:
				nm := i.Args[0].(*Named)
				if m.id == 0 {
					nextSize++
					m.id = nextSize
				}
				nextLabels[nm] = m.id

				if nm == m.nm {
					m.inst.CommentOut()
					changed = true
				}
			default:
				cc.ErrorAt(i).With("[BUG] invalid inst kind %v", i)
			}
		}

		if size != nextSize {
			changed = true
		}
	}
}

func (cc *Compiler) optimizeFlow(insts []*Inst) {
	if !cc.EnableOptimizeFlow {
		return
	}

	var end *Inst
	n := 0
	for x, i := range insts {
		if i.IsMisc(KwBeginProc) {
			n = x
			end = i.Args[1].(*Inst)
		} else if i == end {
			chunk := insts[n : x+1]
			cc.removeUnreachableCode(chunk)
			cc.simplifyRedundantJump(chunk)
			cc.removeRedundantCode(chunk)
		}
	}
}

func (cc *Compiler) optimizeLink(insts []*Inst) {
	if len(cc.LinkEntryPoints) == 0 {
		return
	}

	type part struct {
		nm     *Named
		insts  []*Inst
		marked bool
	}

	var traverse func(Value)
	traverse = func(v Value) {
		switch v := v.(type) {
		case *Vec:
			for _, i := range *v {
				traverse(i)
			}
		case *Constexpr:
			EvalAndCacheIfConst(v, cc)
		}
	}

	mark := func(insts []*Inst) {
		for _, i := range insts {
			switch i.Kind {
			case InstMisc:
				switch i.Args[0] {
				case KwComment:
					traverse(NewVec(i.Args[2:]))
				case KwFallthrough:
					i.Args[1].(*Named).Found()
				}
			case InstAssert:
				EvalAndCacheIfConst(i.Args[0], cc)
			case InstConst:
				EvalAndCacheIfConst(GetNamedValue(i.Args[0], ConstexprT), cc)
			case InstBind:
				EvalAndCacheIfConst(GetNamedValue(i.Args[0], LabelT).At, cc)
			case InstCode:
				for _, i := range i.Args[1:] {
					EvalAndCacheIfConst(i.(*Operand).A0, cc)
				}
			case InstData:
				traverse(i.Args[4])
			case InstDS:
				EvalAndCacheIfConst(i.Args[1], cc)
			case InstLabel, InstOrg, InstAlign, InstBlob:
				// OK
			default:
				cc.ErrorAt(i).With("[BUG] invalid inst kind %v", i)
			}
		}
	}

	deps := map[*Named][]*Named{}
	for _, i := range insts {
		if !i.IsMisc(KwBeginDep) {
			continue
		}

		from := i.Args[1].(*Named)
		env := i.Args[2].(*Env)
		for _, id := range *i.Args[3].(*Vec) {
			id := id.(*Identifier)
			to := cc.LookupNamed(env, id)
			if to == nil || to.Kind != NmLabel {
				cc.ErrorAt(id, i).With("invalid dependency '%s'", id)
			}
			deps[to] = append(deps[to], from)
		}
	}

	parts := []*part{}
	others := []*Inst{}
	register := func(nm *Named, insts []*Inst) {
		m := &part{nm: nm, insts: insts}
		nm.OnFound = func(*Named) {
			if !m.marked {
				m.marked = true
				for _, i := range deps[nm] {
					i.Found()
				}
				cc.EnterContext(CtModule)
				mark(m.insts)
				cc.LeaveContext()
			}
		}
		parts = append(parts, m)
	}

	label := &Label{}
	cc.initConstvals(PhLink)
	for x := 0; x < len(insts); x++ {
		i := insts[x]
		switch i.Kind {
		case InstMisc:
			switch i.Args[0] {
			case KwComment, KwFallthrough:
				others = append(others, i)
			case KwBeginProc:
				n := x
				nm := insts[x+1].Args[0].(*Named)
				for ; x < len(insts); x++ {
					i := insts[x]
					if i.IsMisc(KwEndProc) && i.Args[1] == nm {
						break
					} else if i.IsMisc(KwPatchAnchor) {
						labelNm := insts[x+1].Args[0].(*Named)
						labelNm.OnFound = func(*Named) {
							nm.OnFound(nm)
						}
						parts = append(parts, &part{nm: labelNm, marked: true})
					}
				}
				register(nm, insts[n:x+1])
			case KwBeginDep:
				n := x
				nm := i.Args[1].(*Named)
				for ; x < len(insts); x++ {
					if insts[x].IsMisc(KwEndDep) {
						break
					}
				}
				register(nm, insts[n:x+1])
			}
		case InstConst:
			register(i.Args[0].(*Named), insts[x:x+1])
		case InstBind:
			register(i.Args[0].(*Named), insts[x:x+1])
		case InstLabel:
			label = GetNamedValue(i.Args[0], LabelT)
		case InstData, InstDS:
			if label.Link == i {
				register(insts[x-1].Args[0].(*Named), insts[x-1:x+1])
				continue
			}
			others = append(others, i)
		case InstAssert, InstCode:
			others = append(others, i)
		case InstOrg, InstAlign, InstBlob:
			// OK
		default:
			cc.ErrorAt(i).With("[BUG] invalid inst kind %v", i)
		}
	}

	mark(others)
	for _, id := range cc.LinkEntryPoints {
		nm := cc.LookupNamed(cc.Toplevel, id)
		if nm == nil || nm.OnFound == nil {
			cc.ErrorAt(id).With("invalid entry point '%s'", id)
		}
		nm.Found()
	}

	for _, i := range parts {
		i.nm.OnFound = nil
		delete(deps, i.nm)
		if i.marked {
			continue
		}

		name := i.nm.AsmName.String()
		for _, i := range i.insts {
			i.Kind = InstMisc
			i.Args = []Value{KwUNDER}
		}
		if len(i.insts) > 1 {
			i.insts[0].Args = []Value{KwComment, NewStr("// UNLINKED: " + name)}
		}
	}

	for k, v := range deps {
		for _, i := range v {
			cc.WarnAt(i, k).With("the label cannot be used as a dependency")
		}
	}
}

package core

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func _sizeof(v Value) int {
	switch v {
	case KwByte:
		return 1
	case KwWord:
		return 2
	}
	panic("invalid type")
}

func _asValueSlice(v Value) []Value {
	if vec, ok := v.(*Vec); ok {
		return *vec
	}
	return []Value{v}
}

func (g *Generator) validateProcTail(inst *Inst) {
	if (inst.Kind == InstCode && g.cc.IsValidProcTail(g.cc, inst)) ||
		(inst.Kind == InstMisc && inst.Args[0] == KwFallthrough) {
		return
	}
	g.cc.RaiseCompileError(inst.ExprTag(),
		"the last instruction must be a return/fallthrough within the proc")
}

func (g *Generator) validateFallthrough(insts []*Inst) {
	for _, i := range insts {
		if i.Kind == InstMisc && i.Args[0] == KwEndProc {
			continue
		} else if i.Kind == InstLabel && i.Args[1] == KwProc {
			break // OK
		}
		g.cc.RaiseCompileError(i.ExprTag(),
			"the fallthrough must be followed by a proc")
	}
}

func (g *Generator) validateCallproc(inst *Inst) {
	c := inst.Args[1].(*Constexpr)
	id := c.Body.(*Identifier)
	nm := g.cc.LookupNamed(c.Env, id)
	if nm == nil {
		g.cc.RaiseCompileError(inst.ExprTag(), "undefined proc %s", id)
	}

	var a *Sig
	switch nm.Kind {
	case NmLabel:
		a = nm.Value.(*Label).Sig
	case NmInline:
		a = nm.Value.(*Inline).Sig
	}
	if a == nil {
		g.cc.RaiseCompileError(inst.ExprTag(), "%s is not a proc", id)
	}

	b := inst.Args[2].(*Sig)
	if !a.Equals(b) {
		g.cc.RaiseCompileError(inst.ExprTag(), "proc signature mismatch: %s.\n"+
			"  expected %v,\n"+
			"  given    %v", id, a, b)
	}
}

func (g *Generator) validateInsts(insts []*Inst) {
	ci := insts[0]
	state := []*Inst{}

	for x, i := range insts {
		switch i.Kind {
		case InstCode:
			ci = i
		case InstLabel:
			GetNamedValue(i.Args[0]).(*Label).Addr = 0
			switch i.Args[1] {
			case KwProc:
				state, ci = append(state, ci), i
			case KwLabel:
				ci = i
			}
		case InstMisc:
			switch i.Args[0] {
			case KwCallproc:
				g.validateCallproc(i)
			case KwEndProc:
				g.validateProcTail(ci)
				ci, state = state[len(state)-1], state[:len(state)-1]
			case KwFallthrough:
				g.validateFallthrough(insts[x+1:])
				ci = i
			}
		}
	}
}

func (g *Generator) findInstBody(inst *Inst, adjust bool) []BCode {
	etag := inst.ExprTag()
	op := inst.Args[0].(*Keyword)

	m, ok := g.cc.InstMap[op]
	if !ok {
		g.cc.RaiseCompileError(etag, "[BUG] unknown instruction: %s", op)
	}

	p := m.(InstPat)
	for x, i := range inst.Args[1:] {
		i := i.(*Operand)
		if adjust {
			g.cc.AdjustOperand(g.cc, i, etag)
		}
		p, ok = p[i.Kind].(InstPat)
		if !ok {
			g.cc.RaiseCompileError(etag, "cannot use %s as operand#%d for %s", i.Kind, x+1, op)
		}
	}

	body, ok := p[nil].(InstDat)
	if !ok || (!adjust && body[0].Kind == BcTemp) {
		g.cc.RaiseCompileError(etag, "invalid operands for %s. "+
			"some operand values may be out of range", op)
	}

	if g.Optimizer.OptimizeBCode != nil {
		body = g.Optimizer.OptimizeBCode(g.cc, inst, body, !adjust)
	}
	return body
}

func (g *Generator) resolveInstsWithoutOptimize(insts []*Inst, verify bool) int {
	optimize := g.Optimizer.OptimizeBCode
	g.Optimizer.OptimizeBCode = nil
	codeSize := g.resolveInsts(insts, false)
	g.Optimizer.OptimizeBCode = optimize
	return codeSize
}

func (g *Generator) resolveInsts(insts []*Inst, verify bool) int {
	g.cc.Pc = 0
	g.changes = 0
	g.cc.initConstvals(false)

	org := 0
	limit := 0
	codeSize := 0
	for _, i := range insts {
		size := i.size

		switch i.Kind {
		case InstMisc:
			switch i.Args[0] {
			case KwComment:
				for _, i := range i.Args[2:] {
					EvalAndCacheIfConst(i, g.cc)
				}
			}
		case InstAssert:
			EvalAndCacheIfConst(i.Args[0], g.cc)
		case InstLabel:
			v := GetNamedValue(i.Args[0]).(*Label)
			if v.Addr != g.cc.Pc {
				v.Addr = g.cc.Pc
				g.Changed()
			}
		case InstConst:
			EvalAndCacheIfConst(GetNamedValue(i.Args[0]), g.cc)
		case InstOrg:
			if n := g.cc.Pc - org; verify && limit > 0 && n > limit {
				g.cc.RaiseCompileError(i.ExprTag(), "size limit exceeded(%d/%d)", n, limit)
			}
			if addr := int(i.Args[0].(Int)); addr >= 0 {
				org = addr
				limit = int(i.Args[1].(Int))
				g.cc.Pc = org
				g.cc.Org = org
			}
		case InstCode:
			i.size = len(g.findInstBody(i, true))
			g.cc.Pc += i.size
			codeSize += i.size
		case InstData:
			n := _sizeof(i.Args[0])
			pc := g.cc.Pc // increment the field directly on each iteration for __PC__
			for _, v := range *i.Args[2].(*Vec) {
				for _, v := range _asValueSlice(EvalAndCacheIfConst(v, g.cc)) {
					switch v := v.(type) {
					case Int:
						g.cc.Pc += n
					case *Str:
						if n != 1 {
							g.cc.RaiseCompileError(i.ExprTag(), "word type data cannot contain byte strings")
						}
						g.cc.Pc += len(*v)
					default:
						g.cc.RaiseCompileError(i.ExprTag(), "invalid data value %T", v)
					}
				}
			}
			i.size = (g.cc.Pc - pc) * int(i.Args[1].(Int))
			g.cc.Pc = pc + i.size
		case InstDS:
			i.size = _sizeof(i.Args[0]) * int(i.Args[1].(Int))
			g.cc.Pc += i.size
		case InstAlign:
			n := int(i.Args[0].(Int)) - 1
			if n < 0 || (n&(n+1)) != 0 {
				g.cc.RaiseCompileError(i.ExprTag(), "the alignment size must be power of 2")
			}
			i.size = ((g.cc.Pc + n) & ^n) - g.cc.Pc
			g.cc.Pc += i.size
		case InstFile:
			i.size = len(*i.Args[1].(*Binary))
			g.cc.Pc += i.size
		default:
			g.cc.RaiseCompileError(i.ExprTag(), "[BUG] invalid inst kind %v", i)
		}

		if size != i.size {
			g.Changed()
		}

		if verify && g.changes > 0 {
			fmt.Println(size, i.size)
			g.cc.RaiseCompileError(i.ExprTag(), "cannot determine the size of the instruction")
		}
	}

	return codeSize
}

func (g *Generator) expandBCode(c BCode, ab []*Operand, pc int) byte {
	switch c.Kind {
	case BcByte:
		return c.A0
	case BcLow:
		e := ab[c.A0].A0.(*Constexpr)
		v := g.cc.Constvals[e].(Int)
		return byte(v)
	case BcHigh:
		e := ab[c.A0].A0.(*Constexpr)
		v := g.cc.Constvals[e].(Int)
		return byte(uint16(v) >> 8)
	case BcRlow:
		e := ab[c.A0].A0.(*Constexpr)
		v := g.cc.Constvals[e].(Int)
		d := int(v) - pc + int(int8(c.A1))
		if (c.A2 == 1 && (d < -0x80 || d > 0x7f)) ||
			(c.A2 == 2 && (d < -0x8000 || d > 0x7fff)) {
			g.cc.RaiseCompileError(e.Token.TagId(), "the relative address is out of range")
		}
		return byte(d)
	case BcRhigh:
		e := ab[c.A0].A0.(*Constexpr)
		v := g.cc.Constvals[e].(Int)
		d := int(v) - pc + int(int8(c.A1))
		if (c.A2 == 1 && (d < -0x80 || d > 0x7f)) ||
			(c.A2 == 2 && (d < -0x8000 || d > 0x7fff)) {
			g.cc.RaiseCompileError(e.Token.TagId(), "the relative address is out of range")
		}
		return byte(d >> 8)
	case BcImp:
		e := ab[c.A0].A0.(*Constexpr)
		v := g.cc.Constvals[e].(Int)
		base, mask, shift := c.A1, c.A2, c.A3
		if v < 0 || v > Int(mask) {
			g.cc.RaiseCompileError(e.Token.TagId(), "the operand only accepts 0..%d", mask)
		}
		return base | ((byte(v) & mask) << shift)
	case BcMap:
		e := ab[c.A0].A0.(*Constexpr)
		v := g.cc.Constvals[e].(Int)
		m := g.cc.BMaps[c.A1]
		mask, min, max, items := m[0], m[1], m[2], m[3:]
		if v < Int(min) || v > Int(max) {
			g.cc.RaiseCompileError(e.Token.TagId(), "the operand only accepts %d..%d", min, max)
		}
		return items[(byte(v)-min)&mask]
	}
	panic("[BUG] cannot happen")
}

func (g *Generator) ValueToAsm(env *Env, v Value) string {
	switch v := v.(type) {
	case Int:
		return fmt.Sprintf("%d", v)
	case *Str:
		return v.Inspect()
	case *Constexpr:
		s := g.ValueToAsm(v.Env, v.Body)
		if _, ok := v.Body.(*Vec); ok && s[0] == '(' {
			return s[1 : len(s)-1]
		}
		return s
	case *Vec:
		op := v.At(0).(*Identifier)
		a := []string{}
		for _, i := range (*v)[1:] {
			a = append(a, g.ValueToAsm(env, i))
		}

		if g.cc.BinOps[op.String()] > 0 && len(a) == 2 {
			return fmt.Sprintf("(%s %s %s)", a[0], op.String(), a[1])
		}
		return fmt.Sprintf("%s(%s)", op.String(), strings.Join(a, " "))
	case *Operand:
		return g.cc.OperandToAsm(g, v)
	case *Keyword:
		return v.String()
	case *Identifier:
		if nm := g.cc.LookupNamed(env, v); nm != nil {
			return nm.AsmName.String()
		}
		return v.String()
	}
	g.RaiseGenerateError("[BUG] invalid value: %T %s", v, v.Inspect())
	return "!"
}

var listPadding = bytes.Repeat([]byte{' '}, 40)

func (g *Generator) generateList(insts []*Inst, code []byte) []byte {
	list := &bytes.Buffer{}
	enabled := true
	mode := 1
	pc := 0
	pos := 0
	last := ""

	writes := func(s string, v ...any) {
		if enabled {
			if s == "" && last == "" {
				return
			}
			if s != "" {
				list.Write(listPadding)
				fmt.Fprintf(list, s, v...)
			}
			list.WriteByte('\n')
			last = s
		}
	}
	writeh := func(n, m int, s string, v ...any) {
		if enabled {
			for x, i := 0, 0; i == 0 || i < (n+7)/8; i++ {
				if i > 0 {
					fmt.Fprintf(list, "%11s", ":")
				} else if mode > 0 && m > 0 {
					fmt.Fprintf(list, "%06x %04x", pos, pc)
				} else {
					fmt.Fprintf(list, "%6s %04x", "-", pc)
				}

				hexlen := 0
				if n == 0 && m > 0 { // ds, align, file
					hexlen, _ = fmt.Fprintf(list, "[%d] .. ", m)
				} else if mode > 0 && x < n {
					n := min(8, n-x)
					hexlen = n*3 + 3
					fmt.Fprintf(list, "[%d]", n)
					for i := 0; i < n; i, x = i+1, x+1 {
						fmt.Fprintf(list, " %02x", code[pos+x])
					}
				}

				if i == 0 {
					list.Write(bytes.Repeat([]byte{' '}, 29-hexlen))
					fmt.Fprintf(list, s, v...)
				}
				list.WriteByte('\n')
			}
			last = "+"
		}
		pc += m
		if mode > 0 {
			pos += m
		}
	}

	writes("    ; generated by ocala")
	for _, i := range insts {
		switch i.Kind {
		case InstMisc:
			switch i.Args[0] {
			case KwListConstants:
				enabled = int(i.Args[1].(Int)) != 0
			case KwComment:
				s := ""
				for _, i := range i.Args[2:] {
					i := GetCachedValueIfConst(i, g.cc)
					s += " " + g.ValueToAsm(nil, i)
				}
				writes("    ; %s%s", i.Args[1].(*Str).String(), s)
			case KwEndProc:
				writes("")
			}
		case InstAssert: // NOP
		case InstOrg:
			pos = int(i.Args[3].(Int))
			mode = int(i.Args[2].(Int))
			if mode >= 3 {
				mode -= 3
			}

			if addr := int(i.Args[0].(Int)); addr >= 0 {
				pc = addr
				writes("")
				writeh(0, 0, "    .org %s", g.ValueToAsm(nil, Int(pc)))
			}
		case InstLabel:
			nm := i.Args[0].(*Named)
			s := g.ValueToAsm(nil, nm.AsmName)
			if i.Args[1] == KwProc {
				writes("")
			}
			writeh(0, 0, "%s:", s)
		case InstConst:
			nm := i.Args[0].(*Named)
			a := g.ValueToAsm(nil, nm.AsmName)
			b := g.ValueToAsm(nm.Env, nm.Value)
			if !strings.HasPrefix(a, ".__") {
				writes("    %s = %s", a, b)
			}
		case InstCode:
			kw := i.Args[0]
			if len(i.Args) > 2 {
				a := g.ValueToAsm(nil, i.Args[1])
				b := g.ValueToAsm(nil, i.Args[2])
				writeh(i.size, i.size, "    %-6s %s, %s", kw, a, b)
			} else if len(i.Args) == 2 {
				a := g.ValueToAsm(nil, i.Args[1])
				writeh(i.size, i.size, "    %-6s %s", kw, a)
			} else {
				writeh(i.size, i.size, "    %s", kw)
			}
		case InstData:
			p := pc
			d := fmt.Sprintf("    .%s ", i.Args[0].(*Keyword).String())
			r := _sizeof(i.Args[0])
			n := i.Args[2].(*Vec).Size()
			a := []string{}
			for x, i := range *i.Args[2].(*Vec) {
				s := g.ValueToAsm(nil, i)
				for _, i := range _asValueSlice(GetCachedValueIfConst(i, g.cc)) {
					switch i := i.(type) {
					case Int:
						a = append(a, s)
						if m := r * len(a); m == 8 || x == n-1 {
							writeh(m, m, "%s", d+strings.Join(a, ", "))
							a = a[:0]
						}
					case *Str:
						if m := r * len(a); m > 0 {
							writeh(m, m, "%s", d+strings.Join(a, ", "))
							a = a[:0]
						}
						writeh(len(*i), len(*i), "    .byte "+s)
					}
				}
			}
			if repeat := int(i.Args[1].(Int)); repeat > 1 {
				m := pc - p
				for i := 1; i < repeat; i++ {
					writeh(m, m, "    ; ... repeat %d/%d", i+1, repeat)
				}
			}
		case InstDS:
			if i.Args[0] == KwWord {
				writeh(0, i.size, "    .defw %d", int(i.Args[1].(Int)))
			} else {
				writeh(0, i.size, "    .defb %d", int(i.Args[1].(Int)))
			}
		case InstAlign:
			writeh(0, i.size, "    .align %d ; (.defb %d)", int(i.Args[0].(Int)), i.size)
		case InstFile:
			writeh(0, len(*i.Args[1].(*Binary)), "    .incbin %s", g.ValueToAsm(nil, i.Args[0]))
		default:
			panic("[BUG] cannot happen")
		}
	}

	return list.Bytes()
}

func (g *Generator) prepareToGenerateBin(insts []*Inst) {
	g.validateInsts(insts)

	codeSizes := []int{}
	oldCodeSize := g.resolveInstsWithoutOptimize(insts, false)
	for {
		codeSizes = append(codeSizes, oldCodeSize)
		codeSize := g.resolveInsts(insts, false)
		if codeSize >= oldCodeSize {
			break
		}
		oldCodeSize = codeSize
	}

	if g.DebugMode && g.Optimizer.OptimizeBCode != nil {
		v := []string{}
		for _, i := range codeSizes {
			v = append(v, strconv.Itoa(i))
		}
		fmt.Fprintln(g.ErrWriter, "code size:", strings.Join(v, " -> "))
	}
	g.resolveInsts(insts, true) // verify
}

func (g *Generator) GenerateBin(insts []*Inst) ([]byte, []byte) {
	g.prepareToGenerateBin(insts)

	pushdat := func(s []byte, v int, w bool) []byte {
		if w {
			v := uint16(v)
			return append(s, byte(v), byte((v&0xff00)>>8))
		}
		return append(s, byte(v))
	}

	g.cc.Pc = 0
	code := []byte{}
	seq := []byte{}
	limit := 0
	mode := 1
	anchor := 0
	for _, i := range insts {
		switch i.Kind {
		case InstAssert:
			v := GetCachedValueIfConst(i.Args[0], g.cc)
			if v == Int(0) {
				g.cc.RaiseCompileError(i.ExprTag(), i.Args[1].(*Str).String())
			}
		case InstOrg:
			// mode -- 0: noload, 1: load, 2: fill, 3: +noload, 4: +load
			if mode > 0 {
				code = append(code, seq...)
			}
			seq = []byte{}

			size := len(code) - anchor
			if limit > 0 && size > limit {
				g.cc.RaiseCompileError(i.ExprTag(), "size limit exceeded(%d/%d)", size, limit)
			}

			addr := int(i.Args[0].(Int))
			mode = int(i.Args[2].(Int))
			if mode < 0 || mode >= 5 {
				g.cc.RaiseCompileError(i.ExprTag(), "invalid org mode")
			} else if mode >= 3 { // mode 3/4
				mode -= 3
				if limit == 0 {
					g.cc.RaiseCompileError(i.ExprTag(), "a pack mode must follow a fill/pack mode")
				}
			} else { // mode 0/1/2
				if limit > 0 {
					code = append(code, make([]byte, limit-size)...)
					if addr < 0 {
						g.cc.RaiseCompileError(i.ExprTag(), "the orgin address required")
					}
				}
				limit = 0
				anchor = len(code)
			}

			if addr >= 0 {
				g.cc.Pc = addr
				if mode == 2 {
					limit = int(i.Args[1].(Int))
				}
			}
			i.Args[3] = Int(len(code))
		case InstCode:
			ab := []*Operand{}
			for _, i := range i.Args[1:] {
				ab = append(ab, i.(*Operand))
			}
			for _, i := range g.findInstBody(i, false) {
				seq = append(seq, g.expandBCode(i, ab, g.cc.Pc))
			}
		case InstData:
			p := len(seq)
			w := i.Args[0] == KwWord
			for _, v := range *i.Args[2].(*Vec) {
				for _, v := range _asValueSlice(GetCachedValueIfConst(v, g.cc)) {
					switch v := v.(type) {
					case Int:
						seq = pushdat(seq, int(v), w)
					case *Str:
						seq = append(seq, []byte(*v)...)
					}
				}
			}
			if repeat := int(i.Args[1].(Int)) - 1; repeat > 0 {
				chunk := seq[p:]
				for range repeat {
					seq = append(seq, chunk...)
				}
			}
		case InstDS, InstAlign:
			seq = append(seq, make([]byte, i.size)...)
		case InstFile:
			seq = append(seq, *i.Args[1].(*Binary)...)
		case InstLabel, InstMisc, InstConst: // NOP
		default:
			panic("[BUG] cannot happen")
		}
		g.cc.Pc += i.size
	}

	list := []byte{}
	if g.GenList {
		list = g.generateList(insts[:len(insts)-1], code)
	}
	return code, list
}

package core

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

const (
	asmPassEstimate = iota
	asmPassAdjust
	asmPassVerify
	asmPassCommit
)

func (g *Generator) validateProcTail(inst *Inst) {
	if (inst.Kind == InstCode && g.cc.IsValidProcTail(g.cc, inst)) ||
		inst.IsMisc(KwFallthrough) {
		return
	}
	g.cc.ErrorAt(inst).
		With("the last instruction must be a return/fallthrough within the proc")
}

func (g *Generator) validateFallthrough(insts []*Inst) *Named {
	for x, i := range insts {
		if i.IsMisc(KwEndProc) {
			continue
		} else if i.IsMisc(KwBeginProc) {
			return insts[x+1].Args[0].(*Named)
		}
		g.cc.ErrorAt(i).With("the fallthrough must be followed by a proc")
	}
	panic("[BUG] cannot happen")
}

func (g *Generator) validateCallproc(inst *Inst) {
	c := inst.Args[1].(*Constexpr)
	id := c.Body.(*Identifier)
	nm := g.cc.LookupNamed(c.Env, id)
	if nm == nil {
		g.cc.ErrorAt(inst).With("undefined proc %s", id)
	}

	var a *Sig
	switch nm.Kind {
	case NmLabel:
		a = nm.Value.(*Label).Sig
	case NmInline:
		a = nm.Value.(*Inline).Sig
	}
	if a == nil {
		g.cc.ErrorAt(inst).With("%s is not a proc", id)
	}

	b := inst.Args[2].(*Sig)
	if !a.Equals(b) {
		g.cc.ErrorAt(inst).With("proc signature mismatch: %s.\n"+
			"  expected %v,\n"+
			"  given    %v", id, a, b)
	}
}

func (g *Generator) validateInsts(insts []*Inst) {
	type state struct {
		from  int
		to    int
		inner int
	}

	s := &state{}
	states := []*state{}

	for x, i := range insts {
		switch i.Kind {
		case InstCode:
			s.to = x
		case InstLabel:
			label := GetNamedValue(i.Args[0], LabelT)
			label.Addr = 0
			if !label.LinkedToData() {
				s.to = x
			}
		case InstMisc:
			switch i.Args[0] {
			case KwCallproc:
				g.validateCallproc(i)
			case KwBeginProc:
				if s.inner == 0 {
					s.inner = x
				}
				states, s = append(states, s), &state{from: x, to: x}
			case KwEndProc:
				tail := insts[s.to]
				if s.inner > 0 && s.to > s.inner {
					g.cc.ErrorAt(tail).
						With("the last instruction must be placed before inner procs")
				}

				g.validateProcTail(tail)
				insts[s.from].Args[1] = tail // update begin-proc
				s, states = states[len(states)-1], states[:len(states)-1]
			case KwFallthrough:
				i.Args[1] = g.validateFallthrough(insts[x+1:])
				s.to = x
			}
		}
	}
}

func (g *Generator) findInstBody(inst *Inst, pass int) []BCode {
	etag := inst.ExprTag()
	op := inst.Args[0].(*Keyword)

	m, ok := g.cc.InstMap[op]
	if !ok {
		g.cc.ErrorAt(etag).With("[BUG] unknown instruction: %s", op)
	}

	adjust := pass == asmPassAdjust || pass == asmPassVerify
	p := m.(InstPat)
	for x, i := range inst.Args[1:] {
		i := i.(*Operand)
		if adjust {
			if c, ok := i.A0.(*Constexpr); ok {
				n := EvalConstAs(c, c.Env, IntT, "operand", etag, g.cc)
				g.cc.Constvals[c] = n
				g.cc.AdjustOperand(g.cc, i, int(n), etag)
			}
		}

		p, ok = p[i.Kind].(InstPat)
		if !ok {
			g.cc.ErrorAt(i, etag).With("cannot use %s as operand#%d for %s", i.Kind, x+1, op)
		}
	}

	body, ok := p[nil].(InstDat)
	if !ok {
		g.cc.ErrorAt(etag).With("too few operands for %s", op)
	}

	if pass == asmPassVerify && body[0].Kind == BcTemp {
		g.cc.ErrorAt(etag).With("invalid operands for %s. "+
			"some operand values may be out of range", op)
	}

	if pass != asmPassEstimate {
		body = g.cc.OptimizeBCode(g.cc, inst, body, pass == asmPassCommit)
	}
	return body
}

// Resolve and flatten nested data using the datatype.
// The flattened data contains the following elements:
//
//	*Datatype  change single value size(".byte" or ".word")
//	Int        padding byte size
//	*Constexpr value(int or string)
//
// for examples:
//
//	byte [0 1 2] ==> {ByteType (0) (1) (2)}
//	[4]byte [0 1 2] ==> {ByteType (0) (1) (2) 1}
//	[4]struct{ a byte; b word } [{0 1} {2}] ==> {ByteType (0) WordType (1) ByteType (2) 8}
func (g *Generator) resolveInstData(inst *Inst) {
	acc := []Value{NIL}
	s := &Datatype{}

	checkAndPad := func(size int, next int, etag *Identifier) {
		if size > 0 && g.cc.Pc > next {
			g.cc.ErrorAt(etag).With("too many elements")
		} else if g.cc.Pc < next {
			d := Int(next - g.cc.Pc) // padding size
			if n, ok := acc[len(acc)-1].(Int); ok {
				d += n // merge padding
				acc = acc[:len(acc)-1]
			}
			acc = append(acc, d)
			g.cc.Pc = next
		}
	}

	var traverse func(Value, *Identifier, *Datatype, int, int)
	traverse = func(v Value, etag *Identifier, t *Datatype, size int, limit int) {
		if size != DataSizeSingle {
			e := KwDataList.MatchExpr(v)
			if e == nil {
				g.cc.ErrorAt(v, etag).With("data list required")
			}

			etag := e.At(0).(*Identifier)
			next := g.cc.Pc + t.Size*size
			for _, i := range (*e)[1:] {
				traverse(i, etag, t, DataSizeSingle, limit)
			}
			checkAndPad(size, next, etag)
		} else if t.IsArray() {
			field := t.Fields[0]
			traverse(v, etag, field.Datatype, field.Size, field.Size)
		} else if t.IsStruct() {
			e := KwStructData.MatchExpr(v)
			if e == nil {
				g.cc.ErrorAt(v, etag).With("struct data required")
			}

			etag := e.At(0).(*Identifier)
			next := g.cc.Pc + t.Size
			for x, i := range (*e)[1:] {
				field := t.Fields[x]
				limit := field.Size
				if limit == DataSizeSingle {
					limit = 1
				}
				traverse(i, etag, field.Datatype, field.Size, limit)
			}
			checkAndPad(size, next, etag)
		} else {
			e := CheckValue(v, ConstexprT, "element", etag, g.cc)
			if t != s {
				s = t
				acc = append(acc, s)
			}

			switch v := EvalAndCacheIfConst(e, g.cc).(type) {
			case Int:
				g.cc.Pc += t.Size
			case *Str:
				if t != ByteType {
					g.cc.ErrorAt(v, etag).With("strings are only allowed as byte data")
				}
				if limit > -1 && len(*v) > limit {
					g.cc.ErrorAt(v, etag).With("the string too long")
				}
				g.cc.Pc += len(*v)
			default:
				g.cc.ErrorAt(v, etag).With("invalid data value")
			}
			acc = append(acc, e)
		}
	}

	t := inst.Args[0].(*Datatype)
	size := int(inst.Args[3].(Int))
	traverse(inst.Args[4], inst.ExprTag(), t, size, size)
	inst.Args[2] = NewVec(acc[1:])
}

func (g *Generator) resolveInsts(insts []*Inst, pass int) int {
	g.cc.Pc = 0
	g.changes = 0
	g.cc.initConstvals(PhLink)

	org := 0
	limit := 0
	codeSize := 0
	for _, i := range insts {
		size := i.Size

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
			v := GetNamedValue(i.Args[0], LabelT)
			if v.Addr != g.cc.Pc {
				v.Addr = g.cc.Pc
				g.Changed()
			}
		case InstConst:
			EvalAndCacheIfConst(GetNamedValue(i.Args[0], ConstexprT), g.cc)
		case InstBind:
			EvalAndCacheIfConst(GetNamedValue(i.Args[0], LabelT).At, g.cc)
		case InstOrg:
			if n := g.cc.Pc - org; pass == asmPassVerify && limit > 0 && n > limit {
				g.cc.ErrorAt(i).With("size limit exceeded(%d/%d)", n, limit)
			}
			if addr := int(i.Args[0].(Int)); addr >= 0 {
				org = addr
				limit = int(i.Args[1].(Int))
				g.cc.Pc = org
				g.cc.Org = org
			}
		case InstCode:
			i.Size = len(g.findInstBody(i, pass))
			g.cc.Pc += i.Size
			codeSize += i.Size
		case InstData:
			pc := g.cc.Pc // increment the field directly on each iteration for __PC__
			g.resolveInstData(i)
			i.Size = (g.cc.Pc - pc) * int(i.Args[1].(Int))
			g.cc.Pc = pc + i.Size
		case InstDS:
			n, ok := EvalAndCacheIfConst(i.Args[1], g.cc).(Int)
			if !ok || n < 0 {
				g.cc.ErrorAt(i).With("invalid fill size")
			}
			i.Size = i.Args[0].(*Datatype).Size * int(n)
			g.cc.Pc += i.Size
		case InstAlign:
			n := int(i.Args[0].(Int)) - 1
			i.Size = ((g.cc.Pc + n) & ^n) - g.cc.Pc
			g.cc.Pc += i.Size
		case InstBlob:
			i.Size = len(i.Args[2].(*Blob).data)
			g.cc.Pc += i.Size
		default:
			g.cc.ErrorAt(i).With("[BUG] invalid inst kind %v", i)
		}

		if size != i.Size {
			g.Changed()
		}

		if pass == asmPassVerify && g.changes > 0 {
			g.cc.ErrorAt(i).With("cannot determine the size of the instruction(%d or %d)", size, i.Size)
		}
	}

	return codeSize
}

func (g *Generator) expandBCode(inst *Inst, c BCode, ab []*Operand, pc int) byte {
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
			g.cc.ErrorAt(e).With("the relative address is out of range")
		}
		return byte(d)
	case BcRhigh:
		e := ab[c.A0].A0.(*Constexpr)
		v := g.cc.Constvals[e].(Int)
		d := int(v) - pc + int(int8(c.A1))
		if (c.A2 == 1 && (d < -0x80 || d > 0x7f)) ||
			(c.A2 == 2 && (d < -0x8000 || d > 0x7fff)) {
			g.cc.ErrorAt(e).With("the relative address is out of range")
		}
		return byte(d >> 8)
	case BcImp:
		e := ab[c.A0].A0.(*Constexpr)
		v := g.cc.Constvals[e].(Int)
		base, mask, shift := c.A1, c.A2, c.A3
		if v < 0 || v > Int(mask) {
			g.cc.ErrorAt(e).With("the operand only accepts 0..%d", mask)
		}
		return base | ((byte(v) & mask) << shift)
	case BcMap:
		e := ab[c.A0].A0.(*Constexpr)
		v := g.cc.Constvals[e].(Int)
		m := g.cc.BMaps[c.A1]
		mask, min, max, items := m[0], m[1], m[2], m[3:]
		if v < Int(min) || v > Int(max) {
			g.cc.ErrorAt(e).With("the operand only accepts %d..%d", min, max)
		}
		return items[(byte(v)-min)&mask]
	case BcUnsupported:
		g.cc.ErrorAt(inst).With("unsupported instruction for %s", g.cc.FullArchName())
	}
	panic("[BUG] cannot happen")
}

func (g *Generator) OperandToAsm(e *Operand) string {
	a := g.cc.AsmOperands[e.Kind]
	s := a.Base
	if a.Expand {
		s = strings.Replace(a.Base, "%", g.ValueToAsm(nil, e.A0), 1)
	}
	if g.cc.TrimAsmOperand && s[0] == '0' && s[1] == '+' && s[2] == ' ' && s[3] != '(' {
		s = s[3:]
	}
	return s
}

func (g *Generator) ValueToAsm(env *Env, v Value) string {
	switch v := v.(type) {
	case Int:
		return fmt.Sprintf("%d", v)
	case *Str:
		return v.Inspect()
	case *Blob:
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

		if op.Name == KwField {
			return fmt.Sprintf("%s.%s", a[0], a[1])
		} else if g.cc.Precs[op.Name] > 0 && len(a) == 2 {
			return fmt.Sprintf("(%s %s %s)", a[0], op, a[1])
		}
		return fmt.Sprintf("%s(%s)", op, strings.Join(a, " "))
	case *Operand:
		return g.OperandToAsm(v)
	case *Keyword:
		return v.String()
	case *Identifier:
		if nm := g.cc.LookupNamed(env, v); nm != nil {
			return nm.AsmName.String()
		}
		return v.String()
	}
	g.ErrorAt(v).With("[BUG] invalid value: %T %s", v, v.Inspect())
	return "!"
}

var listPadding = bytes.Repeat([]byte{' '}, 40)

func (g *Generator) generateList(insts []*Inst, code []byte) {
	list := &bytes.Buffer{}
	enabled := true
	issub := g.IsSub
	mode := 1
	pc := 0
	pos := 0
	last := "+"

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
				} else if mode > 0 && m > 0 && !issub {
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
	flushdataif := func(cond bool, a []string, t *Datatype) []string {
		if cond {
			m := t.Size * len(a)
			writeh(m, m, "    .%s %s", t.Name, strings.Join(a, ", "))
			a = a[:0]
		}
		return a
	}

	if issub {
		writes("")
		writes("    ; ----------------")
	} else {
		writes("    ; generated by ocala")
	}
	for _, i := range insts {
		switch i.Kind {
		case InstMisc:
			switch i.Args[0] {
			case KwListConstants:
				enabled = int(i.Args[1].(Int)) != 0
			case KwComment:
				if comment, ok := i.Args[1].(*Str); ok {
					s := ""
					for _, i := range i.Args[2:] {
						i := GetCachedValueIfConst(i, g.cc)
						s += " " + g.ValueToAsm(nil, i)
					}
					writes("    ; %s%s", comment.String(), s)
				}
			case KwBeginProc, KwEndProc:
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
			if !strings.HasPrefix(s, ".__") {
				writeh(0, 0, "%s:", s)
			}
		case InstConst:
			nm := i.Args[0].(*Named)
			a := g.ValueToAsm(nil, nm.AsmName)
			b := g.ValueToAsm(nm.Env, nm.Value)
			if !strings.HasPrefix(a, ".__") {
				writes("    %s = %s", a, b)
			}
		case InstBind:
			nm := i.Args[0].(*Named)
			a := g.ValueToAsm(nil, nm.AsmName)
			b := g.ValueToAsm(nil, nm.Value.(*Label).At)
			writes("    %s = %s", a, b)
		case InstCode:
			kw := i.Args[0]
			comment := ""
			if i.Comment != "" {
				comment = " ; " + i.Comment
			}
			if len(i.Args) == 1 {
				writeh(i.Size, i.Size, "    %s%s", kw, comment)
			} else {
				s := []string{}
				for _, a := range i.Args[1:] {
					s = append(s, g.ValueToAsm(nil, a))
				}
				writeh(i.Size, i.Size, "    %-6s %s%s", kw, strings.Join(s, ", "), comment)
			}
		case InstData:
			p := pc
			t := ByteType
			a := []string{}
			for _, i := range *i.Args[2].(*Vec) {
				switch v := i.(type) {
				case Int:
					a = flushdataif(len(a) > 0, a, t)
					writeh(0, int(v), "    .defb %d", v)
				case *Datatype:
					a = flushdataif(len(a) > 0, a, t)
					t = v
				case *Constexpr:
					s := g.ValueToAsm(nil, i)
					switch i := GetCachedValueIfConst(i, g.cc).(type) {
					case Int:
						a = append(a, s)
						a = flushdataif(len(a)*t.Size == 8, a, t)
					case *Str:
						a = flushdataif(len(a) > 0, a, t)
						writeh(len(*i), len(*i), "    .byte "+s)
					}
				}
			}
			flushdataif(len(a) > 0, a, t)

			if repeat := int(i.Args[1].(Int)); repeat > 1 {
				m := pc - p
				for i := 1; i < repeat; i++ {
					writeh(m, m, "    ; ... repeat %d/%d", i+1, repeat)
				}
			}
		case InstDS:
			if i.Args[0] == WordType {
				writeh(0, i.Size, "    .defw %d", int(i.Size/2))
			} else {
				writeh(0, i.Size, "    .defb %d", int(i.Size))
			}
		case InstAlign:
			writeh(0, i.Size, "    .align %d ; (.defb %d)", int(i.Args[0].(Int)), i.Size)
		case InstBlob:
			blob := i.Args[2].(*Blob)
			s := g.ValueToAsm(nil, NewStr(blob.origPath))
			if blob.compiled {
				s = "\"(compiled):" + s[1:]
			}
			writeh(0, len(blob.data), "    .incbin %s", s)
		default:
			panic("[BUG] cannot happen")
		}
	}

	g.prependList(list.Bytes())
}

func (g *Generator) prepareToGenerateBin(insts []*Inst) {
	g.validateInsts(insts)
	g.cc.optimizeFlow(insts)
	g.cc.optimizeLink(insts)

	oldCodeSize := g.resolveInsts(insts, asmPassEstimate)
	codeSizes := []int{oldCodeSize}
	for {
		codeSize := g.resolveInsts(insts, asmPassAdjust)
		codeSizes = append(codeSizes, codeSize)
		if codeSize >= oldCodeSize {
			break
		}
		oldCodeSize = codeSize
	}

	if g.DebugMode && g.cc.OptimizeBCode != nil {
		v := []string{}
		for _, i := range codeSizes {
			v = append(v, strconv.Itoa(i))
		}
		fmt.Fprintln(g.ErrWriter, "code size:", strings.Join(v, " -> "))
	}
	g.resolveInsts(insts, asmPassVerify) // verify
}

func (g *Generator) GenerateBin(insts []*Inst) []byte {
	g.prepareToGenerateBin(insts)

	pushdat := func(s []byte, v int, t *Datatype) []byte {
		if t == WordType {
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
				g.cc.ErrorAt(i).With("%s", i.Args[1].(*Str))
			}
		case InstOrg:
			// mode -- 0: noload, 1: load, 2: fill, 3: +noload, 4: +load
			if mode > 0 {
				code = append(code, seq...)
			}
			seq = []byte{}

			size := len(code) - anchor
			if limit > 0 && size > limit {
				g.cc.ErrorAt(i).With("size limit exceeded(%d/%d)", size, limit)
			}

			addr := int(i.Args[0].(Int))
			mode = int(i.Args[2].(Int))
			if mode < 0 || mode >= 5 {
				g.cc.ErrorAt(i).With("invalid org mode")
			} else if mode >= 3 { // mode 3/4
				mode -= 3
				if limit == 0 {
					g.cc.ErrorAt(i).With("a pack mode must follow a fill/pack mode")
				}
			} else { // mode 0/1/2
				if limit > 0 {
					code = append(code, make([]byte, limit-size)...)
					if addr < 0 {
						g.cc.ErrorAt(i).With("the orgin address required")
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
			for _, bcode := range g.findInstBody(i, asmPassCommit) {
				seq = append(seq, g.expandBCode(i, bcode, ab, g.cc.Pc))
			}
		case InstData:
			p := len(seq)
			t := ByteType
			for _, v := range *i.Args[2].(*Vec) {
				switch v := v.(type) {
				case Int:
					seq = append(seq, make([]byte, int(v))...)
				case *Datatype:
					t = v
				case *Constexpr:
					switch v := GetCachedValueIfConst(v, g.cc).(type) {
					case Int:
						seq = pushdat(seq, int(v), t)
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
			seq = append(seq, make([]byte, i.Size)...)
		case InstBlob:
			seq = append(seq, i.Args[2].(*Blob).data...)
		case InstLabel, InstMisc, InstConst, InstBind: // NOP
		default:
			panic("[BUG] cannot happen")
		}
		g.cc.Pc += i.Size
	}

	if g.GenList {
		g.generateList(insts[:len(insts)-1], code)
	}
	return code
}

package core

import (
	"fmt"
	"maps"
	"runtime"
	"sort"
	"strconv"
)

var KwUNDER = Intern("_")
var KwFILENAME = Intern("__FILE__")
var KwPROCNAME = Intern("__PROC__")
var KwCURLOC = Intern("__PC__")
var KwCURORG = Intern("__ORG__")
var KwUBEG = Intern("_BEG")
var KwUEND = Intern("_END")
var KwUCOND = Intern("_COND")

var KwDot = Intern(".")
var KwDotDot = Intern("..")
var KwRest = Intern("...")
var KwTEXT = Intern("text")
var KwBSS = Intern("bss")
var KwRODATA = Intern("rodata")

var KwAny = Intern("#.ANY")
var KwLabel = Intern("#.label")
var KwModule = Intern("#.module")
var KwMacro = Intern("#.macro")
var KwProc = Intern("#.proc")
var KwConst = Intern("#.const")
var KwData = Intern("#.data")
var KwMem = Intern("#.mem")
var KwJump = Intern("#.jump")
var KwJumpIf = Intern("-jump-if")
var KwJumpUnless = Intern("-jump-unless")
var KwTpl = Intern("#.tpl")
var KwVec = Intern("#.vec")
var KwProg = Intern("#.prog")
var KwBlock = Intern("#.block")
var KwWith = Intern("#.with")
var KwCall = Intern("#.call")
var KwCallproc = Intern("#.callproc")
var KwUndefined = Intern("#.undefined")
var KwEndProc = Intern("#.endproc")
var KwFallthrough = Intern("#.fallthrough")
var KwValueOf = Intern("#.valueof")
var KwExprdata = Intern("#.exprdata")
var KwEndInline = Intern("#.endinline")
var KwField = Intern("#.field")
var KwSetIota = Intern("#.set-iota")
var KwInvalidExpansion = Intern("#.invalid-expansion")
var KwMakeId = Intern("#.make-id")
var KwInline = Intern("#.inline")

var KwLeftArrow = Intern("<-")
var KwPlusOp = Intern("+")
var KwMulOp = Intern("*")
var KwDivOp = Intern("/")
var KwLogicalNotOp = Intern("!")
var KwNotOp = Intern("~")
var KwToplevel = Intern(":")
var KwInclude = Intern("include")
var KwDo = Intern("do")
var KwQuote = Intern("quote")
var KwByte = Intern("byte")
var KwWord = Intern("word")
var KwRecord = Intern("record")
var KwIf = Intern("if")
var KwElse = Intern("else")
var KwOrg = Intern("org")
var KwMerge = Intern("merge")
var KwArch = Intern("arch")
var KwListConstants = Intern("list-constants")
var KwComment = Intern("comment")
var KwReserved = Intern("<reserved>")
var KwWhen = Intern("when")
var KwWhenAnd = Intern("&&-")
var KwWhenOr = Intern("||-")
var KwThen = Intern("then")
var KwOptimize = Intern("optimize")
var KwVolatile = Intern("<volatile>")
var KwCompileFile = Intern("compile-file")

var IdUNDER = InternalId(KwUNDER)

var binOps = map[string]int{
	"*":   3,
	"/":   3,
	"%":   3,
	"+":   4,
	"-":   4,
	"<<":  5,
	">>":  5,
	">>>": 5,
	"<":   6,
	"<=":  6,
	">":   6,
	">=":  6,
	"==":  7,
	"!=":  7,
	"&":   8,
	"|":   9,
	"^":   10,
	"&&":  11,
	"||":  12,
	"#":   101,
}

const (
	CtModule    = 0b001
	CtProc      = 0b010
	CtConstexpr = 0b100
)

const (
	PhCompile = iota
	PhLink
)

var PhaseLabels = []string{"compile", "link"}

type Compiler struct {
	Arch      string
	Toplevel  *Env
	Builtins  *Env
	Contexts  []byte
	CodeStack [][]*Inst
	Module    *Module
	Section   *Section
	InPath    string
	MainPath  string
	Constvals map[*Constexpr]Value
	Pc        int
	Org       int
	Phase     int
	link      *Vec
	loaded    []string
	g         *Generator

	hooks struct {
		beforeLink func(*Compiler)
	}

	InstMap     map[*Keyword]InstTab
	InstAliases map[string][]string
	CtxOpMap    map[*Keyword]map[*Keyword]map[*Keyword][][]Value
	SyntaxMap   map[*Keyword]SyntaxFn
	FunMap      map[*Keyword]SyntaxFn
	CollectRegs func([]*Keyword, *Keyword) []*Keyword

	ExprToOperand   func(*Compiler, Value) *Operand
	AdjustOperand   func(*Compiler, *Operand, *Identifier)
	IsValidProcTail func(*Compiler, *Inst) bool
	AdjustInline    func(*Compiler, []*Inst)
	OperandToAsm    func(*Generator, *Operand) string

	BMaps         [][]byte
	KwRegA        *Keyword
	Precs         map[*Keyword]int
	Operators     map[string]int
	TokenWords    [][]string
	TokenAliases  map[string]string
	ReservedWords map[string]int32
	MacroNesting  int
	InlineInsts   []*Inst
}

func (cc *Compiler) RaiseCompileError(id *Identifier, message string, args ...any) {
	var token *Token
	if id != nil {
		token = id.Token
	}

	for i := 1; ; i++ {
		if _, file, line, ok := runtime.Caller(i); ok && cc.g.DebugMode {
			message += fmt.Sprintf("\n-- %s:%d", file, line)
			continue
		}
		break
	}
	cc.g.raiseError(token, "compile error: ", message, args...)
}

func (cc *Compiler) EnterCodeBlock() {
	cc.CodeStack = append(cc.CodeStack, []*Inst{})
}

func (cc *Compiler) LeaveCodeBlock() []*Inst {
	insts := cc.CodeStack[len(cc.CodeStack)-1]
	cc.CodeStack = cc.CodeStack[:len(cc.CodeStack)-1]
	return insts
}

func (cc *Compiler) EnterContext(ct byte) {
	cc.Contexts = append(cc.Contexts, ct)
}

func (cc *Compiler) LeaveContext() {
	cc.Contexts = cc.Contexts[:len(cc.Contexts)-1]
}

func (cc *Compiler) Context() byte {
	return cc.Contexts[len(cc.Contexts)-1]
}

func (cc *Compiler) EmitCode(code ...*Inst) Value {
	cc.CodeStack[len(cc.CodeStack)-1] = append(cc.CodeStack[len(cc.CodeStack)-1], code...)
	return NIL
}

func (cc *Compiler) EmitCodeTo(x int, code ...*Inst) Value {
	cc.CodeStack[x] = append(cc.CodeStack[x], code...)
	return NIL
}

func (cc *Compiler) EmitCodeByEnv(env *Env, inst *Inst) Value {
	if env.Module != nil {
		return cc.EmitCodeTo(0, inst)
	}
	return cc.EmitCode(inst)
}

func (cc *Compiler) EmitCodeToSection(section *Section, code ...*Inst) {
	section.Insts = append(section.Insts, code...)
}

func (cc *Compiler) resolveNamespace(env *Env, id *Identifier) *Env {
	if id.Namespace != nil {
		mod := cc.FindModule(env, id.Namespace, id)
		if mod == nil {
			return nil
		}
		return mod.Env
	}
	return env
}

func (cc *Compiler) FindModule(env *Env, kw *Keyword, id *Identifier) *Module {
	nm := cc.Toplevel.Lookup(kw)
	if nm == nil {
		return nil
	}
	if nm.Kind != NmModule {
		cc.RaiseCompileError(id, "%s is not a namespace", kw.String())
	}
	return nm.Value.(*Module)
}

func (cc *Compiler) FindNamed(env *Env, id *Identifier) *Named {
	env = cc.resolveNamespace(env, id)
	if env == nil {
		return nil
	}
	return env.FindById(id)
}

func (cc *Compiler) LookupNamed(env *Env, id *Identifier) *Named {
	if id.Namespace != nil {
		return cc.FindNamed(env, id)
	}

	env = cc.resolveNamespace(env, id)
	if env == nil {
		return nil
	}
	return env.LookupById(id)
}

func (cc *Compiler) InstallNamed(env *Env, id *Identifier, kind int32, value Value) *Named {
	if id.Namespace != nil {
		cc.RaiseCompileError(id, "invalid namespace %s", id.Namespace)
	}
	if env.FindById(id) != nil {
		cc.RaiseCompileError(id, "%s is already defined", id)
	}
	if cc.Builtins.FindById(id) != nil {
		cc.RaiseCompileError(id, "%s is a builtin name", id)
	}

	nm := &Named{Token: id.Token, Name: id.Name, Kind: kind, Value: value}
	if env.Module != nil {
		s := id.String()
		if env != cc.Toplevel {
			s = fmt.Sprintf("%s:%s", env.Module.Name, s)
		}
		nm.AsmName = NewKeyword(s)
		nm.Export = true
	}
	return env.Install(nm)
}

func (cc *Compiler) installBuiltins(env *Env, kind int32, tab map[*Keyword]SyntaxFn) {
	for k, v := range tab {
		env.Install(&Named{Name: k, Kind: kind, Value: v})
	}
}

func (cc *Compiler) installInsts(env *Env) {
	for k, v := range cc.InstAliases {
		k := Intern(k)
		for _, i := range v {
			cc.InstMap[Intern(i)] = cc.InstMap[k]
		}
	}

	for k := range cc.InstMap {
		env.Install(&Named{Name: k, Kind: NmInst, Value: NIL})
	}
}

func (cc *Compiler) initHooks() {
	cc.hooks.beforeLink = func(_ *Compiler) {}
}

func (cc *Compiler) initMainPath(path string) {
	cc.MainPath = path
	cc.loaded = []string{path}
}

func (cc *Compiler) initTopLevelEnv() {
	env := NewEnv(nil)

	cc.Builtins = env
	cc.installInsts(env)
	cc.installBuiltins(env, NmSyntax, SyntaxMap)
	cc.installBuiltins(env, NmSyntax, cc.SyntaxMap)
	cc.installBuiltins(env, NmFunc, FunMap)
	cc.installBuiltins(env, NmFunc, cc.FunMap)
	cc.installBuiltins(env, NmSpecial, SpecialMap)
	env.Install(&Named{Name: KwCompileFile, Kind: NmSyntax, Value: SyntaxFn(sCompileFile)})

	cc.Toplevel = env.Enter()
	cc.Module = NewModule(KwToplevel, cc.Toplevel)
	cc.Module.Env.Module = cc.Module
	cc.Section = cc.Module.Sections[KwBSS]
	cc.Toplevel.Install(&Named{Name: KwToplevel, Kind: NmModule, Value: cc.Module})
}

func (cc *Compiler) initConstvals(phase int) {
	cc.Constvals = map[*Constexpr]Value{}
	cc.Phase = phase
}

func (cc *Compiler) initReservedWords() {
	cc.ReservedWords = map[string]int32{}
	maps.Copy(cc.ReservedWords, reservedWords)
	for _, i := range cc.TokenWords[0] {
		cc.ReservedWords[i] = tkREG
	}
	for _, i := range cc.TokenWords[1] {
		cc.ReservedWords[i] = tkCOND
	}

	cc.Operators = map[string]int{}
	for _, i := range cc.TokenWords[2] {
		cc.Operators[i] |= 1
	}

	cc.Precs = map[*Keyword]int{}
	for k, v := range binOps {
		cc.Operators[k] |= 2
		cc.Precs[Intern(k)] = v
	}
	for _, i := range cc.TokenWords[3] {
		cc.Operators[i] |= 2
	}
}

func (cc *Compiler) buildSetupCode() *Vec {
	parser := &Scanner{Path: "<-D>", Text: []byte("(constants defined by -D options)")}
	tag := &Identifier{Name: KwConst}
	tag.Token = &Token{From: parser, Value: tag}

	r := &Vec{tag.Expand(KwBlock)}
	for _, i := range cc.g.Defs {
		id := &Identifier{Name: Intern(i)}
		id.Token = &Token{From: parser, Value: id}
		e := &Vec{tag.Expand(KwConst), id.ToConstexpr(nil), &Vec{}, InternalConstexpr(Int(1))}
		r.Push(e)
	}
	return r
}

func (cc *Compiler) Compile(path string, text []byte) []*Inst {
	cc.initHooks()
	cc.initMainPath(path)
	cc.initTopLevelEnv()
	cc.initConstvals(PhCompile)
	cc.initReservedWords()
	cc.EnterContext(CtModule)
	cc.EnterCodeBlock() // Constants
	cc.EnterCodeBlock() // Toplevel
	cc.CompileExpr(cc.Toplevel, cc.buildSetupCode())
	cc.CompileExpr(cc.Toplevel, cc.Parse(path, text))
	return cc.doLink()
}

func (cc *Compiler) CompileIncluded(path string, text []byte) Value {
	inPath, module, section := cc.InPath, cc.Module, cc.Section
	cc.InPath = path
	cc.Module = cc.Toplevel.Module
	cc.Section = cc.Module.Sections[KwBSS]
	cc.EnterContext(CtModule)
	cc.EnterCodeBlock()
	cc.CompileExpr(cc.Toplevel, cc.Parse(path, text))
	cc.EmitCodeToSection(cc.Section, cc.LeaveCodeBlock()...)
	cc.LeaveContext()
	cc.InPath, cc.Module, cc.Section = inPath, module, section
	return NIL
}

func (cc *Compiler) Parse(path string, text []byte) *Vec {
	cc.InPath = path
	p := &Parser{Scanner: Scanner{Path: path, Text: text, cc: cc}, contexts: []byte{'{'}}
	res, _ := p._parse()
	return res.(*Vec)
}

func (cc *Compiler) GetOptimizer() *Optimizer {
	return &cc.g.Optimizer
}

//

func CheckExpr(e *Vec, min, max int, contexts byte, cc *Compiler) (*Identifier, int) {
	etag := e.ExprTag()
	if etag == nil {
		cc.RaiseCompileError(GetTokenAsId(e), "invalid form")
	} else if contexts != 0 && (cc.Context()&contexts) == 0 {
		cc.RaiseCompileError(etag, "%s is not allowed in this context", etag)
	}

	argc := e.Size()
	if argc < min || (max > -1 && argc > max) {
		s := strconv.Itoa(min - 1)
		if max < 0 {
			s += "+"
		} else if max > min {
			s += ".." + strconv.Itoa(max-1)
		}
		cc.RaiseCompileError(etag, "%s: %s argument(s) required, but given %d", etag, s, argc-1)
	}

	return etag, argc
}

func CheckValue[T Value](v Value, t T, name string, id *Identifier, cc *Compiler) T {
	if w, ok := v.(T); ok {
		return w
	}
	cc.RaiseCompileError(id, "%s is must be %s", name, TypeLabels[t])
	return t
}

func CheckConst[T Value](v Value, t T, name string, id *Identifier, cc *Compiler) T {
	a := CheckValue(v, ConstexprT, name, id, cc)
	return CheckValue(a.Body, t, name, id, cc)
}

func EvalConstAs[T Value](v Value, env *Env, t T, name string, id *Identifier, cc *Compiler) T {
	a := EvalConst(v, env, id, cc)
	return CheckValue(a, t, name, id, cc)
}

func EvalConst(v Value, env *Env, id *Identifier, cc *Compiler) Value {
	if cc.Context()&CtConstexpr != 0 {
		return cc.evaluateConstexpr(env, v, id.Token)
	} else if v, ok := v.(*Constexpr); ok {
		cc.EnterContext(CtConstexpr)
		r := cc.evaluateConstexpr(env, v.Body, v.Token)
		cc.LeaveContext()
		return r
	}
	cc.RaiseCompileError(nil, "[BUG] invalid context for the constexpr, %s", v.Inspect())
	return NIL
}

func CheckAndEvalConstAs[T Value](v Value, env *Env, t T, name string, id *Identifier, cc *Compiler) T {
	a := CheckValue(v, ConstexprT, name, id, cc)
	return EvalConstAs(a, env, t, name, id, cc)
}

func CheckAndEvalConst(v Value, env *Env, name string, id *Identifier, cc *Compiler) Value {
	CheckValue(v, ConstexprT, name, id, cc)
	return EvalConst(v, env, id, cc)
}

func EvalAndCacheIfConst(v Value, cc *Compiler) Value {
	if e, ok := v.(*Constexpr); ok {
		v = EvalConst(e, e.Env, nil, cc)
		cc.Constvals[e] = v
	}
	return v
}

func GetCachedValueIfConst(v Value, cc *Compiler) Value {
	if e, ok := v.(*Constexpr); ok {
		return cc.Constvals[e]
	}
	return v
}

func GetConstBody(v Value) Value {
	return v.(*Constexpr).Body
}

func GetNamedValue[T Value](v Value, t T) T {
	return v.(*Named).Value.(T)
}

func GetTokenAsId(e Value) *Identifier {
	switch e := e.(type) {
	case *Vec:
		return GetTokenAsId(e.AtOrUndef(0))
	case *Identifier:
		return e
	case *Constexpr:
		return e.Token.TagId()
	}
	return nil
}

func CheckToplevelEnv(env *Env, etag *Identifier, cc *Compiler) {
	if env != cc.Toplevel {
		cc.RaiseCompileError(etag, "%s must be in toplevel", etag.String())
	}
}

func CheckToplevelEnvIfCtProc(env *Env, etag *Identifier, cc *Compiler) {
	if (cc.Context() & CtProc) != 0 {
		CheckToplevelEnv(env, etag, cc)
	}
}

func CheckReserved(nm *Named, id *Identifier, cc *Compiler) {
	if nm == nil || nm.Kind != NmLabel {
		cc.RaiseCompileError(id, "unknown label: %s", id)
	}
	if e, ok := nm.Value.(*Label); !ok || !e.IsReserved() {
		cc.RaiseCompileError(id, "%s already defined", id)
	}
}

func CheckConstPlainId(v Value, name string, id *Identifier, cc *Compiler) *Identifier {
	a := CheckConst(v, IdentifierT, name, id, cc)
	if a.Namespace != nil {
		cc.RaiseCompileError(id, "qualified name is not allowed in this context")
	}
	return a
}

func AsTaggedVec(v Value) (*Identifier, *Vec) {
	if v, ok := v.(*Vec); ok {
		if id := v.ExprTag(); id != nil {
			return id, v
		}
	}
	return nil, nil
}

func AsMemForm(v Value) *Vec {
	if etag, v := AsTaggedVec(v); etag != nil && etag.Namespace == nil && etag.Name == KwMem {
		return v
	}
	return nil
}

func AsBlockForm(v Value) *Vec {
	if etag, v := AsTaggedVec(v); etag != nil && etag.Namespace == nil &&
		(etag.Name == KwProg || etag.Name == KwBlock) {
		return v
	}
	return nil
}

func CheckBlockForm(v Value, name string, id *Identifier, cc *Compiler) *Vec {
	if v := AsBlockForm(v); v != nil {
		return v
	}
	cc.RaiseCompileError(id, "%s is must be block-form", name)
	return nil
}

func CheckPhase(phase int, id *Identifier, cc *Compiler) {
	if cc.Phase == phase {
		return
	}
	cc.RaiseCompileError(id, "cannot use %s in %s phase", id.String(), PhaseLabels[cc.Phase])
}

//

func (cc *Compiler) CompileExpr(env *Env, e Value) Value {
	switch e := e.(type) {
	case *Constexpr:
		if e.Env != nil && e.Env != env {
			cc.RaiseCompileError(e.Token.TagId(), "[BUG] constexpr compiled twice")
		}
		e.Env = env
	case *Vec:
		etag := e.ExprTag()
		if etag == nil {
			cc.RaiseCompileError(GetTokenAsId(e), "invalid form")
		}

		nm := cc.LookupNamed(env, etag)
		if nm == nil {
			cc.RaiseCompileError(etag, "unknown form name %s", etag)
		}

		switch nm.Kind {
		case NmInst:
			return cc.compileInst(env, e)
		case NmSyntax:
			return nm.Value.(SyntaxFn)(cc, env, e)
		case NmMacro:
			return cc.expandMacro(env, e, nm.Value.(*Macro))
		default:
			cc.RaiseCompileError(etag, "invalid form %s(%s)", etag, NamedKindLabels[nm.Kind])
		}
	}
	return e
}

func (cc *Compiler) compileInst(env *Env, e *Vec) Value {
	op, _ := CheckExpr(e, 1, -1, CtProc, cc)

	args := []Value{op.Name}
	for _, i := range (*e)[1:] {
		i := cc.ExprToOperand(cc, cc.CompileExpr(env, i))
		args = append(args, i)
	}
	return cc.EmitCode(NewInst(e, InstCode, args...))
}

func (cc *Compiler) unwrapExprdata(v Value, etag *Identifier) Value {
	if v := KwExprdata.MatchExpr(v); v != nil && v.Size() == 2 {
		return v.At(1)
	}
	cc.RaiseCompileError(etag, "the expression itself required. use the %%& placeholder")
	return NIL
}

func (cc *Compiler) unwrapConstexprIf(cond bool, v Value, etag *Identifier) Value {
	if cond {
		if v, ok := v.(*Constexpr); ok {
			return v.Body
		}
		cc.RaiseCompileError(etag, "the placeholder only accepts constant expressions")
	}
	return v
}

func (cc *Compiler) expandMacroBody(env *Env, e Value, r *Vec, unwrap bool, by *Identifier) *Vec {
	switch e := e.(type) {
	case *Identifier:
		e = e.Clone()
		e.ExpandedBy = by

		if e.PlaceHolder == "" {
			return r.Push(e)
		} else if e.PlaceHolder[:2] == "%%" {
			e.PlaceHolder = e.PlaceHolder[1:]
			return r.Push(e)
		}

		nm := cc.FindNamed(env, e)
		if nm == nil {
			cc.RaiseCompileError(e, "unknown placeholder %s in macro body", e)
		}

		v := nm.Value.Dup()
		if e.PlaceHolder == "%=" {
			return r.Push(cc.unwrapConstexprIf(unwrap, v, e))
		} else if e.PlaceHolder == "%&" {
			if !unwrap {
				cc.RaiseCompileError(e, "the %%& placeholder is only allowed within constant expressions")
			}
			return r.Push(&Vec{by.Expand(KwExprdata), v})
		} else if nm.Special {
			v := []Value(*v.(*Vec))
			if e.PlaceHolder == "%#" {
				v := &Constexpr{Token: e.Token, Body: Int(len(v))}
				return r.Push(cc.unwrapConstexprIf(unwrap, v, e))
			}

			for _, i := range e.PlaceHolder[1:] {
				switch i {
				case '*':
					for _, i := range v {
						r.Push(cc.unwrapConstexprIf(unwrap, i, e))
					}
					return r
				case '<':
					if len(v) == 0 {
						v := &Constexpr{Token: e.Token, Body: &Vec{e.Expand(KwInvalidExpansion)}}
						return r.Push(cc.unwrapConstexprIf(unwrap, v, e))
					}
					return r.Push(cc.unwrapConstexprIf(unwrap, v[0], e))
				case '>':
					if len(v) == 0 {
						v := &Constexpr{Token: e.Token, Body: &Vec{e.Expand(KwInvalidExpansion)}}
						return r.Push(cc.unwrapConstexprIf(unwrap, v, e))
					}
					v = v[1:]
				}
			}

			if unwrap {
				cc.RaiseCompileError(e, "the non-constant value cannot be expanded in the constant expression")
			}
			return r.Push(NewVec(v))
		}
		cc.RaiseCompileError(e, "vector operations only allowed on the rest parameter")
	case *Vec:
		v := &Vec{}
		for _, i := range *e {
			cc.expandMacroBody(env, i, v, unwrap, by)
		}
		if r == nil {
			return v
		}
		return r.Push(v)
	case *Constexpr:
		return r.Push(&Constexpr{
			Token: CopyPtr(e.Token),
			Body:  cc.expandMacroBody(env, e.Body, &Vec{}, true, by).At(0),
		})
	}
	return r.Push(e)
}

func (cc *Compiler) expandMacro(env *Env, e *Vec, mac *Macro) Value {
	etag := e.At(0).(*Identifier)
	tenv := NewEnv(env)

	if cc.MacroNesting > 64 {
		cc.RaiseCompileError(etag, "macro expansion too deep")
	}
	cc.MacroNesting++

	as := (*e)[1:]
	m := len(as)
	n := len(mac.Args)
	d := n - len(mac.Opts)
	x := 0
	if mac.Rest {
		CheckExpr(e, d+1, -1, 0, cc)

		for ; x < m && x < n-1; x++ {
			tenv.Install(&Named{Name: mac.Args[x], Kind: NmVar, Value: as[x]})
		}
		tenv.Install(&Named{Name: mac.Args[n-1], Kind: NmVar, Special: true, Value: NewVec(as[x:])})

		for ; x < n-1; x++ {
			tenv.Install(&Named{Name: mac.Args[x], Kind: NmVar, Value: mac.Opts[x-d].Dup()})
		}
	} else {
		CheckExpr(e, d+1, n+1, 0, cc)

		for ; x < m; x++ {
			tenv.Install(&Named{Name: mac.Args[x], Kind: NmVar, Value: as[x]})
		}
		for ; x < n; x++ {
			tenv.Install(&Named{Name: mac.Args[x], Kind: NmVar, Value: mac.Opts[x-d].Dup()})
		}
	}

	for x, n := 0, len(mac.Vars); x < n; x += 2 {
		k := mac.Vars[x].(*Keyword)
		v := mac.Vars[x+1]
		if v == NIL {
			v = Gensym(k.String()).ToId(etag.Token).ToConstexpr(nil)
		} else {
			v = &Constexpr{Token: etag.Token, Body: EvalConst(v, tenv, etag, cc)}
		}
		tenv.Install(&Named{Name: k, Kind: NmVar, Value: v})
	}
	r := cc.expandMacroBody(tenv, mac.Body, nil, false, etag)

	cc.CompileExpr(env, r)
	cc.MacroNesting--
	return NIL
}

func (cc *Compiler) expandInline(env *Env, e *Vec, id *Identifier, inline *Inline) Value {
	etag := e.At(0).(*Identifier)

	env = env.Enter()
	cc.InstallNamed(env, etag.Expand(KwPROCNAME), NmInvalid, NIL)
	nm := cc.InstallNamed(env, etag.Expand(KwEndInline), NmLabel, &Label{})
	r := cc.expandMacroBody(NewEnv(nil), inline.Body, nil, false, etag)

	cc.EnterCodeBlock()
	// at least one element required in the body.
	cc.EmitCode(NewInst(e, InstMisc, KwComment, NewStr("begin-inline "+id.String())))
	cc.CompileExpr(env, r)
	cc.AdjustInline(cc, cc.CodeStack[len(cc.CodeStack)-1])
	cc.EmitCode(NewInst(e, InstLabel, nm))
	cc.EmitCode(cc.LeaveCodeBlock()...)
	return NIL
}

func (cc *Compiler) evaluateConstexpr(env *Env, e Value, token *Token) Value {
	switch e := e.(type) {
	case Int, *Str, *Keyword, *Blob, *Nil:
		return e
	case *Vec:
		id, _ := CheckExpr(e, 1, -1, CtConstexpr, cc)
		nm := cc.LookupNamed(env, id)
		if nm == nil || nm.Value == nil {
			cc.RaiseCompileError(id, "unknown operator %s", id)
		}

		switch fn := nm.Value.(type) {
		case *ConstFn:
			tenv := NewEnv(fn.Env)
			as := (*e)[1:]
			m := len(as)
			n := len(fn.Args)
			d := n - len(fn.Opts)
			x := 0
			CheckExpr(e, d+1, n+1, 0, cc)

			for ; x < m; x++ {
				v := cc.evaluateConstexpr(env, as[x], token)
				tenv.Install(&Named{Name: fn.Args[x], Kind: NmVar, Value: v})
			}
			for ; x < n; x++ {
				i := fn.Opts[x-d].(*Constexpr)
				v := cc.evaluateConstexpr(tenv, i.Body, i.Token)
				tenv.Install(&Named{Name: fn.Args[x], Kind: NmVar, Value: v})
			}
			return cc.evaluateConstexpr(tenv, fn.Body, token)
		case SyntaxFn:
			switch nm.Kind {
			case NmSyntax:
				return fn(cc, env, e)
			case NmFunc:
				args := &Vec{id}
				for _, i := range (*e)[1:] {
					args.Push(cc.evaluateConstexpr(env, i, token))
				}
				return fn(cc, env, args)
			}
		}
		cc.RaiseCompileError(id, "%s is not callable", id)
	case *Identifier:
		nm := cc.LookupNamed(env, e)
		if nm == nil {
			cc.RaiseCompileError(e, "undefined name %s", e)
		}

		switch nm.Kind {
		case NmVar:
			return nm.Value
		case NmLabel:
			if cc.Phase == PhCompile {
				cc.RaiseCompileError(e, "cannot use label address in compile phase")
			}
			v := nm.Value.(*Label)
			if v.At == nil {
				return Int(v.Addr)
			}

			w := cc.Constvals[v.At]
			if w == nil {
				cc.RaiseCompileError(e, "label %s used before declaration", e)
			}
			return w
		case NmSpecial:
			if v, ok := nm.Value.(SyntaxFn); ok {
				return v(cc, env, &Vec{e})
			}
			cc.RaiseCompileError(e, "[BUG] invalid special variable %s", e)
		case NmConst:
			v := nm.Value.(*Constexpr)
			w := cc.Constvals[v]
			if w == nil {
				if cc.Phase == PhCompile {
					w = cc.evaluateConstexpr(v.Env, v.Body, v.Token)
					cc.Constvals[v] = w
				} else {
					cc.RaiseCompileError(e, "constant %s used before declaration", e)
				}
			}
			return w
		default:
			cc.RaiseCompileError(e, "cannot use the %s `%s` within this context", NamedKindLabels[nm.Kind], e)
		}
	}
	cc.RaiseCompileError(token.TagId(), "[BUG] invalid constexpr: %s[%T]", e.Inspect(), e)
	return NIL
}

var idOpLast = InternalId(Intern("#"))

func (cc *Compiler) orderByPrec(s []Value, right *Identifier, v Value) []Value {
	rp, ok := cc.Precs[right.Name]
	if !ok || rp == 0 {
		cc.RaiseCompileError(right.Token.TagId(), "[BUG] Unknown operator '%s'", right.Token)
	}

	for len(s) >= 3 {
		left := s[len(s)-2].(*Identifier)
		lp := cc.Precs[left.Name]

		if lp <= rp {
			s = append(s[:len(s)-3], &Vec{left, s[len(s)-3], s[len(s)-1]})
			continue
		}
		break
	}
	if rp > 100 {
		return s
	} else {
		return append(s, right, v)
	}
}

func (cc *Compiler) sortRegs(v Value, fallback []*Keyword) []*Keyword {
	regs := []*Keyword{}
	if v == NIL {
		return regs
	}

	w := v.(*Vec)
	if len(fallback) > 0 && w.Size() == 0 {
		return append(regs, fallback...)
	}

	for _, i := range *w {
		i := i.(*Identifier).Name
		regs = cc.CollectRegs(regs, i)
	}
	sort.Slice(regs, func(i, j int) bool {
		return string(*regs[i]) < string(*regs[j])
	})
	return regs

}

func (cc *Compiler) newSig(vec *Vec) *Sig {
	required := cc.sortRegs(vec.At(0), nil)
	results := cc.sortRegs(vec.At(1), nil)
	invalidated := cc.sortRegs(vec.At(2), []*Keyword{KwUNDER})
	isinline := vec.At(3) != NIL
	return &Sig{Required: required, Results: results, Invalidated: invalidated, IsInline: isinline}
}

func (cc *Compiler) IsCond(a *Keyword) bool {
	return cc.ReservedWords[a.String()] == tkCOND
}

func (cc *Compiler) IsReg(a *Keyword) bool {
	return cc.ReservedWords[a.String()] == tkREG
}

func (cc *Compiler) emitCodeFromModule(m *Module, s *Keyword) {
	if section := m.Sections[s]; section != nil {
		cc.EmitCode(section.Insts...)
		section.Insts = []*Inst{}
	}
}

func (cc *Compiler) expandAllInlines() {
	cc.EnterContext(CtProc)
	for n := 0; n < 64 && len(cc.InlineInsts) > 0; n++ {
		insts := cc.InlineInsts
		cc.InlineInsts = []*Inst{}
		for _, i := range insts {
			env := i.Args[1].(*Env)
			id := i.Args[2].(*Identifier)
			nm := cc.LookupNamed(env, id)
			if nm == nil {
				cc.RaiseCompileError(i.ExprTag(), "undefined proc %s", id.String())
			} else if nm.Kind != NmInline {
				cc.RaiseCompileError(i.ExprTag(), "%s is not a inline proc", id.String())
			}
			cc.EnterCodeBlock()
			cc.expandInline(nm.Env, i.From, id, nm.Value.(*Inline))
			i.Args[3] = &Section{Insts: cc.LeaveCodeBlock()}
		}
	}
	cc.LeaveContext()

	if len(cc.InlineInsts) > 0 {
		etag := cc.InlineInsts[0].ExprTag()
		cc.RaiseCompileError(etag, "inline proc expansion too deep")
	}
}

func flattenInsts(acc *[]*Inst, insts []*Inst) {
	for _, i := range insts {
		if i.Kind == InstMisc && i.Args[0] == KwInline {
			flattenInsts(acc, i.Args[3].(*Section).Insts)
		} else {
			*acc = append(*acc, i)
		}
	}
}

var defaultLink = []byte("org 0 0 1; merge text _; merge rodata _; merge bss _")

func (cc *Compiler) doLink() []*Inst {
	if cc.link == nil {
		cc.link = cc.Parse("@", defaultLink)
	}

	cc.hooks.beforeLink(cc)
	cc.EmitCodeToSection(cc.Section, cc.LeaveCodeBlock()...) // Toplevel
	cc.expandAllInlines()
	cc.EnterCodeBlock()
	env := cc.Toplevel
	modules := env.Filter(NmModule)
	for _, e := range (*cc.link)[1:] {
		e := e.(*Vec)
		switch e.ExprTagName() {
		case KwOrg:
			etag, _ := CheckExpr(e, 4, 4, 0, cc)
			addr := EvalConstAs(e.At(1), env, IntT, "origin address", etag, cc)
			size := EvalConstAs(e.At(2), env, IntT, "limit size", etag, cc)
			mode := EvalConstAs(e.At(3), env, IntT, "mode", etag, cc)
			cc.EmitCode(NewInst(e, InstOrg, addr, size, mode, Int(0)))
		case KwMerge:
			etag, _ := CheckExpr(e, 2, -1, 0, cc)
			sec := CheckConst(e.At(1), IdentifierT, "section name", etag, cc)
			for _, i := range (*e)[2:] {
				if KwProg.MatchExpr(i) != nil {
					cc.CompileExpr(env, i)
					continue
				}

				id := CheckConst(i, IdentifierT, "module name", etag, cc)
				if id.Name == KwUNDER {
					for _, nm := range modules {
						cc.emitCodeFromModule(nm.Value.(*Module), sec.Name)
					}
				} else {
					nm := cc.Toplevel.FindById(id)
					if nm == nil || nm.Kind != NmModule {
						cc.RaiseCompileError(etag, "unknown module %s", id)
					}
					cc.emitCodeFromModule(nm.Value.(*Module), sec.Name)
				}
			}
		default:
			cc.RaiseCompileError(e.ExprTag(), "invalid link form")
		}
	}

	for _, nm := range modules {
		for name, section := range nm.Value.(*Module).Sections {
			if len(section.Insts) > 0 {
				cc.RaiseCompileError(nil, "the section `%s@%s` is not linked", name.String(), nm.Name.String())
			}
		}
	}
	cc.EmitCode(cc.LeaveCodeBlock()...)
	cc.EmitCode(NewInst(&Vec{IdUNDER}, InstOrg, Int(0), Int(0), Int(0), Int(0)))

	flatten := []*Inst{}
	flattenInsts(&flatten, cc.LeaveCodeBlock())
	return flatten
}

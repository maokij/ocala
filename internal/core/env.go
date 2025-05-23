package core

import (
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"
)

func CopyPtr[T any](v *T) *T {
	w := *v
	return &w
}

type Value interface {
	Inspect() string
	Dup() Value
}

type InstTab interface {
	Find(*Keyword) InstTab
}

type InstPat map[*Keyword]InstTab

func (v InstPat) Find(k *Keyword) InstTab {
	return v[k]
}

type InstDat []BCode

func (v InstDat) Find(k *Keyword) InstTab {
	return nil
}

// //////////////////////////////////////////////////////////
type Nil struct{}

var NIL = &Nil{}

func (v *Nil) Inspect() string {
	return "<NIL>"
}

func (*Nil) Dup() Value {
	return NIL
}

// //////////////////////////////////////////////////////////
type Undefined struct{}

var UNDEFINED = &Undefined{}

func (v *Undefined) Inspect() string {
	return "<UNDEFINED>"
}

func (*Undefined) Dup() Value {
	return UNDEFINED
}

// //////////////////////////////////////////////////////////
type Int int

var IntT = Int(0)

func (v Int) Inspect() string {
	return fmt.Sprint(v)
}

func (v Int) Dup() Value {
	return v
}

func BoolInt(a bool) Int {
	if a {
		return Int(1)
	}
	return Int(0)
}

// //////////////////////////////////////////////////////////
type Binary []byte

var BynaryT = &Binary{}

func (v *Binary) Inspect() string {
	return fmt.Sprint(v)
}

func (v *Binary) Dup() Value {
	w := append(make(Binary, 0, len(*v)), *v...)
	return &w
}

// //////////////////////////////////////////////////////////
type Keyword string

var KeywordT = NewKeyword("")

func (v *Keyword) Inspect() string {
	return ":" + string(*v)
}

func (v *Keyword) Dup() Value {
	return v
}

func (v *Keyword) String() string {
	return string(*v)
}

func (v *Keyword) ToId(token *Token) *Identifier {
	id := &Identifier{Name: v, Token: CopyPtr(token)}
	id.Token.Value = id
	return id
}

func (v *Keyword) MatchId(a Value) *Identifier {
	if a, ok := a.(*Identifier); ok && a.Namespace == nil && v == a.Name {
		return a
	}
	return nil
}

func (v *Keyword) MatchExpr(a Value) *Vec {
	if a, ok := a.(*Vec); ok && v.MatchId(a.AtOrUndef(0)) != nil {
		return a
	}
	return nil
}

var keywords = map[string]*Keyword{}

func Intern(s string) *Keyword {
	v, ok := keywords[s]
	if !ok {
		v = NewKeyword(strings.Clone(s))
		keywords[s] = v
	}
	return v
}

var gensymSerial = 0

func Gensym(s string) *Keyword {
	gensymSerial++
	return NewKeyword(fmt.Sprintf("%s.G%d", s, gensymSerial))
}

func NewKeyword(s string) *Keyword {
	return (*Keyword)(&s)
}

// //////////////////////////////////////////////////////////
type Str string

var StrT = NewStr("")

func (v *Str) Inspect() string {
	s := []byte{'"'}
	for _, c := range []byte(*v) {
		if d, ok := invertedEscapeChars[c]; ok {
			s = append(s, '\\', d)
		} else if c < 0x20 || c > 0x7e {
			s = append(s, fmt.Sprintf("\\x%02x", c)...)
		} else {
			s = append(s, c)
		}
	}
	s = append(s, '"')
	return string(s)
}

func (v *Str) Dup() Value {
	return CopyPtr(v)
}

func (v *Str) String() string {
	return string(*v)
}

func (v *Str) Intern() *Keyword {
	return Intern(v.String())
}

func NewStr(s string) *Str {
	return (*Str)(&s)
}

// //////////////////////////////////////////////////////////
type Vec []Value

var VecT = &Vec{}

func (v *Vec) Inspect() string {
	s := strings.Builder{}
	s.WriteByte('[')
	for x, i := range *v {
		if i != nil {
			s.WriteString(i.Inspect())
		} else {
			s.WriteString("<nil>")
		}
		if x < len(*v)-1 {
			s.WriteByte(' ')
		}
	}
	s.WriteByte(']')
	return s.String()
}

func (v *Vec) Dup() Value {
	w := []Value{}
	for _, i := range *v {
		w = append(w, i.Dup())
	}
	return (*Vec)(&w)
}

func (v *Vec) At(x int) Value {
	return (*v)[x]
}

func (v *Vec) AtOrUndef(x int) Value {
	if x < len(*v) {
		return (*v)[x]
	}
	return UNDEFINED
}

func (v *Vec) SetAt(x int, a Value) {
	(*v)[x] = a
}

func (v *Vec) Size() int {
	return len(*v)
}

func (v *Vec) Push(a ...Value) *Vec {
	*v = append(*v, a...)
	return v
}

func (v *Vec) Flatten() *Vec {
	r := &Vec{}
	for _, i := range *v {
		switch v := i.(type) {
		case *Vec:
			r.Push(*v.Flatten()...)
		default:
			r.Push(v)
		}
	}
	return r
}

func (v *Vec) ExprTag() *Identifier {
	if len(*v) > 0 {
		if v, ok := (*v)[0].(*Identifier); ok {
			return v
		}
	}
	return nil
}

func (v *Vec) ExprTagName() *Keyword {
	if etag := v.ExprTag(); etag != nil && etag.Namespace == nil {
		return etag.Name
	}
	return nil
}

func (v *Vec) OperandAt(x int) *Operand {
	if x < len(*v) {
		return (*v)[x].(*Operand)
	}
	return NoOperand
}

func NewVec(a []Value) *Vec {
	return (*Vec)(&a)
}

// //////////////////////////////////////////////////////////
type Map map[Value]Value

var MapT = &Map{}

func (v *Map) Inspect() string {
	return fmt.Sprint(*v)
}

func (v *Map) Dup() Value {
	return v
}

// //////////////////////////////////////////////////////////
type Identifier struct {
	Name        *Keyword
	Token       *Token
	Namespace   *Keyword
	ExpandedBy  *Identifier
	PlaceHolder string
}

var IdentifierT = &Identifier{}

func (v *Identifier) Inspect() string {
	return v.String()
}

func (v *Identifier) Dup() Value {
	return v.Clone()
}

func (v *Identifier) Clone() *Identifier {
	v = CopyPtr(v)
	v.Token = CopyPtr(v.Token)
	v.Token.Value = v
	return v
}

func (v *Identifier) String() string {
	s := v.PlaceHolder + v.Name.String()
	if v.Namespace != nil {
		return v.Namespace.String() + ":" + s
	}
	return s
}

func (v *Identifier) Expand(kw *Keyword) *Identifier {
	id := &Identifier{Name: kw, ExpandedBy: v}
	id.Token = &Token{From: InternalParser, Value: id}
	return id
}

func (v *Identifier) ToConstexpr(env *Env) *Constexpr {
	return &Constexpr{Token: v.Token, Body: v, Env: env}
}

func InternalId(kw *Keyword) *Identifier {
	id := &Identifier{Name: kw}
	id.Token = &Token{From: InternalParser, Value: id}
	return id
}

// //////////////////////////////////////////////////////////
type Sig struct {
	IsInline    bool
	Required    []*Keyword
	Results     []*Keyword
	Invalidated []*Keyword
}

func (v *Sig) Inspect() string {
	return fmt.Sprintf("<Sig:%v>", v)
}

func (v *Sig) Dup() Value {
	return CopyPtr(v)
}

func (v *Sig) String() string {
	s := ""
	if v.IsInline {
		s = "-* "
	}
	return fmt.Sprintf("%s%v => %v ! %v", s, v.Required, v.Results, v.Invalidated)
}

func (a *Sig) Equals(b *Sig) bool {
	return a.IsInline == b.IsInline &&
		reflect.DeepEqual(a.Required, b.Required) &&
		reflect.DeepEqual(a.Results, b.Results) &&
		reflect.DeepEqual(a.Invalidated, b.Invalidated)
}

// //////////////////////////////////////////////////////////
type Macro struct {
	Args []*Keyword
	Vars []Value
	Opts []Value
	Body *Vec
	Rest bool
}

func (v *Macro) Inspect() string {
	return fmt.Sprintf("<Macro:%v>", v)
}

func (v *Macro) Dup() Value {
	return v
}

// //////////////////////////////////////////////////////////
type Constexpr struct {
	Token *Token
	Body  Value
	Env   *Env
}

var ConstexprT = &Constexpr{}

func (v *Constexpr) Inspect() string {
	return fmt.Sprintf("(%v)", v.Body.Inspect())
}

func (v *Constexpr) Dup() Value {
	return &Constexpr{Token: CopyPtr(v.Token), Body: v.Body.Dup()}
}

func InternalConstexpr(v Value) *Constexpr {
	return &Constexpr{Token: &Token{From: InternalParser}, Body: v}
}

// //////////////////////////////////////////////////////////
type ConstFn struct {
	Args  []*Keyword
	Opts  []Value
	Token *Token
	Body  Value
	Env   *Env
}

var ConstFnT = &Constexpr{}

func (v *ConstFn) Inspect() string {
	return fmt.Sprintf("%v(%v)", v.Args, v.Body.Inspect())
}

func (v *ConstFn) Dup() Value {
	return v
}

func NewConstFn(args []*Keyword, opts []Value, v *Constexpr) *ConstFn {
	return &ConstFn{Args: args, Opts: opts, Token: v.Token, Body: v.Body, Env: v.Env}
}

// //////////////////////////////////////////////////////////
type SyntaxFn func(*Compiler, *Env, *Vec) Value

func (v SyntaxFn) Inspect() string {
	return fmt.Sprintf("<Syntax:%v>", v)
}

func (v SyntaxFn) Dup() Value {
	return v
}

// //////////////////////////////////////////////////////////
type BCode struct {
	Kind byte
	A0   byte
	A1   byte
	A2   byte
	A3   byte
	A4   byte
	A5   byte
	A6   byte
}

const (
	BcByte = iota
	BcLow
	BcHigh
	BcRlow
	BcRhigh
	BcImp
	BcMap
	BcTemp
)

var BCodeT = BCode{}

func (v BCode) Inspect() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x", v.Kind, v.A0, v.A1, v.A2)
}

func (v BCode) Dup() Value {
	return v
}

// //////////////////////////////////////////////////////////
type Operand struct {
	Kind *Keyword
	A0   Value
	A1   Value
	A2   Value
}

var OperandT = &Operand{}
var NoOperand = &Operand{Kind: KwUndefined}
var InvalidOperand = &Operand{Kind: NewKeyword("invalid-operand")}

func (v *Operand) Inspect() string {
	if a, ok := v.A0.(*Constexpr); ok {
		return fmt.Sprintf("{%s:%v}", v.Kind, a.Inspect())
	}
	return fmt.Sprintf("{%s}", v.Kind)
}

func (v *Operand) Dup() Value {
	return CopyPtr(v)
}

// //////////////////////////////////////////////////////////
type Inst struct {
	Kind int
	From *Vec
	Args []Value
	size int
}

const (
	InstLabel = iota
	InstCode
	InstData
	InstDS
	InstOrg
	InstAlign
	InstFile
	InstConst
	InstMisc
	InstAssert
)

func (v *Inst) Inspect() string {
	return fmt.Sprintf("%d %s", byte(v.Kind), (*Vec)(&v.Args).Inspect())
}

func (v *Inst) Dup() Value {
	args := []Value{}
	for _, i := range v.Args {
		args = append(args, i.Dup())
	}
	return &Inst{Kind: v.Kind, From: v.From, Args: args, size: v.size}
}

func (v *Inst) ExprTag() *Identifier {
	return v.From.ExprTag()
}

func (v *Inst) MatchCode(a ...*Keyword) bool {
	return v.Kind == InstCode && slices.Contains(a, v.Args[0].(*Keyword))
}

func NewInst(from *Vec, kind int, args ...Value) *Inst {
	return &Inst{Kind: kind, From: from, Args: args}
}

// //////////////////////////////////////////////////////////
type Module struct {
	Name     *Keyword
	Env      *Env
	Sections map[*Keyword][]*Inst
}

var ModuleT = &Module{}

func (v *Module) Inspect() string {
	return fmt.Sprintf("<Module:%v>", v)
}

func (v *Module) Dup() Value {
	return v
}

func NewModule(name *Keyword, env *Env) *Module {
	return &Module{
		Name:     name,
		Env:      env,
		Sections: map[*Keyword][]*Inst{KwTEXT: {}, KwBSS: {}, KwRODATA: {}},
	}
}

// //////////////////////////////////////////////////////////
type Named struct {
	Token   *Token
	Name    *Keyword
	AsmName *Keyword
	Env     *Env
	Value   Value
	Kind    int32
	Export  bool
	Special bool
	Serial  int
}

const (
	NmModule = iota
	NmSyntax
	NmMacro
	NmInline
	NmFunc
	NmInst
	NmConst
	NmLabel
	NmVar
	NmSpecial
	NmInvalid
)

var NamedKindLabels = []string{
	"module",
	"syntax",
	"macro",
	"inline",
	"func",
	"inst",
	"const",
	"label",
	"var",
	"special",
	"invalid",
}

var NamedT = &Named{}
var namedSerial = 0

func nextNamedSerial() int {
	namedSerial++
	return namedSerial
}

func (v *Named) Inspect() string {
	return fmt.Sprintf("<Named:%s>", v.Name.String())
}

func (v *Named) Dup() Value {
	return CopyPtr(v)
}

// //////////////////////////////////////////////////////////
type Label struct {
	Addr int
	Link *Inst
	Sig  *Sig
}

var LabelT = &Label{}

func (v *Label) Inspect() string {
	return fmt.Sprintf("<Label:%v>", *v)
}

func (v *Label) Dup() Value {
	return CopyPtr(v)
}

func (v *Label) LinkedToData() bool {
	return v.Link != nil && (v.Link.Kind == InstData || v.Link.Kind == InstDS)
}

// //////////////////////////////////////////////////////////
type Inline struct {
	Body *Vec
	Sig  *Sig
}

var InlineT = &Inline{}

func (v *Inline) Inspect() string {
	return fmt.Sprintf("<Inline:%v>", *v)
}

func (v *Inline) Dup() Value {
	return CopyPtr(v)
}

// //////////////////////////////////////////////////////////
type Env struct {
	outer  *Env
	inners []*Env
	names  map[*Keyword]*Named
	serial int
	Module *Module
}

var envSerial = 0

func NewEnv(outer *Env) *Env {
	envSerial++
	return &Env{serial: envSerial, outer: outer, names: map[*Keyword]*Named{}}
}

func (env *Env) Inspect() string {
	return fmt.Sprintf("<Env:%d>", env.serial)
}

func (env *Env) Dup() Value {
	return env
}

func (env *Env) Outer() *Env {
	return env.outer
}

func (env *Env) Enter() *Env {
	inner := NewEnv(env)
	env.inners = append(env.inners, inner)
	return inner
}

// var reAsmUnsafeChars = regexp.MustCompile(`[^_a-zA-Z0-9]`)

func (env *Env) Install(nm *Named) *Named {
	if nm.Env == nil {
		nm.Env = env
	}
	if nm.Serial == 0 {
		nm.Serial = nextNamedSerial()
	}
	if nm.AsmName == nil {
		nm.AsmName = NewKeyword(fmt.Sprintf(".%s.#%d", nm.Name.String(), nm.Serial))
	}
	if nm.Token == nil {
		nm.Token = &Token{From: InternalParser}
	}
	env.names[nm.Name] = nm
	return nm
}

func (env *Env) LookupById(id *Identifier) *Named {
	return env.Lookup(id.Name)
}

func (env *Env) Lookup(k *Keyword) *Named {
	nm := env.Find(k)
	if nm != nil {
		return nm
	}
	if env.outer != nil {
		return env.outer.Lookup(k)
	}
	return nil
}

func (env *Env) FindById(id *Identifier) *Named {
	return env.Find(id.Name)
}

func (env *Env) Find(k *Keyword) *Named {
	return env.names[k]
}

func (env *Env) Filter(kind int32) []*Named {
	filtered := []*Named{}
	for _, i := range env.names {
		if i.Kind == kind {
			filtered = append(filtered, i)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Serial < filtered[j].Serial
	})
	return filtered
}

func (env *Env) InsertEnv(e *Env) {
	env.outer = NewEnv(env.outer)
	env.outer.names = e.names
}

package file

import (
	"fmt"
	_ "ocala"
	"ocala/core"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

const TkDTMI = tkDTMI

var invalidToken = &core.Token{From: core.InternalParser, Kind: tkEOF, Value: core.NIL}

var compilersCache = map[string]*core.Compiler{}

var AppRoot = ""

func Init(path string) error {
	var err error
	path, err = core.FindAppRoot(path)
	if err != nil {
		return err
	}
	AppRoot = path
	return nil
}

func BuildCompilerFromSource(text []byte) (*core.Compiler, error) {
	arch := core.FindArchDirective(text)
	if arch == "" {
		return nil, fmt.Errorf("arch required")
	}

	cc := BuildCompiler(arch)
	if cc == nil {
		return nil, fmt.Errorf("unknown arch: %s", arch)
	}
	return cc, nil
}

func BuildCompiler(arch string) *core.Compiler {
	if cc := compilersCache[arch]; cc != nil {
		return cc
	}
	if cc := core.NewCompiler(arch); cc != nil {
		cc.Init("-")
		env := cc.Builtins
		for i := range reservedWords {
			env.Install(&core.Named{
				Name:  core.Intern(i),
				Kind:  core.NmSpecial,
				Value: core.NIL,
			})
		}
		compilersCache[arch] = cc
		return cc
	}
	return nil
}

type ErrorInfo struct {
	Message string
	Token   *core.Token
}

func ErrorInfoOf(err *core.InternalError) []*ErrorInfo {
	message := err.Error()
	result := []*ErrorInfo{}
	for _, token := range err.Tokens() {
		message := fmt.Sprintf("%s\n[error #%d]\n", message, len(result))
		message += token.FormatAsErrorLine("at")
		if id, ok := token.Value.(*core.Identifier); ok {
			for ; id.ExpandedBy != nil; id = id.ExpandedBy {
				from := id.ExpandedBy.Token
				if token.IsInternal() {
					token = from
				}
				message += from.FormatAsErrorLine("from")
			}
		}
		if !token.IsInternal() {
			result = append(result, &ErrorInfo{Message: message, Token: token})
		}
	}
	if len(result) == 0 {
		result = append(result, &ErrorInfo{Message: message, Token: invalidToken})
	}
	return result
}

func newInternalError(s string) *core.InternalError {
	err := &core.InternalError{}
	err.SetMessage(s)
	return err
}

func CheckCode(path string, text []byte, options *CompileOptions) []*core.InternalError {
	cc, err := BuildCompilerFromSource(text)
	if err != nil {
		err := newInternalError(err.Error())
		return []*core.InternalError{err}
	}
	if r := checkSyntax(cc, path, text, options); r != nil {
		return r
	}
	if r := checkSemantics(cc, path, text, options); r != nil {
		return r
	}
	return nil
}

func checkSyntax(cc *core.Compiler, path string, text []byte, options *CompileOptions) []*core.InternalError {
	f := &File{
		Parser:         cc.NewParser(path, text),
		followStack:    make([]uint64, 0, 64),
		enableCheck:    true,
		Node:           &BlockNode{},
		CompileOptions: options,
	}
	f.Scanner.OnError = func(err *core.InternalError) {
		f.Errors = append(f.Errors, err)
	}
	f._parse()
	return f.Errors
}

func checkSemantics(cc *core.Compiler, path string, text []byte, options *CompileOptions) []*core.InternalError {
	cc = core.NewCompiler(cc.FullArchName())
	g := &core.Generator{}
	g.IncPaths = options.IncPaths
	g.Defs = options.Defs
	g.SetCompiler(cc)
	func() {
		defer g.HandlePanic()
		g.CheckLink(g.Compile(path, text))
	}()

	if g.Err == nil {
		return nil
	}
	return []*core.InternalError{g.Err}
}

type CompileOptions struct {
	IncPaths []string `json:"incPaths"`
	Defs     []string `json:"defs"`
}

func NewCompileOptions() *CompileOptions {
	return &CompileOptions{IncPaths: AdjustIncPaths(nil)}
}

type File struct {
	*core.Parser
	*CompileOptions
	mu sync.Mutex

	Errors   []*core.InternalError
	Tokens   []*core.Token
	Depends  []*File
	Includes []string
	Node     Node
	Env      *core.Env
	Version  int
	Sequence int
	Analyzed int

	enableCheck bool
	followStack []uint64
	scope       *core.Token
	pending     Node
}

func NewFile(path string, version int, options *CompileOptions) *File {
	f := &File{Parser: &core.Parser{}, Version: version, CompileOptions: options}
	f.Path = path
	f.followStack = make([]uint64, 0, 64)
	return f
}

func (f *File) ConsumeToken() *core.Token {
	token := f.Scanner.ConsumeToken()
	token.Parent = f.scope
	f.Tokens = append(f.Tokens, token)
	return token
}

func (f *File) ensureEOF() {
	for {
		if f.PeekToken(); f.ConsumeToken().Kind == tkEOF {
			break
		}
	}
}

func (f *File) ErrorUnexpected(expected string) {
	if !f.enableCheck || f.Recovering {
		return
	}

	token := f.PeekToken()
	if token.Kind == tkSCANERROR {
		f.Trace("[SCANERROR]")
		return
	}
	err := f.Parser.ErrorAt(token)
	label := tokenLabels[token.Kind]
	err.SetMessage("unexpected %s, expected %s", label, expected)
	f.Trace("[ERROR] ", err.Error())
	f.Errors = append(f.Errors, err)
}

func (f *File) HandleError(v core.Value) (recovered bool) {
	size := len(f.followStack)
	n := f.followStack[size-1]
	f.followStack = f.followStack[:size-1]
	if v != SYNTAX_ERROR {
		return true
	}

	m := uint64(tpEOF)
	for _, i := range f.followStack {
		m |= i
	}
	f.Recovering = true
	for {
		token := f.PeekToken()
		mask := uint64(1) << token.Kind
		if n&mask != 0 {
			f.Trace("recovered:", token)
			recovered = true
			break
		}
		if m&mask != 0 {
			f.Trace("catched:", token)
			break
		}
		f.Trace("not recovered:", token)
		f.ConsumeToken()
	}
	return
}

func (f *File) Follow(n uint64) {
	f.followStack = append(f.followStack, n)
}

func (f *File) Trace(s ...any) {
	if core.Debug.Enabled {
		fmt.Fprintln(os.Stderr, s...)
	}
}

func (f *File) HoldNode(node Node, name core.Value) {
	if f.AddNode(node, name) {
		f.pending = node
	}
}

func (f *File) AddNode(node Node, name core.Value) bool {
	if id, ok := name.(*core.Identifier); ok && !f.Recovering {
		f.addNamedNode(node, id)
		return true
	}
	return false
}

func (f *File) addNamedNode(node Node, id *core.Identifier) {
	child := node.Base()
	child.Name = id
	child.parent = f.Node

	parent := f.Node.Base()
	parent.AddChildNode(node)
}

func (f *File) EnterNode(token *core.Token) {
	node := f.pending
	if node == nil {
		node = &BlockNode{NodeBase{Token: token}}
		f.addNamedNode(node, core.IdUNDER)
	}
	f.Node = node
	f.scope = token
	f.pending = nil
	token.Value = node
}

func (f *File) LeaveNode() {
	f.Node = f.Node.Base().parent
	if f.scope != nil {
		f.scope = f.scope.Parent
	}
}

func (f *File) Update(text []byte) {
	f.mu.Lock()
	f.Text = text
	f.mu.Unlock()
	f.Tokens = f.Tokens[:0]
	f.Scanner.Init()
	f.Sequence++
}

func (f *File) SafeGetText() []byte {
	f.mu.Lock()
	text := f.Text
	f.mu.Unlock()
	return text
}

func (f *File) Parse() {
	cc, err := BuildCompilerFromSource(f.Text)
	if err != nil {
		cc = BuildCompiler("z80")
	}

	f.Parser.SetCompiler(cc)
	f.followStack = f.followStack[:0]
	f.Tokens = f.Tokens[:0]
	f.Scanner.OnError = func(err *core.InternalError) {}

	rootNode := &BlockNode{
		NodeBase: NodeBase{Token: &core.Token{Kind: tkEQLC}, Name: core.IdUNDER},
	}
	f.scope = rootNode.Token
	f.scope.Value = rootNode
	f.Node = rootNode
	f.Includes = f.Includes[:0]
	f._parse()
	f.ensureEOF()
	f.buildEnv(core.NewEnv(cc.Toplevel), rootNode)
	f.Node = rootNode
	f.Env = rootNode.Env
}

func (f *File) Analyze(open func(string) *File) {
	if f.Analyzed == f.Sequence {
		return
	}

	f.Analyzed = f.Sequence
	f.Parse()
	for _, i := range f.Includes {
		dir := filepath.Dir(f.Path)
		path, err := core.RegularizePath(i, dir, f.IncPaths)
		if err != nil {
			continue
		}
		if dep := open(path); dep != nil {
			f.Depends = append(f.Depends, dep)
		}
	}
	f.MergeEnv()
}

func (f *File) MergeEnv() {
	var all func(*File, map[*File]bool)
	all = func(f *File, merged map[*File]bool) {
		if merged[f] {
			return
		}
		merged[f] = true

		for _, i := range f.Depends {
			all(i, merged)
		}
	}
	deps := map[*File]bool{}
	all(f, deps)
	delete(deps, f)

	env := core.NewEnv(nil)
	for i := range deps {
		if i.Env != nil {
			env.MergeEnv(i.Env)
		}
	}
	f.Env.InsertEnv(env)
}

func (f *File) buildEnv(env *core.Env, e Node) {
	switch e := e.(type) {
	case *MacroNode:
		f.Install(env, e.Name, core.NmMacro, e)
		e.Env = env.Enter()
		for _, i := range e.Args {
			f.Install(e.Env, i, core.NmVar, nil)
		}
		for _, i := range e.Vars {
			f.Install(e.Env, i, core.NmVar, nil)
		}
		for _, i := range e.children {
			f.buildEnv(e.Env, i)
		}
	case *ProcNode:
		f.Install(env, e.Name, core.NmLabel, e)
		e.Env = env.Enter()
		for _, i := range e.children {
			f.buildEnv(e.Env, i)
		}
	case *ConstNode:
		f.Install(env, e.Name, core.NmConst, e)
	case *ConstFnNode:
		f.Install(env, e.Name, core.NmConst, e)
	case *DataNode:
		f.Install(env, e.Name, core.NmLabel, e)
		e.Env = env
	case *ModuleNode:
		f.Install(env, e.Name, core.NmModule, e)
		e.Env = env.Enter()
		for _, i := range e.children {
			f.buildEnv(e.Env, i)
		}
	case *LabelNode:
		if env.Lookup(e.Name.Name) == nil {
			f.Install(env, e.Name, core.NmLabel, e)
		}
	case *StructNode:
		f.Install(env, e.Name, core.NmDatatype, e)
		e.Env = env.Enter()
		for _, i := range e.children {
			f.buildEnv(e.Env, i)
		}
	case *StructFieldNode:
		f.Install(env, e.Name, core.NmVar, e)
		e.Env = env
	case *BlockNode:
		if e.Token.Kind == tkLC {
			env = env.Enter()
		}
		e.Env = env
		for _, i := range e.children {
			f.buildEnv(env, i)
		}
	}
}

func (f *File) TokenIndexAt(pos int32) int {
	n := len(f.Tokens)
	x := sort.Search(n, func(x int) bool {
		return f.Tokens[x].Pos >= pos
	})
	if x == n {
		return -1
	}
	return x - 1 // x == 0: invalid(-1)
}

func (f *File) Install(env *core.Env, id *core.Identifier, kind int32, value core.Value) *core.Named {
	nm := &core.Named{Name: id.Name, Token: id.Token, Kind: kind, Value: value}
	return env.Install(nm)
}

func ParentNodeAt(token *core.Token) Node {
	return token.Parent.Value.(Node)
}

func EnvAt(token *core.Token) *core.Env {
	return ParentNodeAt(token).Base().Env
}

func IsDotToken(token *core.Token) bool {
	return token.Kind == tkDTMI
}

func IsDotLikeToken(token *core.Token) bool {
	return token.Kind == tkDTMI || token.Kind == tkDOT_OPERATOR
}

func TokenString(token *core.Token) string {
	return string(token.From.Text[token.Pos:token.End.Pos])
}

func AdjustIncPaths(paths []string) []string {
	base := []string{filepath.Join(AppRoot, "share/ocala/include")}
	return append(base, paths...)
}

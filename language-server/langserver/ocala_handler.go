package langserver

import (
	"bytes"
	"cmp"
	"fmt"
	"ocala/core"
	"ocala/language-server/file"
	"slices"
	"strings"
)

type DocumentSymbol struct {
	Name           string           `json:"name"`
	Detail         string           `json:"detail"`
	Kind           int64            `json:"kind"`
	Tags           []int            `json:"tags"`
	Deprecated     bool             `json:"deprecated"`
	Range          Range            `json:"range"`
	SelectionRange Range            `json:"selectionRange"`
	Children       []DocumentSymbol `json:"children"`
}

const (
	SymbolKindFile          = 1
	SymbolKindModule        = 2
	SymbolKindNamespace     = 3
	SymbolKindPackage       = 4
	SymbolKindClass         = 5
	SymbolKindMethod        = 6
	SymbolKindProperty      = 7
	SymbolKindField         = 8
	SymbolKindConstructor   = 9
	SymbolKindEnum          = 10
	SymbolKindInterface     = 11
	SymbolKindFunction      = 12
	SymbolKindVariable      = 13
	SymbolKindConstant      = 14
	SymbolKindString        = 15
	SymbolKindNumber        = 16
	SymbolKindBoolean       = 17
	SymbolKindArray         = 18
	SymbolKindObject        = 19
	SymbolKindKey           = 20
	SymbolKindNull          = 21
	SymbolKindEnumMember    = 22
	SymbolKindStruct        = 23
	SymbolKindEvent         = 24
	SymbolKindOperator      = 25
	SymbolKindTypeparameter = 26
)

func namesStartWith(env *core.Env, s string) []*core.Named {
	names := []*core.Named{}
	for nm := range env.Names() {
		if strings.HasPrefix(string(*nm.Name), s) {
			names = append(names, nm)
		}
	}
	return names
}

func lookupNamed(env, toplevel *core.Env, token *core.Token) *core.Named {
	name := file.TokenString(token)
	if x := strings.IndexByte(name, ':'); x > -1 {
		nm := toplevel.Lookup(core.Intern(name[:x]))
		if nm != nil && nm.Kind == core.NmModule {
			env = nm.Value.(*file.ModuleNode).Env
			return env.Find(core.Intern(name[x+1:]))
		}
		return nil
	}
	return env.Lookup(core.Intern(name))
}

func chainEnv(nm *core.Named) *core.Env {
	if nm == nil {
		return nil
	}

	switch v := nm.Value.(type) {
	case *file.StructNode:
		return v.Env
	case *file.DataNode:
		if v.Type != nil {
			return chainEnv(v.Env.Lookup(v.Type.Name))
		}
	case *file.StructFieldNode:
		if v.Type != nil {
			return chainEnv(v.Env.Lookup(v.Type.Name))
		}
	}
	return nil
}

func resolveEnv(f *file.File, i int) *core.Env {
	j := i - 2
	for j >= 1 && file.IsDotToken(f.Tokens[j]) {
		j -= 2
	}
	j += 1

	nm := lookupNamed(file.EnvAt(f.Tokens[i]), f.Env, f.Tokens[j])
	env := chainEnv(nm)
	j += 2

	for j < i && env != nil {
		nm := env.Lookup(core.Intern(file.TokenString(f.Tokens[j])))
		env = chainEnv(nm)
		j += 2
	}
	return env
}

func (h *langHandler) completion(uri DocumentURI, params *CompletionParams) ([]CompletionItem, error) {
	f, ok := h.files[uri]
	if !ok {
		return nil, fmt.Errorf("document not found: %v", uri)
	}

	f.Analyze(h.addIncludedFile)
	pt := positionToPt(f, params.Position)
	x := f.TokenIndexAt(pt.Pos)
	if x < 0 {
		return nil, nil
	}

	token := f.Tokens[x]
	prefix := string(f.Text[token.Pos:pt.Pos])
	founds := []*core.Named{}
	if x >= 1 && file.IsDotLikeToken(f.Tokens[x]) {
		if env := resolveEnv(f, x); env != nil {
			founds = namesStartWith(env, "")
		}
	} else if x >= 2 && file.IsDotToken(f.Tokens[x-1]) {
		if env := resolveEnv(f, x-1); env != nil {
			founds = namesStartWith(env, prefix)
		}
	} else if x := strings.IndexByte(prefix, ':'); x > -1 {
		env := f.Env
		ns := core.Intern(prefix[:x])
		nm := env.Lookup(ns)
		if nm != nil && nm.Kind == core.NmLabel { // shadowed by label?
			nm = env.Outer().Lookup(ns)
		}
		if nm != nil && nm.Kind == core.NmModule {
			env := nm.Value.(*file.ModuleNode).Env
			founds = namesStartWith(env, prefix[x+1:])
		}
	} else {
		for env := file.EnvAt(token); env != nil; env = env.Outer() {
			founds = append(founds, namesStartWith(env, prefix)...)
		}
	}

	slices.SortStableFunc(founds, func(a, b *core.Named) int {
		return cmp.Compare(*a.Name, *b.Name)
	})
	names := slices.CompactFunc(founds, func(a, b *core.Named) bool {
		return a.Name == b.Name
	})

	result := []CompletionItem{}
	for _, nm := range names {
		label := string(*nm.Name)
		insert := label
		kind := CompletionItemKind(0)

		switch v := nm.Value.(type) {
		case *file.MacroNode:
			label = fmt.Sprintf("%s(%s)", label, file.JoinStringer(v.Args, " "))
			kind = KeywordCompletion
		case *file.ProcNode:
			label = fmt.Sprintf("%s(%s)", label, file.SigDetail(v.Sig))
			insert = label
			kind = FunctionCompletion
		case *file.ConstNode:
			kind = ConstantCompletion
		case *file.ConstFnNode:
			label = fmt.Sprintf("%s(%s)", label, file.JoinStringer(v.Args, " "))
			kind = ConstantCompletion
		case *file.DataNode:
			kind = ConstantCompletion
		case *file.ModuleNode:
			kind = ModuleCompletion
		case *file.LabelNode:
			kind = ConstantCompletion
		case *file.StructNode:
			kind = ClassCompletion
		default: // builtins
			switch nm.Kind {
			case core.NmVar:
				kind = VariableCompletion
			case core.NmConst, core.NmFunc, core.NmLabel:
				kind = ConstantCompletion
			case core.NmSyntax, core.NmInst, core.NmSpecial:
				kind = KeywordCompletion
			}
		}
		result = append(result, CompletionItem{
			Label:      label,
			InsertText: insert,
			Kind:       kind,
		})
	}
	return result, nil
}

func (h *langHandler) symbol(uri DocumentURI) ([]DocumentSymbol, error) {
	var traverse func(file.Node, *[]DocumentSymbol)
	traverse = func(node file.Node, acc *[]DocumentSymbol) {
		name := core.IdUNDER
		detail := ""
		kind := int64(0)
		switch node := node.(type) {
		case *file.MacroNode:
			name = node.Name
			kind = SymbolKindOperator
			detail = node.Detail()
		case *file.ProcNode:
			name = node.Name
			kind = SymbolKindFunction
			detail = node.Detail()
		case *file.ConstNode:
			name = node.Name
			kind = SymbolKindConstant
		case *file.ConstFnNode:
			name = node.Name
			kind = SymbolKindConstant
			detail = node.Detail()
		case *file.DataNode:
			name = node.Name
			kind = SymbolKindVariable
		case *file.ModuleNode:
			name = node.Name
			kind = SymbolKindModule
		case *file.LabelNode:
			name = node.Name
			kind = SymbolKindVariable
		case *file.StructNode:
			name = node.Name
			kind = SymbolKindVariable
		case *file.StructFieldNode:
			name = node.Name
			kind = SymbolKindField
		case *file.BlockNode:
			name = node.Name
		}

		if name.Name != core.KwUNDER {
			token := name.Token
			*acc = append(*acc, DocumentSymbol{
				Range:          tokenToRange(token),
				SelectionRange: tokenToRange(token),
				Kind:           kind,
				Name:           string(*name.Name),
				Detail:         detail,
			})
			acc = &(*acc)[len(*acc)-1].Children
		}
		for _, i := range node.Base().Children() {
			traverse(i, acc)
		}
	}

	f, ok := h.files[uri]
	if !ok {
		return nil, fmt.Errorf("document not found: %v", uri)
	}

	f.Analyze(h.addIncludedFile)
	acc := &[]DocumentSymbol{}
	traverse(f.Node, acc)
	return *acc, nil
}

var noadjust = []byte(";\n:.()[]")

func (h *langHandler) definition(uri DocumentURI, params *DocumentDefinitionParams) ([]Location, error) {
	f, ok := h.files[uri]
	if !ok {
		return nil, fmt.Errorf("document not found: %v", uri)
	}

	f.Analyze(h.addIncludedFile)
	pt := positionToPt(f, params.Position)
	if int(pt.Pos) < len(f.Text) {
		if c := f.Text[pt.Pos]; bytes.IndexByte(noadjust, c) < 0 {
			pt.Pos++
		}
	}
	x := f.TokenIndexAt(pt.Pos)
	if x < 0 {
		return nil, nil
	}

	var nm *core.Named
	token := f.Tokens[x]
	if x >= 2 && file.IsDotToken(f.Tokens[x-1]) {
		if env := resolveEnv(f, x-1); env != nil {
			nm = env.Find(core.Intern(file.TokenString(token)))
		}
	} else {
		nm = lookupNamed(file.EnvAt(token), f.Env, token)
	}
	if nm == nil {
		return nil, nil
	}

	found := nm.Token
	if found == nil || found.From == core.InternalParser {
		return nil, nil
	}
	loc := Location{
		URI:   toURI(found.From.Path),
		Range: tokenToRange(found),
	}
	return []Location{loc}, nil
}

func (h *langHandler) didChangeConfiguration(config *Config) (any, error) {
	h.mu.Lock()
	h.compileOptions.IncPaths = file.AdjustIncPaths(config.IncPaths)
	h.compileOptions.Defs = config.Defs
	h.mu.Unlock()
	for uri := range h.files {
		h.lintRequest(uri, eventTypeChange)
	}
	return nil, nil
}

func (h *langHandler) didChangeWorkspaceFolders(params *DidChangeWorkspaceFoldersParams) (result any, err error) {
	folders := h.folders
	for _, removed := range params.Event.Removed {
		folders = slices.DeleteFunc(folders, func(i string) bool {
			return toURI(i) == removed.URI
		})
	}
	for _, added := range params.Event.Added {
		x := slices.IndexFunc(folders, func(i string) bool {
			return toURI(i) == added.URI
		})
		if x < 0 {
			if folder, err := fromURI(added.URI); err == nil {
				folders = append(folders, folder)
			}
		}
	}
	h.folders = folders
	return nil, nil
}

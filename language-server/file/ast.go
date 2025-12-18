package file

import (
	"fmt"
	"ocala/core"
	"strings"
)

func SigDetail(sig *core.Sig) string {
	s := []string{}
	if sig != nil {
		if sig.IsInline {
			s = []string{"-*"}
		}
		if len(sig.Required) > 0 {
			s = append(s, JoinStringer(sig.Required, " "))
		}
		if len(sig.Results) > 0 {
			s = append(s, "=> "+JoinStringer(sig.Results, " "))
		}
		if len(sig.Invalidated) > 0 {
			if sig.Invalidated[0] == core.KwUNDER {
				s = append(s, "!")
			} else {
				s = append(s, "! "+JoinStringer(sig.Invalidated, " "))
			}
		}
	}
	return strings.Join(s, " ")
}

func JoinStringer[T fmt.Stringer](s []T, sep string) string {
	t := make([]string, 0, len(s))
	for _, i := range s {
		t = append(t, i.String())
	}
	return strings.Join(t, sep)
}

type Node interface {
	core.Value
	Base() *NodeBase
	Parent() Node
	Children() []Node
}

type NodeBase struct {
	Token *core.Token
	Name  *core.Identifier
	Env   *core.Env

	parent   Node
	children []Node
}

func (node *NodeBase) Dup() core.Value {
	return core.CopyPtr(node)
}

func (node *NodeBase) Inspect() string {
	r := ""
	for _, i := range node.children {
		r += i.Inspect()
	}
	return "<Node:" + r + ">"
}

func (node *NodeBase) Base() *NodeBase {
	return node
}

func (node *NodeBase) Parent() Node {
	return node.parent
}

func (node *NodeBase) Children() []Node {
	return node.children
}

func (v *NodeBase) AddChildNode(node Node) {
	v.children = append(v.children, node)
}

type MacroNode struct {
	NodeBase
	Args []*core.Identifier
	Vars []*core.Identifier
}

func (v *MacroNode) AddArg(arg core.Value) {
	if id, ok := arg.(*core.Identifier); ok {
		v.Args = append(v.Args, id)
	}
}

func (v *MacroNode) AddVar(arg core.Value) {
	if id, ok := arg.(*core.Identifier); ok {
		v.Vars = append(v.Vars, id)
	}
}

func (v *MacroNode) Detail() string {
	args := JoinStringer(v.Args, ", ")
	return fmt.Sprintf("%s(%s)", v.Name, args)
}

type ProcNode struct {
	NodeBase
	Sig *core.Sig
}

func (v *ProcNode) SetSignature(sig core.Value) {
	if sig, ok := sig.(*core.Sig); ok {
		v.Sig = sig
	}
}

func (v *ProcNode) Detail() string {
	return fmt.Sprintf("%s(%s)", v.Name, SigDetail(v.Sig))
}

type ConstNode struct {
	NodeBase
}

type ConstFnNode struct {
	NodeBase
	Args []*core.Identifier
}

func (v *ConstFnNode) AddArg(arg core.Value) {
	if id, ok := arg.(*core.Identifier); ok {
		v.Args = append(v.Args, id)
	}
}

func (v *ConstFnNode) Detail() string {
	args := JoinStringer(v.Args, ", ")
	return fmt.Sprintf("%s(%s)", v.Name, args)
}

type DataNode struct {
	NodeBase
	Type *core.Identifier
}

type ModuleNode struct {
	NodeBase
}

type LabelNode struct {
	NodeBase
}

type StructNode struct {
	NodeBase
}

type StructFieldNode struct {
	NodeBase
	Type *core.Identifier
}

type BlockNode struct {
	NodeBase
}

func inspectNode[T Node](node T) string {
	r := []string{}
	for _, i := range node.Children() {
		for _, j := range strings.Split(i.Inspect(), "\n") {
			r = append(r, "    "+j)
		}
	}
	return fmt.Sprintf("<%T:\n%s>", node, strings.Join(r, "\n"))
}

func (node *MacroNode) Inspect() string       { return inspectNode(node) }
func (node *ProcNode) Inspect() string        { return inspectNode(node) }
func (node *ConstNode) Inspect() string       { return inspectNode(node) }
func (node *ConstFnNode) Inspect() string     { return inspectNode(node) }
func (node *DataNode) Inspect() string        { return inspectNode(node) }
func (node *ModuleNode) Inspect() string      { return inspectNode(node) }
func (node *LabelNode) Inspect() string       { return inspectNode(node) }
func (node *StructNode) Inspect() string      { return inspectNode(node) }
func (node *StructFieldNode) Inspect() string { return inspectNode(node) }
func (node *BlockNode) Inspect() string       { return inspectNode(node) }
